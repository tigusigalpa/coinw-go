// Package errors defines the CoinW SDK error taxonomy.
//
// Every failure raised by the SDK is one of the typed errors below (or wraps
// one). Use errors.As / errors.Is to classify. Error messages are redacted —
// credentials, signatures and tokens never appear in them.
package errors

import (
	"fmt"
	"strings"
)

// ValidationError — client-side input rejected before any network call.
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string { return e.Message }

// AuthenticationError — missing/invalid credentials or a documented auth
// business code (6000/6001/6003/-103/1004).
type AuthenticationError struct {
	Message string
}

func (e *AuthenticationError) Error() string { return e.Message }

// RateLimitError — HTTP 429 or documented code 29001.
type RateLimitError struct {
	Code    string
	Message string
}

func (e *RateLimitError) Error() string {
	return fmt.Sprintf("coinw rate limit (%s): %s", e.Code, e.Message)
}

// HTTPStatusError — non-2xx transport status.
type HTTPStatusError struct {
	StatusCode int
	Message    string
}

func (e *HTTPStatusError) Error() string {
	return fmt.Sprintf("coinw http %d: %s", e.StatusCode, e.Message)
}

// APIError — business-level failure code returned by the API envelope.
type APIError struct {
	EndpointID string
	Code       string
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("coinw api error %s at %s: %s", e.Code, e.EndpointID, e.Message)
}

// ProtocolError — malformed JSON envelope / WebSocket frame.
type ProtocolError struct {
	Message string
}

func (e *ProtocolError) Error() string { return e.Message }

// TransportError — connection/DNS/TLS/body-read failures and timeouts.
type TransportError struct {
	Message string
	Err     error
}

func (e *TransportError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *TransportError) Unwrap() error { return e.Err }

// CancelledError — a local cancellation bound was reached (reconnect
// attempts, cancellation token). Context cancellation surfaces ctx.Err()
// directly.
type CancelledError struct {
	Message string
}

func (e *CancelledError) Error() string { return e.Message }

// AmbiguousDocumentationError — the endpoint is registered but the official
// documentation contradicts itself, so no request is sent until verified.
type AmbiguousDocumentationError struct {
	EndpointID string
	Note       string
}

func (e *AmbiguousDocumentationError) Error() string {
	return fmt.Sprintf("coinw endpoint %s blocked by conflicting documentation: %s",
		e.EndpointID, e.Note)
}

// Redact strips credential/signature values from a string before it reaches
// logs or error messages.
func Redact(s string) string {
	for _, key := range []string{
		"api_key", "secret_key", "passphrase", "sign", "signature",
		"secretKey", "apiKey", "secret", "authorization", "token",
	} {
		for _, marker := range []string{key + "=", `"` + key + `":"`} {
			for {
				idx := strings.Index(s, marker)
				if idx < 0 {
					break
				}
				start := idx + len(marker)
				end := start
				for end < len(s) && !strings.ContainsRune("&\"' \t\n", rune(s[end])) {
					end++
				}
				s = s[:start] + "***" + s[end:]
			}
		}
	}
	return s
}
