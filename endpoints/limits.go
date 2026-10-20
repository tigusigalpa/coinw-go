package endpoints

// Frequency Limit sections checked on 2026-10-20 at
// https://www.coinw.com/api-doc/en/. Rates below are requests per second;
// e.g. 5 requests/2s = 2.5, and 1 request/5s = 0.2.
// Explicit scopes preserve UID-only Spot limits and asymmetric Futures limits.
func applyDocumentedLimits(defs map[string]Definition) {
	set := func(id string, ip, uid float64) {
		d := defs[id]
		d.PerIPRPS, d.PerUIDRPS = ip, uid
		d.PerInterfaceRPS = ip
		if ip == 0 || (uid > 0 && uid < ip) {
			d.PerInterfaceRPS = uid
		}
		defs[id] = d
	}
	for id, rate := range map[string]float64{
		"common.account.cancelWithdrawal":            3,
		"common.account.getDepositWithdrawalAddress": 3,
		"common.account.getDepositWithdrawalHistory": 3,
		"common.account.initiateWithdrawal":          3,
		"common.account.transferAssets":              1,
		"common.subAccount.getSubAccountList":        3,
		"common.subAccount.getSubAccountApiKeys":     3,
		"common.subAccount.resetSubAccountApiKey":    3,
		"common.subAccount.getMaxTransferableAmount": 3,
		"common.subAccount.getTransferHistory":       3,
		"common.subAccount.accountTransfer":          3,
		"spot.account.getSpotAccountBalance":         3,
		"spot.account.getFullSpotAccountBalance":     2.5,
		"spot.check.getBatchHistoricalOrders":        0.5,
		"spot.check.getCurrentOrders":                15,
		"spot.check.getHistoricalOrders":             5,
		"spot.check.getOrderDetails":                 15,
		"spot.check.getOrderStatus":                  15,
		"spot.check.getTransactionHistory":           2.5,
		"spot.trade.placeOrder":                      15,
		"spot.trade.cancelOrder":                     15,
		"spot.trade.cancelAllOrders":                 15,
	} {
		set(id, 0, rate)
	}
	for id, rate := range map[string]float64{
		"spot.market.getInstrumentInformation":            80,
		"spot.market.get24hTradeSummaryForAllInstruments": 80,
		"spot.market.getOrderBook":                        10,
		"spot.market.getRecentTrades":                     10,
		"spot.market.getKLine":                            10,
		"spot.market.get24hVolumeForPopularInstruments":   80,
		"spot.market.getCoinDepositWithdrawalLimits":      80,
	} {
		set(id, rate, 0)
	}
	for id, rate := range map[string]float64{
		"futures.market.getInstrumentInfo":                       5,
		"futures.market.getKLineOfAnInstrument":                  10,
		"futures.market.getBatchInstrumentsInfo":                 2,
		"futures.market.getLastTradeSummaryOfAnInstrument":       5,
		"futures.market.getBatchLastTradeSummary":                0.2,
		"futures.market.getLastTradeSummaryOfAllInstruments":     5,
		"futures.market.getLastSettlementFundingFeeRate":         8,
		"futures.market.getOrderBookOfAnInstrument":              5,
		"futures.market.getTradeDataOfAnInstrument":              5,
		"futures.market.getMarginRequirementsOfAllInstruments":   5,
		"futures.market.getHistoricalPublicTrades":               5,
		"futures.assets.coinConversion":                          5,
		"futures.assets.enableDisableMegaCoupon":                 0.5,
		"futures.assets.getFuturesAccountAssets":                 2.5,
		"futures.assets.getFuturesAccountFees":                   1,
		"futures.assets.getMarginMode":                           5,
		"futures.assets.getMaxTransferableBalance":               5,
		"futures.assets.getMaxUsersContractSize":                 5,
		"futures.assets.getMegaCouponBalance":                    0.5,
		"futures.assets.getTransactionDetails3Days":              10,
		"futures.assets.getTransactionDetails3Months":            5,
		"futures.assets.setMarginMode":                           1,
		"futures.common.getLeverageInformation":                  10,
		"futures.check.getBatchHistoricalOrders3Months":          1,
		"futures.check.getBatchHistoricalOrders7Days":            1,
		"futures.check.getCurrentOrders":                         10,
		"futures.check.getHistoricalOrders3Months":               5,
		"futures.check.getHistoricalOrders7Days":                 2.5,
		"futures.check.getOrderInformation":                      10,
		"futures.check.getPendingOrderCount":                     5,
		"futures.check.getStopLossTakeProfitInformation":         5,
		"futures.check.getTrailingStopLossTakeProfitInformation": 5,
		"futures.position.getCurrentPositionInformation":         10,
		"futures.position.getMaxContractSize":                    8,
		"futures.position.getPositionMarginRate":                 5,
		"futures.trade.adjustMargin":                             5,
		"futures.trade.cancelAnOrder":                            10,
		"futures.trade.cancelBatchOrders":                        10,
		"futures.trade.closeBatchPositions":                      5,
		"futures.trade.closePositionsAtMarketPrice":              5,
		"futures.trade.modifyAnOrder":                            10,
		"futures.trade.modifyStopLossTakeProfitInBatches":        5,
		"futures.trade.placeBatchOrders":                         10,
		"futures.trade.reverseAPosition":                         5,
		"futures.trade.setStopLossTakeProfit":                    5,
		"futures.trade.setStopLossTakeProfitInBatches":           5,
		"futures.trade.setTrailingStopLossTakeProfit":            5,
	} {
		set(id, rate, rate)
	}
	set("futures.position.getAllCurrentPositions", 50000, 30)
	set("futures.position.getHistoricalPositionInformation", 100000, 15)
	set("futures.trade.closeAPosition", 100000, 30)
	set("futures.trade.placeAnOrder", 100000, 30)
}
