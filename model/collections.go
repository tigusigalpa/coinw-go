package model

import (
	"encoding/json"
	"fmt"
)

// RowsPage retains published cursor tokens without inferring page arithmetic.
type RowsPage[T any] struct {
	NextID *StringOrInteger `json:"nextId,omitempty"`
	PrevID *StringOrInteger `json:"prevId,omitempty"`
	Rows   []T              `json:"rows" required:"true"`
	Total  StringOrInteger  `json:"total" required:"true"`
}

type ListPage[T any] struct {
	List    []T             `json:"list" required:"true"`
	Total   StringOrInteger `json:"total" required:"true"`
	Context *OpaqueNull     `json:"context,omitempty"`
}

// OpaqueNull records an undocumented field shown only as null. A non-null
// value requires a verified schema update; the SDK assigns it no meaning.
type OpaqueNull struct{}

func (*OpaqueNull) UnmarshalJSON(raw []byte) error {
	if string(raw) == "null" {
		return nil
	}
	return fmt.Errorf("undocumented non-null field")
}
func (OpaqueNull) MarshalJSON() ([]byte, error) { return []byte("null"), nil }

type SpotHistoryPage struct {
	List   []SpotHistoricalOrder `json:"list" required:"true"`
	Before *StringOrInteger      `json:"before,omitempty"`
	After  *StringOrInteger      `json:"after,omitempty"`
}

// PopularVolumes has three documented totals and dynamic pair/currency keys.
// Pair totals are deliberately separate from the overall currency totals.
type PopularVolumes struct {
	TotalETH  WireDecimal
	TotalUSDT WireDecimal
	TotalBTC  WireDecimal
	Pairs     map[string]map[string]WireDecimal
}

func (v *PopularVolumes) UnmarshalJSON(raw []byte) error {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(raw, &obj); err != nil || obj == nil {
		return fmt.Errorf("expected volume object")
	}
	var next PopularVolumes
	next.Pairs = map[string]map[string]WireDecimal{}
	for key, d := range map[string]*WireDecimal{"totalETH": &next.TotalETH, "totalUSDT": &next.TotalUSDT, "totalBTC": &next.TotalBTC} {
		b, ok := obj[key]
		if !ok {
			return fmt.Errorf("missing total volume")
		}
		if err := Decode(b, d); err != nil {
			return err
		}
		delete(obj, key)
	}
	for key, b := range obj {
		var pairs map[string]WireDecimal
		if err := Decode(b, &pairs); err != nil {
			return err
		}
		next.Pairs[key] = pairs
	}
	*v = next
	return nil
}
func (v PopularVolumes) MarshalJSON() ([]byte, error) {
	obj := map[string]any{"totalETH": v.TotalETH, "totalUSDT": v.TotalUSDT, "totalBTC": v.TotalBTC}
	for k, p := range v.Pairs {
		if _, exists := obj[k]; exists {
			return nil, fmt.Errorf("pair key conflicts with total")
		}
		obj[k] = p
	}
	return json.Marshal(obj)
}

// PriceQuantity is the exact two-element Spot price/quantity tuple.
type PriceQuantity struct {
	Price    WireDecimal
	Quantity WireDecimal
}

func (v *PriceQuantity) UnmarshalJSON(raw []byte) error {
	a, err := tuple(raw, 2)
	if err != nil {
		return err
	}
	var n PriceQuantity
	if err = Decode(a[0], &n.Price); err != nil {
		return err
	}
	if err = Decode(a[1], &n.Quantity); err != nil {
		return err
	}
	*v = n
	return nil
}
func (v PriceQuantity) MarshalJSON() ([]byte, error) {
	return json.Marshal([2]WireDecimal{v.Price, v.Quantity})
}

type SequencedPriceQuantity struct {
	Price    WireDecimal
	Quantity WireDecimal
	Sequence StringOrInteger
}

func (v *SequencedPriceQuantity) UnmarshalJSON(raw []byte) error {
	a, err := tuple(raw, 3)
	if err != nil {
		return err
	}
	var n SequencedPriceQuantity
	if err = Decode(a[0], &n.Price); err != nil {
		return err
	}
	if err = Decode(a[1], &n.Quantity); err != nil {
		return err
	}
	if err = Decode(a[2], &n.Sequence); err != nil {
		return err
	}
	*v = n
	return nil
}
func (v SequencedPriceQuantity) MarshalJSON() ([]byte, error) {
	return json.Marshal([3]any{v.Price, v.Quantity, v.Sequence})
}

// FuturesCandleTuple retains positional slots: the REST table and example
// disagree about high/open ordering. No OHLC interpretation is guessed.
type FuturesCandleTuple struct {
	Timestamp StringOrInteger
	Values    [5]WireDecimal
}

func (v *FuturesCandleTuple) UnmarshalJSON(raw []byte) error {
	a, err := tuple(raw, 6)
	if err != nil {
		return err
	}
	var n FuturesCandleTuple
	if err = Decode(a[0], &n.Timestamp); err != nil {
		return err
	}
	for i := range n.Values {
		if err = Decode(a[i+1], &n.Values[i]); err != nil {
			return err
		}
	}
	*v = n
	return nil
}
func (v FuturesCandleTuple) MarshalJSON() ([]byte, error) {
	a := []any{v.Timestamp}
	for _, d := range v.Values {
		a = append(a, d)
	}
	return json.Marshal(a)
}

// SpotCandleTuple follows the documented Spot stream order.
type SpotCandleTuple struct {
	Timestamp                                  StringOrInteger
	Open, Close, High, Low, Volume, QuoteValue WireDecimal
}

func (v *SpotCandleTuple) UnmarshalJSON(raw []byte) error {
	a, err := tuple(raw, 7)
	if err != nil {
		return err
	}
	var n SpotCandleTuple
	if err = Decode(a[0], &n.Timestamp); err != nil {
		return err
	}
	for i, p := range []*WireDecimal{&n.Open, &n.Close, &n.High, &n.Low, &n.Volume, &n.QuoteValue} {
		if err = Decode(a[i+1], p); err != nil {
			return err
		}
	}
	*v = n
	return nil
}
func (v SpotCandleTuple) MarshalJSON() ([]byte, error) {
	return json.Marshal([7]any{v.Timestamp, v.Open, v.Close, v.High, v.Low, v.Volume, v.QuoteValue})
}

func tuple(raw []byte, length int) ([]json.RawMessage, error) {
	var a []json.RawMessage
	if err := json.Unmarshal(raw, &a); err != nil || len(a) != length {
		return nil, fmt.Errorf("expected a %d-element tuple", length)
	}
	return a, nil
}

// FuturesStreamCandle follows the separate WS table's timestamp/OHLC/volume order.
type FuturesStreamCandle struct {
	Timestamp                      StringOrInteger
	Open, High, Low, Close, Volume WireDecimal
}

func (v *FuturesStreamCandle) UnmarshalJSON(raw []byte) error {
	a, err := tuple(raw, 6)
	if err != nil {
		return err
	}
	var n FuturesStreamCandle
	if err = Decode(a[0], &n.Timestamp); err != nil {
		return err
	}
	for i, p := range []*WireDecimal{&n.Open, &n.High, &n.Low, &n.Close, &n.Volume} {
		if err = Decode(a[i+1], p); err != nil {
			return err
		}
	}
	*v = n
	return nil
}
func (v FuturesStreamCandle) MarshalJSON() ([]byte, error) {
	return json.Marshal([6]any{v.Timestamp, v.Open, v.High, v.Low, v.Close, v.Volume})
}
