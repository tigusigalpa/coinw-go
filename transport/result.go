// Package transport contains the request pipeline: ambiguity gate →
// rate limiter → signer → HTTP → retry (read-only only) → validation.
package transport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"strings"

	"github.com/tigusigalpa/coinw-go/endpoints"
	coinwerr "github.com/tigusigalpa/coinw-go/errors"
)

// Result is the original REST envelope retained by Response[T] and the
// low-level executor. Data and Envelope can contain private values; never log
// them. Public services validate endpoint-specific typed Data before returning.
type Result struct {
	Code       json.RawMessage // business code ("200", "0", 29001, …)
	Msg        string          // message field, when present
	Data       json.RawMessage // `data` field verbatim (may be null)
	Envelope   []byte          // complete decoded body verbatim
	HTTPStatus int
	EndpointID string
}

// IsSuccess mirrors the documented success codes across envelopes:
// spot "200", futures 0 (also "0"/200 on partner API).
func (r *Result) IsSuccess() bool {
	var s string
	if err := json.Unmarshal(r.Code, &s); err == nil {
		return s == "200" || s == "0"
	}
	var n int
	if err := json.Unmarshal(r.Code, &n); err == nil {
		return n == 200 || n == 0
	}
	return false
}

// DecodeData decodes data into v. Interface-valued numeric leaves use
// json.Number to preserve decimal precision and IDs larger than 2^53.
func (r *Result) DecodeData(v interface{}) error {
	if len(r.Data) == 0 {
		return nil
	}
	dec := json.NewDecoder(bytes.NewReader(r.Data))
	dec.UseNumber()
	return dec.Decode(v)
}

var authErrorCodes = map[string]bool{
	"6000": true, "6001": true, "6003": true, "-103": true, "1004": true,
}

var rateLimitCodes = map[string]bool{
	"29001": true, "429": true,
}

// validateResponse maps HTTP status → decoded envelope → business code to a
// Result or a typed error.
func validateResponse(resp *http.Response, body []byte, ep endpoints.Definition) (*Result, error) {
	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, &coinwerr.RateLimitError{Code: "429", Message: "Too Many Requests"}
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &coinwerr.HTTPStatusError{
			StatusCode: resp.StatusCode,
			Message:    "request failed at " + ep.ID,
		}
	}
	if ct := resp.Header.Get("Content-Type"); ct != "" {
		mediaType, _, err := mime.ParseMediaType(ct)
		if err != nil || (mediaType != "application/json" && !strings.HasSuffix(mediaType, "+json")) {
			return nil, &coinwerr.ProtocolError{Message: "expected a JSON Content-Type for " + ep.ID}
		}
	}

	var envelope map[string]json.RawMessage
	if err := json.Unmarshal(body, &envelope); err != nil || envelope == nil {
		return nil, &coinwerr.ProtocolError{
			Message: "expected a JSON object response for " + ep.ID,
		}
	}

	var codeRaw json.RawMessage
	if c, ok := envelope["code"]; ok {
		codeRaw = c
	}
	msg := readStringField(envelope, "msg")
	if msg == "" {
		msg = readStringField(envelope, "message")
	}
	codeStr := codeToString(codeRaw)
	errCodeStr := codeToString(envelope["errorCode"])
	if codeStr == "" && errCodeStr == "" {
		return nil, &coinwerr.ProtocolError{Message: "response missing a valid business code for " + ep.ID}
	}

	isSuccess := (codeStr == "200" || codeStr == "0") &&
		!rawBoolFalse(envelope["success"]) &&
		!rawBoolFalse(envelope["result"]) &&
		(len(envelope["errorCode"]) == 0 || bytes.Equal(bytes.TrimSpace(envelope["errorCode"]), []byte("null")))

	if !isSuccess {
		effCode := codeStr
		if errCodeStr != "" {
			effCode = errCodeStr
		}
		if effCode == "" {
			effCode = "unknown"
		}
		effMsg := msg
		if effMsg == "" {
			effMsg = readStringField(envelope, "errorMsg")
		}
		if effMsg == "" {
			effMsg = "unknown API error"
		}
		switch {
		case rateLimitCodes[effCode]:
			return nil, &coinwerr.RateLimitError{Code: effCode, Message: effMsg}
		case authErrorCodes[effCode]:
			return nil, &coinwerr.AuthenticationError{
				Message: fmt.Sprintf("coinw auth error %s: %s", effCode, effMsg),
			}
		default:
			return nil, &coinwerr.APIError{EndpointID: ep.ID, Code: effCode, Message: effMsg}
		}
	}

	return &Result{
		Code:       codeRaw,
		Msg:        msg,
		Data:       envelope["data"],
		Envelope:   body,
		HTTPStatus: resp.StatusCode,
		EndpointID: ep.ID,
	}, nil
}

func readStringField(env map[string]json.RawMessage, key string) string {
	raw, ok := env[key]
	if !ok {
		return ""
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return s
	}
	return ""
}

func rawBoolFalse(raw json.RawMessage) bool {
	return bytes.Equal(bytes.TrimSpace(raw), []byte("false"))
}

func codeToString(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return s
	}
	var n json.Number
	if err := json.Unmarshal(raw, &n); err == nil {
		if _, err := n.Int64(); err == nil {
			return n.String()
		}
	}
	return ""
}
