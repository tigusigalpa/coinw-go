package transport

import (
	"context"
	"crypto/rand"
	"io"
	"math"
	"math/big"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/tigusigalpa/coinw-go/auth"
	"github.com/tigusigalpa/coinw-go/endpoints"
	coinwerr "github.com/tigusigalpa/coinw-go/errors"
	"github.com/tigusigalpa/coinw-go/ratelimit"
)

// HTTPDoer is satisfied by *http.Client and any test double.
type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Clock supplies signing timestamps; injectable for deterministic tests.
type Clock interface {
	NowMs() int64
}

type systemClock struct{}

func (systemClock) NowMs() int64 { return time.Now().UnixMilli() }

// Config carries the executor's immutable settings.
type Config struct {
	BaseURL          string
	SpotAPIKey       string
	SpotSecretKey    string
	FuturesAPIKey    string
	FuturesSecretKey string
	UID              string
	HTTPClient       HTTPDoer // required
	Limiter          ratelimit.Limiter
	Clock            Clock
	MaxReadRetries   int                                        // automatic retries for read-only endpoints only
	RequestTimeout   time.Duration                              // per-attempt timeout; 0 = no extra timeout
	RetryBackoff     func(attempt int) time.Duration            // nil uses bounded exponential backoff with jitter
	RetrySleep       func(context.Context, time.Duration) error // nil uses context-aware timers
}

// Executor runs the central request pipeline. Construct via NewExecutor.
type Executor struct {
	cfg Config
}

// NewExecutor builds the pipeline. HTTPClient is required.
func NewExecutor(cfg Config) (*Executor, error) {
	u, err := url.Parse(cfg.BaseURL)
	if err != nil || u.Host == "" || (u.Scheme != "http" && u.Scheme != "https") || u.RawQuery != "" || u.Fragment != "" || u.User != nil {
		return nil, &coinwerr.ValidationError{Message: "BaseURL must be an absolute HTTP(S) URL without credentials, query or fragment"}
	}
	if cfg.HTTPClient == nil {
		return nil, &coinwerr.ValidationError{
			Message: "an HTTPClient is required (e.g. http.DefaultClient or a tuned *http.Client)",
		}
	}
	if cfg.Clock == nil {
		cfg.Clock = systemClock{}
	}
	if cfg.RetryBackoff == nil {
		cfg.RetryBackoff = backoff
	}
	if cfg.RetrySleep == nil {
		cfg.RetrySleep = sleepCtx
	}
	if cfg.MaxReadRetries < 0 {
		cfg.MaxReadRetries = 0
	}
	if cfg.RequestTimeout < 0 {
		return nil, &coinwerr.ValidationError{Message: "RequestTimeout must not be negative"}
	}
	return &Executor{cfg: cfg}, nil
}

// Execute performs a registered endpoint call. Params:
//   - map[string]interface{} — params for all verb/auth combos
//   - a slice/struct — JSON body for futures non-GET (e.g. batch arrays)
//   - nil — no params (futures non-GET then sends documented "{}" body)
func (e *Executor) Execute(ctx context.Context, endpointID string, params interface{}) (*Result, error) {
	ep, ok := endpoints.Get(endpointID)
	if !ok {
		return nil, &coinwerr.ValidationError{Message: "unknown endpoint id: " + endpointID}
	}
	if ep.Ambiguous {
		return nil, &coinwerr.AmbiguousDocumentationError{
			EndpointID: ep.ID,
			Note:       orDefault(ep.AmbiguityNote, "conflicting official documentation"),
		}
	}

	maxAttempts := 1
	if ep.ReadOnly && ep.Method == http.MethodGet {
		maxAttempts += e.cfg.MaxReadRetries
	}

	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		if e.cfg.Limiter != nil {
			if err := e.cfg.Limiter.Acquire(ctx, ep, e.cfg.UID); err != nil {
				return nil, err
			}
		}

		result, err := e.attempt(ctx, ep, params)
		if err == nil {
			return result, nil
		}
		lastErr = err

		if !retryable(ep, err) || attempt == maxAttempts {
			return nil, err
		}
		delay := e.cfg.RetryBackoff(attempt)
		if delay < 0 {
			delay = 0
		}
		if delay > 5*time.Second {
			delay = 5 * time.Second
		}
		if err := e.cfg.RetrySleep(ctx, delay); err != nil {
			return nil, err
		}
	}
	return nil, lastErr
}

func (e *Executor) attempt(ctx context.Context, ep endpoints.Definition, params interface{}) (*Result, error) {
	if e.cfg.RequestTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, e.cfg.RequestTimeout)
		defer cancel()
	}
	req, err := e.buildRequest(ctx, ep, params)
	if err != nil {
		return nil, err
	}
	return e.do(req, ep)
}

func (e *Executor) do(req *http.Request, ep endpoints.Definition) (*Result, error) {
	resp, err := e.cfg.HTTPClient.Do(req)
	if err != nil {
		return nil, &coinwerr.TransportError{Message: "http transport error at " + ep.ID, Err: err}
	}
	defer resp.Body.Close()
	const maxBody = 16 << 20
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxBody+1))
	if err != nil {
		return nil, &coinwerr.TransportError{Message: "failed to read response body", Err: err}
	}
	if len(body) > maxBody {
		return nil, &coinwerr.ProtocolError{Message: "response body exceeds 16 MiB"}
	}
	return validateResponse(resp, body, ep)
}

// buildRequest constructs the *http.Request for the endpoint per its auth
// profile. Spot signs into the URL; futures signs headers (+JSON body).
func (e *Executor) buildRequest(ctx context.Context, ep endpoints.Definition, params interface{}) (*http.Request, error) {
	base := strings.TrimRight(e.cfg.BaseURL, "/")
	method := strings.ToUpper(ep.Method)

	switch ep.Auth {
	case endpoints.AuthPublic:
		u := base + ep.Path
		if method == http.MethodGet {
			if m, ok := params.(map[string]interface{}); ok && len(m) > 0 {
				sep := "?"
				if strings.Contains(ep.Path, "?") {
					sep = "&"
				}
				u += sep + queryString(m)
			}
		}
		return http.NewRequestWithContext(ctx, method, u, nil)

	case endpoints.AuthSpot:
		if e.cfg.SpotAPIKey == "" || e.cfg.SpotSecretKey == "" {
			return nil, &coinwerr.AuthenticationError{
				Message: "spot credentials are required for " + ep.ID + " (Spot private profile)",
			}
		}
		m, _ := params.(map[string]interface{})
		if m == nil {
			m = map[string]interface{}{}
		}
		u := auth.SpotSignedURL(base, ep.Path, m, e.cfg.SpotAPIKey, e.cfg.SpotSecretKey)
		var body io.Reader
		if method == http.MethodPost {
			body = strings.NewReader("")
		}
		req, err := http.NewRequestWithContext(ctx, method, u, body)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
		return req, nil

	case endpoints.AuthFutures:
		if e.cfg.FuturesAPIKey == "" || e.cfg.FuturesSecretKey == "" {
			return nil, &coinwerr.AuthenticationError{
				Message: "futures credentials are required for " + ep.ID + " (Futures private profile)",
			}
		}
		ts := e.cfg.Clock.NowMs()
		if method == http.MethodGet {
			m, _ := params.(map[string]interface{})
			query, hdr := auth.FuturesSignedGet(ep.Path, m, e.cfg.FuturesAPIKey, e.cfg.FuturesSecretKey, ts)
			u := base + ep.Path
			if query != "" {
				u += "?" + query
			}
			req, err := http.NewRequestWithContext(ctx, method, u, nil)
			if err != nil {
				return nil, err
			}
			applyFuturesHeaders(req, hdr, false)
			return req, nil
		}
		payload, hdr, err := auth.FuturesSignedBody(method, ep.Path, params,
			e.cfg.FuturesAPIKey, e.cfg.FuturesSecretKey, ts)
		if err != nil {
			return nil, &coinwerr.SerializationError{Message: "failed to encode futures body", Err: err}
		}
		req, err := http.NewRequestWithContext(ctx, method, base+ep.Path, strings.NewReader(payload))
		if err != nil {
			return nil, err
		}
		applyFuturesHeaders(req, hdr, true)
		return req, nil
	}
	return nil, &coinwerr.ValidationError{Message: "unknown auth profile: " + ep.Auth}
}

func applyFuturesHeaders(req *http.Request, h auth.FuturesHeaders, jsonBody bool) {
	req.Header.Set("sign", h.Sign)
	req.Header.Set("api_key", h.APIKey)
	req.Header.Set("timestamp", h.Timestamp)
	if jsonBody {
		req.Header.Set("Content-Type", "application/json")
	}
}

// queryString encodes params as RFC-3986 query preserving sorted key order.
func queryString(m map[string]interface{}) string {
	v := url.Values{}
	for k, val := range m {
		if val == nil {
			continue
		}
		v.Set(k, auth.StringifyParam(val))
	}
	return v.Encode()
}

// retryable decides whether an error may be retried automatically. Only
// read-only endpoints are ever retried; auth/validation/ambiguity/business
// errors are never retried.
func retryable(ep endpoints.Definition, err error) bool {
	if ep.IsMutation() || ep.Method != http.MethodGet {
		return false
	}
	switch err.(type) {
	case *coinwerr.TransportError, *coinwerr.RateLimitError:
		return true
	}
	return false
}

func backoff(attempt int) time.Duration {
	base := float64(200) * math.Pow(2, float64(attempt-1)) // 200ms, 400ms, 800ms…
	jitterN, _ := rand.Int(rand.Reader, big.NewInt(int64(math.Max(1, base*0.25))))
	d := time.Duration(base+float64(jitterN.Int64())) * time.Millisecond
	if d > 5*time.Second {
		d = 5 * time.Second
	}
	return d
}

func sleepCtx(ctx context.Context, d time.Duration) error {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}

func orDefault(s, def string) string {
	if s == "" {
		return def
	}
	return s
}
