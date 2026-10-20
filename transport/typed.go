package transport

import (
	coinwerr "github.com/tigusigalpa/coinw-go/errors"
	"github.com/tigusigalpa/coinw-go/model"
)

// Response carries validated endpoint-specific Data and the original envelope.
// Result.Data remains the original JSON; Response.Data is the typed payload.
type Response[T any] struct {
	*Result
	Data T
}

// DecodeResponse enforces the endpoint's complete payload contract.
func DecodeResponse[T any](raw *Result) (*Response[T], error) {
	if raw == nil {
		return nil, &coinwerr.ProtocolError{Message: "nil REST result"}
	}
	data := raw.Data
	if _, ok := any(new(T)).(*model.NoData); ok && len(data) == 0 {
		data = []byte("null")
	}
	var typed T
	if err := model.Decode(data, &typed); err != nil {
		return nil, &coinwerr.ProtocolError{Message: "response contract for " + raw.EndpointID + ": " + err.Error()}
	}
	return &Response[T]{Result: raw, Data: typed}, nil
}
