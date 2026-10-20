package coinw

import (
	"encoding/json"
	"testing"
)

func TestDecimal(t *testing.T) {
	valid := []string{"100.5", "-0.00000001", "0", "42", "  100.5  ", "-7"}
	for _, v := range valid {
		if _, err := NewDecimal(v); err != nil {
			t.Errorf("expected %q valid: %v", v, err)
		}
	}
	invalid := []string{"abc", "1e5", "0x10", "1,5", "", "01.5", ".5", "5.", "--1", "NaN"}
	for _, v := range invalid {
		if _, err := NewDecimal(v); err == nil {
			t.Errorf("expected %q rejected", v)
		}
	}
	if !MustDecimal("0.0000").IsZero() || !MustDecimal("-0").IsZero() {
		t.Error("zero variants should be IsZero")
	}
	if MustDecimal("0.0001").IsZero() {
		t.Error("non-zero must not be IsZero")
	}
	raw, _ := MustDecimal("1.5").MarshalJSON()
	if string(raw) != `"1.5"` {
		t.Errorf("marshal: %s", raw)
	}
}

func TestDecimalJSONBoundary(t *testing.T) {
	for _, raw := range []string{`"invalid"`, `"1e5"`, `1.5`, `null`} {
		d := MustDecimal("1.000")
		if err := json.Unmarshal([]byte(raw), &d); err == nil {
			t.Errorf("invalid decimal JSON accepted: %s", raw)
		}
		if d.String() != "1.000" {
			t.Fatal("failed decoding changed the destination")
		}
	}
	if _, err := json.Marshal(Decimal("invalid")); err == nil {
		t.Fatal("invalid cast Decimal was serialized")
	}
}
