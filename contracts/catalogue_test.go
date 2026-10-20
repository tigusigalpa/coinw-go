package contracts_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	coinw "github.com/tigusigalpa/coinw-go"
	"github.com/tigusigalpa/coinw-go/endpoints"
	coinwerr "github.com/tigusigalpa/coinw-go/errors"
	"github.com/tigusigalpa/coinw-go/model"
)

type contractEntry struct {
	ID, Source, Fixture, TypedMethod, Method, Path, Auth, ResponseType, RequestType string
	Ambiguous, ReadOnly                                                             bool
}
type catalogue struct {
	REST      []contractEntry
	Websocket []struct{ ID, Source, Fixture, AckFixture, Adapter, Method string }
}

func loadCatalogue(t *testing.T) catalogue {
	t.Helper()
	raw, err := os.ReadFile("catalogue.json")
	if err != nil {
		t.Fatal(err)
	}
	var c catalogue
	if err = json.Unmarshal(raw, &c); err != nil {
		t.Fatal(err)
	}
	return c
}

// This frozen source mapping is deliberately independent of the runtime
// registry. Adding/changing an operation requires reviewing its contract too.
func TestCatalogueIsExhaustive(t *testing.T) {
	c := loadCatalogue(t)
	registry := endpoints.All()
	seen := map[string]bool{}
	for _, e := range c.REST {
		if seen[e.ID] {
			t.Fatalf("duplicate REST contract %s", e.ID)
		}
		seen[e.ID] = true
		d, ok := registry[e.ID]
		if !ok {
			t.Fatalf("unregistered contract %s", e.ID)
		}
		if d.Method != e.Method || d.Path != e.Path || d.Auth != e.Auth || d.Ambiguous != e.Ambiguous || d.ReadOnly != e.ReadOnly {
			t.Fatalf("registry differs from verified contract %s", e.ID)
		}
		if e.TypedMethod == "" || e.ResponseType == "" || !strings.HasPrefix(e.Source, "https://www.coinw.com/api-doc/en/") {
			t.Fatalf("incomplete mapping %s", e.ID)
		}
		delete(registry, e.ID)
	}
	if len(registry) > 0 {
		t.Fatalf("uncovered REST endpoints: %d", len(registry))
	}
	files, err := os.ReadDir("fixtures/rest")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != len(c.REST) {
		t.Fatal("unmapped REST fixtures")
	}
	seen = map[string]bool{}
	for _, e := range c.Websocket {
		if seen[e.ID] || e.Method == "" || e.Source == "" {
			t.Fatalf("incomplete WS contract %s", e.ID)
		}
		seen[e.ID] = true
		if _, err = os.Stat(e.Fixture); err != nil {
			t.Fatal(err)
		}
		if e.Adapter == "raw" {
			if _, err = os.Stat(e.AckFixture); err != nil {
				t.Fatal(err)
			}
		}
	}
	if len(c.REST) != 86 || len(c.Websocket) != 26 {
		t.Fatal("catalogue snapshot changed; review source coverage")
	}
}

func TestEveryRESTMethodOverHTTP(t *testing.T) {
	for _, entry := range loadCatalogue(t).REST {
		t.Run(entry.ID, func(t *testing.T) {
			fixture, err := os.ReadFile(entry.Fixture)
			if err != nil {
				t.Fatal(err)
			}
			calls := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls++
				expected, err := url.Parse(entry.Path)
				if err != nil {
					t.Error(err)
					return
				}
				if r.Method != entry.Method || r.URL.Path != expected.Path || r.URL.Query().Get("command") != expected.Query().Get("command") {
					t.Error("wrong documented method/path/command")
				}
				q := r.URL.Query()
				switch entry.Auth {
				case "public":
					if q.Has("api_key") || q.Has("sign") || r.Header.Get("api_key") != "" || r.Header.Get("sign") != "" {
						t.Error("public endpoint received credentials")
					}
				case "spot":
					if q.Get("api_key") != "spot-fixture-key" || q.Get("sign") == "" || r.Header.Get("api_key") != "" {
						t.Error("Spot profile missing or mixed")
					}
				case "futures":
					if r.Header.Get("api_key") != "futures-fixture-key" || r.Header.Get("sign") == "" || r.Header.Get("timestamp") == "" || q.Has("api_key") {
						t.Error("Futures profile missing or mixed")
					}
				}
				body, _ := io.ReadAll(r.Body)
				if r.Method == "GET" && len(body) != 0 {
					t.Error("GET parameters moved into body")
				}
				if entry.ID == "common.affiliate.getInvitedUsersTradingVolume" && string(body) != "{}" {
					t.Error("unfiltered volume list must sign and send empty params object")
				}
				if entry.ID == "futures.market.getKLineOfAnInstrument" && (q.Get("granuality") != "0" || q.Has("granularity")) {
					t.Error("K-line parameter spelling changed")
				}
				if strings.HasPrefix(entry.ID, "spot.") && (strings.HasSuffix(entry.ID, "getOrderStatus") || strings.HasSuffix(entry.ID, "getOrderDetails") || strings.HasSuffix(entry.ID, "cancelOrder")) {
					if q.Get("orderNumber") == "" || q.Has("orderId") {
						t.Error("Spot orderNumber contract changed")
					}
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write(fixture)
			}))
			defer server.Close()
			client, err := coinw.NewClient(coinw.WithRestBaseURL(server.URL), coinw.WithHTTPClient(server.Client()), coinw.WithSpotCredentials("spot-fixture-key", "spot-fixture-secret"), coinw.WithFuturesCredentials("futures-fixture-key", "futures-fixture-secret"), coinw.DisableRateLimiter())
			if err != nil {
				t.Fatal(err)
			}
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			err = invokeContract(ctx, client, entry)
			if entry.Ambiguous {
				var ambiguous *coinwerr.AmbiguousDocumentationError
				if !errors.As(err, &ambiguous) || calls != 0 {
					t.Fatalf("ambiguous route was emitted: %v", err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if calls != 1 {
				t.Fatalf("expected one explicit call, got %d", calls)
			}
		})
	}
}

func invokeContract(ctx context.Context, c *coinw.Client, e contractEntry) error {
	services := map[string]any{"AccountService": c.Account(), "SubAccountService": c.SubAccounts(), "AffiliateService": c.Affiliate(), "SpotMarketService": c.SpotMarket(), "SpotAccountService": c.SpotAccount(), "SpotCheckService": c.SpotCheck(), "SpotTradeService": c.SpotTrade(), "FuturesMarketService": c.FuturesMarket(), "FuturesAssetsService": c.FuturesAssets(), "FuturesCommonService": c.FuturesCommon(), "FuturesPositionService": c.FuturesPosition(), "FuturesCheckService": c.FuturesCheck(), "FuturesTradeService": c.FuturesTrade()}
	parts := strings.Split(e.TypedMethod, ".")
	method := reflect.ValueOf(services[parts[0]]).MethodByName(parts[1])
	if !method.IsValid() {
		return errors.New("mapped typed method missing")
	}
	if parts[0] == "FuturesTradeService" {
		return invokeTrade(ctx, c, e.ID)
	}
	request := reflect.New(method.Type().In(1)).Elem()
	fillRequest(request, "")
	out := method.Call([]reflect.Value{reflect.ValueOf(ctx), request})
	if out[1].IsNil() {
		return nil
	}
	return out[1].Interface().(error)
}

func fillRequest(v reflect.Value, key string) {
	if v.Type() == reflect.TypeFor[model.Instrument]() {
		instrument, _ := model.USDTInstrument("BTC")
		v.Set(reflect.ValueOf(instrument))
		return
	}
	if v.Type() == reflect.TypeFor[model.Decimal]() {
		v.SetString("1.000000000000000001")
		return
	}
	if v.Type() == reflect.TypeFor[model.InstrumentList]() {
		instrument, _ := model.USDTInstrument("BTC")
		v.Set(reflect.ValueOf(model.InstrumentList{instrument}))
		return
	}
	if v.Kind() == reflect.Pointer {
		switch key {
		case "currencyPair":
			v.Set(reflect.New(v.Type().Elem()))
			v.Elem().SetString("BTC_USDT")
		case "amount", "rate", "dealPiece":
			v.Set(reflect.New(v.Type().Elem()))
			v.Elem().SetString("1.000")
		case "openId", "orderId":
			v.Set(reflect.New(v.Type().Elem()))
			v.Elem().SetInt(9007199254740993)
		}
		return
	}
	switch v.Kind() {
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			if v.Field(i).Kind() == reflect.Slice && strings.Contains(v.Type().Field(i).Tag.Get("json"), "omitempty") {
				continue
			}
			fillRequest(v.Field(i), strings.Split(v.Type().Field(i).Tag.Get("json"), ",")[0])
		}
	case reflect.String:
		value := "fixture"
		switch key {
		case "symbol", "currencyPair":
			value = "BTC_USDT"
		case "instrument", "currencyCode", "base", "symbolId", "coinCode", "currency":
			value = "BTC"
		case "chain":
			value = "BTC"
		case "side":
			value = "OUT"
		case "type", "granuality", "status":
			value = "0"
		case "accountType":
			value = "SPOT"
		case "targetAccountType":
			value = "WEALTH"
		case "bizType":
			value = "SPOT_TO_WEALTH"
		case "id", "orderNumber", "subUid", "subAccount_apiKey":
			value = "9007199254740993"
		case "symbols":
			value = "BTC,ETH"
		case "originType", "positionType":
			value = "execute"
		}
		v.SetString(value)
	case reflect.Int, reflect.Int64:
		n := int64(1)
		switch key {
		case "type", "positionModel", "layout":
			n = 0
		case "period":
			n = 60
		case "to":
			n = 2
		}
		v.SetInt(n)
	case reflect.Slice:
		v.Set(reflect.MakeSlice(v.Type(), 1, 1))
		v.Index(0).SetString("9007199254740993")
	}
}

func invokeTrade(ctx context.Context, c *coinw.Client, id string) error {
	s := c.FuturesTrade()
	instrument, _ := model.USDTInstrument("BTC")
	order := model.FuturesOrderRequest{Instrument: instrument, Direction: model.Long, Leverage: 1, QuantityUnit: model.ContractQuantity, Quantity: model.MustDecimal("1"), PositionModel: model.IsolatedMargin, PositionType: model.MarketOrder}
	var err error
	switch id {
	case "futures.trade.adjustMargin":
		_, err = s.AdjustMargin(ctx, "9007199254740993", "1", "", "")
	case "futures.trade.cancelAnOrder":
		_, err = s.CancelAnOrder(ctx, "9007199254740993")
	case "futures.trade.cancelBatchOrders":
		_, err = s.CancelBatchOrders(ctx, []string{"9007199254740993"}, "")
	case "futures.trade.closeAPosition":
		_, err = s.CloseAPosition(ctx, "9007199254740993", "execute", "1", "", "")
	case "futures.trade.closeBatchPositions":
		_, err = s.CloseBatchPositions(ctx, []model.CloseBatchPositionRequest{{ThirdOrderID: "fixture-order"}})
	case "futures.trade.closePositionsAtMarketPrice":
		_, err = s.ClosePositionsAtMarketPrice(ctx, "BTC")
	case "futures.trade.modifyAnOrder":
		_, err = s.ModifyAnOrder(ctx, model.ModifyOrderRequest{ID: "9007199254740993", FuturesOrderRequest: order})
	case "futures.trade.modifyStopLossTakeProfitInBatches":
		_, err = s.ModifyStopLossTakeProfitInBatches(ctx, model.ModifyBatchSLTPRequest{ID: "9007199254740993", StopFrom: 2})
	case "futures.trade.placeAnOrder":
		_, err = s.PlaceAnOrder(ctx, order)
	case "futures.trade.placeBatchOrders":
		_, err = s.PlaceBatchOrders(ctx, []model.FuturesOrderRequest{order})
	case "futures.trade.reverseAPosition":
		_, err = s.ReverseAPosition(ctx, "9007199254740993")
	case "futures.trade.setStopLossTakeProfit":
		_, err = s.SetStopLossTakeProfit(ctx, "9007199254740993", "", model.SLTPPrices{})
	case "futures.trade.setStopLossTakeProfitInBatches":
		_, err = s.SetStopLossTakeProfitInBatches(ctx, "9007199254740993", "", model.BatchSLTPRequest{StopFrom: 2})
	case "futures.trade.setTrailingStopLossTakeProfit":
		_, err = s.SetTrailingStopLossTakeProfit(ctx, "9007199254740993", "0.1", "1", 1, "")
	default:
		return errors.New("unmapped trade contract")
	}
	return err
}
