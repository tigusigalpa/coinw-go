package contracts_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	coinw "github.com/tigusigalpa/coinw-go"
	coinwerr "github.com/tigusigalpa/coinw-go/errors"
)

func TestEveryRESTErrorContract(t *testing.T) {
	for _, entry := range loadCatalogue(t).REST {
		if entry.Ambiguous {
			continue
		}
		t.Run(entry.ID, func(t *testing.T) {
			for _, tc := range []struct {
				name, file string
				status     int
			}{{"API", "fixtures/api_error.json", 200}, {"RateLimit", "fixtures/rate_limit.json", 200}, {"Malformed", "fixtures/malformed.json", 200}, {"HTTP429", "fixtures/rate_limit.json", 429}} {
				t.Run(tc.name, func(t *testing.T) {
					body, err := os.ReadFile(tc.file)
					if err != nil {
						t.Fatal(err)
					}
					calls := 0
					server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
						calls++
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(tc.status)
						_, _ = w.Write(body)
					}))
					defer server.Close()
					client, err := coinw.NewClient(coinw.WithRestBaseURL(server.URL), coinw.WithHTTPClient(server.Client()), coinw.WithMaxReadRetries(0), coinw.DisableRateLimiter(), coinw.WithSpotCredentials("fixture-key", "fixture-secret"), coinw.WithFuturesCredentials("fixture-key", "fixture-secret"))
					if err != nil {
						t.Fatal(err)
					}
					err = invokeContract(context.Background(), client, entry)
					switch tc.name {
					case "API":
						var target *coinwerr.APIError
						if !errors.As(err, &target) {
							t.Fatalf("wrong API error class: %v", err)
						}
					case "RateLimit", "HTTP429":
						var target *coinwerr.RateLimitError
						if !errors.As(err, &target) {
							t.Fatalf("wrong rate-limit class: %v", err)
						}
					case "Malformed":
						var target *coinwerr.ProtocolError
						if !errors.As(err, &target) {
							t.Fatalf("wrong decoder class: %v", err)
						}
					}
					if calls != 1 {
						t.Fatal("error contract unexpectedly retried")
					}
				})
			}
		})
	}
}
