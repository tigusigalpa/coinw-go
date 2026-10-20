package coinw

import "github.com/tigusigalpa/coinw-go/model"

// Decimal is a validated decimal string.
type Decimal = model.Decimal

// NewDecimal validates a monetary string.
func NewDecimal(v string) (Decimal, error) { return model.NewDecimal(v) }

// MustDecimal validates a constant or panics.
func MustDecimal(v string) Decimal { return model.MustDecimal(v) }
