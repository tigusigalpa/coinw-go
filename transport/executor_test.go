package transport

import (
	"context"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"

	coinwerr "github.com/tigusigalpa/coinw-go/errors"
)

type fakeDoer struct {
	handler func(req *http.Request, body []byte) *http.Response
	calls   int32
}

func (f *fakeDoer) Do(req *http.Request) (*http.Response, error) {
	atomic.AddInt32(&f.calls, 1)
	var body []byte
	if req.Body != nil {
		body, _ = io.ReadAll(req.Body)
	}
	return f.handler(req, body), nil
}

type failOnceDoer struct {
	calls int32
}

func (f *failOnceDoer) Do(req *http.Request) (*http.Response, error) {
	n := atomic.AddInt32(&f.calls, 1)
	if n == 1 {
		return nil, errors.New("simulated transport failure")
	}
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(`{"code":"200","data":{}}`)),
	}, nil
}

func okJSON(body string) *http.Response {
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

type fixedClock struct{}

func (fixedClock) NowMs() int64 { return 1700000000123 }

func newExec(doer HTTPDoer) *Executor {
	e, err := NewExecutor(Config{
		BaseURL:          "https://api.coinw.com",
		SpotAPIKey:       "SPOT_KEY",
		SpotSecretKey:    "SPOT_SECRET",
		FuturesAPIKey:    "FUT_KEY",
		FuturesSecretKey: "FUT_SECRET",
		UID:              "12345",
		HTTPClient:       doer,
		Clock:            fixedClock{},
		MaxReadRetries:   2,
	})
	if err != nil {
		panic(err)
	}
	return e
}

func TestSpotPublicGet(t *testing.T) {
	var got *http.Request
	doer := &fakeDoer{handler: func(req *http.Request, _ []byte) *http.Response {
		got = req
		return okJSON(`{"code":"200","data":{"x":1}}`)
	}}
	r, err := newExec(doer).Execute(context.Background(),
		"spot.market.getRecentTrades", map[string]interface{}{"symbol": "BTC_USDT"})
	if err != nil {
		t.Fatal(err)
	}
	if !r.IsSuccess() {
		t.Fatal("expected success result")
	}
	if got.Method != http.MethodGet {
		t.Fatalf("method %s", got.Method)
	}
	if !strings.Contains(got.URL.String(), "command=returnTradeHistory") ||
		!strings.Contains(got.URL.String(), "symbol=BTC_USDT") {
		t.Fatalf("bad url: %s", got.URL)
	}
	if strings.Contains(got.URL.String(), "sign=") {
		t.Fatal("public request must not be signed")
	}
}

func TestSpotPrivateSignedURL(t *testing.T) {
	var got *http.Request
	doer := &fakeDoer{handler: func(req *http.Request, _ []byte) *http.Response {
		got = req
		return okJSON(`{"code":"200","data":{}}`)
	}}
	_, err := newExec(doer).Execute(context.Background(),
		"spot.account.getSpotAccountBalance", nil)
	if err != nil {
		t.Fatal(err)
	}
	u := got.URL.String()
	if !strings.Contains(u, "command=returnBalances") ||
		!strings.Contains(u, "api_key=SPOT_KEY") {
		t.Fatalf("bad url: %s", u)
	}
	if !strings.Contains(u, "&sign=") {
		t.Fatal("missing sign")
	}
	if strings.Contains(u, "SPOT_SECRET") {
		t.Fatal("secret leaked into URL")
	}
}

func TestFuturesPrivateHeaders(t *testing.T) {
	var got *http.Request
	doer := &fakeDoer{handler: func(req *http.Request, _ []byte) *http.Response {
		got = req
		return okJSON(`{"code":0,"data":{}}`)
	}}
	_, err := newExec(doer).Execute(context.Background(),
		"futures.assets.getFuturesAccountAssets", nil)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got.URL.String(), "/v1/perpum/account/getUserAssets") {
		t.Fatalf("bad url: %s", got.URL)
	}
	if got.Header.Get("api_key") != "FUT_KEY" ||
		got.Header.Get("timestamp") != "1700000000123" {
		t.Fatalf("bad headers: %v", got.Header)
	}
	if _, err := base64.StdEncoding.DecodeString(got.Header.Get("sign")); err != nil {
		t.Fatal("sign must be base64")
	}
}

func TestFuturesPostJSONBody(t *testing.T) {
	var body []byte
	doer := &fakeDoer{handler: func(req *http.Request, b []byte) *http.Response {
		body = b
		return okJSON(`{"code":0,"data":{}}`)
	}}
	_, err := newExec(doer).Execute(context.Background(), "futures.trade.placeAnOrder",
		map[string]interface{}{"instrument": "BTC_USDT", "direction": "long"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(body), `"instrument":"BTC_USDT"`) {
		t.Fatalf("bad body: %s", body)
	}
}

func TestEmptyBodyIsJSONObject(t *testing.T) {
	var body []byte
	doer := &fakeDoer{handler: func(req *http.Request, b []byte) *http.Response {
		body = b
		return okJSON(`{"code":0,"data":{}}`)
	}}
	_, err := newExec(doer).Execute(context.Background(),
		"common.affiliate.getInvitedUsersTradingVolume", nil)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != "{}" {
		t.Fatalf("documented empty body must be {}, got %q", body)
	}
}

func TestMissingCredentials(t *testing.T) {
	e, _ := NewExecutor(Config{BaseURL: "https://api.coinw.com", HTTPClient: &fakeDoer{
		handler: func(*http.Request, []byte) *http.Response { return okJSON(`{"code":"200"}`) },
	}})
	_, err := e.Execute(context.Background(), "spot.account.getSpotAccountBalance", nil)
	var authErr *coinwerr.AuthenticationError
	if !errors.As(err, &authErr) {
		t.Fatalf("expected AuthenticationError, got %v", err)
	}
}

func TestAmbiguousGateBlocksRequest(t *testing.T) {
	doer := &fakeDoer{handler: func(*http.Request, []byte) *http.Response {
		return okJSON(`{"code":"200"}`)
	}}
	_, err := newExec(doer).Execute(context.Background(),
		"common.subAccount.getMaxTransferableAmount", map[string]interface{}{"coinCode": "BTC"})
	var amb *coinwerr.AmbiguousDocumentationError
	if !errors.As(err, &amb) {
		t.Fatalf("expected AmbiguousDocumentationError, got %v", err)
	}
	if atomic.LoadInt32(&doer.calls) != 0 {
		t.Fatal("gated endpoint must not send any request")
	}
}

func TestReadOnlyRetriesOnTransportError(t *testing.T) {
	doer := &failOnceDoer{}
	_, err := newExec(doer).Execute(context.Background(),
		"spot.market.getRecentTrades", map[string]interface{}{"symbol": "BTC_USDT"})
	if err != nil {
		t.Fatal(err)
	}
	if atomic.LoadInt32(&doer.calls) != 2 {
		t.Fatalf("expected 2 attempts, got %d", doer.calls)
	}
}

func TestMutationNeverRetried(t *testing.T) {
	doer := &failOnceDoer{}
	_, err := newExec(doer).Execute(context.Background(),
		"spot.trade.placeOrder", map[string]interface{}{"symbol": "BTC_USDT"})
	var te *coinwerr.TransportError
	if !errors.As(err, &te) {
		t.Fatalf("expected TransportError, got %v", err)
	}
	if atomic.LoadInt32(&doer.calls) != 1 {
		t.Fatalf("mutation must not be retried, got %d calls", doer.calls)
	}
}

func TestBusinessErrorMapping(t *testing.T) {
	execOnce := func(body string) error {
		doer := &fakeDoer{handler: func(*http.Request, []byte) *http.Response {
			return okJSON(body)
		}}
		_, err := newExec(doer).Execute(context.Background(),
			"spot.market.getRecentTrades", map[string]interface{}{"symbol": "X"})
		return err
	}

	var rl *coinwerr.RateLimitError
	if err := execOnce(`{"code":"29001","msg":"frequent"}`); !errors.As(err, &rl) {
		t.Errorf("29001 → expected RateLimitError, got %v", err)
	}
	var au *coinwerr.AuthenticationError
	if err := execOnce(`{"code":"6001","msg":"sig"}`); !errors.As(err, &au) {
		t.Errorf("6001 → expected AuthenticationError, got %v", err)
	}
	var api *coinwerr.APIError
	if err := execOnce(`{"code":"-200","msg":"biz"}`); !errors.As(err, &api) {
		t.Errorf("-200 → expected APIError, got %v", err)
	}
}
