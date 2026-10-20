package auth

import (
	"bytes"
	"encoding/json"
	"net/url"
	"os"
	"testing"
)

func TestPublishedAlgorithmGoldenVectors(t *testing.T) {
	raw, err := os.ReadFile("testdata/golden.json")
	if err != nil {
		t.Fatal(err)
	}
	type vector struct {
		Path, Canonical, Signature, Payload, Query string
		Params, Body                               map[string]interface{}
	}
	var fixture struct {
		APIKey, SecretKey          string
		Timestamp                  int64
		Spot, Get, Post, EmptyPost vector
	}
	d := json.NewDecoder(bytes.NewReader(raw))
	d.UseNumber()
	if err = d.Decode(&fixture); err != nil {
		t.Fatal(err)
	}
	if got := SpotCanonicalString(fixture.Spot.Params, fixture.APIKey, fixture.SecretKey); got != fixture.Spot.Canonical {
		t.Fatal("Spot canonical bytes differ from golden")
	}
	if SpotSign(fixture.Spot.Params, fixture.APIKey, fixture.SecretKey) != fixture.Spot.Signature {
		t.Fatal("Spot MD5 golden mismatch")
	}
	query, headers := FuturesSignedGet(fixture.Get.Path, fixture.Get.Params, fixture.APIKey, fixture.SecretKey, fixture.Timestamp)
	if headers.Sign != fixture.Get.Signature {
		t.Fatal("Futures GET HMAC golden mismatch")
	}
	wire, err := url.ParseQuery(query)
	if err != nil || wire.Get("referralCode") != fixture.Get.Params["referralCode"] {
		t.Fatal("escaped wire query changed the signed values")
	}
	for _, v := range []vector{fixture.Post, fixture.EmptyPost} {
		payload, headers, err := FuturesSignedBody("POST", v.Path, v.Body, fixture.APIKey, fixture.SecretKey, fixture.Timestamp)
		if err != nil || payload != v.Payload || headers.Sign != v.Signature {
			t.Fatalf("Futures body golden mismatch at %s: %v", v.Path, err)
		}
	}
}
