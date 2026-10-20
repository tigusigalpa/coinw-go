package services

import (
	"context"
	"errors"
	"testing"

	coinwerr "github.com/tigusigalpa/coinw-go/errors"
	"github.com/tigusigalpa/coinw-go/model"
)

func TestTypedRequestBoundariesFailBeforeNetwork(t *testing.T) {
	s, _, d := testServices(t)
	instrument, _ := model.USDTInstrument("BTC")
	base := model.FuturesOrderRequest{Instrument: instrument, Direction: model.Long, Leverage: 1, QuantityUnit: model.ContractQuantity, Quantity: model.MustDecimal("1"), PositionType: model.MarketOrder}
	for _, change := range []func(*model.FuturesOrderRequest){func(r *model.FuturesOrderRequest) { r.Quantity = model.Decimal("1e-1") }, func(r *model.FuturesOrderRequest) { r.Direction = "opposite" }, func(r *model.FuturesOrderRequest) { r.PositionType = model.LimitOrder }, func(r *model.FuturesOrderRequest) { r.QuantityUnit = 3 }, func(r *model.FuturesOrderRequest) { r.ThirdOrderID = "invalid id" }, func(r *model.FuturesOrderRequest) { r.Instrument = model.Instrument{} }} {
		r := base
		change(&r)
		if _, err := s.PlaceAnOrder(context.Background(), r); err == nil {
			t.Fatal("invalid typed order reached transport")
		}
	}
	partial := 1
	if _, err := s.SetStopLossTakeProfitInBatches(context.Background(), "1", "", model.BatchSLTPRequest{StopFrom: 2, StopType: &partial}); err == nil {
		t.Fatal("partial SL/TP accepted without closePiece")
	}
	if _, err := s.CancelBatchOrders(context.Background(), make([]string, 21), ""); err == nil {
		t.Fatal("21 cancellations accepted")
	}
	if _, err := s.CloseAPosition(context.Background(), "1", "execute", "", "1.1", ""); err == nil {
		t.Fatal("closeRate outside [0,1] accepted")
	}
	if d.calls != 0 {
		t.Fatal("invalid typed request performed network I/O")
	}
}

func TestTypedCancellationCannotBypassInstrumentGate(t *testing.T) {
	trade, _, doer := testServices(t)
	spot := NewSpotTradeService(trade.ex)
	blank := " "
	for _, filter := range []*string{nil, &blank} {
		_, err := spot.CancelAllOrdersWithRequest(context.Background(), model.SpotTradeCancelAllOrdersRequest{CurrencyPair: filter})
		var gate *coinwerr.AmbiguousDocumentationError
		if !errors.As(err, &gate) {
			t.Fatalf("typed unfiltered cancellation bypassed its gate: %v", err)
		}
	}
	if doer.calls != 0 {
		t.Fatal("unfiltered cancellation reached the network")
	}
}
