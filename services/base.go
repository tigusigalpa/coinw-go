// Package services groups the typed REST clients for every documented
// CoinW domain: spot, futures and common (accounts/sub-accounts/affiliate).
// All service I/O methods take a context and return a typed Response envelope.
package services

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
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
	if err := validateContractRequest(id, params); err != nil {
		return nil, err
	}
	if err := validateRequestValues(reflect.ValueOf(params), "", 0); err != nil {
		return nil, err
	}
	return b.ex.Execute(ctx, id, params)
}

func callTyped[T any](b base, ctx context.Context, id string, params interface{}) (*transport.Response[T], error) {
	raw, err := b.call(ctx, id, params)
	if err != nil {
		return nil, err
	}
	return transport.DecodeResponse[T](raw)
}

// Validate flexible batch/extra fields too: they must not bypass the same
// decimal-string boundary as the explicitly typed method arguments.
func validateRequestValues(v reflect.Value, field string, depth int) error {
	if depth > 100 {
		return validationErr("request parameters contain a cycle or excessive nesting")
	}
	for v.IsValid() && (v.Kind() == reflect.Interface || v.Kind() == reflect.Pointer) {
		depth++
		if depth > 100 {
			return validationErr("request parameters contain a cycle or excessive nesting")
		}
		if v.IsNil() {
			return validationErr("parameter %s must not be null", field)
		}
		v = v.Elem()
	}
	if !v.IsValid() {
		return nil
	}
	if monetaryField(field) {
		if v.Type() == reflect.TypeOf(json.Number("")) || v.Kind() != reflect.String || !decimalPattern.MatchString(v.String()) {
			return validationErr("parameter %s must be a decimal string", field)
		}
		return nil
	}
	switch v.Kind() {
	case reflect.Float32, reflect.Float64:
		return validationErr("parameter %s must not be a binary float", field)
	case reflect.Map:
		iter := v.MapRange()
		for iter.Next() {
			if err := validateRequestValues(iter.Value(), fmt.Sprint(iter.Key().Interface()), depth+1); err != nil {
				return err
			}
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			if err := validateRequestValues(v.Index(i), field, depth+1); err != nil {
				return err
			}
		}
	}
	return nil
}

func monetaryField(field string) bool {
	switch field {
	case "amount", "rate", "funds", "quantity", "openPrice", "orderPrice",
		"addMargin", "reduceMargin", "closeNum", "closeRate", "closePiece",
		"callbackRate", "triggerPrice", "dealPiece", "faceValue", "baseSize",
		"stopLossOrderPrice", "stopProfitOrderPrice", "stopLossPrice", "stopProfitPrice",
		"stopLossRate", "stopProfitRate", "stopLossRates", "stopProfitRates":
		return true
	}
	return false
}

func mergeExtras(params, extra map[string]interface{}) error {
	for k, v := range extra {
		if _, exists := params[k]; exists {
			return validationErr("extra parameter %s conflicts with an explicit argument", k)
		}
		params[k] = v
	}
	return nil
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
