package services

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	coinwerr "github.com/tigusigalpa/coinw-go/errors"
	"github.com/tigusigalpa/coinw-go/model"
	"github.com/tigusigalpa/coinw-go/transport"
)

type recordingDoer struct {
	calls int
	body  string
	query string
}

func (d *recordingDoer) Do(req *http.Request) (*http.Response, error) {
	d.calls++
	d.query = req.URL.RawQuery
	if req.Body != nil {
		b, _ := io.ReadAll(req.Body)
		d.body = string(b)
	}
	data := `{}`
	switch req.URL.Path {
	case "/v1/perpum/order":
		data = `{"value":"9007199254740993"}`
	case "/v1/perpum/positions/history", "/v1/perpum/orders/deals", "/v1/perpum/orders/deals/history":
		data = `{"rows":[],"total":0}`
	}
	return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(`{"code":"200","data":` + data + `}`))}, nil
}

func TestIsolatedMarginFilterIsSent(t *testing.T) {
	d := &recordingDoer{}
	e, err := transport.NewExecutor(transport.Config{BaseURL: "https://example.test", HTTPClient: d, FuturesAPIKey: "K", FuturesSecretKey: "S"})
	if err != nil {
		t.Fatal(err)
	}
	assets := NewFuturesAssetsService(e)
	positions := NewFuturesPositionService(e)
	for _, call := range []func() error{
		func() error {
			_, err := assets.GetTransactionDetails3Days(context.Background(), "BTC", 0, 0, 0)
			return err
		},
		func() error {
			_, err := assets.GetTransactionDetails3Months(context.Background(), "BTC", 0, 0, 0)
			return err
		},
		func() error {
			_, err := positions.GetHistoricalPositionInformation(context.Background(), "", 0)
			return err
		},
	} {
		if err := call(); err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(d.query, "positionModel=0") {
			t.Fatalf("isolated margin filter omitted: %s", d.query)
		}
	}
	if _, err := positions.GetHistoricalPositionInformation(context.Background(), "", -1); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(d.query, "positionModel") {
		t.Fatalf("negative sentinel must omit filter: %s", d.query)
	}
}

func testServices(t *testing.T) (*FuturesTradeService, *AccountService, *recordingDoer) {
	t.Helper()
	d := &recordingDoer{}
	e, err := transport.NewExecutor(transport.Config{BaseURL: "https://example.test", HTTPClient: d, FuturesAPIKey: "K", FuturesSecretKey: "S", SpotAPIKey: "K", SpotSecretKey: "S"})
	if err != nil {
		t.Fatal(err)
	}
	return NewFuturesTradeService(e), NewAccountService(e), d
}

func TestExtraCannotReplaceValidatedArguments(t *testing.T) {
	_, account, d := testServices(t)
	err := mergeExtras(map[string]interface{}{"quantity": "1"}, map[string]interface{}{"quantity": 2.5})
	if err == nil {
		t.Fatal("extra quantity bypassed validation")
	}
	_, err = account.InitiateWithdrawal(context.Background(), WithdrawalRequest{Amount: "1", Currency: "BTC", Address: "original"})
	if err == nil {
		t.Fatal("withdrawal without mandatory chain reached network")
	}
	if d.calls != 0 {
		t.Fatal("invalid arguments reached the network")
	}
}

func TestBatchAndTPSLRejectInvalidMoney(t *testing.T) {
	s, _, d := testServices(t)
	for _, fields := range []map[string]interface{}{
		{"stopType": "1", "closePiece": "1"},
		{"stopType": 1, "closePiece": nil},
		{"stopType": 1, "closePiece": ""},
		{"stopType": 1, "closePiece": 1},
		{"stopLossPrice": 12.5},
		{"stopLossPrice": json.Number("12.5")},
	} {
		_, err := s.setStopLossTakeProfit(context.Background(), "123", "BTC", fields)
		var validation *coinwerr.ValidationError
		if !errors.As(err, &validation) {
			t.Errorf("%v: expected validation error, got %v", fields, err)
		}
	}
	instrument, _ := model.USDTInstrument("BTC")
	_, err := s.PlaceBatchOrders(context.Background(), []model.FuturesOrderRequest{{Instrument: instrument, Quantity: model.Decimal("1e-1")}})
	if err == nil {
		t.Fatal("batch monetary float accepted")
	}
	if d.calls != 0 {
		t.Fatal("invalid arguments reached the network")
	}
	_, err = s.setStopLossTakeProfit(context.Background(), "123", "BTC", map[string]interface{}{"stopLossPrice": "1.000"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(d.body, `"stopLossPrice":"1.000"`) {
		t.Fatalf("decimal precision lost: %s", d.body)
	}
}

func TestDocumentedOrderTypes(t *testing.T) {
	s, _, d := testServices(t)
	for _, typ := range []string{"bbo", "IOC", "FOK", "PostOnly"} {
		instrument, _ := model.USDTInstrument("BTC")
		_, err := s.PlaceAnOrder(context.Background(), model.FuturesOrderRequest{Instrument: instrument, Direction: model.Long, Leverage: 1, QuantityUnit: model.ContractQuantity, Quantity: model.MustDecimal("1"), PositionType: model.OrderType(typ)})
		if err != nil {
			t.Fatalf("documented order type %s rejected: %v", typ, err)
		}
		if !strings.Contains(d.body, `"positionType":"`+typ+`"`) {
			t.Fatalf("wrong type in request: %s", d.body)
		}
	}
}

func TestUnfilteredCancellationIsGated(t *testing.T) {
	d := &recordingDoer{}
	e, err := transport.NewExecutor(transport.Config{BaseURL: "https://example.test", HTTPClient: d, SpotAPIKey: "K", SpotSecretKey: "S"})
	if err != nil {
		t.Fatal(err)
	}
	_, err = NewSpotTradeService(e).CancelAllOrders(context.Background(), "")
	var ambiguity *coinwerr.AmbiguousDocumentationError
	if !errors.As(err, &ambiguity) || d.calls != 0 {
		t.Fatalf("unverified cancellation scope was sent: %v", err)
	}
}

func TestBatchSLTPRequiresStopFromAndAllowsOptionalInstrument(t *testing.T) {
	s, _, d := testServices(t)
	if _, err := s.setStopLossTakeProfitInBatches(context.Background(), "1", "", nil); err == nil || d.calls != 0 {
		t.Fatal("batch SL/TP without stopFrom reached network")
	}
	if _, err := s.setStopLossTakeProfitInBatches(context.Background(), "1", "", map[string]interface{}{"stopFrom": 2, "stopType": 1, "closePiece": "1"}); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(d.body, `"instrument"`) {
		t.Fatalf("optional instrument must be omitted: %s", d.body)
	}
}
