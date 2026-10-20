package coinw

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

type tokenTransport struct{ body string }

func (r tokenTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(r.body)), Header: http.Header{}, Request: req}, nil
}

func TestPublicTokenValidation(t *testing.T) {
	for _, tc := range []struct {
		body  string
		valid bool
	}{
		{`{"data":{"token":"PUBLIC_PLACEHOLDER"}}`, true},
		{`{"code":"200","data":{"token":"PUBLIC_PLACEHOLDER"}}`, true},
		{`{"code":"6001","data":{"token":"PUBLIC_PLACEHOLDER"}}`, false},
		{`{"code":"200","data":{}}`, false},
		{`null`, false},
		{`{"data":{"token":"PUBLIC_PLACEHOLDER"}} {}`, false},
	} {
		client, err := NewClient(WithHTTPClient(&http.Client{Transport: tokenTransport{body: tc.body}}))
		if err != nil {
			t.Fatal(err)
		}
		token, _, err := client.fetchPublicToken(context.Background())
		if (err == nil) != tc.valid {
			t.Errorf("valid=%t: unexpected token validation result: %v", tc.valid, err)
		}
		if tc.valid && token != "PUBLIC_PLACEHOLDER" {
			t.Fatal("missing public token")
		}
	}
}
