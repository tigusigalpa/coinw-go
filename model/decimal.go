package model

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/tigusigalpa/coinw-go/errors"
)

// Decimal is a validated decimal string. Every monetary quantity in the SDK
// (price, quantity, amount, fee, PnL, rate, margin, funding rate) is carried
// as a decimal string — binary floats are never emitted for money.
type Decimal string

var decimalPattern = regexp.MustCompile(`^-?(?:0|[1-9]\d*)(?:\.\d+)?$`)

// NewDecimal validates a decimal string.
func NewDecimal(v string) (Decimal, error) {
	v = strings.TrimSpace(v)
	if !decimalPattern.MatchString(v) {
		return "", &errors.ValidationError{
			Message: "invalid decimal string " + v + " — monetary values must be decimal strings, not floats",
		}
	}
	return Decimal(v), nil
}

// MustDecimal validates or panics — intended for tests and constants only.
func MustDecimal(v string) Decimal {
	d, err := NewDecimal(v)
	if err != nil {
		panic(err)
	}
	return d
}

// String returns the canonical decimal string.
func (d Decimal) String() string { return string(d) }

// IsZero reports whether the value equals zero in any decimal spelling.
func (d Decimal) IsZero() bool {
	t := strings.TrimLeft(string(d), "-")
	for _, c := range t {
		if c != '0' && c != '.' {
			return false
		}
	}
	return true
}

// MarshalJSON emits the decimal as a JSON string.
func (d Decimal) MarshalJSON() ([]byte, error) {
	if !decimalPattern.MatchString(string(d)) {
		return nil, &errors.ValidationError{Message: "invalid Decimal value"}
	}
	return json.Marshal(string(d))
}

// UnmarshalJSON accepts only validated decimal strings. Numeric API fields
// require an explicitly chosen compatibility adapter, never a float coercion.
func (d *Decimal) UnmarshalJSON(raw []byte) error {
	var v string
	if err := json.Unmarshal(raw, &v); err != nil || !decimalPattern.MatchString(v) {
		return &errors.ValidationError{Message: "Decimal JSON value must be a decimal string"}
	}
	*d = Decimal(v)
	return nil
}
