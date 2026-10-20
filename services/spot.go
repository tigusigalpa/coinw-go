package services

import (
	"context"
	"sort"

	coinwerr "github.com/tigusigalpa/coinw-go/errors"
	"github.com/tigusigalpa/coinw-go/transport"
)

// ---------------------------------------------------------------------------
// Spot > Market (public)
// ---------------------------------------------------------------------------

// SpotMarketService — public spot market data.
type SpotMarketService struct{ base }

func NewSpotMarketService(ex *transport.Executor) *SpotMarketService {
	return &SpotMarketService{newBase(ex)}
}

// GetInstrumentInformation — command=returnSymbol.
func (s *SpotMarketService) GetInstrumentInformation(ctx context.Context) (*transport.Result, error) {
	return s.call(ctx, "spot.market.getInstrumentInformation", nil)
}

// Get24hTradeSummaryForAllInstruments — command=returnTicker.
func (s *SpotMarketService) Get24hTradeSummaryForAllInstruments(ctx context.Context) (*transport.Result, error) {
	return s.call(ctx, "spot.market.get24hTradeSummaryForAllInstruments", nil)
}

// GetCoinDepositWithdrawalLimits — command=returnCurrencies.
func (s *SpotMarketService) GetCoinDepositWithdrawalLimits(ctx context.Context) (*transport.Result, error) {
	return s.call(ctx, "spot.market.getCoinDepositWithdrawalLimits", nil)
}

// GetOrderBook — command=returnOrderBook. size=0 omits the documented size.
func (s *SpotMarketService) GetOrderBook(ctx context.Context, symbol string, size int) (*transport.Result, error) {
	if _, err := requireNonEmpty(symbol, "symbol"); err != nil {
		return nil, err
	}
	params := map[string]interface{}{"symbol": symbol}
	if size > 0 {
		params["size"] = size
	}
	return s.call(ctx, "spot.market.getOrderBook", params)
}

// GetRecentTrades — command=returnTradeHistory. end=0 omits the documented
// end timestamp (ms).
func (s *SpotMarketService) GetRecentTrades(ctx context.Context, symbol string, end int64) (*transport.Result, error) {
	if _, err := requireNonEmpty(symbol, "symbol"); err != nil {
		return nil, err
	}
	params := map[string]interface{}{"symbol": symbol}
	if end != 0 {
		params["end"] = end
	}
	return s.call(ctx, "spot.market.getRecentTrades", params)
}

// SpotKLinePeriods are the documented K-line intervals in SECONDS.
// The doc page's hyperlink example shows period=30, contradicting the table —
// the SDK enforces the table contract and does not rewrite the value.
var SpotKLinePeriods = []int{60, 180, 300, 900, 1800, 7200, 14400}

// GetKLine — command=returnChartData.
// currencyPairs: e.g. ["BTC_USDT","ETH_USDT"], max 20. period must be one of
// SpotKLinePeriods. start/end are documented ms timestamps; 0 omits them.
func (s *SpotMarketService) GetKLine(ctx context.Context, currencyPairs []string, period int, start, end int64) (*transport.Result, error) {
	ok := sort.SearchInts(SpotKLinePeriods, period) < len(SpotKLinePeriods) &&
		SpotKLinePeriods[sort.SearchInts(SpotKLinePeriods, period)] == period
	if !ok {
		return nil, validationErr("period must be one of the documented K-line intervals (seconds): %v", SpotKLinePeriods)
	}
	if _, err := requireMaxInstruments(currencyPairs, 20); err != nil {
		return nil, err
	}
	params := map[string]interface{}{
		"currencyPair": csv(currencyPairs),
		"period":       period,
	}
	if start != 0 {
		params["start"] = start
	}
	if end != 0 {
		params["end"] = end
	}
	return s.call(ctx, "spot.market.getKLine", params)
}

// Get24hVolumeForPopularInstruments — command=return24hVolume.
func (s *SpotMarketService) Get24hVolumeForPopularInstruments(ctx context.Context) (*transport.Result, error) {
	return s.call(ctx, "spot.market.get24hVolumeForPopularInstruments", nil)
}

// ---------------------------------------------------------------------------
// Spot > Account (private)
// ---------------------------------------------------------------------------

// SpotAccountService — private spot account reads.
type SpotAccountService struct{ base }

func NewSpotAccountService(ex *transport.Executor) *SpotAccountService {
	return &SpotAccountService{newBase(ex)}
}

// GetSpotAccountBalance — command=returnBalances.
func (s *SpotAccountService) GetSpotAccountBalance(ctx context.Context) (*transport.Result, error) {
	return s.call(ctx, "spot.account.getSpotAccountBalance", nil)
}

// GetFullSpotAccountBalance — command=returnCompleteBalances.
func (s *SpotAccountService) GetFullSpotAccountBalance(ctx context.Context) (*transport.Result, error) {
	return s.call(ctx, "spot.account.getFullSpotAccountBalance", nil)
}

// ---------------------------------------------------------------------------
// Spot > Check (private reads)
// ---------------------------------------------------------------------------

// SpotCheckService — order/transaction state reads.
type SpotCheckService struct{ base }

func NewSpotCheckService(ex *transport.Executor) *SpotCheckService {
	return &SpotCheckService{newBase(ex)}
}

// GetBatchHistoricalOrders — command=getBatchHistoryOrders at /v1/private.
//
// GATED: the documented endpoint omits the /api prefix used by every other
// spot private route; host/path is unverified. Returns
// AmbiguousDocumentationError until resolved. orderIds must contain 1..50
// entries.
func (s *SpotCheckService) GetBatchHistoricalOrders(ctx context.Context, orderIDs []string) (*transport.Result, error) {
	if len(orderIDs) == 0 || len(orderIDs) > 50 {
		return nil, validationErr("orderIds must contain 1..50 entries")
	}
	return nil, &coinwerr.AmbiguousDocumentationError{
		EndpointID: "spot.check.getBatchHistoricalOrders",
		Note: "documented endpoint /v1/private omits the /api prefix used by " +
			"all other spot private routes; the host/path is unverified",
	}
}

// GetCurrentOrders — command=returnOpenOrders. Empty currencyPair lists all.
func (s *SpotCheckService) GetCurrentOrders(ctx context.Context, currencyPair string) (*transport.Result, error) {
	return s.call(ctx, "spot.check.getCurrentOrders", nonNull(map[string]interface{}{
		"currencyPair": emptyToNil(currencyPair),
	}))
}

// GetHistoricalOrders — command=getUserTrades.
func (s *SpotCheckService) GetHistoricalOrders(ctx context.Context, symbol string) (*transport.Result, error) {
	return s.call(ctx, "spot.check.getHistoricalOrders", nonNull(map[string]interface{}{
		"symbol": emptyToNil(symbol),
	}))
}

// GetOrderDetails — command=returnOrderTrades.
func (s *SpotCheckService) GetOrderDetails(ctx context.Context, orderID string) (*transport.Result, error) {
	return s.call(ctx, "spot.check.getOrderDetails", nonNull(map[string]interface{}{
		"orderId": emptyToNil(orderID),
	}))
}

// GetOrderStatus — command=returnOrderStatus.
func (s *SpotCheckService) GetOrderStatus(ctx context.Context, orderID string) (*transport.Result, error) {
	return s.call(ctx, "spot.check.getOrderStatus", nonNull(map[string]interface{}{
		"orderId": emptyToNil(orderID),
	}))
}

// GetTransactionHistory — command=returnUTradeHistory.
func (s *SpotCheckService) GetTransactionHistory(ctx context.Context, currencyPair string) (*transport.Result, error) {
	return s.call(ctx, "spot.check.getTransactionHistory", nonNull(map[string]interface{}{
		"currencyPair": emptyToNil(currencyPair),
	}))
}

// ---------------------------------------------------------------------------
// Spot > Trade (private mutations — never auto-retried)
// ---------------------------------------------------------------------------

const (
	SpotSideBuy  = 0
	SpotSideSell = 1
)

// SpotOrderOptions carries optional documented placeOrder fields.
type SpotOrderOptions struct {
	Rate       string // limit price — decimal string (omit for market)
	IsMarket   *bool  // documented market-order flag
	OutTradeNo string // client order reference
	Funds      string // quote amount for market orders — decimal string
}

// SpotTradeService — order placement/cancellation. MUTATING: the SDK never
// invokes these implicitly and never retries them; after a transport
// timeout verify state via SpotCheckService before resubmitting.
type SpotTradeService struct{ base }

func NewSpotTradeService(ex *transport.Executor) *SpotTradeService {
	return &SpotTradeService{newBase(ex)}
}

// PlaceOrder — command=doTrade. MUTATION.
// symbol e.g. "BTC_USDT"; typ is SpotSideBuy (0) or SpotSideSell (1);
// amount is a decimal string.
func (s *SpotTradeService) PlaceOrder(ctx context.Context, symbol string, typ int, amount string, opts SpotOrderOptions) (*transport.Result, error) {
	if typ != SpotSideBuy && typ != SpotSideSell {
		return nil, validationErr("type must be 0 (buy) or 1 (sell)")
	}
	if _, err := requireNonEmpty(symbol, "symbol"); err != nil {
		return nil, err
	}
	if _, err := decimalArg(amount, "amount"); err != nil {
		return nil, err
	}
	params := map[string]interface{}{
		"symbol": symbol,
		"type":   typ,
		"amount": amount,
	}
	if opts.Rate != "" {
		if _, err := decimalArg(opts.Rate, "rate"); err != nil {
			return nil, err
		}
		params["rate"] = opts.Rate
	}
	if opts.IsMarket != nil {
		params["isMarket"] = *opts.IsMarket
	}
	if opts.OutTradeNo != "" {
		params["out_trade_no"] = opts.OutTradeNo
	}
	if opts.Funds != "" {
		if _, err := decimalArg(opts.Funds, "funds"); err != nil {
			return nil, err
		}
		params["funds"] = opts.Funds
	}
	return s.call(ctx, "spot.trade.placeOrder", params)
}

// CancelOrder — command=cancelOrder. MUTATION.
func (s *SpotTradeService) CancelOrder(ctx context.Context, orderID string) (*transport.Result, error) {
	return s.call(ctx, "spot.trade.cancelOrder", nonNull(map[string]interface{}{
		"orderId": emptyToNil(orderID),
	}))
}

// CancelAllOrders — command=cancelAllOrder. MUTATION.
// currencyPair is sent only when provided; the SDK does not widen the
// documented semantics to "all markets" on its own.
func (s *SpotTradeService) CancelAllOrders(ctx context.Context, currencyPair string) (*transport.Result, error) {
	return s.call(ctx, "spot.trade.cancelAllOrders", nonNull(map[string]interface{}{
		"currencyPair": emptyToNil(currencyPair),
	}))
}

func emptyToNil(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

func decimalArg(v, name string) (string, error) {
	if _, err := requireNonEmpty(v, name); err != nil {
		return "", err
	}
	// validate canonical decimal form without depending on the root package
	if !decimalPattern.MatchString(v) {
		return "", validationErr("invalid decimal string %q for %s — monetary values must be decimal strings", v, name)
	}
	return v, nil
}
