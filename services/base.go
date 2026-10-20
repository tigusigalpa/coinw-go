// Package services groups the typed REST clients for every documented
// CoinW domain: spot, futures and common (accounts/sub-accounts/affiliate).
// All methods take a context and return the typed Result envelope.
package services

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	coinwerr "github.com/tigusigalpa/coinw-go/errors"
	"github.com/tigusigalpa/coinw-go/transport"
)

// decimalPattern matches canonical decimal strings — monetary params are
// always validated at the client boundary; floats are never accepted.
var decimalPattern = regexp.MustCompile(`^-?(?:0|[1-9]\d*)(?:\.\d+)?$`)

type base struct {
	ex *transport.Executor
}

func newBase(ex *transport.Executor) base { return base{ex: ex} }

func (b base) call(ctx context.Context, id string, params interface{}) (*transport.Result, error) {
	return b.ex.Execute(ctx, id, params)
}

// nonNull drops nil values, keeping false/0.
func nonNull(m map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(m))
	for k, v := range m {
		if v != nil {
			out[k] = v
		}
	}
	return out
}

func requireNonEmpty(v, name string) (string, error) {
	if strings.TrimSpace(v) == "" {
		return "", &coinwerr.ValidationError{
			Message: fmt.Sprintf("parameter %s must be a non-empty string", name),
		}
	}
	return v, nil
}

// requireMaxInstruments enforces the documented 20-entry instrument cap.
func requireMaxInstruments(items []string, max int) ([]string, error) {
	if len(items) == 0 || len(items) > max {
		return nil, &coinwerr.ValidationError{
			Message: fmt.Sprintf("instrument list must contain 1..%d entries", max),
		}
	}
	return items, nil
}

func requirePageSize(size, max int) (int, error) {
	if size < 1 || size > max {
		return 0, &coinwerr.ValidationError{
			Message: fmt.Sprintf("pageSize must be within 1..%d", max),
		}
	}
	return size, nil
}

func csv(items []string) string { return strings.Join(items, ",") }

func validationErr(format string, args ...interface{}) error {
	return &coinwerr.ValidationError{Message: fmt.Sprintf(format, args...)}
}
