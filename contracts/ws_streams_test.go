package contracts_test

import (
	"github.com/tigusigalpa/coinw-go/model"
	ws "github.com/tigusigalpa/coinw-go/websocket"
	"testing"
)

func TestWSStreamContracts(t *testing.T) {
	t.Run("spot.market.ticker.socketio", func(t *testing.T) {
		checkStream(t, ws.SpotTickerSocketIO("BTC-USDT"), "fixtures/ws/spot.market.ticker.socketio.json", "", true, false)
	})
	t.Run("spot.market.ticker.raw", func(t *testing.T) {
		checkStream(t, ws.SpotTicker("78"), "fixtures/ws/spot.market.ticker.raw.json", "fixtures/ws/spot.market.ticker.raw.ack.json", false, false)
	})
	t.Run("spot.market.allTickers.raw", func(t *testing.T) {
		checkStream(t, ws.SpotAllTickers(), "fixtures/ws/spot.market.allTickers.raw.json", "fixtures/ws/spot.market.allTickers.raw.ack.json", false, false)
	})
	t.Run("spot.market.depth.raw", func(t *testing.T) {
		checkStream(t, ws.SpotIncrementalDepth("78"), "fixtures/ws/spot.market.depth.raw.json", "fixtures/ws/spot.market.depth.raw.ack.json", false, false)
	})
	t.Run("spot.market.candles.socketio", func(t *testing.T) {
		checkStream(t, ws.SpotCandlesSocketIO("BTC-USDT", "1m"), "fixtures/ws/spot.market.candles.socketio.json", "", true, false)
	})
	t.Run("spot.market.candles.raw", func(t *testing.T) {
		checkStream(t, ws.SpotCandles("78", "1m"), "fixtures/ws/spot.market.candles.raw.json", "fixtures/ws/spot.market.candles.raw.ack.json", false, false)
	})
	t.Run("spot.market.depthSnapshot.socketio", func(t *testing.T) {
		checkStream(t, ws.SpotDepthSocketIO("BTC-USDT"), "fixtures/ws/spot.market.depthSnapshot.socketio.json", "", true, false)
	})
	t.Run("spot.market.depthSnapshot.raw", func(t *testing.T) {
		checkStream(t, ws.SpotDepthSnapshot("78"), "fixtures/ws/spot.market.depthSnapshot.raw.json", "fixtures/ws/spot.market.depthSnapshot.raw.ack.json", false, false)
	})
	t.Run("spot.market.trades.socketio", func(t *testing.T) {
		checkStream(t, ws.SpotTradesSocketIO("BTC-USDT"), "fixtures/ws/spot.market.trades.socketio.json", "", true, false)
	})
	t.Run("spot.market.trades.raw", func(t *testing.T) {
		checkStream(t, ws.SpotTrades("78"), "fixtures/ws/spot.market.trades.raw.json", "fixtures/ws/spot.market.trades.raw.ack.json", false, false)
	})
	t.Run("spot.account.assets.raw", func(t *testing.T) {
		checkStream(t, ws.SpotAssets(), "fixtures/ws/spot.account.assets.raw.json", "fixtures/ws/spot.account.assets.raw.ack.json", false, true)
	})
	t.Run("spot.check.orders.raw", func(t *testing.T) {
		checkStream(t, ws.SpotOrders(), "fixtures/ws/spot.check.orders.raw.json", "fixtures/ws/spot.check.orders.raw.ack.json", false, true)
	})
	t.Run("futures.market.ticker.raw", func(t *testing.T) {
		instrument, _ := model.USDTInstrument("BTC")
		checkStream(t, ws.FuturesTicker(instrument), "fixtures/ws/futures.market.ticker.raw.json", "fixtures/ws/futures.market.ticker.raw.ack.json", false, false)
	})
	t.Run("futures.market.depth.raw", func(t *testing.T) {
		instrument, _ := model.USDTInstrument("BTC")
		checkStream(t, ws.FuturesDepth(instrument), "fixtures/ws/futures.market.depth.raw.json", "fixtures/ws/futures.market.depth.raw.ack.json", false, false)
	})
	t.Run("futures.market.trades.raw", func(t *testing.T) {
		instrument, _ := model.USDTInstrument("BTC")
		checkStream(t, ws.FuturesTrades(instrument), "fixtures/ws/futures.market.trades.raw.json", "fixtures/ws/futures.market.trades.raw.ack.json", false, false)
	})
	t.Run("futures.market.candlesUTC8.raw", func(t *testing.T) {
		instrument, _ := model.USDTInstrument("BTC")
		checkStream(t, ws.FuturesCandlesUTC8(instrument, "1"), "fixtures/ws/futures.market.candlesUTC8.raw.json", "fixtures/ws/futures.market.candlesUTC8.raw.ack.json", false, false)
	})
	t.Run("futures.market.candlesUTC0.raw", func(t *testing.T) {
		instrument, _ := model.USDTInstrument("BTC")
		checkStream(t, ws.FuturesCandlesUTC0(instrument, "1"), "fixtures/ws/futures.market.candlesUTC0.raw.json", "fixtures/ws/futures.market.candlesUTC0.raw.ack.json", false, false)
	})
	t.Run("futures.market.indexPrice.raw", func(t *testing.T) {
		instrument, _ := model.USDTInstrument("BTC")
		checkStream(t, ws.FuturesIndexPrice(instrument), "fixtures/ws/futures.market.indexPrice.raw.json", "fixtures/ws/futures.market.indexPrice.raw.ack.json", false, false)
	})
	t.Run("futures.market.markPrice.raw", func(t *testing.T) {
		instrument, _ := model.USDTInstrument("BTC")
		checkStream(t, ws.FuturesMarkPrice(instrument), "fixtures/ws/futures.market.markPrice.raw.json", "fixtures/ws/futures.market.markPrice.raw.ack.json", false, false)
	})
	t.Run("futures.market.funding.raw", func(t *testing.T) {
		instrument, _ := model.USDTInstrument("BTC")
		checkStream(t, ws.FuturesFundingRate(instrument), "fixtures/ws/futures.market.funding.raw.json", "fixtures/ws/futures.market.funding.raw.ack.json", false, false)
	})
	t.Run("futures.assets.assets.raw", func(t *testing.T) {
		checkStream(t, ws.FuturesAssets(), "fixtures/ws/futures.assets.assets.raw.json", "fixtures/ws/futures.assets.assets.raw.ack.json", false, true)
	})
	t.Run("futures.assets.assets_ag.raw", func(t *testing.T) {
		checkStream(t, ws.FuturesMegaCoupon(), "fixtures/ws/futures.assets.assets_ag.raw.json", "fixtures/ws/futures.assets.assets_ag.raw.ack.json", false, true)
	})
	t.Run("futures.assets.user_setting.raw", func(t *testing.T) {
		checkStream(t, ws.FuturesMarginMode(), "fixtures/ws/futures.assets.user_setting.raw.json", "fixtures/ws/futures.assets.user_setting.raw.ack.json", false, true)
	})
	t.Run("futures.check.order.raw", func(t *testing.T) {
		checkStream(t, ws.FuturesOrders(), "fixtures/ws/futures.check.order.raw.json", "fixtures/ws/futures.check.order.raw.ack.json", false, true)
	})
	t.Run("futures.position.position.raw", func(t *testing.T) {
		checkStream(t, ws.FuturesPositions(), "fixtures/ws/futures.position.position.raw.json", "fixtures/ws/futures.position.position.raw.ack.json", false, true)
	})
	t.Run("futures.position.position_change.raw", func(t *testing.T) {
		checkStream(t, ws.FuturesPositionChanges(), "fixtures/ws/futures.position.position_change.raw.json", "fixtures/ws/futures.position.position_change.raw.ack.json", false, true)
	})
}
