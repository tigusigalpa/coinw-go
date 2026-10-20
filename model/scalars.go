// Package model contains the wire contracts from the CoinW documentation.
// Response compatibility types retain the exact JSON token when the published
// tables and examples disagree. No money passes through binary floating point.
package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var numberPattern = regexp.MustCompile(`^(-?)(0|[1-9][0-9]*)(?:\.([0-9]+))?(?:[eE]([+-]?[0-9]+))?$`)
var integerPattern = regexp.MustCompile(`^-?(?:0|[1-9][0-9]*)$`)

// WireDecimal explicitly accepts the documented string/JSON-number conflict.
// String and MarshalJSON always expose validated plain decimal strings.
// RawJSON returns a copy of the original token, including its original scale.
type WireDecimal struct {
	value string
	raw   json.RawMessage
}

func (d WireDecimal) String() string           { return d.value }
func (d WireDecimal) RawJSON() json.RawMessage { return bytes.Clone(d.raw) }
func (d WireDecimal) MarshalJSON() ([]byte, error) {
	if !decimalPattern.MatchString(d.value) {
		return nil, fmt.Errorf("uninitialized WireDecimal")
	}
	return json.Marshal(d.value)
}
func (d *WireDecimal) UnmarshalJSON(raw []byte) error {
	text := string(raw)
	if len(raw) > 0 && raw[0] == '"' {
		if err := json.Unmarshal(raw, &text); err != nil {
			return fmt.Errorf("invalid WireDecimal string")
		}
	}
	normalized, err := expandDecimal(text)
	if err != nil {
		return err
	}
	d.value, d.raw = normalized, bytes.Clone(raw)
	return nil
}

func expandDecimal(s string) (string, error) {
	// Bound work on untrusted exponent tokens; no numeric precision is lost.
	if len(s) > 8192 {
		return "", fmt.Errorf("decimal token exceeds 8192 bytes")
	}
	m := numberPattern.FindStringSubmatch(s)
	if m == nil {
		return "", fmt.Errorf("expected a decimal string or number")
	}
	exponent := 0
	if m[4] != "" {
		n, err := strconv.Atoi(m[4])
		if err != nil || n < -4096 || n > 4096 {
			return "", fmt.Errorf("decimal exponent outside supported range")
		}
		exponent = n
	}
	if exponent == 0 {
		return m[1] + m[2] + func() string {
			if m[3] != "" {
				return "." + m[3]
			}
			return ""
		}(), nil
	}
	digits, point := m[2]+m[3], len(m[2])+exponent
	var result string
	switch {
	case point <= 0:
		result = "0." + strings.Repeat("0", -point) + digits
	case point >= len(digits):
		result = digits + strings.Repeat("0", point-len(digits))
	default:
		result = digits[:point] + "." + digits[point:]
	}
	parts := strings.SplitN(result, ".", 2)
	parts[0] = strings.TrimLeft(parts[0], "0")
	if parts[0] == "" {
		parts[0] = "0"
	}
	return m[1] + strings.Join(parts, "."), nil
}

// StringOrInteger preserves opaque identifiers, timestamps and status fields
// whose documented examples alternate between strings and integer tokens.
// It does not infer timestamp units, date formats, or identifier arithmetic.
type StringOrInteger struct {
	text string
	raw  json.RawMessage
}

func (s StringOrInteger) String() string           { return s.text }
func (s StringOrInteger) RawJSON() json.RawMessage { return bytes.Clone(s.raw) }
func (s StringOrInteger) MarshalJSON() ([]byte, error) {
	if len(s.raw) == 0 {
		return nil, fmt.Errorf("uninitialized StringOrInteger")
	}
	return bytes.Clone(s.raw), nil
}
func (s *StringOrInteger) UnmarshalJSON(raw []byte) error {
	text := string(raw)
	if len(raw) > 0 && raw[0] == '"' {
		if err := json.Unmarshal(raw, &text); err != nil {
			return fmt.Errorf("invalid string token")
		}
	} else if !integerPattern.Match(raw) {
		return fmt.Errorf("expected a string or integer token")
	}
	s.text, s.raw = text, bytes.Clone(raw)
	return nil
}

// MaskedDecimal is used only for affiliate contractBalance, which the docs
// allow to be hidden as "-" or "". A masked balance has no numeric meaning.
type MaskedDecimal struct {
	Decimal *WireDecimal
	Mask    *string
}

func (m *MaskedDecimal) UnmarshalJSON(raw []byte) error {
	var s string
	if json.Unmarshal(raw, &s) == nil && (s == "-" || s == "") {
		m.Mask, m.Decimal = &s, nil
		return nil
	}
	var d WireDecimal
	if err := json.Unmarshal(raw, &d); err != nil {
		return err
	}
	m.Decimal, m.Mask = &d, nil
	return nil
}
func (m MaskedDecimal) MarshalJSON() ([]byte, error) {
	if m.Decimal != nil && m.Mask == nil {
		return json.Marshal(m.Decimal)
	}
	if m.Mask != nil && m.Decimal == nil && (*m.Mask == "-" || *m.Mask == "") {
		return json.Marshal(*m.Mask)
	}
	return nil, fmt.Errorf("invalid masked decimal union")
}

// PercentageText retains the literal percentage suffix; no fraction conversion
// is implied by historical-position profitRate.
type PercentageText string

func (p *PercentageText) UnmarshalJSON(raw []byte) error {
	var s string
	if json.Unmarshal(raw, &s) != nil || !strings.HasSuffix(s, "%") {
		return fmt.Errorf("expected percentage text")
	}
	if _, err := expandDecimal(strings.TrimSuffix(s, "%")); err != nil {
		return err
	}
	*p = PercentageText(s)
	return nil
}

func (p PercentageText) MarshalJSON() ([]byte, error) {
	s := string(p)
	if !strings.HasSuffix(s, "%") {
		return nil, fmt.Errorf("invalid percentage text")
	}
	if _, err := expandDecimal(strings.TrimSuffix(s, "%")); err != nil {
		return nil, err
	}
	return json.Marshal(s)
}

// ValueOrScalar models pages where the table declares a scalar but the example
// wraps it in {"value":...,"ts":...}. Raw form remains available to the caller.
type ValueOrScalar[T any] struct {
	Value T
	TS    *StringOrInteger
	raw   json.RawMessage
}

func (v ValueOrScalar[T]) RawJSON() json.RawMessage { return bytes.Clone(v.raw) }
func (v *ValueOrScalar[T]) UnmarshalJSON(raw []byte) error {
	var value T
	var ts *StringOrInteger
	if len(raw) > 0 && raw[0] == '{' {
		var obj struct {
			Value *T               `json:"value" required:"true"`
			TS    *StringOrInteger `json:"ts,omitempty"`
		}
		if err := Decode(raw, &obj); err != nil {
			return err
		}
		if obj.Value == nil {
			return fmt.Errorf("value must not be null")
		}
		value, ts = *obj.Value, obj.TS
	} else if err := Decode(raw, &value); err != nil {
		return err
	}
	v.Value, v.TS, v.raw = value, ts, bytes.Clone(raw)
	return nil
}
func (v ValueOrScalar[T]) MarshalJSON() ([]byte, error) {
	if len(v.raw) > 0 && v.raw[0] == '{' {
		return json.Marshal(struct {
			Value T                `json:"value"`
			TS    *StringOrInteger `json:"ts,omitempty"`
		}{v.Value, v.TS})
	}
	return json.Marshal(v.Value)
}

// ObjectOrArray is an explicit object/array conflict, not an invented page.
type ObjectOrArray[T any] struct {
	Object *T
	Array  []T
	raw    json.RawMessage
}

func (v ObjectOrArray[T]) RawJSON() json.RawMessage { return bytes.Clone(v.raw) }
func (v *ObjectOrArray[T]) UnmarshalJSON(raw []byte) error {
	if len(raw) > 0 && raw[0] == '[' {
		var a []T
		if err := Decode(raw, &a); err != nil {
			return err
		}
		v.Array, v.Object = a, nil
	} else {
		var t T
		if err := Decode(raw, &t); err != nil {
			return err
		}
		v.Object, v.Array = &t, nil
	}
	v.raw = bytes.Clone(raw)
	return nil
}
func (v ObjectOrArray[T]) MarshalJSON() ([]byte, error) {
	if v.Object != nil && v.Array == nil {
		return json.Marshal(v.Object)
	}
	if v.Object == nil && v.Array != nil {
		return json.Marshal(v.Array)
	}
	return nil, fmt.Errorf("invalid object/array union")
}

// JSONEncoded accepts the JSON-string payload used in Spot streams as well as
// the object/array shape printed in their tables. It never evaluates text.
type JSONEncoded[T any] struct {
	Value T
	raw   json.RawMessage
}

func (v JSONEncoded[T]) RawJSON() json.RawMessage { return bytes.Clone(v.raw) }
func (v *JSONEncoded[T]) UnmarshalJSON(raw []byte) error {
	b := raw
	if len(raw) > 0 && raw[0] == '"' {
		var s string
		if json.Unmarshal(raw, &s) != nil {
			return fmt.Errorf("invalid encoded JSON")
		}
		b = []byte(s)
	}
	var t T
	if err := Decode(b, &t); err != nil {
		return err
	}
	v.Value, v.raw = t, bytes.Clone(raw)
	return nil
}
func (v JSONEncoded[T]) MarshalJSON() ([]byte, error) {
	raw, err := json.Marshal(v.Value)
	if err != nil {
		return nil, err
	}
	if len(v.raw) > 0 && v.raw[0] == '"' {
		return json.Marshal(string(raw))
	}
	return raw, nil
}

// LeverageOptions preserves the documented string versus integer-array forms.
type LeverageOptions struct {
	Text   *string
	Values []int
	raw    json.RawMessage
}

func (v *LeverageOptions) UnmarshalJSON(raw []byte) error {
	if len(raw) > 0 && raw[0] == '"' {
		var s string
		if err := json.Unmarshal(raw, &s); err != nil {
			return err
		}
		v.Text, v.Values = &s, nil
	} else {
		var a []int
		if err := Decode(raw, &a); err != nil {
			return err
		}
		v.Text, v.Values = nil, a
	}
	v.raw = bytes.Clone(raw)
	return nil
}
func (v LeverageOptions) RawJSON() json.RawMessage { return bytes.Clone(v.raw) }
func (v LeverageOptions) MarshalJSON() ([]byte, error) {
	if v.Text != nil {
		return json.Marshal(v.Text)
	}
	if v.Values != nil {
		return json.Marshal(v.Values)
	}
	return nil, fmt.Errorf("uninitialized leverage options")
}
