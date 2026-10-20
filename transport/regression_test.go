package transport

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/tigusigalpa/coinw-go/endpoints"
	coinwerr "github.com/tigusigalpa/coinw-go/errors"
)

func TestRequestTimeoutOverridesLongerHTTPClientTimeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { <-r.Context().Done() }))
	defer srv.Close()
	e, err := NewExecutor(Config{BaseURL: srv.URL, HTTPClient: &http.Client{Timeout: 5 * time.Second}, RequestTimeout: 30 * time.Millisecond})
	if err != nil {
		t.Fatal(err)
	}
	start := time.Now()
	_, err = e.Execute(context.Background(), "spot.market.getRecentTrades", nil)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected timeout, got %v", err)
	}
	if time.Since(start) > time.Second {
		t.Fatal("per-attempt timeout was ignored")
	}
}

func TestDecodeDataPreservesNumericPrecision(t *testing.T) {
	r := Result{Data: json.RawMessage(`{"id":9007199254740993,"amount":0.1234567890123456789}`)}
	var data map[string]interface{}
	if err := r.DecodeData(&data); err != nil {
		t.Fatal(err)
	}
	if data["id"] != json.Number("9007199254740993") || data["amount"] != json.Number("0.1234567890123456789") {
		t.Fatalf("numeric precision lost: %#v", data)
	}
}

func TestInjectedRetryWaitHonorsCancellationAndSafety(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	d := &failOnceDoer{}
	waits := 0
	e, err := NewExecutor(Config{BaseURL: "https://fixture.test", HTTPClient: d, MaxReadRetries: 3, RetryBackoff: func(int) time.Duration { return time.Hour }, RetrySleep: func(ctx context.Context, delay time.Duration) error {
		waits++
		if delay != 5*time.Second {
			t.Error("retry delay escaped its bound")
		}
		cancel()
		return ctx.Err()
	}})
	if err != nil {
		t.Fatal(err)
	}
	_, err = e.Execute(ctx, "spot.market.getInstrumentInformation", nil)
	if !errors.Is(err, context.Canceled) || waits != 1 || d.calls != 1 {
		t.Fatal("retry wait ignored cancellation")
	}
}

func TestSerializationErrorIsTypedAndDoesNotReachNetwork(t *testing.T) {
	d := &failOnceDoer{}
	e, err := NewExecutor(Config{BaseURL: "https://fixture.test", HTTPClient: d, FuturesAPIKey: "fixture-key", FuturesSecretKey: "fixture-secret"})
	if err != nil {
		t.Fatal(err)
	}
	_, err = e.Execute(context.Background(), "futures.trade.placeAnOrder", map[string]interface{}{"invalid": func() {}})
	var target *coinwerr.SerializationError
	if !errors.As(err, &target) || d.calls != 0 {
		t.Fatalf("invalid JSON reached network or lost serialization class: %v", err)
	}
}

func TestPostReadIsNotRetried(t *testing.T) {
	d := &failOnceDoer{}
	_, err := newExec(d).Execute(context.Background(), "spot.account.getSpotAccountBalance", nil)
	if err == nil || d.calls != 1 {
		t.Fatalf("POST must have exactly one attempt, calls=%d err=%v", d.calls, err)
	}
}

func TestResponseEnvelopeValidation(t *testing.T) {
	ep, _ := endpoints.Get("spot.market.getRecentTrades")
	for _, body := range []string{`null`, `[]`, `{}`, `{"code":null}`, `{"code":{}}`, `invalid`} {
		_, err := validateResponse(okJSON(body), []byte(body), ep)
		var protocol *coinwerr.ProtocolError
		if !errors.As(err, &protocol) {
			t.Errorf("%s: expected protocol error, got %v", body, err)
		}
	}
	body := `{"code":"200","success":null,"errorCode":null,"data":{}}`
	if _, err := validateResponse(okJSON(body), []byte(body), ep); err != nil {
		t.Fatal(err)
	}
	resp := okJSON(`{"code":"200"}`)
	resp.Header = http.Header{"Content-Type": {"text/html"}}
	if _, err := validateResponse(resp, []byte(`{"code":"200"}`), ep); err == nil {
		t.Fatal("expected JSON content type validation")
	}
}

func TestHTTPErrorDoesNotEchoPrivateBody(t *testing.T) {
	ep, _ := endpoints.Get("spot.account.getSpotAccountBalance")
	resp := okJSON("")
	resp.StatusCode = 500
	_, err := validateResponse(resp, []byte("arbitrary-private-payload"), ep)
	if err == nil || strings.Contains(err.Error(), "arbitrary-private-payload") {
		t.Fatal("expected a safe HTTP error")
	}
}
