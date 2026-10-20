package auth

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"strings"
	"testing"
)

func TestSpotCanonicalString(t *testing.T) {
	got := SpotCanonicalString(
		map[string]interface{}{"quantity": "1", "price": "100.5", "instrument": "BTC_USDT"},
		"TEST_API_KEY", "TEST_SECRET",
	)
	want := "api_key=TEST_API_KEY&instrument=BTC_USDT&price=100.5&quantity=1&secret_key=TEST_SECRET"
	if got != want {
		t.Fatalf("canonical mismatch:\n got %s\nwant %s", got, want)
	}
}

func TestSpotSignIsUppercaseMD5(t *testing.T) {
	canonical := "api_key=K&instrument=BTC_USDT&secret_key=S"
	sum := md5.Sum([]byte(canonical))
	want := strings.ToUpper(hex.EncodeToString(sum[:]))

	got := SpotSign(map[string]interface{}{"instrument": "BTC_USDT"}, "K", "S")
	if got != want {
		t.Fatalf("sign mismatch: got %s want %s", got, want)
	}
	if len(got) != 32 {
		t.Fatalf("sign must be 32 hex chars")
	}
}

func TestSpotSignedURL(t *testing.T) {
	u := SpotSignedURL("https://api.coinw.com", "/api/v1/private?command=doTrade",
		map[string]interface{}{"instrument": "BTC USDT", "price": "100.5"}, "K", "S")

	if !strings.HasPrefix(u, "https://api.coinw.com/api/v1/private?command=doTrade&sign=") {
		t.Fatalf("unexpected url prefix: %s", u)
	}
	if !strings.Contains(u, "api_key=K") {
		t.Fatalf("api_key missing: %s", u)
	}
	if !strings.Contains(u, "instrument=BTC+USDT") && !strings.Contains(u, "instrument=BTC%20USDT") {
		t.Fatalf("instrument encoding missing: %s", u)
	}
}

func TestFuturesSignedGet(t *testing.T) {
	const ts = int64(1700000000123)
	q, h := FuturesSignedGet("/v1/perpum/user/assets",
		map[string]interface{}{"a": "1", "b": "x"}, "KEY", "SECRET", ts)

	if q != "a=1&b=x" {
		t.Fatalf("query mismatch: %s", q)
	}
	payload := "1700000000123GET/v1/perpum/user/assets?a=1&b=x"
	mac := hmac.New(sha256.New, []byte("SECRET"))
	mac.Write([]byte(payload))
	want := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	if h.Sign != want {
		t.Fatalf("sign mismatch:\n got %s\nwant %s", h.Sign, want)
	}
	if h.APIKey != "KEY" || h.Timestamp != "1700000000123" {
		t.Fatalf("headers mismatch: %+v", h)
	}
}

func TestFuturesSignedGetEmptyQuery(t *testing.T) {
	const ts = int64(1700000000123)
	q, h := FuturesSignedGet("/v1/perpum/x", nil, "KEY", "SECRET", ts)
	if q != "" {
		t.Fatalf("expected empty query, got %s", q)
	}
	payload := "1700000000123GET/v1/perpum/x"
	mac := hmac.New(sha256.New, []byte("SECRET"))
	mac.Write([]byte(payload))
	if h.Sign != base64.StdEncoding.EncodeToString(mac.Sum(nil)) {
		t.Fatal("sign payload must not contain '?' when query empty")
	}
}

func TestFuturesSignedBody(t *testing.T) {
	const ts = int64(1700000000123)
	body, h, err := FuturesSignedBody("POST", "/v1/perpum/order",
		map[string]interface{}{"instrument": "BTC_USDT", "quantity": "1"}, "KEY", "SECRET", ts)
	if err != nil {
		t.Fatal(err)
	}
	if body != `{"instrument":"BTC_USDT","quantity":"1"}` {
		t.Fatalf("body mismatch: %s", body)
	}
	payload := "1700000000123POST/v1/perpum/order" + body
	mac := hmac.New(sha256.New, []byte("SECRET"))
	mac.Write([]byte(payload))
	if h.Sign != base64.StdEncoding.EncodeToString(mac.Sum(nil)) {
		t.Fatal("body sign mismatch")
	}
}

func TestFuturesEmptyBodyIsJSONObject(t *testing.T) {
	body, err := FuturesCanonicalBody(nil)
	if err != nil {
		t.Fatal(err)
	}
	if body != "{}" {
		t.Fatalf("documented empty body must be {}, got %s", body)
	}
}
