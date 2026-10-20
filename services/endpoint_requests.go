package services

import (
	"context"
	"github.com/tigusigalpa/coinw-go/model"
	"github.com/tigusigalpa/coinw-go/transport"
)

// CancelWithdrawalWithRequest exposes every filter published on the operation's request page.
func (s *AccountService) CancelWithdrawalWithRequest(ctx context.Context, r model.CommonAccountCancelWithdrawalRequest) (*transport.Response[model.NoData], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[model.NoData](s.base, ctx, "common.account.cancelWithdrawal", p)
}

// GetDepositWithdrawalAddressWithRequest exposes every filter published on the operation's request page.
func (s *AccountService) GetDepositWithdrawalAddressWithRequest(ctx context.Context, r model.CommonAccountGetDepositWithdrawalAddressRequest) (*transport.Response[[]model.DepositAddress], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[[]model.DepositAddress](s.base, ctx, "common.account.getDepositWithdrawalAddress", p)
}

// GetDepositWithdrawalHistoryWithRequest exposes every filter published on the operation's request page.
func (s *AccountService) GetDepositWithdrawalHistoryWithRequest(ctx context.Context, r model.CommonAccountGetDepositWithdrawalHistoryRequest) (*transport.Response[[]model.DepositWithdrawal], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[[]model.DepositWithdrawal](s.base, ctx, "common.account.getDepositWithdrawalHistory", p)
}

// InitiateWithdrawalWithRequest exposes every filter published on the operation's request page.
func (s *AccountService) InitiateWithdrawalWithRequest(ctx context.Context, r model.CommonAccountInitiateWithdrawalRequest) (*transport.Response[model.WithdrawalAcknowledgement], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[model.WithdrawalAcknowledgement](s.base, ctx, "common.account.initiateWithdrawal", p)
}

// TransferAssetsWithRequest exposes every filter published on the operation's request page.
func (s *AccountService) TransferAssetsWithRequest(ctx context.Context, r model.CommonAccountTransferAssetsRequest) (*transport.Response[bool], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[bool](s.base, ctx, "common.account.transferAssets", p)
}

// GetSubAccountListWithRequest exposes every filter published on the operation's request page.
func (s *SubAccountService) GetSubAccountListWithRequest(ctx context.Context, r model.CommonSubAccountGetSubAccountListRequest) (*transport.Response[[]model.SubAccount], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[[]model.SubAccount](s.base, ctx, "common.subAccount.getSubAccountList", p)
}

// GetSubAccountApiKeysWithRequest exposes every filter published on the operation's request page.
func (s *SubAccountService) GetSubAccountApiKeysWithRequest(ctx context.Context, r model.CommonSubAccountGetSubAccountApiKeysRequest) (*transport.Response[[]model.SubAccountAPIKey], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[[]model.SubAccountAPIKey](s.base, ctx, "common.subAccount.getSubAccountApiKeys", p)
}

// ResetSubAccountApiKeyWithRequest exposes every filter published on the operation's request page.
func (s *SubAccountService) ResetSubAccountApiKeyWithRequest(ctx context.Context, r model.CommonSubAccountResetSubAccountApiKeyRequest) (*transport.Response[model.NoData], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[model.NoData](s.base, ctx, "common.subAccount.resetSubAccountApiKey", p)
}

// GetMaxTransferableAmountWithRequest exposes every filter published on the operation's request page.
func (s *SubAccountService) GetMaxTransferableAmountWithRequest(ctx context.Context, r model.CommonSubAccountGetMaxTransferableAmountRequest) (*transport.Response[model.WireDecimal], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[model.WireDecimal](s.base, ctx, "common.subAccount.getMaxTransferableAmount", p)
}

// GetTransferHistoryWithRequest exposes every filter published on the operation's request page.
func (s *SubAccountService) GetTransferHistoryWithRequest(ctx context.Context, r model.CommonSubAccountGetTransferHistoryRequest) (*transport.Response[[]model.Transfer], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[[]model.Transfer](s.base, ctx, "common.subAccount.getTransferHistory", p)
}

// AccountTransferWithRequest exposes every filter published on the operation's request page.
func (s *SubAccountService) AccountTransferWithRequest(ctx context.Context, r model.CommonSubAccountAccountTransferRequest) (*transport.Response[model.TransferAcknowledgement], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[model.TransferAcknowledgement](s.base, ctx, "common.subAccount.accountTransfer", p)
}

// GetInvitedUserListWithRequest exposes every filter published on the operation's request page.
func (s *AffiliateService) GetInvitedUserListWithRequest(ctx context.Context, r model.CommonAffiliateGetInvitedUserListRequest) (*transport.Response[model.ListPage[model.InvitedUser]], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[model.ListPage[model.InvitedUser]](s.base, ctx, "common.affiliate.getInvitedUserList", p)
}

// GetInvitedUsersKycResultsWithRequest exposes every filter published on the operation's request page.
func (s *AffiliateService) GetInvitedUsersKycResultsWithRequest(ctx context.Context, r model.CommonAffiliateGetInvitedUsersKycResultsRequest) (*transport.Response[model.ListPage[model.InvitedUserKYC]], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[model.ListPage[model.InvitedUserKYC]](s.base, ctx, "common.affiliate.getInvitedUsersKycResults", p)
}

// GetInvitedUsersTradingVolumeWithRequest exposes every filter published on the operation's request page.
func (s *AffiliateService) GetInvitedUsersTradingVolumeWithRequest(ctx context.Context, r model.CommonAffiliateGetInvitedUsersTradingVolumeRequest) (*transport.Response[model.ListPage[model.InvitedUserVolume]], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[model.ListPage[model.InvitedUserVolume]](s.base, ctx, "common.affiliate.getInvitedUsersTradingVolume", p)
}

// GetAffiliateUserActivityDataWithRequest exposes every filter published on the operation's request page.
func (s *AffiliateService) GetAffiliateUserActivityDataWithRequest(ctx context.Context, r model.CommonAffiliateGetAffiliateUserActivityDataRequest) (*transport.Response[model.ListPage[model.AffiliateActivity]], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[model.ListPage[model.AffiliateActivity]](s.base, ctx, "common.affiliate.getAffiliateUserActivityData", p)
}

// GetInvitedUserAssetsWithRequest exposes every filter published on the operation's request page.
func (s *AffiliateService) GetInvitedUserAssetsWithRequest(ctx context.Context, r model.CommonAffiliateGetInvitedUserAssetsRequest) (*transport.Response[model.ListPage[model.InvitedUserAssets]], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[model.ListPage[model.InvitedUserAssets]](s.base, ctx, "common.affiliate.getInvitedUserAssets", p)
}

// GetAffiliateTeamDirectUserRangeWithRequest exposes every filter published on the operation's request page.
func (s *AffiliateService) GetAffiliateTeamDirectUserRangeWithRequest(ctx context.Context, r model.CommonAffiliateGetAffiliateTeamDirectUserRangeRequest) (*transport.Response[model.AffiliateTeam], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[model.AffiliateTeam](s.base, ctx, "common.affiliate.getAffiliateTeamDirectUserRange", p)
}

// GetInstrumentInformationWithRequest exposes every filter published on the operation's request page.
func (s *SpotMarketService) GetInstrumentInformationWithRequest(ctx context.Context, r model.SpotMarketGetInstrumentInformationRequest) (*transport.Response[[]model.SpotInstrument], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[[]model.SpotInstrument](s.base, ctx, "spot.market.getInstrumentInformation", p)
}

// Get24hTradeSummaryForAllInstrumentsWithRequest exposes every filter published on the operation's request page.
func (s *SpotMarketService) Get24hTradeSummaryForAllInstrumentsWithRequest(ctx context.Context, r model.SpotMarketGet24hTradeSummaryForAllInstrumentsRequest) (*transport.Response[map[string]model.SpotTicker], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[map[string]model.SpotTicker](s.base, ctx, "spot.market.get24hTradeSummaryForAllInstruments", p)
}

// GetCoinDepositWithdrawalLimitsWithRequest exposes every filter published on the operation's request page.
func (s *SpotMarketService) GetCoinDepositWithdrawalLimitsWithRequest(ctx context.Context, r model.SpotMarketGetCoinDepositWithdrawalLimitsRequest) (*transport.Response[map[string]model.SpotCurrency], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[map[string]model.SpotCurrency](s.base, ctx, "spot.market.getCoinDepositWithdrawalLimits", p)
}

// GetOrderBookWithRequest exposes every filter published on the operation's request page.
func (s *SpotMarketService) GetOrderBookWithRequest(ctx context.Context, r model.SpotMarketGetOrderBookRequest) (*transport.Response[[]model.SpotOrderBook], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[[]model.SpotOrderBook](s.base, ctx, "spot.market.getOrderBook", p)
}

// GetRecentTradesWithRequest exposes every filter published on the operation's request page.
func (s *SpotMarketService) GetRecentTradesWithRequest(ctx context.Context, r model.SpotMarketGetRecentTradesRequest) (*transport.Response[[]model.SpotTrade], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[[]model.SpotTrade](s.base, ctx, "spot.market.getRecentTrades", p)
}

// GetKLineWithRequest exposes every filter published on the operation's request page.
func (s *SpotMarketService) GetKLineWithRequest(ctx context.Context, r model.SpotMarketGetKLineRequest) (*transport.Response[[]model.SpotCandle], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[[]model.SpotCandle](s.base, ctx, "spot.market.getKLine", p)
}

// Get24hVolumeForPopularInstrumentsWithRequest exposes every filter published on the operation's request page.
func (s *SpotMarketService) Get24hVolumeForPopularInstrumentsWithRequest(ctx context.Context, r model.SpotMarketGet24hVolumeForPopularInstrumentsRequest) (*transport.Response[model.PopularVolumes], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[model.PopularVolumes](s.base, ctx, "spot.market.get24hVolumeForPopularInstruments", p)
}

// GetSpotAccountBalanceWithRequest exposes every filter published on the operation's request page.
func (s *SpotAccountService) GetSpotAccountBalanceWithRequest(ctx context.Context, r model.SpotAccountGetSpotAccountBalanceRequest) (*transport.Response[map[string]model.WireDecimal], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[map[string]model.WireDecimal](s.base, ctx, "spot.account.getSpotAccountBalance", p)
}

// GetFullSpotAccountBalanceWithRequest exposes every filter published on the operation's request page.
func (s *SpotAccountService) GetFullSpotAccountBalanceWithRequest(ctx context.Context, r model.SpotAccountGetFullSpotAccountBalanceRequest) (*transport.Response[map[string]model.SpotFullBalance], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[map[string]model.SpotFullBalance](s.base, ctx, "spot.account.getFullSpotAccountBalance", p)
}

// GetBatchHistoricalOrdersWithRequest exposes every filter published on the operation's request page.
func (s *SpotCheckService) GetBatchHistoricalOrdersWithRequest(ctx context.Context, r model.SpotCheckGetBatchHistoricalOrdersRequest) (*transport.Response[model.ObjectOrArray[model.SpotBatchHistoricalOrder]], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[model.ObjectOrArray[model.SpotBatchHistoricalOrder]](s.base, ctx, "spot.check.getBatchHistoricalOrders", p)
}

// GetCurrentOrdersWithRequest exposes every filter published on the operation's request page.
func (s *SpotCheckService) GetCurrentOrdersWithRequest(ctx context.Context, r model.SpotCheckGetCurrentOrdersRequest) (*transport.Response[[]model.SpotCurrentOrder], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[[]model.SpotCurrentOrder](s.base, ctx, "spot.check.getCurrentOrders", p)
}

// GetHistoricalOrdersWithRequest exposes every filter published on the operation's request page.
func (s *SpotCheckService) GetHistoricalOrdersWithRequest(ctx context.Context, r model.SpotCheckGetHistoricalOrdersRequest) (*transport.Response[model.SpotHistoryPage], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[model.SpotHistoryPage](s.base, ctx, "spot.check.getHistoricalOrders", p)
}

// GetOrderDetailsWithRequest exposes every filter published on the operation's request page.
func (s *SpotCheckService) GetOrderDetailsWithRequest(ctx context.Context, r model.SpotCheckGetOrderDetailsRequest) (*transport.Response[model.SpotOrderDetails], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[model.SpotOrderDetails](s.base, ctx, "spot.check.getOrderDetails", p)
}

// GetOrderStatusWithRequest exposes every filter published on the operation's request page.
func (s *SpotCheckService) GetOrderStatusWithRequest(ctx context.Context, r model.SpotCheckGetOrderStatusRequest) (*transport.Response[model.SpotOrderStatus], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[model.SpotOrderStatus](s.base, ctx, "spot.check.getOrderStatus", p)
}

// GetTransactionHistoryWithRequest exposes every filter published on the operation's request page.
func (s *SpotCheckService) GetTransactionHistoryWithRequest(ctx context.Context, r model.SpotCheckGetTransactionHistoryRequest) (*transport.Response[[]model.SpotTransaction], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[[]model.SpotTransaction](s.base, ctx, "spot.check.getTransactionHistory", p)
}

// PlaceOrderWithRequest exposes every filter published on the operation's request page.
func (s *SpotTradeService) PlaceOrderWithRequest(ctx context.Context, r model.SpotTradePlaceOrderRequest) (*transport.Response[model.SpotOrderAcknowledgement], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[model.SpotOrderAcknowledgement](s.base, ctx, "spot.trade.placeOrder", p)
}

// CancelOrderWithRequest exposes every filter published on the operation's request page.
func (s *SpotTradeService) CancelOrderWithRequest(ctx context.Context, r model.SpotTradeCancelOrderRequest) (*transport.Response[model.SpotCancelAcknowledgement], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[model.SpotCancelAcknowledgement](s.base, ctx, "spot.trade.cancelOrder", p)
}

// CancelAllOrdersWithRequest exposes every filter published on the operation's request page.
func (s *SpotTradeService) CancelAllOrdersWithRequest(ctx context.Context, r model.SpotTradeCancelAllOrdersRequest) (*transport.Response[model.NoData], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[model.NoData](s.base, ctx, "spot.trade.cancelAllOrders", p)
}

// GetInstrumentInfoWithRequest exposes every filter published on the operation's request page.
func (s *FuturesMarketService) GetInstrumentInfoWithRequest(ctx context.Context, r model.FuturesMarketGetInstrumentInfoRequest) (*transport.Response[[]model.FuturesInstrument], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[[]model.FuturesInstrument](s.base, ctx, "futures.market.getInstrumentInfo", p)
}

// GetBatchInstrumentsInfoWithRequest exposes every filter published on the operation's request page.
func (s *FuturesMarketService) GetBatchInstrumentsInfoWithRequest(ctx context.Context, r model.FuturesMarketGetBatchInstrumentsInfoRequest) (*transport.Response[[]model.FuturesInstrument], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[[]model.FuturesInstrument](s.base, ctx, "futures.market.getBatchInstrumentsInfo", p)
}

// GetLastTradeSummaryOfAnInstrumentWithRequest exposes every filter published on the operation's request page.
func (s *FuturesMarketService) GetLastTradeSummaryOfAnInstrumentWithRequest(ctx context.Context, r model.FuturesMarketGetLastTradeSummaryOfAnInstrumentRequest) (*transport.Response[[]model.FuturesTicker], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[[]model.FuturesTicker](s.base, ctx, "futures.market.getLastTradeSummaryOfAnInstrument", p)
}

// GetBatchLastTradeSummaryWithRequest exposes every filter published on the operation's request page.
func (s *FuturesMarketService) GetBatchLastTradeSummaryWithRequest(ctx context.Context, r model.FuturesMarketGetBatchLastTradeSummaryRequest) (*transport.Response[[]model.FuturesTicker], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[[]model.FuturesTicker](s.base, ctx, "futures.market.getBatchLastTradeSummary", p)
}

// GetLastTradeSummaryOfAllInstrumentsWithRequest exposes every filter published on the operation's request page.
func (s *FuturesMarketService) GetLastTradeSummaryOfAllInstrumentsWithRequest(ctx context.Context, r model.FuturesMarketGetLastTradeSummaryOfAllInstrumentsRequest) (*transport.Response[[]model.FuturesTicker], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[[]model.FuturesTicker](s.base, ctx, "futures.market.getLastTradeSummaryOfAllInstruments", p)
}

// GetKLineOfAnInstrumentWithRequest exposes every filter published on the operation's request page.
func (s *FuturesMarketService) GetKLineOfAnInstrumentWithRequest(ctx context.Context, r model.FuturesMarketGetKLineOfAnInstrumentRequest) (*transport.Response[[]model.FuturesCandleTuple], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[[]model.FuturesCandleTuple](s.base, ctx, "futures.market.getKLineOfAnInstrument", p)
}

// GetLastSettlementFundingFeeRateWithRequest exposes every filter published on the operation's request page.
func (s *FuturesMarketService) GetLastSettlementFundingFeeRateWithRequest(ctx context.Context, r model.FuturesMarketGetLastSettlementFundingFeeRateRequest) (*transport.Response[model.DecimalAtTime], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[model.DecimalAtTime](s.base, ctx, "futures.market.getLastSettlementFundingFeeRate", p)
}

// GetOrderBookOfAnInstrumentWithRequest exposes every filter published on the operation's request page.
func (s *FuturesMarketService) GetOrderBookOfAnInstrumentWithRequest(ctx context.Context, r model.FuturesMarketGetOrderBookOfAnInstrumentRequest) (*transport.Response[model.FuturesOrderBook], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[model.FuturesOrderBook](s.base, ctx, "futures.market.getOrderBookOfAnInstrument", p)
}

// GetTradeDataOfAnInstrumentWithRequest exposes every filter published on the operation's request page.
func (s *FuturesMarketService) GetTradeDataOfAnInstrumentWithRequest(ctx context.Context, r model.FuturesMarketGetTradeDataOfAnInstrumentRequest) (*transport.Response[[]model.FuturesTrade], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[[]model.FuturesTrade](s.base, ctx, "futures.market.getTradeDataOfAnInstrument", p)
}

// GetMarginRequirementsOfAllInstrumentsWithRequest exposes every filter published on the operation's request page.
func (s *FuturesMarketService) GetMarginRequirementsOfAllInstrumentsWithRequest(ctx context.Context, r model.FuturesMarketGetMarginRequirementsOfAllInstrumentsRequest) (*transport.Response[model.FuturesLadders], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[model.FuturesLadders](s.base, ctx, "futures.market.getMarginRequirementsOfAllInstruments", p)
}

// GetHistoricalPublicTradesWithRequest exposes every filter published on the operation's request page.
func (s *FuturesMarketService) GetHistoricalPublicTradesWithRequest(ctx context.Context, r model.FuturesMarketGetHistoricalPublicTradesRequest) (*transport.Response[model.RowsPage[model.HistoricalPublicTrade]], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[model.RowsPage[model.HistoricalPublicTrade]](s.base, ctx, "futures.market.getHistoricalPublicTrades", p)
}

// CoinConversionWithRequest exposes every filter published on the operation's request page.
func (s *FuturesAssetsService) CoinConversionWithRequest(ctx context.Context, r model.FuturesAssetsCoinConversionRequest) (*transport.Response[model.DecimalValue], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[model.DecimalValue](s.base, ctx, "futures.assets.coinConversion", p)
}

// EnableDisableMegaCouponWithRequest exposes every filter published on the operation's request page.
func (s *FuturesAssetsService) EnableDisableMegaCouponWithRequest(ctx context.Context, r model.FuturesAssetsEnableDisableMegaCouponRequest) (*transport.Response[model.NoData], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[model.NoData](s.base, ctx, "futures.assets.enableDisableMegaCoupon", p)
}

// GetFuturesAccountAssetsWithRequest exposes every filter published on the operation's request page.
func (s *FuturesAssetsService) GetFuturesAccountAssetsWithRequest(ctx context.Context, r model.FuturesAssetsGetFuturesAccountAssetsRequest) (*transport.Response[model.FuturesAccountAssets], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[model.FuturesAccountAssets](s.base, ctx, "futures.assets.getFuturesAccountAssets", p)
}

// GetFuturesAccountFeesWithRequest exposes every filter published on the operation's request page.
func (s *FuturesAssetsService) GetFuturesAccountFeesWithRequest(ctx context.Context, r model.FuturesAssetsGetFuturesAccountFeesRequest) (*transport.Response[model.FuturesAccountFees], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[model.FuturesAccountFees](s.base, ctx, "futures.assets.getFuturesAccountFees", p)
}

// GetMarginModeWithRequest exposes every filter published on the operation's request page.
func (s *FuturesAssetsService) GetMarginModeWithRequest(ctx context.Context, r model.FuturesAssetsGetMarginModeRequest) (*transport.Response[model.MarginMode], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[model.MarginMode](s.base, ctx, "futures.assets.getMarginMode", p)
}

// GetMaxTransferableBalanceWithRequest exposes every filter published on the operation's request page.
func (s *FuturesAssetsService) GetMaxTransferableBalanceWithRequest(ctx context.Context, r model.FuturesAssetsGetMaxTransferableBalanceRequest) (*transport.Response[model.ValueOrScalar[model.WireDecimal]], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[model.ValueOrScalar[model.WireDecimal]](s.base, ctx, "futures.assets.getMaxTransferableBalance", p)
}

// GetMaxUsersContractSizeWithRequest exposes every filter published on the operation's request page.
func (s *FuturesAssetsService) GetMaxUsersContractSizeWithRequest(ctx context.Context, r model.FuturesAssetsGetMaxUsersContractSizeRequest) (*transport.Response[model.AvailableContractSize], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[model.AvailableContractSize](s.base, ctx, "futures.assets.getMaxUsersContractSize", p)
}

// GetMegaCouponBalanceWithRequest exposes every filter published on the operation's request page.
func (s *FuturesAssetsService) GetMegaCouponBalanceWithRequest(ctx context.Context, r model.FuturesAssetsGetMegaCouponBalanceRequest) (*transport.Response[[]model.MegaCoupon], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[[]model.MegaCoupon](s.base, ctx, "futures.assets.getMegaCouponBalance", p)
}

// GetTransactionDetails3DaysWithRequest exposes every filter published on the operation's request page.
func (s *FuturesAssetsService) GetTransactionDetails3DaysWithRequest(ctx context.Context, r model.FuturesAssetsGetTransactionDetails3DaysRequest) (*transport.Response[model.RowsPage[model.FuturesTransaction]], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[model.RowsPage[model.FuturesTransaction]](s.base, ctx, "futures.assets.getTransactionDetails3Days", p)
}

// GetTransactionDetails3MonthsWithRequest exposes every filter published on the operation's request page.
func (s *FuturesAssetsService) GetTransactionDetails3MonthsWithRequest(ctx context.Context, r model.FuturesAssetsGetTransactionDetails3MonthsRequest) (*transport.Response[model.RowsPage[model.FuturesTransaction]], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[model.RowsPage[model.FuturesTransaction]](s.base, ctx, "futures.assets.getTransactionDetails3Months", p)
}

// SetMarginModeWithRequest exposes every filter published on the operation's request page.
func (s *FuturesAssetsService) SetMarginModeWithRequest(ctx context.Context, r model.FuturesAssetsSetMarginModeRequest) (*transport.Response[string], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[string](s.base, ctx, "futures.assets.setMarginMode", p)
}

// GetLeverageInformationWithRequest exposes every filter published on the operation's request page.
func (s *FuturesCommonService) GetLeverageInformationWithRequest(ctx context.Context, r model.FuturesCommonGetLeverageInformationRequest) (*transport.Response[model.ValueOrScalar[model.WireDecimal]], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[model.ValueOrScalar[model.WireDecimal]](s.base, ctx, "futures.common.getLeverageInformation", p)
}

// GetCurrentPositionInformationWithRequest exposes every filter published on the operation's request page.
func (s *FuturesPositionService) GetCurrentPositionInformationWithRequest(ctx context.Context, r model.FuturesPositionGetCurrentPositionInformationRequest) (*transport.Response[[]model.FuturesPosition], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[[]model.FuturesPosition](s.base, ctx, "futures.position.getCurrentPositionInformation", p)
}

// GetAllCurrentPositionsWithRequest exposes every filter published on the operation's request page.
func (s *FuturesPositionService) GetAllCurrentPositionsWithRequest(ctx context.Context, r model.FuturesPositionGetAllCurrentPositionsRequest) (*transport.Response[[]model.FuturesPosition], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[[]model.FuturesPosition](s.base, ctx, "futures.position.getAllCurrentPositions", p)
}

// GetHistoricalPositionInformationWithRequest exposes every filter published on the operation's request page.
func (s *FuturesPositionService) GetHistoricalPositionInformationWithRequest(ctx context.Context, r model.FuturesPositionGetHistoricalPositionInformationRequest) (*transport.Response[model.RowsPage[model.HistoricalPosition]], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[model.RowsPage[model.HistoricalPosition]](s.base, ctx, "futures.position.getHistoricalPositionInformation", p)
}

// GetMaxContractSizeWithRequest exposes every filter published on the operation's request page.
func (s *FuturesPositionService) GetMaxContractSizeWithRequest(ctx context.Context, r model.FuturesPositionGetMaxContractSizeRequest) (*transport.Response[model.MaximumContractSize], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[model.MaximumContractSize](s.base, ctx, "futures.position.getMaxContractSize", p)
}

// GetPositionMarginRateWithRequest exposes every filter published on the operation's request page.
func (s *FuturesPositionService) GetPositionMarginRateWithRequest(ctx context.Context, r model.FuturesPositionGetPositionMarginRateRequest) (*transport.Response[model.ValueOrScalar[model.WireDecimal]], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[model.ValueOrScalar[model.WireDecimal]](s.base, ctx, "futures.position.getPositionMarginRate", p)
}

// GetBatchHistoricalOrders7DaysWithRequest exposes every filter published on the operation's request page.
func (s *FuturesCheckService) GetBatchHistoricalOrders7DaysWithRequest(ctx context.Context, r model.FuturesCheckGetBatchHistoricalOrders7DaysRequest) (*transport.Response[model.RowsPage[model.FuturesOrder]], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[model.RowsPage[model.FuturesOrder]](s.base, ctx, "futures.check.getBatchHistoricalOrders7Days", p)
}

// GetBatchHistoricalOrders3MonthsWithRequest exposes every filter published on the operation's request page.
func (s *FuturesCheckService) GetBatchHistoricalOrders3MonthsWithRequest(ctx context.Context, r model.FuturesCheckGetBatchHistoricalOrders3MonthsRequest) (*transport.Response[model.RowsPage[model.FuturesOrder]], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[model.RowsPage[model.FuturesOrder]](s.base, ctx, "futures.check.getBatchHistoricalOrders3Months", p)
}

// GetCurrentOrdersWithRequest exposes every filter published on the operation's request page.
func (s *FuturesCheckService) GetCurrentOrdersWithRequest(ctx context.Context, r model.FuturesCheckGetCurrentOrdersRequest) (*transport.Response[model.RowsPage[model.FuturesOrder]], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[model.RowsPage[model.FuturesOrder]](s.base, ctx, "futures.check.getCurrentOrders", p)
}

// GetHistoricalOrders7DaysWithRequest exposes every filter published on the operation's request page.
func (s *FuturesCheckService) GetHistoricalOrders7DaysWithRequest(ctx context.Context, r model.FuturesCheckGetHistoricalOrders7DaysRequest) (*transport.Response[model.RowsPage[model.FuturesOrder]], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[model.RowsPage[model.FuturesOrder]](s.base, ctx, "futures.check.getHistoricalOrders7Days", p)
}

// GetHistoricalOrders3MonthsWithRequest exposes every filter published on the operation's request page.
func (s *FuturesCheckService) GetHistoricalOrders3MonthsWithRequest(ctx context.Context, r model.FuturesCheckGetHistoricalOrders3MonthsRequest) (*transport.Response[model.RowsPage[model.FuturesOrder]], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[model.RowsPage[model.FuturesOrder]](s.base, ctx, "futures.check.getHistoricalOrders3Months", p)
}

// GetOrderInformationWithRequest exposes every filter published on the operation's request page.
func (s *FuturesCheckService) GetOrderInformationWithRequest(ctx context.Context, r model.FuturesCheckGetOrderInformationRequest) (*transport.Response[[]model.FuturesOrder], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[[]model.FuturesOrder](s.base, ctx, "futures.check.getOrderInformation", p)
}

// GetPendingOrderCountWithRequest exposes every filter published on the operation's request page.
func (s *FuturesCheckService) GetPendingOrderCountWithRequest(ctx context.Context, r model.FuturesCheckGetPendingOrderCountRequest) (*transport.Response[model.ValueOrScalar[model.StringOrInteger]], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[model.ValueOrScalar[model.StringOrInteger]](s.base, ctx, "futures.check.getPendingOrderCount", p)
}

// GetStopLossTakeProfitInformationWithRequest exposes every filter published on the operation's request page.
func (s *FuturesCheckService) GetStopLossTakeProfitInformationWithRequest(ctx context.Context, r model.FuturesCheckGetStopLossTakeProfitInformationRequest) (*transport.Response[[]model.StopLossTakeProfit], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[[]model.StopLossTakeProfit](s.base, ctx, "futures.check.getStopLossTakeProfitInformation", p)
}

// GetTrailingStopLossTakeProfitInformationWithRequest exposes every filter published on the operation's request page.
func (s *FuturesCheckService) GetTrailingStopLossTakeProfitInformationWithRequest(ctx context.Context, r model.FuturesCheckGetTrailingStopLossTakeProfitInformationRequest) (*transport.Response[[]model.TrailingStopLossTakeProfit], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return callTyped[[]model.TrailingStopLossTakeProfit](s.base, ctx, "futures.check.getTrailingStopLossTakeProfitInformation", p)
}
