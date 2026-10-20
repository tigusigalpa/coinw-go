package coinw_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	coinw "github.com/tigusigalpa/coinw-go"
	"github.com/tigusigalpa/coinw-go/model"
)

type exampleTransport struct{}

func (exampleTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": []string{"application/json"}}, Body: io.NopCloser(strings.NewReader(`{"code":"200","data":[{"id":9007199254740993,"price":"83000.00000001","amount":"0.000000000000000001","time":"2026-10-20 12:00:00"}]}`)), Request: r}, nil
}

// This executable public-usage example uses an injected offline transport.
func ExampleClient() {
	c, err := coinw.NewClient(coinw.WithHTTPClient(&http.Client{Transport: exampleTransport{}}))
	if err != nil {
		panic(err)
	}
	r, err := c.SpotMarket().GetRecentTrades(context.Background(), "BTC_USDT", 0)
	if err != nil {
		panic(err)
	}
	fmt.Println(r.Data[0].ID.String(), r.Data[0].Price.String(), r.Data[0].Amount.String())
	// Output: 9007199254740993 83000.00000001 0.000000000000000001
}

func ExampleDecimal() {
	quantity, err := coinw.NewDecimal("0.000000000000000001")
	if err != nil {
		panic(err)
	}
	fmt.Println(quantity.String())
	// Output: 0.000000000000000001
}

func ExampleClient_FuturesTrade() {
	instrument, err := model.USDCInstrument("BTC_USDC")
	if err != nil {
		panic(err)
	}
	request := model.FuturesOrderRequest{Instrument: instrument, Direction: model.Long, Leverage: 1, QuantityUnit: model.ContractQuantity, Quantity: model.MustDecimal("1"), PositionModel: model.IsolatedMargin, PositionType: model.MarketOrder}
	// Only an explicit client.FuturesTrade().PlaceAnOrder(ctx, request) submits it.
	fmt.Println(request.Instrument.String(), request.Quantity.String())
	// Output: BTC_USDC 1
}
