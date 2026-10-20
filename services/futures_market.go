package services

import (
	"context"

	"github.com/tigusigalpa/coinw-go/model"
	"github.com/tigusigalpa/coinw-go/transport"
)

// FuturesMarketService — futures market data.
type FuturesMarketService struct{ base }

func NewFuturesMarketService(ex *transport.Executor) *FuturesMarketService {
	return &FuturesMarketService{newBase(ex)}
}

// GetInstrumentInfo — GET /v1/perpum/instruments (public).
func (s *FuturesMarketService) GetInstrumentInfo(ctx context.Context, name string) (*transport.Response[[]model.FuturesInstrument], error) {
	return callTyped[[]model.FuturesInstrument](s.base, ctx, "futures.market.getInstrumentInfo", nonNull(map[string]interface{}{"name": emptyToNil(name)}))
}

// GetBatchInstrumentsInfo — GET /v1/perpum/instrumentList (public, max 20).
func (s *FuturesMarketService) GetBatchInstrumentsInfo(ctx context.Context, symbols []string) (*transport.Response[[]model.FuturesInstrument], error) {
	if _, err := requireMaxInstruments(symbols, 20); err != nil {
		return nil, err
	}
	return callTyped[[]model.FuturesInstrument](s.base, ctx, "futures.market.getBatchInstrumentsInfo",
		map[string]interface{}{"symbols": csv(symbols)})
}

// GetLastTradeSummaryOfAnInstrument — GET /v1/perpumPublic/ticker (public).
func (s *FuturesMarketService) GetLastTradeSummaryOfAnInstrument(ctx context.Context, instrument string) (*transport.Response[[]model.FuturesTicker], error) {
	if _, err := requireNonEmpty(instrument, "instrument"); err != nil {
		return nil, err
	}
	return callTyped[[]model.FuturesTicker](s.base, ctx, "futures.market.getLastTradeSummaryOfAnInstrument",
		map[string]interface{}{"instrument": instrument})
}

// GetBatchLastTradeSummary — GET /v1/perpumPublic/ticker/list (public).
func (s *FuturesMarketService) GetBatchLastTradeSummary(ctx context.Context, symbols []string) (*transport.Response[[]model.FuturesTicker], error) {
	if _, err := requireMaxInstruments(symbols, 20); err != nil {
		return nil, err
	}
	return callTyped[[]model.FuturesTicker](s.base, ctx, "futures.market.getBatchLastTradeSummary",
		map[string]interface{}{"symbols": csv(symbols)})
}

// GetLastTradeSummaryOfAllInstruments — GET /v1/perpumPublic/tickers (public).
func (s *FuturesMarketService) GetLastTradeSummaryOfAllInstruments(ctx context.Context) (*transport.Response[[]model.FuturesTicker], error) {
	return callTyped[[]model.FuturesTicker](s.base, ctx, "futures.market.getLastTradeSummaryOfAllInstruments", nil)
}

// GetKLineOfAnInstrument — GET /v1/perpumPublic/klines (public).
// granularity is the documented interval; limit=0 omits it.
func (s *FuturesMarketService) GetKLineOfAnInstrument(ctx context.Context, currencyCode, granularity string, limit int) (*transport.Response[[]model.FuturesCandleTuple], error) {
	if _, err := requireNonEmpty(currencyCode, "currencyCode"); err != nil {
		return nil, err
	}
	params := map[string]interface{}{"currencyCode": currencyCode, "granuality": granularity}
	if limit > 0 {
		params["limit"] = limit
	}
	return callTyped[[]model.FuturesCandleTuple](s.base, ctx, "futures.market.getKLineOfAnInstrument", params)
}

// GetLastSettlementFundingFeeRate — GET /v1/perpum/fundingRate (public).
func (s *FuturesMarketService) GetLastSettlementFundingFeeRate(ctx context.Context, instrument string) (*transport.Response[model.DecimalAtTime], error) {
	if _, err := requireNonEmpty(instrument, "instrument"); err != nil {
		return nil, err
	}
	return callTyped[model.DecimalAtTime](s.base, ctx, "futures.market.getLastSettlementFundingFeeRate",
		map[string]interface{}{"instrument": instrument})
}

// GetOrderBookOfAnInstrument — GET /v1/perpumPublic/depth (public).
func (s *FuturesMarketService) GetOrderBookOfAnInstrument(ctx context.Context, base string) (*transport.Response[model.FuturesOrderBook], error) {
	if _, err := requireNonEmpty(base, "base"); err != nil {
		return nil, err
	}
	return callTyped[model.FuturesOrderBook](s.base, ctx, "futures.market.getOrderBookOfAnInstrument",
		map[string]interface{}{"base": base})
}

// GetTradeDataOfAnInstrument — GET /v1/perpumPublic/trades (public).
func (s *FuturesMarketService) GetTradeDataOfAnInstrument(ctx context.Context, base string) (*transport.Response[[]model.FuturesTrade], error) {
	if _, err := requireNonEmpty(base, "base"); err != nil {
		return nil, err
	}
	return callTyped[[]model.FuturesTrade](s.base, ctx, "futures.market.getTradeDataOfAnInstrument",
		map[string]interface{}{"base": base})
}

// GetMarginRequirementsOfAllInstruments — GET /v1/perpum/ladders.
// Documented under market info but requires Futures private auth.
func (s *FuturesMarketService) GetMarginRequirementsOfAllInstruments(ctx context.Context) (*transport.Response[model.FuturesLadders], error) {
	return callTyped[model.FuturesLadders](s.base, ctx, "futures.market.getMarginRequirementsOfAllInstruments", nil)
}

// GetHistoricalPublicTrades — GET /v1/perpum/orders/trades.
// Its doc page declares authentication required — resolved to the Futures
// HMAC profile. page=0 omits; pageSize<=500 documented cap.
func (s *FuturesMarketService) GetHistoricalPublicTrades(ctx context.Context, instrument string, page, pageSize int) (*transport.Response[model.RowsPage[model.HistoricalPublicTrade]], error) {
	if _, err := requireNonEmpty(instrument, "instrument"); err != nil {
		return nil, err
	}
	params := map[string]interface{}{"instrument": instrument}
	if page > 0 {
		params["page"] = page
	}
	if pageSize > 0 {
		if _, err := requirePageSize(pageSize, 500); err != nil {
			return nil, err
		}
		params["pageSize"] = pageSize
	}
	return callTyped[model.RowsPage[model.HistoricalPublicTrade]](s.base, ctx, "futures.market.getHistoricalPublicTrades", params)
}
