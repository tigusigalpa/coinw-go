// Package coinw is the entry point of the CoinW Go SDK.
//
// The SDK covers the entire documented surface of the CoinW API
// (https://www.coinw.com/api-doc/en, sitemap snapshot 2026-10-20):
// Spot & Futures REST, Common (accounts/sub-accounts/affiliate), and all
// five documented WebSocket variants.
//
// Two isolated authentication profiles are implemented exactly per docs —
// Spot MD5 and Futures HMAC-SHA256 — and never mixed. Endpoints whose
// documentation contradicts itself are registered but gated behind
// errors.AmbiguousDocumentationError instead of being guessed.
package coinw

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	coinwerr "github.com/tigusigalpa/coinw-go/errors"
	"github.com/tigusigalpa/coinw-go/ratelimit"
	"github.com/tigusigalpa/coinw-go/services"
	"github.com/tigusigalpa/coinw-go/transport"
	ws "github.com/tigusigalpa/coinw-go/websocket"
)

// Documented endpoints.
const (
	DefaultRestBaseURL        = "https://api.coinw.com"
	DefaultSpotWSURL          = "wss://ws.futurescw.com"
	DefaultFuturesWSURL       = "wss://ws.futurescw.com/perpum"
	DefaultSpotPublicTokenURL = "https://www.coinw.com/pusher/public-token"
	DefaultSpotSocketIOWSURL  = "wss://ws.futurescw.info"
)

// Client groups all REST services and WS session factories.
type Client struct {
	cfg    ClientConfig
	exec   *transport.Executor
	httpDo *http.Client

	spotMarket  *services.SpotMarketService
	spotAccount *services.SpotAccountService
	spotCheck   *services.SpotCheckService
	spotTrade   *services.SpotTradeService

	futuresMarket   *services.FuturesMarketService
	futuresAssets   *services.FuturesAssetsService
	futuresCommon   *services.FuturesCommonService
	futuresCheck    *services.FuturesCheckService
	futuresPosition *services.FuturesPositionService
	futuresTrade    *services.FuturesTradeService

	account     *services.AccountService
	subAccounts *services.SubAccountService
	affiliate   *services.AffiliateService
}

// ClientConfig holds immutable client settings, populated via ClientOption.
type ClientConfig struct {
	RestBaseURL        string
	SpotWSURL          string
	FuturesWSURL       string
	SpotPublicTokenURL string
	SpotSocketIOWSURL  string

	SpotAPIKey       string
	SpotSecretKey    string
	FuturesAPIKey    string
	FuturesSecretKey string
	UID              string

	HTTPClient     *http.Client
	RequestTimeout time.Duration

	Limiter        ratelimit.Limiter // nil → default in-memory limiter; use DisableRateLimiter() to turn off
	LimiterSet     bool
	Clock          transport.Clock
	Dialer         ws.Dialer
	MaxReadRetries int
	RetryBackoff   func(int) time.Duration
	RetrySleep     func(context.Context, time.Duration) error
}

// ClientOption mutates ClientConfig during NewClient.
type ClientOption func(*ClientConfig)

// WithSpotCredentials sets the Spot private profile keys.
func WithSpotCredentials(apiKey, secretKey string) ClientOption {
	return func(c *ClientConfig) { c.SpotAPIKey, c.SpotSecretKey = apiKey, secretKey }
}

// WithFuturesCredentials sets the Futures private profile keys.
func WithFuturesCredentials(apiKey, secretKey string) ClientOption {
	return func(c *ClientConfig) { c.FuturesAPIKey, c.FuturesSecretKey = apiKey, secretKey }
}

// WithUID enables the documented UID-scoped rate-limit buckets.
func WithUID(uid string) ClientOption {
	return func(c *ClientConfig) { c.UID = uid }
}

// WithRestBaseURL overrides the REST base (testing/proxy).
func WithRestBaseURL(u string) ClientOption {
	return func(c *ClientConfig) { c.RestBaseURL = u }
}

// WithHTTPClient injects a custom *http.Client (timeouts, proxies, tracing).
func WithHTTPClient(h *http.Client) ClientOption {
	return func(c *ClientConfig) { c.HTTPClient = h }
}

// WithRequestTimeout sets the per-request transport timeout.
func WithRequestTimeout(d time.Duration) ClientOption {
	return func(c *ClientConfig) { c.RequestTimeout = d }
}

// WithRateLimiter replaces the default limiter with a custom one.
func WithRateLimiter(l ratelimit.Limiter) ClientOption {
	return func(c *ClientConfig) { c.Limiter, c.LimiterSet = l, true }
}

// DisableRateLimiter turns throttling off entirely.
func DisableRateLimiter() ClientOption {
	return func(c *ClientConfig) { c.Limiter, c.LimiterSet = nil, true }
}

// WithClock injects a signing clock (tests).
func WithClock(cl transport.Clock) ClientOption {
	return func(c *ClientConfig) { c.Clock = cl }
}

// WithDialer injects a custom WS dialer for all sessions created afterwards.
func WithDialer(d ws.Dialer) ClientOption {
	return func(c *ClientConfig) { c.Dialer = d }
}

// WithMaxReadRetries bounds automatic retries of read-only endpoints.
// Mutations are never retried.
func WithMaxReadRetries(n int) ClientOption {
	return func(c *ClientConfig) { c.MaxReadRetries = n }
}

// WithRetryPolicy injects concurrency-safe deterministic backoff/sleep hooks.
// Delays are bounded to 0..5 seconds; the endpoint retry safety rule still applies.
func WithRetryPolicy(backoff func(int) time.Duration, sleep func(context.Context, time.Duration) error) ClientOption {
	return func(c *ClientConfig) { c.RetryBackoff, c.RetrySleep = backoff, sleep }
}

// WithWebSocketURLs overrides documented WS endpoints (testing/proxy).
func WithWebSocketURLs(spot, futures, socketIO, tokenURL string) ClientOption {
	return func(c *ClientConfig) {
		if spot != "" {
			c.SpotWSURL = spot
		}
		if futures != "" {
			c.FuturesWSURL = futures
		}
		if socketIO != "" {
			c.SpotSocketIOWSURL = socketIO
		}
		if tokenURL != "" {
			c.SpotPublicTokenURL = tokenURL
		}
	}
}

// NewClient builds the SDK entry point.
func NewClient(opts ...ClientOption) (*Client, error) {
	cfg := ClientConfig{
		RestBaseURL:        DefaultRestBaseURL,
		SpotWSURL:          DefaultSpotWSURL,
		FuturesWSURL:       DefaultFuturesWSURL,
		SpotPublicTokenURL: DefaultSpotPublicTokenURL,
		SpotSocketIOWSURL:  DefaultSpotSocketIOWSURL,
		RequestTimeout:     15 * time.Second,
		MaxReadRetries:     2,
	}
	for _, o := range opts {
		o(&cfg)
	}

	httpDo := cfg.HTTPClient
	if httpDo == nil {
		httpDo = &http.Client{Timeout: cfg.RequestTimeout}
	} else {
		clone := *httpDo
		if clone.Timeout == 0 && cfg.RequestTimeout > 0 {
			clone.Timeout = cfg.RequestTimeout
		}
		httpDo = &clone
	}

	limiter := cfg.Limiter
	if !cfg.LimiterSet {
		limiter = ratelimit.NewInMemoryLimiter()
	}

	exec, err := transport.NewExecutor(transport.Config{
		BaseURL:          cfg.RestBaseURL,
		SpotAPIKey:       cfg.SpotAPIKey,
		SpotSecretKey:    cfg.SpotSecretKey,
		FuturesAPIKey:    cfg.FuturesAPIKey,
		FuturesSecretKey: cfg.FuturesSecretKey,
		UID:              cfg.UID,
		HTTPClient:       httpDo,
		Limiter:          limiter,
		Clock:            cfg.Clock,
		MaxReadRetries:   cfg.MaxReadRetries,
		RetryBackoff:     cfg.RetryBackoff,
		RetrySleep:       cfg.RetrySleep,
		RequestTimeout:   cfg.RequestTimeout,
	})
	if err != nil {
		return nil, err
	}

	c := &Client{cfg: cfg, exec: exec, httpDo: httpDo}
	c.spotMarket = services.NewSpotMarketService(exec)
	c.spotAccount = services.NewSpotAccountService(exec)
	c.spotCheck = services.NewSpotCheckService(exec)
	c.spotTrade = services.NewSpotTradeService(exec)
	c.futuresMarket = services.NewFuturesMarketService(exec)
	c.futuresAssets = services.NewFuturesAssetsService(exec)
	c.futuresCommon = services.NewFuturesCommonService(exec)
	c.futuresCheck = services.NewFuturesCheckService(exec)
	c.futuresPosition = services.NewFuturesPositionService(exec)
	c.futuresTrade = services.NewFuturesTradeService(exec)
	c.account = services.NewAccountService(exec)
	c.subAccounts = services.NewSubAccountService(exec)
	c.affiliate = services.NewAffiliateService(exec)
	return c, nil
}

// REST services.
func (c *Client) SpotMarket() *services.SpotMarketService       { return c.spotMarket }
func (c *Client) SpotAccount() *services.SpotAccountService     { return c.spotAccount }
func (c *Client) SpotCheck() *services.SpotCheckService         { return c.spotCheck }
func (c *Client) SpotTrade() *services.SpotTradeService         { return c.spotTrade }
func (c *Client) FuturesMarket() *services.FuturesMarketService { return c.futuresMarket }
func (c *Client) FuturesAssets() *services.FuturesAssetsService { return c.futuresAssets }
func (c *Client) FuturesCommon() *services.FuturesCommonService { return c.futuresCommon }
func (c *Client) FuturesCheck() *services.FuturesCheckService   { return c.futuresCheck }
func (c *Client) FuturesPosition() *services.FuturesPositionService {
	return c.futuresPosition
}
func (c *Client) FuturesTrade() *services.FuturesTradeService { return c.futuresTrade }
func (c *Client) Account() *services.AccountService           { return c.account }
func (c *Client) SubAccounts() *services.SubAccountService    { return c.subAccounts }
func (c *Client) Affiliate() *services.AffiliateService       { return c.affiliate }

// ---------------------------------------------------------------------------
// WebSocket session factories
// ---------------------------------------------------------------------------

// SpotPublicWS — Spot raw-JSON public stream (Method 2), no auth.
// Incremental depth requires DepthSequenceTracker + snapshot resync.
func (c *Client) SpotPublicWS(opts ...ws.RawOption) *ws.RawJsonSession {
	return ws.NewRawJsonSession(c.cfg.SpotWSURL, false, "", "",
		append([]ws.RawOption{ws.WithDialer(c.dialerOrDefault()), ws.WithProfile(ws.SpotProfile)}, opts...)...)
}

// SpotPrivateWS — same URL + documented login frame. Requires Spot creds.
func (c *Client) SpotPrivateWS(opts ...ws.RawOption) (*ws.RawJsonSession, error) {
	if c.cfg.SpotAPIKey == "" || c.cfg.SpotSecretKey == "" {
		return nil, &coinwerr.AuthenticationError{
			Message: "spot credentials are required for the private WebSocket stream",
		}
	}
	return ws.NewRawJsonSession(c.cfg.SpotWSURL, true, c.cfg.SpotAPIKey, c.cfg.SpotSecretKey,
		append([]ws.RawOption{ws.WithDialer(c.dialerOrDefault()), ws.WithProfile(ws.SpotProfile)}, opts...)...), nil
}

// SpotSocketIOWS — Spot public Method 1 (Socket.IO + token fetch).
func (c *Client) SpotSocketIOWS(opts ...ws.SocketIOOption) *ws.SocketIOSession {
	return ws.NewSocketIOSession(c.cfg.SpotSocketIOWSURL, c.fetchPublicToken,
		append([]ws.SocketIOOption{ws.WithSocketIODialer(c.dialerOrDefault())}, opts...)...)
}

// FuturesPublicWS — Futures public stream on /perpum.
func (c *Client) FuturesPublicWS(opts ...ws.RawOption) *ws.RawJsonSession {
	return ws.NewRawJsonSession(c.cfg.FuturesWSURL, false, "", "",
		append([]ws.RawOption{ws.WithDialer(c.dialerOrDefault()), ws.WithProfile(ws.FuturesProfile)}, opts...)...)
}

// FuturesPrivateWS — /perpum + documented login. Requires Futures creds.
func (c *Client) FuturesPrivateWS(opts ...ws.RawOption) (*ws.RawJsonSession, error) {
	if c.cfg.FuturesAPIKey == "" || c.cfg.FuturesSecretKey == "" {
		return nil, &coinwerr.AuthenticationError{
			Message: "futures credentials are required for the private WebSocket stream",
		}
	}
	return ws.NewRawJsonSession(c.cfg.FuturesWSURL, true, c.cfg.FuturesAPIKey, c.cfg.FuturesSecretKey,
		append([]ws.RawOption{ws.WithDialer(c.dialerOrDefault()), ws.WithProfile(ws.FuturesProfile)}, opts...)...), nil
}

func (c *Client) dialerOrDefault() ws.Dialer {
	if c.cfg.Dialer != nil {
		return c.cfg.Dialer
	}
	return ws.DefaultDialer()
}

// fetchPublicToken implements the documented token retrieval:
// GET {spotPublicTokenUrl} → {"code":"200","data":{"endpoint":…,"token":…}}.
func (c *Client) fetchPublicToken(ctx context.Context) (string, string, error) {
	if c.cfg.RequestTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.cfg.RequestTimeout)
		defer cancel()
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.cfg.SpotPublicTokenURL, nil)
	if err != nil {
		return "", "", err
	}
	resp, err := c.httpDo.Do(req)
	if err != nil {
		return "", "", &coinwerr.TransportError{Message: "failed to fetch public WS token", Err: err}
	}
	defer resp.Body.Close()
	var payload struct {
		Code json.RawMessage `json:"code"`
		Data struct {
			Token    string `json:"token"`
			Endpoint string `json:"endpoint"`
		} `json:"data"`
	}
	if resp.StatusCode != http.StatusOK {
		return "", "", &coinwerr.HTTPStatusError{StatusCode: resp.StatusCode, Message: "public WS token endpoint"}
	}
	raw, err := io.ReadAll(io.LimitReader(resp.Body, (1<<20)+1))
	if err != nil {
		return "", "", &coinwerr.TransportError{Message: "failed reading public WS token response", Err: err}
	}
	if len(raw) > 1<<20 || json.Unmarshal(raw, &payload) != nil {
		return "", "", &coinwerr.ProtocolError{Message: "malformed public WS token response"}
	}
	result := transport.Result{Code: payload.Code}
	if (len(payload.Code) > 0 && !result.IsSuccess()) || payload.Data.Token == "" {
		return "", "", &coinwerr.ProtocolError{Message: "public WS token response must contain data.token and no failure code"}
	}
	return payload.Data.Token, payload.Data.Endpoint, nil
}
