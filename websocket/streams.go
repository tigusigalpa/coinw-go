package websocket

import "github.com/tigusigalpa/coinw-go/model"

// SpotTickerSocketIO — https://www.coinw.com/api-doc/en/spot-trading/market/subscribe-24h-trade-summary
func SpotTickerSocketIO(pair string) Stream[model.JSONEncoded[model.SpotTickerEvent]] {
	return Stream[model.JSONEncoded[model.SpotTickerEvent]]{id: "spot.market.ticker.socketio", kind: "ticker", profile: SpotProfile, private: false, pair: pair, interval: "", args: "spot/market-api-ticker" + ":" + pair, needsPair: true, needsInterval: false}
}

// SpotTicker — https://www.coinw.com/api-doc/en/spot-trading/market/subscribe-24h-trade-summary
func SpotTicker(pair string) Stream[model.JSONEncoded[model.SpotTickerEvent]] {
	return Stream[model.JSONEncoded[model.SpotTickerEvent]]{id: "spot.market.ticker.raw", kind: "ticker", profile: SpotProfile, private: false, pair: pair, interval: "", args: "", needsPair: true, needsInterval: false}
}

// SpotAllTickers — https://www.coinw.com/api-doc/en/spot-trading/market/subscribe-24h-trades-summary-for-all-instruments
func SpotAllTickers() Stream[model.JSONEncoded[[]model.SpotAllTickerEvent]] {
	return Stream[model.JSONEncoded[[]model.SpotAllTickerEvent]]{id: "spot.market.allTickers.raw", kind: "ticker_all", profile: SpotProfile, private: false, pair: "", interval: "", args: "", needsPair: false, needsInterval: false}
}

// SpotIncrementalDepth — https://www.coinw.com/api-doc/en/spot-trading/market/subscribe-incremental-order-book
func SpotIncrementalDepth(pair string) Stream[model.JSONEncoded[model.SpotDepthIncrement]] {
	return Stream[model.JSONEncoded[model.SpotDepthIncrement]]{id: "spot.market.depth.raw", kind: "depth", profile: SpotProfile, private: false, pair: pair, interval: "", args: "", needsPair: true, needsInterval: false}
}

// SpotCandlesSocketIO — https://www.coinw.com/api-doc/en/spot-trading/market/subscribe-k-line
func SpotCandlesSocketIO(pair string, interval string) Stream[model.JSONEncoded[model.SpotCandleTuple]] {
	return Stream[model.JSONEncoded[model.SpotCandleTuple]]{id: "spot.market.candles.socketio", kind: "candles", profile: SpotProfile, private: false, pair: pair, interval: interval, args: "spot/candle-" + interval + ":" + pair, needsPair: true, needsInterval: true}
}

// SpotCandles — https://www.coinw.com/api-doc/en/spot-trading/market/subscribe-k-line
func SpotCandles(pair string, interval string) Stream[model.JSONEncoded[model.SpotCandleTuple]] {
	return Stream[model.JSONEncoded[model.SpotCandleTuple]]{id: "spot.market.candles.raw", kind: "candles", profile: SpotProfile, private: false, pair: pair, interval: interval, args: "", needsPair: true, needsInterval: true}
}

// SpotDepthSocketIO — https://www.coinw.com/api-doc/en/spot-trading/market/subscribe-order-book
func SpotDepthSocketIO(pair string) Stream[model.JSONEncoded[model.SpotDepthSnapshot]] {
	return Stream[model.JSONEncoded[model.SpotDepthSnapshot]]{id: "spot.market.depthSnapshot.socketio", kind: "depth_snapshot", profile: SpotProfile, private: false, pair: pair, interval: "", args: "spot/level2_20" + ":" + pair, needsPair: true, needsInterval: false}
}

// SpotDepthSnapshot — https://www.coinw.com/api-doc/en/spot-trading/market/subscribe-order-book
func SpotDepthSnapshot(pair string) Stream[model.JSONEncoded[model.SpotDepthSnapshot]] {
	return Stream[model.JSONEncoded[model.SpotDepthSnapshot]]{id: "spot.market.depthSnapshot.raw", kind: "depth_snapshot", profile: SpotProfile, private: false, pair: pair, interval: "", args: "", needsPair: true, needsInterval: false}
}

// SpotTradesSocketIO — https://www.coinw.com/api-doc/en/spot-trading/market/subscribe-trades
func SpotTradesSocketIO(pair string) Stream[model.JSONEncoded[[]model.SpotStreamTrade]] {
	return Stream[model.JSONEncoded[[]model.SpotStreamTrade]]{id: "spot.market.trades.socketio", kind: "fills", profile: SpotProfile, private: false, pair: pair, interval: "", args: "spot/match" + ":" + pair, needsPair: true, needsInterval: false}
}

// SpotTrades — https://www.coinw.com/api-doc/en/spot-trading/market/subscribe-trades
func SpotTrades(pair string) Stream[model.JSONEncoded[[]model.SpotStreamTrade]] {
	return Stream[model.JSONEncoded[[]model.SpotStreamTrade]]{id: "spot.market.trades.raw", kind: "fills", profile: SpotProfile, private: false, pair: pair, interval: "", args: "", needsPair: true, needsInterval: false}
}

// SpotAssets — https://www.coinw.com/api-doc/en/spot-trading/account/subscribe-assets
func SpotAssets() Stream[model.ObjectOrArray[model.SpotAssetEvent]] {
	return Stream[model.ObjectOrArray[model.SpotAssetEvent]]{id: "spot.account.assets.raw", kind: "assets", profile: SpotProfile, private: true, pair: "", interval: "", args: "", needsPair: false, needsInterval: false}
}

// SpotOrders — https://www.coinw.com/api-doc/en/spot-trading/check/subscribe-current-orders
func SpotOrders() Stream[model.ObjectOrArray[model.SpotOrderEvent]] {
	return Stream[model.ObjectOrArray[model.SpotOrderEvent]]{id: "spot.check.orders.raw", kind: "order", profile: SpotProfile, private: true, pair: "", interval: "", args: "", needsPair: false, needsInterval: false}
}

// FuturesTicker — https://www.coinw.com/api-doc/en/futures-trading/market/subscribe-24h-trade-summary
func FuturesTicker(instrument model.Instrument) Stream[model.FuturesTickerEvent] {
	return Stream[model.FuturesTickerEvent]{id: "futures.market.ticker.raw", kind: "ticker_swap", profile: FuturesProfile, private: false, pair: instrument.String(), interval: "", args: "", needsPair: true, needsInterval: false}
}

// FuturesDepth — https://www.coinw.com/api-doc/en/futures-trading/market/subscribe-order-book
func FuturesDepth(instrument model.Instrument) Stream[model.FuturesDepthEvent] {
	return Stream[model.FuturesDepthEvent]{id: "futures.market.depth.raw", kind: "depth", profile: FuturesProfile, private: false, pair: instrument.String(), interval: "", args: "", needsPair: true, needsInterval: false}
}

// FuturesTrades — https://www.coinw.com/api-doc/en/futures-trading/market/subscribe-trade-data
func FuturesTrades(instrument model.Instrument) Stream[[]model.FuturesTrade] {
	return Stream[[]model.FuturesTrade]{id: "futures.market.trades.raw", kind: "fills", profile: FuturesProfile, private: false, pair: instrument.String(), interval: "", args: "", needsPair: true, needsInterval: false}
}

// FuturesCandlesUTC8 — https://www.coinw.com/api-doc/en/futures-trading/market/subscribe-k-line-utc+8
func FuturesCandlesUTC8(instrument model.Instrument, interval string) Stream[model.FuturesStreamCandle] {
	return Stream[model.FuturesStreamCandle]{id: "futures.market.candlesUTC8.raw", kind: "candles_swap", profile: FuturesProfile, private: false, pair: instrument.String(), interval: interval, args: "", needsPair: true, needsInterval: true}
}

// FuturesCandlesUTC0 — https://www.coinw.com/api-doc/en/futures-trading/market/subscribe-k-line-utc+0
func FuturesCandlesUTC0(instrument model.Instrument, interval string) Stream[model.FuturesStreamCandle] {
	return Stream[model.FuturesStreamCandle]{id: "futures.market.candlesUTC0.raw", kind: "candles_swap_utc", profile: FuturesProfile, private: false, pair: instrument.String(), interval: interval, args: "", needsPair: true, needsInterval: true}
}

// FuturesIndexPrice — https://www.coinw.com/api-doc/en/futures-trading/market/subscribe-index-price
func FuturesIndexPrice(instrument model.Instrument) Stream[model.FuturesPriceEvent] {
	return Stream[model.FuturesPriceEvent]{id: "futures.market.indexPrice.raw", kind: "index_price", profile: FuturesProfile, private: false, pair: instrument.String(), interval: "", args: "", needsPair: true, needsInterval: false}
}

// FuturesMarkPrice — https://www.coinw.com/api-doc/en/futures-trading/market/subscribe-mark-price
func FuturesMarkPrice(instrument model.Instrument) Stream[model.FuturesPriceEvent] {
	return Stream[model.FuturesPriceEvent]{id: "futures.market.markPrice.raw", kind: "mark_price", profile: FuturesProfile, private: false, pair: instrument.String(), interval: "", args: "", needsPair: true, needsInterval: false}
}

// FuturesFundingRate — https://www.coinw.com/api-doc/en/futures-trading/market/subscribe-funding-fee-rate
func FuturesFundingRate(instrument model.Instrument) Stream[model.FuturesFundingEvent] {
	return Stream[model.FuturesFundingEvent]{id: "futures.market.funding.raw", kind: "funding_rate", profile: FuturesProfile, private: false, pair: instrument.String(), interval: "", args: "", needsPair: true, needsInterval: false}
}

// FuturesAssets — https://www.coinw.com/api-doc/en/futures-trading/assets/subscribe-assets
func FuturesAssets() Stream[[]model.FuturesAssetEvent] {
	return Stream[[]model.FuturesAssetEvent]{id: "futures.assets.assets.raw", kind: "assets", profile: FuturesProfile, private: true, pair: "", interval: "", args: "", needsPair: false, needsInterval: false}
}

// FuturesMegaCoupon — https://www.coinw.com/api-doc/en/futures-trading/assets/subscribe-mega-coupon
func FuturesMegaCoupon() Stream[[]model.MegaCoupon] {
	return Stream[[]model.MegaCoupon]{id: "futures.assets.assets_ag.raw", kind: "assets_ag", profile: FuturesProfile, private: true, pair: "", interval: "", args: "", needsPair: false, needsInterval: false}
}

// FuturesMarginMode — https://www.coinw.com/api-doc/en/futures-trading/assets/subscribe-margin-mode
func FuturesMarginMode() Stream[[]model.MarginMode] {
	return Stream[[]model.MarginMode]{id: "futures.assets.user_setting.raw", kind: "user_setting", profile: FuturesProfile, private: true, pair: "", interval: "", args: "", needsPair: false, needsInterval: false}
}

// FuturesOrders — https://www.coinw.com/api-doc/en/futures-trading/check/subscribe-current-orders
func FuturesOrders() Stream[[]model.FuturesOrder] {
	return Stream[[]model.FuturesOrder]{id: "futures.check.order.raw", kind: "order", profile: FuturesProfile, private: true, pair: "", interval: "", args: "", needsPair: false, needsInterval: false}
}

// FuturesPositions — https://www.coinw.com/api-doc/en/futures-trading/position/subscribe-position
func FuturesPositions() Stream[[]model.FuturesPosition] {
	return Stream[[]model.FuturesPosition]{id: "futures.position.position.raw", kind: "position", profile: FuturesProfile, private: true, pair: "", interval: "", args: "", needsPair: false, needsInterval: false}
}

// FuturesPositionChanges — https://www.coinw.com/api-doc/en/futures-trading/position/subscribe-position-change
func FuturesPositionChanges() Stream[[]model.FuturesPositionChange] {
	return Stream[[]model.FuturesPositionChange]{id: "futures.position.position_change.raw", kind: "position_change", profile: FuturesProfile, private: true, pair: "", interval: "", args: "", needsPair: false, needsInterval: false}
}
