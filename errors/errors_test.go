package errors

import (
	stderrors "errors"
	"net/url"
	"strings"
	"testing"
)

func TestRedact(t *testing.T) {
	for _, input := range []string{
		`https://example.test?api_key=PRIVATE&sign=SIGNATURE&symbol=BTC`,
		`{"api_key": "PRIVATE", "passphrase":"SEC\"RET", "token": "TOKEN"}`,
		"Authorization: Bearer PRIVATE\r\nX-Request-ID: 123",
		`subAccount_apiKey=PRIVATE&secret_key=SECRET&token=TOKEN`,
	} {
		t.Run(input, func(t *testing.T) {
			got := Redact(input)
			for _, secret := range []string{"PRIVATE", "SIGNATURE", `SEC\"RET`, "SECRET", "TOKEN"} {
				if strings.Contains(got, secret) {
					t.Fatalf("secret was not redacted: %s", got)
				}
			}
			if Redact(got) != got {
				t.Fatalf("redaction must be idempotent: %s", got)
			}
		})
	}
}

func TestTransportErrorRedactsURLAndUnwraps(t *testing.T) {
	cause := &url.Error{Op: "Post", URL: "https://example.test?api_key=PRIVATE&sign=SIGNATURE&address=PRIVATE_DESTINATION&amount=123.456789", Err: stderrors.New("failed")}
	err := &TransportError{Message: "request failed", Err: cause}
	if strings.Contains(err.Error(), "PRIVATE") || strings.Contains(err.Error(), "SIGNATURE") {
		t.Fatal("credentials leaked")
	}
	if strings.Contains(err.Error(), "123.456789") {
		t.Fatal("private request parameters leaked")
	}
	if !stderrors.Is(err, cause) {
		t.Fatal("cause must be preserved for errors.Is/As")
	}
	for _, err := range []error{
		&AuthenticationError{Message: `api_key=PRIVATE`},
		&APIError{Message: `{"token": "PRIVATE"}`},
		&RateLimitError{Message: `sign=PRIVATE`},
	} {
		if strings.Contains(err.Error(), "PRIVATE") {
			t.Fatal("error message leaked credentials")
		}
	}
}
