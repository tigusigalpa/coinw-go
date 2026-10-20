package websocket

import (
	"strconv"

	coinwerr "github.com/tigusigalpa/coinw-go/errors"
)

// AsError classifies an explicitly rejected WS operation without returning
// its private payload. ReadMessage callers can also use this outside Run.
func (m *Message) AsError() error {
	if m == nil || !m.IsError() {
		return nil
	}
	message := m.ErrorMsg
	if message == "" {
		message = "WebSocket operation rejected"
	}
	code := ""
	if m.ErrorCode != 0 {
		code = strconv.Itoa(m.ErrorCode)
	}
	if m.IsLoginAck() {
		return &coinwerr.AuthenticationError{Message: message}
	}
	switch m.ErrorCode {
	case 6000, 6001, 6003, -103, 1004:
		return &coinwerr.AuthenticationError{Message: message}
	case 29001, 429:
		return &coinwerr.RateLimitError{Code: code, Message: message}
	default:
		return &coinwerr.APIError{EndpointID: "websocket", Code: code, Message: message}
	}
}
