package model

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func TestExactDecimalCompatibility(t *testing.T) {
	for _, tc := range []struct{ raw, want string }{{`"0E-18"`, "0.000000000000000000"}, {`2.964e-05`, "0.00002964"}, {`"1.044E+5"`, "104400"}, {`9007199254740993123456789.00000001`, "9007199254740993123456789.00000001"}, {`"-0.00100"`, "-0.00100"}} {
		var d WireDecimal
		if err := Decode([]byte(tc.raw), &d); err != nil {
			t.Fatal(err)
		}
		if d.String() != tc.want || string(d.RawJSON()) != tc.raw {
			t.Fatalf("decimal lost exact representation for %s", tc.raw)
		}
		encoded, err := json.Marshal(d)
		if err != nil || string(encoded) != `"`+tc.want+`"` {
			t.Fatalf("decimal serialized as a non-string: %s", encoded)
		}
		raw := d.RawJSON()
		raw[0] = '!'
		if string(d.RawJSON()) != tc.raw {
			t.Fatal("raw token can be mutated")
		}
	}
	for _, raw := range []string{`null`, `true`, `[]`, `{}`, `"NaN"`, `"1e999999"`, `"1e-999999"`, `"01"`, `" 1 "`, `"+1"`} {
		var d WireDecimal
		if Decode([]byte(raw), &d) == nil {
			t.Fatalf("invalid decimal accepted: %s", raw)
		}
	}
	if _, err := json.Marshal(WireDecimal{}); err == nil {
		t.Fatal("uninitialized response decimal serialized")
	}
}

func TestCompatibilityUnions(t *testing.T) {
	var id StringOrInteger
	if err := Decode([]byte(`9007199254740993123`), &id); err != nil || id.String() != "9007199254740993123" {
		t.Fatal("ID precision lost")
	}
	for _, raw := range []string{`1.2`, `true`, `null`, `[]`} {
		if Decode([]byte(raw), &id) == nil {
			t.Fatal("invalid identifier token")
		}
	}
	for _, raw := range []string{`"1.000"`, `{"value":"1.000","ts":1775443378356}`} {
		var d ValueOrScalar[WireDecimal]
		if err := Decode([]byte(raw), &d); err != nil || d.Value.String() != "1.000" || string(d.RawJSON()) != raw {
			t.Fatalf("scalar/wrapper conflict lost: %v", err)
		}
	}
	for _, raw := range []string{`{}`, `{"value":null}`, `{"value":"1","unexpected":true}`} {
		var d ValueOrScalar[WireDecimal]
		if Decode([]byte(raw), &d) == nil {
			t.Fatal("invalid wrapper accepted")
		}
	}
	var masked MaskedDecimal
	for _, raw := range []string{`"-"`, `""`, `"0.00"`} {
		if err := Decode([]byte(raw), &masked); err != nil {
			t.Fatal(err)
		}
	}
}

func TestStrictResponseShape(t *testing.T) {
	for _, raw := range []string{`{}`, `{"makerFee":"0"}`, `{"makerFee":"0","takerFee":null}`, `{"makerFee":false,"takerFee":"0"}`, `{"makerFee":"0","takerFee":"0","surprise":"PRIVATE_VALUE"}`, `{"makerFee":"PRIVATE_VALUE","takerFee":"0"}`, `null`, `[]`, `{"makerFee":"0","takerFee":"0"} {}`} {
		var fees FuturesAccountFees
		err := Decode([]byte(raw), &fees)
		if err == nil {
			t.Fatalf("malformed fees accepted: %s", raw)
		}
		if strings.Contains(err.Error(), "PRIVATE_VALUE") {
			t.Fatal("private payload leaked in decoder error")
		}
	}
	var depth PriceQuantity
	for _, raw := range []string{`[]`, `["1"]`, `["1","2","3"]`, `["1",null]`} {
		if Decode([]byte(raw), &depth) == nil {
			t.Fatal("malformed tuple accepted")
		}
	}
}

func TestModelsContainNoBinaryFloats(t *testing.T) {
	visited := map[reflect.Type]bool{}
	var check func(reflect.Type)
	check = func(typ reflect.Type) {
		if visited[typ] {
			return
		}
		visited[typ] = true
		switch typ.Kind() {
		case reflect.Float32, reflect.Float64:
			t.Fatalf("binary float in model %v", typ)
		case reflect.Pointer, reflect.Slice, reflect.Array, reflect.Map:
			check(typ.Elem())
		case reflect.Struct:
			for i := 0; i < typ.NumField(); i++ {
				check(typ.Field(i).Type)
			}
		}
	}
	for _, value := range []any{FuturesOrderRequest{}, FuturesOrder{}, FuturesPosition{}, FuturesTransaction{}, SpotInstrument{}, SpotTicker{}, SpotCandle{}, FuturesInstrument{}, FuturesTicker{}, FuturesFundingEvent{}, FuturesAssetEvent{}, AffiliateTeam{}, HistoricalPosition{}, FuturesStreamCandle{}} {
		check(reflect.TypeOf(value))
	}
}

func TestInstrumentAndUnitBoundaries(t *testing.T) {
	usdt, err := USDTInstrument("1000PEPE")
	if err != nil || usdt.Settlement() != USDT {
		t.Fatal(err)
	}
	usdc, err := USDCInstrument("BTC_USDC")
	if err != nil || usdc.Settlement() != USDC {
		t.Fatal(err)
	}
	if _, err = USDTInstrument("BTC_USDC"); err == nil {
		t.Fatal("USDC accepted as USDT")
	}
	if _, err = USDCInstrument("BTC"); err == nil {
		t.Fatal("base accepted as USDC pair")
	}
}
