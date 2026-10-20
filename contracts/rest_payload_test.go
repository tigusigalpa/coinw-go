package contracts_test

import (
	"encoding/json"
	"github.com/tigusigalpa/coinw-go/model"
	"github.com/tigusigalpa/coinw-go/transport"
	"os"
	"testing"
)

func TestRESTPayloadContracts(t *testing.T) {
	t.Run("common.account.cancelWithdrawal", func(t *testing.T) {
		checkPayload[model.NoData](t, "fixtures/rest/common.account.cancelWithdrawal.json")
	})
	t.Run("common.account.getDepositWithdrawalAddress", func(t *testing.T) {
		checkPayload[[]model.DepositAddress](t, "fixtures/rest/common.account.getDepositWithdrawalAddress.json")
	})
	t.Run("common.account.getDepositWithdrawalHistory", func(t *testing.T) {
		checkPayload[[]model.DepositWithdrawal](t, "fixtures/rest/common.account.getDepositWithdrawalHistory.json")
	})
	t.Run("common.account.initiateWithdrawal", func(t *testing.T) {
		checkPayload[model.WithdrawalAcknowledgement](t, "fixtures/rest/common.account.initiateWithdrawal.json")
	})
	t.Run("common.account.transferAssets", func(t *testing.T) { checkPayload[bool](t, "fixtures/rest/common.account.transferAssets.json") })
	t.Run("common.subAccount.getSubAccountList", func(t *testing.T) {
		checkPayload[[]model.SubAccount](t, "fixtures/rest/common.subAccount.getSubAccountList.json")
	})
	t.Run("common.subAccount.getSubAccountApiKeys", func(t *testing.T) {
		checkPayload[[]model.SubAccountAPIKey](t, "fixtures/rest/common.subAccount.getSubAccountApiKeys.json")
	})
	t.Run("common.subAccount.resetSubAccountApiKey", func(t *testing.T) {
		checkPayload[model.NoData](t, "fixtures/rest/common.subAccount.resetSubAccountApiKey.json")
	})
	t.Run("common.subAccount.getMaxTransferableAmount", func(t *testing.T) {
		checkPayload[model.WireDecimal](t, "fixtures/rest/common.subAccount.getMaxTransferableAmount.json")
	})
	t.Run("common.subAccount.getTransferHistory", func(t *testing.T) {
		checkPayload[[]model.Transfer](t, "fixtures/rest/common.subAccount.getTransferHistory.json")
	})
	t.Run("common.subAccount.accountTransfer", func(t *testing.T) {
		checkPayload[model.TransferAcknowledgement](t, "fixtures/rest/common.subAccount.accountTransfer.json")
	})
	t.Run("common.affiliate.getInvitedUserList", func(t *testing.T) {
		checkPayload[model.ListPage[model.InvitedUser]](t, "fixtures/rest/common.affiliate.getInvitedUserList.json")
	})
	t.Run("common.affiliate.getInvitedUsersKycResults", func(t *testing.T) {
		checkPayload[model.ListPage[model.InvitedUserKYC]](t, "fixtures/rest/common.affiliate.getInvitedUsersKycResults.json")
	})
	t.Run("common.affiliate.getInvitedUsersTradingVolume", func(t *testing.T) {
		checkPayload[model.ListPage[model.InvitedUserVolume]](t, "fixtures/rest/common.affiliate.getInvitedUsersTradingVolume.json")
	})
	t.Run("common.affiliate.getAffiliateUserActivityData", func(t *testing.T) {
		checkPayload[model.ListPage[model.AffiliateActivity]](t, "fixtures/rest/common.affiliate.getAffiliateUserActivityData.json")
	})
	t.Run("common.affiliate.getInvitedUserAssets", func(t *testing.T) {
		checkPayload[model.ListPage[model.InvitedUserAssets]](t, "fixtures/rest/common.affiliate.getInvitedUserAssets.json")
	})
	t.Run("common.affiliate.getAffiliateTeamDirectUserRange", func(t *testing.T) {
		checkPayload[model.AffiliateTeam](t, "fixtures/rest/common.affiliate.getAffiliateTeamDirectUserRange.json")
	})
	t.Run("spot.market.getInstrumentInformation", func(t *testing.T) {
		checkPayload[[]model.SpotInstrument](t, "fixtures/rest/spot.market.getInstrumentInformation.json")
	})
	t.Run("spot.market.get24hTradeSummaryForAllInstruments", func(t *testing.T) {
		checkPayload[map[string]model.SpotTicker](t, "fixtures/rest/spot.market.get24hTradeSummaryForAllInstruments.json")
	})
	t.Run("spot.market.getCoinDepositWithdrawalLimits", func(t *testing.T) {
		checkPayload[map[string]model.SpotCurrency](t, "fixtures/rest/spot.market.getCoinDepositWithdrawalLimits.json")
	})
	t.Run("spot.market.getOrderBook", func(t *testing.T) {
		checkPayload[[]model.SpotOrderBook](t, "fixtures/rest/spot.market.getOrderBook.json")
	})
	t.Run("spot.market.getRecentTrades", func(t *testing.T) {
		checkPayload[[]model.SpotTrade](t, "fixtures/rest/spot.market.getRecentTrades.json")
	})
	t.Run("spot.market.getKLine", func(t *testing.T) { checkPayload[[]model.SpotCandle](t, "fixtures/rest/spot.market.getKLine.json") })
	t.Run("spot.market.get24hVolumeForPopularInstruments", func(t *testing.T) {
		checkPayload[model.PopularVolumes](t, "fixtures/rest/spot.market.get24hVolumeForPopularInstruments.json")
	})
	t.Run("spot.account.getSpotAccountBalance", func(t *testing.T) {
		checkPayload[map[string]model.WireDecimal](t, "fixtures/rest/spot.account.getSpotAccountBalance.json")
	})
	t.Run("spot.account.getFullSpotAccountBalance", func(t *testing.T) {
		checkPayload[map[string]model.SpotFullBalance](t, "fixtures/rest/spot.account.getFullSpotAccountBalance.json")
	})
	t.Run("spot.check.getBatchHistoricalOrders", func(t *testing.T) {
		checkPayload[model.ObjectOrArray[model.SpotBatchHistoricalOrder]](t, "fixtures/rest/spot.check.getBatchHistoricalOrders.json")
	})
	t.Run("spot.check.getCurrentOrders", func(t *testing.T) {
		checkPayload[[]model.SpotCurrentOrder](t, "fixtures/rest/spot.check.getCurrentOrders.json")
	})
	t.Run("spot.check.getHistoricalOrders", func(t *testing.T) {
		checkPayload[model.SpotHistoryPage](t, "fixtures/rest/spot.check.getHistoricalOrders.json")
	})
	t.Run("spot.check.getOrderDetails", func(t *testing.T) {
		checkPayload[model.SpotOrderDetails](t, "fixtures/rest/spot.check.getOrderDetails.json")
	})
	t.Run("spot.check.getOrderStatus", func(t *testing.T) {
		checkPayload[model.SpotOrderStatus](t, "fixtures/rest/spot.check.getOrderStatus.json")
	})
	t.Run("spot.check.getTransactionHistory", func(t *testing.T) {
		checkPayload[[]model.SpotTransaction](t, "fixtures/rest/spot.check.getTransactionHistory.json")
	})
	t.Run("spot.trade.placeOrder", func(t *testing.T) {
		checkPayload[model.SpotOrderAcknowledgement](t, "fixtures/rest/spot.trade.placeOrder.json")
	})
	t.Run("spot.trade.cancelOrder", func(t *testing.T) {
		checkPayload[model.SpotCancelAcknowledgement](t, "fixtures/rest/spot.trade.cancelOrder.json")
	})
	t.Run("spot.trade.cancelAllOrders", func(t *testing.T) { checkPayload[model.NoData](t, "fixtures/rest/spot.trade.cancelAllOrders.json") })
	t.Run("futures.market.getInstrumentInfo", func(t *testing.T) {
		checkPayload[[]model.FuturesInstrument](t, "fixtures/rest/futures.market.getInstrumentInfo.json")
	})
	t.Run("futures.market.getBatchInstrumentsInfo", func(t *testing.T) {
		checkPayload[[]model.FuturesInstrument](t, "fixtures/rest/futures.market.getBatchInstrumentsInfo.json")
	})
	t.Run("futures.market.getLastTradeSummaryOfAnInstrument", func(t *testing.T) {
		checkPayload[[]model.FuturesTicker](t, "fixtures/rest/futures.market.getLastTradeSummaryOfAnInstrument.json")
	})
	t.Run("futures.market.getBatchLastTradeSummary", func(t *testing.T) {
		checkPayload[[]model.FuturesTicker](t, "fixtures/rest/futures.market.getBatchLastTradeSummary.json")
	})
	t.Run("futures.market.getLastTradeSummaryOfAllInstruments", func(t *testing.T) {
		checkPayload[[]model.FuturesTicker](t, "fixtures/rest/futures.market.getLastTradeSummaryOfAllInstruments.json")
	})
	t.Run("futures.market.getKLineOfAnInstrument", func(t *testing.T) {
		checkPayload[[]model.FuturesCandleTuple](t, "fixtures/rest/futures.market.getKLineOfAnInstrument.json")
	})
	t.Run("futures.market.getLastSettlementFundingFeeRate", func(t *testing.T) {
		checkPayload[model.DecimalAtTime](t, "fixtures/rest/futures.market.getLastSettlementFundingFeeRate.json")
	})
	t.Run("futures.market.getOrderBookOfAnInstrument", func(t *testing.T) {
		checkPayload[model.FuturesOrderBook](t, "fixtures/rest/futures.market.getOrderBookOfAnInstrument.json")
	})
	t.Run("futures.market.getTradeDataOfAnInstrument", func(t *testing.T) {
		checkPayload[[]model.FuturesTrade](t, "fixtures/rest/futures.market.getTradeDataOfAnInstrument.json")
	})
	t.Run("futures.market.getMarginRequirementsOfAllInstruments", func(t *testing.T) {
		checkPayload[model.FuturesLadders](t, "fixtures/rest/futures.market.getMarginRequirementsOfAllInstruments.json")
	})
	t.Run("futures.market.getHistoricalPublicTrades", func(t *testing.T) {
		checkPayload[model.RowsPage[model.HistoricalPublicTrade]](t, "fixtures/rest/futures.market.getHistoricalPublicTrades.json")
	})
	t.Run("futures.assets.coinConversion", func(t *testing.T) {
		checkPayload[model.DecimalValue](t, "fixtures/rest/futures.assets.coinConversion.json")
	})
	t.Run("futures.assets.enableDisableMegaCoupon", func(t *testing.T) {
		checkPayload[model.NoData](t, "fixtures/rest/futures.assets.enableDisableMegaCoupon.json")
	})
	t.Run("futures.assets.getFuturesAccountAssets", func(t *testing.T) {
		checkPayload[model.FuturesAccountAssets](t, "fixtures/rest/futures.assets.getFuturesAccountAssets.json")
	})
	t.Run("futures.assets.getFuturesAccountFees", func(t *testing.T) {
		checkPayload[model.FuturesAccountFees](t, "fixtures/rest/futures.assets.getFuturesAccountFees.json")
	})
	t.Run("futures.assets.getMarginMode", func(t *testing.T) {
		checkPayload[model.MarginMode](t, "fixtures/rest/futures.assets.getMarginMode.json")
	})
	t.Run("futures.assets.getMaxTransferableBalance", func(t *testing.T) {
		checkPayload[model.ValueOrScalar[model.WireDecimal]](t, "fixtures/rest/futures.assets.getMaxTransferableBalance.json")
	})
	t.Run("futures.assets.getMaxUsersContractSize", func(t *testing.T) {
		checkPayload[model.AvailableContractSize](t, "fixtures/rest/futures.assets.getMaxUsersContractSize.json")
	})
	t.Run("futures.assets.getMegaCouponBalance", func(t *testing.T) {
		checkPayload[[]model.MegaCoupon](t, "fixtures/rest/futures.assets.getMegaCouponBalance.json")
	})
	t.Run("futures.assets.getTransactionDetails3Days", func(t *testing.T) {
		checkPayload[model.RowsPage[model.FuturesTransaction]](t, "fixtures/rest/futures.assets.getTransactionDetails3Days.json")
	})
	t.Run("futures.assets.getTransactionDetails3Months", func(t *testing.T) {
		checkPayload[model.RowsPage[model.FuturesTransaction]](t, "fixtures/rest/futures.assets.getTransactionDetails3Months.json")
	})
	t.Run("futures.assets.setMarginMode", func(t *testing.T) { checkPayload[string](t, "fixtures/rest/futures.assets.setMarginMode.json") })
	t.Run("futures.common.getLeverageInformation", func(t *testing.T) {
		checkPayload[model.ValueOrScalar[model.WireDecimal]](t, "fixtures/rest/futures.common.getLeverageInformation.json")
	})
	t.Run("futures.position.getCurrentPositionInformation", func(t *testing.T) {
		checkPayload[[]model.FuturesPosition](t, "fixtures/rest/futures.position.getCurrentPositionInformation.json")
	})
	t.Run("futures.position.getAllCurrentPositions", func(t *testing.T) {
		checkPayload[[]model.FuturesPosition](t, "fixtures/rest/futures.position.getAllCurrentPositions.json")
	})
	t.Run("futures.position.getHistoricalPositionInformation", func(t *testing.T) {
		checkPayload[model.RowsPage[model.HistoricalPosition]](t, "fixtures/rest/futures.position.getHistoricalPositionInformation.json")
	})
	t.Run("futures.position.getMaxContractSize", func(t *testing.T) {
		checkPayload[model.MaximumContractSize](t, "fixtures/rest/futures.position.getMaxContractSize.json")
	})
	t.Run("futures.position.getPositionMarginRate", func(t *testing.T) {
		checkPayload[model.ValueOrScalar[model.WireDecimal]](t, "fixtures/rest/futures.position.getPositionMarginRate.json")
	})
	t.Run("futures.check.getBatchHistoricalOrders7Days", func(t *testing.T) {
		checkPayload[model.RowsPage[model.FuturesOrder]](t, "fixtures/rest/futures.check.getBatchHistoricalOrders7Days.json")
	})
	t.Run("futures.check.getBatchHistoricalOrders3Months", func(t *testing.T) {
		checkPayload[model.RowsPage[model.FuturesOrder]](t, "fixtures/rest/futures.check.getBatchHistoricalOrders3Months.json")
	})
	t.Run("futures.check.getCurrentOrders", func(t *testing.T) {
		checkPayload[model.RowsPage[model.FuturesOrder]](t, "fixtures/rest/futures.check.getCurrentOrders.json")
	})
	t.Run("futures.check.getHistoricalOrders7Days", func(t *testing.T) {
		checkPayload[model.RowsPage[model.FuturesOrder]](t, "fixtures/rest/futures.check.getHistoricalOrders7Days.json")
	})
	t.Run("futures.check.getHistoricalOrders3Months", func(t *testing.T) {
		checkPayload[model.RowsPage[model.FuturesOrder]](t, "fixtures/rest/futures.check.getHistoricalOrders3Months.json")
	})
	t.Run("futures.check.getOrderInformation", func(t *testing.T) {
		checkPayload[[]model.FuturesOrder](t, "fixtures/rest/futures.check.getOrderInformation.json")
	})
	t.Run("futures.check.getPendingOrderCount", func(t *testing.T) {
		checkPayload[model.ValueOrScalar[model.StringOrInteger]](t, "fixtures/rest/futures.check.getPendingOrderCount.json")
	})
	t.Run("futures.check.getStopLossTakeProfitInformation", func(t *testing.T) {
		checkPayload[[]model.StopLossTakeProfit](t, "fixtures/rest/futures.check.getStopLossTakeProfitInformation.json")
	})
	t.Run("futures.check.getTrailingStopLossTakeProfitInformation", func(t *testing.T) {
		checkPayload[[]model.TrailingStopLossTakeProfit](t, "fixtures/rest/futures.check.getTrailingStopLossTakeProfitInformation.json")
	})
	t.Run("futures.trade.adjustMargin", func(t *testing.T) { checkPayload[model.NoData](t, "fixtures/rest/futures.trade.adjustMargin.json") })
	t.Run("futures.trade.cancelAnOrder", func(t *testing.T) { checkPayload[model.NoData](t, "fixtures/rest/futures.trade.cancelAnOrder.json") })
	t.Run("futures.trade.cancelBatchOrders", func(t *testing.T) {
		checkPayload[model.NoData](t, "fixtures/rest/futures.trade.cancelBatchOrders.json")
	})
	t.Run("futures.trade.closeAPosition", func(t *testing.T) {
		checkPayload[model.ValueOrScalar[model.StringOrInteger]](t, "fixtures/rest/futures.trade.closeAPosition.json")
	})
	t.Run("futures.trade.closeBatchPositions", func(t *testing.T) {
		checkPayload[[]model.BatchOrderAcknowledgement](t, "fixtures/rest/futures.trade.closeBatchPositions.json")
	})
	t.Run("futures.trade.closePositionsAtMarketPrice", func(t *testing.T) {
		checkPayload[model.NoData](t, "fixtures/rest/futures.trade.closePositionsAtMarketPrice.json")
	})
	t.Run("futures.trade.modifyAnOrder", func(t *testing.T) {
		checkPayload[model.ModifiedOrderAcknowledgement](t, "fixtures/rest/futures.trade.modifyAnOrder.json")
	})
	t.Run("futures.trade.modifyStopLossTakeProfitInBatches", func(t *testing.T) {
		checkPayload[model.NoData](t, "fixtures/rest/futures.trade.modifyStopLossTakeProfitInBatches.json")
	})
	t.Run("futures.trade.placeAnOrder", func(t *testing.T) {
		checkPayload[model.ValueOrScalar[model.StringOrInteger]](t, "fixtures/rest/futures.trade.placeAnOrder.json")
	})
	t.Run("futures.trade.placeBatchOrders", func(t *testing.T) {
		checkPayload[[]model.BatchOrderAcknowledgement](t, "fixtures/rest/futures.trade.placeBatchOrders.json")
	})
	t.Run("futures.trade.reverseAPosition", func(t *testing.T) {
		checkPayload[model.ValueOrScalar[model.StringOrInteger]](t, "fixtures/rest/futures.trade.reverseAPosition.json")
	})
	t.Run("futures.trade.setStopLossTakeProfit", func(t *testing.T) {
		checkPayload[model.NoData](t, "fixtures/rest/futures.trade.setStopLossTakeProfit.json")
	})
	t.Run("futures.trade.setStopLossTakeProfitInBatches", func(t *testing.T) {
		checkPayload[model.NoData](t, "fixtures/rest/futures.trade.setStopLossTakeProfitInBatches.json")
	})
	t.Run("futures.trade.setTrailingStopLossTakeProfit", func(t *testing.T) {
		checkPayload[model.NoData](t, "fixtures/rest/futures.trade.setTrailingStopLossTakeProfit.json")
	})
}
func checkPayload[T any](t *testing.T, path string) {
	t.Helper()
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var env struct {
		Data json.RawMessage `json:"data"`
	}
	if err = json.Unmarshal(body, &env); err != nil {
		t.Fatal(err)
	}
	r, err := transport.DecodeResponse[T](&transport.Result{Data: env.Data, Envelope: body, EndpointID: t.Name()})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = json.Marshal(r.Data); err != nil {
		t.Fatalf("typed payload cannot serialize: %v", err)
	}
}
