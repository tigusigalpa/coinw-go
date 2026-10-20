package model

// Request DTOs mirror the request tables at the source links in contracts/catalogue.json.
// Pointer fields distinguish omitted filters from explicit zero values.

// CommonAccountCancelWithdrawalRequest — https://www.coinw.com/api-doc/en/common/account/cancel-withdrawal
type CommonAccountCancelWithdrawalRequest struct {
	ID string `json:"id"`
}

// CommonAccountGetDepositWithdrawalAddressRequest — https://www.coinw.com/api-doc/en/common/account/get-deposit-withdrawal-address
type CommonAccountGetDepositWithdrawalAddressRequest struct {
	SymbolID string `json:"symbolId"`
	Chain    string `json:"chain"`
}

// CommonAccountGetDepositWithdrawalHistoryRequest — https://www.coinw.com/api-doc/en/common/account/get-deposit-withdrawal-history
type CommonAccountGetDepositWithdrawalHistoryRequest struct {
	Symbol        string  `json:"symbol"`
	DepositNumber *string `json:"depositNumber,omitempty"`
}

// CommonAccountInitiateWithdrawalRequest — https://www.coinw.com/api-doc/en/common/account/initiate-withdrawal
type CommonAccountInitiateWithdrawalRequest struct {
	Memo        string  `json:"memo"`
	Type        *string `json:"type,omitempty"`
	Amount      Decimal `json:"amount"`
	Currency    string  `json:"currency"`
	Address     string  `json:"address"`
	Chain       string  `json:"chain"`
	InnerToType int64   `json:"innerToType"`
}

// CommonAccountTransferAssetsRequest — https://www.coinw.com/api-doc/en/common/account/transfer-assets
type CommonAccountTransferAssetsRequest struct {
	AccountType       string  `json:"accountType"`
	TargetAccountType string  `json:"targetAccountType"`
	BizType           string  `json:"bizType"`
	CoinCode          string  `json:"coinCode"`
	Amount            Decimal `json:"amount"`
}

// CommonSubAccountGetSubAccountListRequest — https://www.coinw.com/api-doc/en/common/sub-account/get-sub-account-list
type CommonSubAccountGetSubAccountListRequest struct {
}

// CommonSubAccountGetSubAccountApiKeysRequest — https://www.coinw.com/api-doc/en/common/sub-account/get-sub-account-api-keys
type CommonSubAccountGetSubAccountApiKeysRequest struct {
}

// CommonSubAccountResetSubAccountApiKeyRequest — https://www.coinw.com/api-doc/en/common/sub-account/reset-sub-account-api-key
type CommonSubAccountResetSubAccountApiKeyRequest struct {
	SubAccountApiKey string `json:"subAccount_apiKey"`
}

// CommonSubAccountGetMaxTransferableAmountRequest — https://www.coinw.com/api-doc/en/common/sub-account/get-max-transferable-amount
type CommonSubAccountGetMaxTransferableAmountRequest struct {
	CoinCode string `json:"coinCode"`
	SubUID   int64  `json:"subUid"`
}

// CommonSubAccountGetTransferHistoryRequest — https://www.coinw.com/api-doc/en/common/sub-account/get-transfer-history
type CommonSubAccountGetTransferHistoryRequest struct {
	Currency string `json:"currency"`
	From     int64  `json:"from"`
	To       int64  `json:"to"`
	PageNo   int64  `json:"pageNo"`
	PageSize int64  `json:"pageSize"`
}

// CommonSubAccountAccountTransferRequest — https://www.coinw.com/api-doc/en/common/sub-account/account-transfer
type CommonSubAccountAccountTransferRequest struct {
	Side     string  `json:"side"`
	SubUID   int64   `json:"subUid"`
	Currency string  `json:"currency"`
	Amount   Decimal `json:"amount"`
}

// CommonAffiliateGetInvitedUserListRequest — https://www.coinw.com/api-doc/en/common/affiliate/get-invited-user-list
type CommonAffiliateGetInvitedUserListRequest struct {
	StartTime    *string `json:"startTime,omitempty"`
	EndTime      *string `json:"endTime,omitempty"`
	PageNo       *int64  `json:"pageNo,omitempty"`
	PageSize     *int64  `json:"pageSize,omitempty"`
	UID          *int64  `json:"uid,omitempty"`
	ReferralCode *string `json:"referralCode,omitempty"`
}

// CommonAffiliateGetInvitedUsersKycResultsRequest — https://www.coinw.com/api-doc/en/common/affiliate/get-invited-users-kyc-results
type CommonAffiliateGetInvitedUsersKycResultsRequest struct {
	StartTime *string `json:"startTime,omitempty"`
	EndTime   *string `json:"endTime,omitempty"`
	PageNo    *int64  `json:"pageNo,omitempty"`
	PageSize  *int64  `json:"pageSize,omitempty"`
	UID       *int64  `json:"uid,omitempty"`
}

// CommonAffiliateGetInvitedUsersTradingVolumeRequest — https://www.coinw.com/api-doc/en/common/affiliate/get-invited-users-trading-volume
type CommonAffiliateGetInvitedUsersTradingVolumeRequest struct {
	StartTime *string `json:"startTime,omitempty"`
	EndTime   *string `json:"endTime,omitempty"`
	PageNo    *int64  `json:"pageNo,omitempty"`
	PageSize  *int64  `json:"pageSize,omitempty"`
	UID       *int64  `json:"uid,omitempty"`
}

// CommonAffiliateGetAffiliateUserActivityDataRequest — https://www.coinw.com/api-doc/en/common/affiliate/get-affiliate-user-activity-data
type CommonAffiliateGetAffiliateUserActivityDataRequest struct {
	ActivityID string  `json:"activityId"`
	StartTime  *string `json:"startTime,omitempty"`
	EndTime    *string `json:"endTime,omitempty"`
	Page       *int64  `json:"page,omitempty"`
	PageSize   *int64  `json:"pageSize,omitempty"`
	UID        *int64  `json:"uid,omitempty"`
	ShowDetail *string `json:"showDetail,omitempty"`
}

// CommonAffiliateGetInvitedUserAssetsRequest — https://www.coinw.com/api-doc/en/common/affiliate/get-invited-user-assets
type CommonAffiliateGetInvitedUserAssetsRequest struct {
	UID       *int64  `json:"uid,omitempty"`
	StartTime *string `json:"startTime,omitempty"`
	EndTime   *string `json:"endTime,omitempty"`
	Quote     *string `json:"quote,omitempty"`
	PageNo    *int64  `json:"pageNo,omitempty"`
	PageSize  *int64  `json:"pageSize,omitempty"`
}

// CommonAffiliateGetAffiliateTeamDirectUserRangeRequest — https://www.coinw.com/api-doc/en/common/affiliate/get-affiliate-team-direct-user-range
type CommonAffiliateGetAffiliateTeamDirectUserRangeRequest struct {
	StartTime    *string `json:"startTime,omitempty"`
	EndTime      *string `json:"endTime,omitempty"`
	PageNo       *int64  `json:"pageNo,omitempty"`
	PageSize     *int64  `json:"pageSize,omitempty"`
	UID          *int64  `json:"uid,omitempty"`
	ReferralCode *string `json:"referralCode,omitempty"`
	Type         *string `json:"type,omitempty"`
	PartnerUids  []int64 `json:"partnerUids,omitempty"`
}

// SpotMarketGetInstrumentInformationRequest — https://www.coinw.com/api-doc/en/spot-trading/market/get-trading-pair-information
type SpotMarketGetInstrumentInformationRequest struct {
}

// SpotMarketGet24hTradeSummaryForAllInstrumentsRequest — https://www.coinw.com/api-doc/en/spot-trading/market/get-24h-trade-summary-for-all-instruments
type SpotMarketGet24hTradeSummaryForAllInstrumentsRequest struct {
}

// SpotMarketGetCoinDepositWithdrawalLimitsRequest — https://www.coinw.com/api-doc/en/spot-trading/market/get-coin-deposits-withdrawal-limits
type SpotMarketGetCoinDepositWithdrawalLimitsRequest struct {
}

// SpotMarketGetOrderBookRequest — https://www.coinw.com/api-doc/en/spot-trading/market/get-order-book
type SpotMarketGetOrderBookRequest struct {
	Size   *int64 `json:"size,omitempty"`
	Symbol string `json:"symbol"`
}

// SpotMarketGetRecentTradesRequest — https://www.coinw.com/api-doc/en/spot-trading/market/get-recent-trades
type SpotMarketGetRecentTradesRequest struct {
	Symbol string  `json:"symbol"`
	Start  *string `json:"start,omitempty"`
	End    *string `json:"end,omitempty"`
}

// SpotMarketGetKLineRequest — https://www.coinw.com/api-doc/en/spot-trading/market/get-k-line
type SpotMarketGetKLineRequest struct {
	Period       int    `json:"period"`
	CurrencyPair string `json:"currencyPair"`
	Start        *int64 `json:"start,omitempty"`
	End          *int64 `json:"end,omitempty"`
}

// SpotMarketGet24hVolumeForPopularInstrumentsRequest — https://www.coinw.com/api-doc/en/spot-trading/market/get-24h-volume-for-popular-instruments
type SpotMarketGet24hVolumeForPopularInstrumentsRequest struct {
}

// SpotAccountGetSpotAccountBalanceRequest — https://www.coinw.com/api-doc/en/spot-trading/account/get-spot-account-balance
type SpotAccountGetSpotAccountBalanceRequest struct {
}

// SpotAccountGetFullSpotAccountBalanceRequest — https://www.coinw.com/api-doc/en/spot-trading/account/get-total-spot-account-balance
type SpotAccountGetFullSpotAccountBalanceRequest struct {
}

// SpotCheckGetBatchHistoricalOrdersRequest — https://www.coinw.com/api-doc/en/spot-trading/check/get-batch-historical-orders
type SpotCheckGetBatchHistoricalOrdersRequest struct {
	OrderIds []string `json:"orderIds"`
}

// SpotCheckGetCurrentOrdersRequest — https://www.coinw.com/api-doc/en/spot-trading/check/get-current-orders
type SpotCheckGetCurrentOrdersRequest struct {
	CurrencyPair string `json:"currencyPair"`
	StartAt      *int64 `json:"startAt,omitempty"`
	EndAt        *int64 `json:"endAt,omitempty"`
}

// SpotCheckGetHistoricalOrdersRequest — https://www.coinw.com/api-doc/en/spot-trading/check/get-historical-orders
type SpotCheckGetHistoricalOrdersRequest struct {
	Symbol  *string `json:"symbol,omitempty"`
	StartAt *int64  `json:"startAt,omitempty"`
	EndAt   *int64  `json:"endAt,omitempty"`
	Limit   *int64  `json:"limit,omitempty"`
	Before  *string `json:"before,omitempty"`
	After   *string `json:"after,omitempty"`
}

// SpotCheckGetOrderDetailsRequest — https://www.coinw.com/api-doc/en/spot-trading/check/get-order-details
type SpotCheckGetOrderDetailsRequest struct {
	OrderNumber string `json:"orderNumber"`
}

// SpotCheckGetOrderStatusRequest — https://www.coinw.com/api-doc/en/spot-trading/check/get-order-status
type SpotCheckGetOrderStatusRequest struct {
	OrderNumber string `json:"orderNumber"`
}

// SpotCheckGetTransactionHistoryRequest — https://www.coinw.com/api-doc/en/spot-trading/check/get-transaction-history
type SpotCheckGetTransactionHistoryRequest struct {
	CurrencyPair string  `json:"currencyPair"`
	StartAt      *string `json:"startAt,omitempty"`
	EndAt        *string `json:"endAt,omitempty"`
}

// SpotTradePlaceOrderRequest — https://www.coinw.com/api-doc/en/spot-trading/trade/place-order
type SpotTradePlaceOrderRequest struct {
	Symbol     string   `json:"symbol"`
	Type       int      `json:"type"`
	Amount     *Decimal `json:"amount,omitempty"`
	Rate       *Decimal `json:"rate,omitempty"`
	Funds      *Decimal `json:"funds,omitempty"`
	IsMarket   bool     `json:"isMarket"`
	OutTradeNo string   `json:"out_trade_no"`
}

// SpotTradeCancelOrderRequest — https://www.coinw.com/api-doc/en/spot-trading/trade/cancel-order
type SpotTradeCancelOrderRequest struct {
	OrderNumber string `json:"orderNumber"`
}

// SpotTradeCancelAllOrdersRequest — https://www.coinw.com/api-doc/en/spot-trading/trade/cancel-all-orders
type SpotTradeCancelAllOrdersRequest struct {
	CurrencyPair *string `json:"currencyPair,omitempty"`
}

// FuturesMarketGetInstrumentInfoRequest — https://www.coinw.com/api-doc/en/futures-trading/market/get-instrument-information
type FuturesMarketGetInstrumentInfoRequest struct {
	Name *string `json:"name,omitempty"`
}

// FuturesMarketGetBatchInstrumentsInfoRequest — https://www.coinw.com/api-doc/en/futures-trading/market/get-batch-instrument-information
type FuturesMarketGetBatchInstrumentsInfoRequest struct {
	Symbols InstrumentList `json:"symbols"`
}

// FuturesMarketGetLastTradeSummaryOfAnInstrumentRequest — https://www.coinw.com/api-doc/en/futures-trading/market/get-last-trade-summary-of-an-instrument
type FuturesMarketGetLastTradeSummaryOfAnInstrumentRequest struct {
	Instrument Instrument `json:"instrument"`
}

// FuturesMarketGetBatchLastTradeSummaryRequest — https://www.coinw.com/api-doc/en/futures-trading/market/get-batch-last-trade-summary
type FuturesMarketGetBatchLastTradeSummaryRequest struct {
	Symbols InstrumentList `json:"symbols"`
}

// FuturesMarketGetLastTradeSummaryOfAllInstrumentsRequest — https://www.coinw.com/api-doc/en/futures-trading/market/get-last-trade-summary-of-all-instruments
type FuturesMarketGetLastTradeSummaryOfAllInstrumentsRequest struct {
}

// FuturesMarketGetKLineOfAnInstrumentRequest — https://www.coinw.com/api-doc/en/futures-trading/market/get-k-line-of-an-instrument
type FuturesMarketGetKLineOfAnInstrumentRequest struct {
	CurrencyCode Instrument `json:"currencyCode"`
	Granuality   string     `json:"granuality"`
	Limit        *int64     `json:"limit,omitempty"`
	KlineType    *string    `json:"klineType,omitempty"`
	SinceStr     *string    `json:"sinceStr,omitempty"`
	SinceEndStr  *string    `json:"sinceEndStr,omitempty"`
}

// FuturesMarketGetLastSettlementFundingFeeRateRequest — https://www.coinw.com/api-doc/en/futures-trading/market/get-last-settelment-funding-fee-rate
type FuturesMarketGetLastSettlementFundingFeeRateRequest struct {
	Instrument Instrument `json:"instrument"`
}

// FuturesMarketGetOrderBookOfAnInstrumentRequest — https://www.coinw.com/api-doc/en/futures-trading/market/get-order-book-of-an-instrument
type FuturesMarketGetOrderBookOfAnInstrumentRequest struct {
	Base Instrument `json:"base"`
}

// FuturesMarketGetTradeDataOfAnInstrumentRequest — https://www.coinw.com/api-doc/en/futures-trading/market/get-trade-data-of-an-instrument
type FuturesMarketGetTradeDataOfAnInstrumentRequest struct {
	Base Instrument `json:"base"`
}

// FuturesMarketGetMarginRequirementsOfAllInstrumentsRequest — https://www.coinw.com/api-doc/en/futures-trading/market/get-margin-requirements-of-all-instruments
type FuturesMarketGetMarginRequirementsOfAllInstrumentsRequest struct {
}

// FuturesMarketGetHistoricalPublicTradesRequest — https://www.coinw.com/api-doc/en/futures-trading/market/get-historical-public-trades
type FuturesMarketGetHistoricalPublicTradesRequest struct {
	Instrument Instrument `json:"instrument"`
	Page       *int64     `json:"page,omitempty"`
	PageSize   *int64     `json:"pageSize,omitempty"`
}

// FuturesAssetsCoinConversionRequest — https://www.coinw.com/api-doc/en/futures-trading/assets/coin-conversion
type FuturesAssetsCoinConversionRequest struct {
	ConvertType int64    `json:"convertType"`
	DealPiece   *Decimal `json:"dealPiece,omitempty"`
	FaceValue   Decimal  `json:"faceValue"`
	BaseSize    *Decimal `json:"baseSize,omitempty"`
}

// FuturesAssetsEnableDisableMegaCouponRequest — https://www.coinw.com/api-doc/en/futures-trading/assets/enable-disable-mega-coupon
type FuturesAssetsEnableDisableMegaCouponRequest struct {
	Status string `json:"status"`
}

// FuturesAssetsGetFuturesAccountAssetsRequest — https://www.coinw.com/api-doc/en/futures-trading/assets/get-futures-account-assets
type FuturesAssetsGetFuturesAccountAssetsRequest struct {
	Quote *string `json:"quote,omitempty"`
}

// FuturesAssetsGetFuturesAccountFeesRequest — https://www.coinw.com/api-doc/en/futures-trading/assets/get-futures-account-fees
type FuturesAssetsGetFuturesAccountFeesRequest struct {
}

// FuturesAssetsGetMarginModeRequest — https://www.coinw.com/api-doc/en/futures-trading/assets/get-margin-mode
type FuturesAssetsGetMarginModeRequest struct {
}

// FuturesAssetsGetMaxTransferableBalanceRequest — https://www.coinw.com/api-doc/en/futures-trading/assets/get-max-transferable-balance
type FuturesAssetsGetMaxTransferableBalanceRequest struct {
}

// FuturesAssetsGetMaxUsersContractSizeRequest — https://www.coinw.com/api-doc/en/futures-trading/assets/get-max-users-contract-size
type FuturesAssetsGetMaxUsersContractSizeRequest struct {
	Instrument Instrument `json:"instrument"`
}

// FuturesAssetsGetMegaCouponBalanceRequest — https://www.coinw.com/api-doc/en/futures-trading/assets/get-mega-coupon-balance
type FuturesAssetsGetMegaCouponBalanceRequest struct {
	StartTime *int64 `json:"startTime,omitempty"`
	EndTime   *int64 `json:"endTime,omitempty"`
	Type      int64  `json:"type"`
}

// FuturesAssetsGetTransactionDetails3DaysRequest — https://www.coinw.com/api-doc/en/futures-trading/assets/get-transaction-details-3days
type FuturesAssetsGetTransactionDetails3DaysRequest struct {
	Page          *int64         `json:"page,omitempty"`
	PageSize      *int64         `json:"pageSize,omitempty"`
	OriginType    *string        `json:"originType,omitempty"`
	Instrument    Instrument     `json:"instrument"`
	PositionModel *PositionModel `json:"positionModel,omitempty"`
}

// FuturesAssetsGetTransactionDetails3MonthsRequest — https://www.coinw.com/api-doc/en/futures-trading/assets/get-transaction-details-3months
type FuturesAssetsGetTransactionDetails3MonthsRequest struct {
	Page          *int64         `json:"page,omitempty"`
	PageSize      *int64         `json:"pageSize,omitempty"`
	OriginType    *string        `json:"originType,omitempty"`
	Instrument    Instrument     `json:"instrument"`
	PositionModel *PositionModel `json:"positionModel,omitempty"`
}

// FuturesAssetsSetMarginModeRequest — https://www.coinw.com/api-doc/en/futures-trading/assets/set-margin-mode
type FuturesAssetsSetMarginModeRequest struct {
	PositionModel PositionModel `json:"positionModel"`
	Layout        int           `json:"layout"`
}

// FuturesCommonGetLeverageInformationRequest — https://www.coinw.com/api-doc/en/futures-trading/common/get-leverage-information
type FuturesCommonGetLeverageInformationRequest struct {
	PositionID *int64 `json:"positionId,omitempty"`
	OrderID    *int64 `json:"orderId,omitempty"`
}

// FuturesPositionGetCurrentPositionInformationRequest — https://www.coinw.com/api-doc/en/futures-trading/position/get-current-position-information
type FuturesPositionGetCurrentPositionInformationRequest struct {
	Instrument Instrument `json:"instrument"`
	OpenIds    *string    `json:"openIds,omitempty"`
}

// FuturesPositionGetAllCurrentPositionsRequest — https://www.coinw.com/api-doc/en/futures-trading/position/get-current-positions
type FuturesPositionGetAllCurrentPositionsRequest struct {
}

// FuturesPositionGetHistoricalPositionInformationRequest — https://www.coinw.com/api-doc/en/futures-trading/position/get-historical-position-information
type FuturesPositionGetHistoricalPositionInformationRequest struct {
	Instrument    *Instrument    `json:"instrument,omitempty"`
	PositionModel *PositionModel `json:"positionModel,omitempty"`
}

// FuturesPositionGetMaxContractSizeRequest — https://www.coinw.com/api-doc/en/futures-trading/position/get-max-contract-size
type FuturesPositionGetMaxContractSizeRequest struct {
	Leverage      int64         `json:"leverage"`
	Instrument    Instrument    `json:"instrument"`
	OrderPrice    *Decimal      `json:"orderPrice,omitempty"`
	PositionModel PositionModel `json:"positionModel"`
}

// FuturesPositionGetPositionMarginRateRequest — https://www.coinw.com/api-doc/en/futures-trading/position/get-position-margin-rate
type FuturesPositionGetPositionMarginRateRequest struct {
	PositionID *int64 `json:"positionId,omitempty"`
}

// FuturesCheckGetBatchHistoricalOrders7DaysRequest — https://www.coinw.com/api-doc/en/futures-trading/check/get-batch-historical-orders-7days
type FuturesCheckGetBatchHistoricalOrders7DaysRequest struct {
	Page       *int64         `json:"page,omitempty"`
	PageSize   *int64         `json:"pageSize,omitempty"`
	OriginType string         `json:"originType"`
	Symbols    InstrumentList `json:"symbols"`
}

// FuturesCheckGetBatchHistoricalOrders3MonthsRequest — https://www.coinw.com/api-doc/en/futures-trading/check/get-batch-historical-orders-3months
type FuturesCheckGetBatchHistoricalOrders3MonthsRequest struct {
	Page       *int64         `json:"page,omitempty"`
	PageSize   *int64         `json:"pageSize,omitempty"`
	OriginType string         `json:"originType"`
	Symbols    InstrumentList `json:"symbols"`
}

// FuturesCheckGetCurrentOrdersRequest — https://www.coinw.com/api-doc/en/futures-trading/check/get-current-orders
type FuturesCheckGetCurrentOrdersRequest struct {
	Instrument   Instrument `json:"instrument"`
	PositionType OrderType  `json:"positionType"`
	Page         *int64     `json:"page,omitempty"`
	PageSize     *int64     `json:"pageSize,omitempty"`
}

// FuturesCheckGetHistoricalOrders7DaysRequest — https://www.coinw.com/api-doc/en/futures-trading/check/get-historical-orders-7days
type FuturesCheckGetHistoricalOrders7DaysRequest struct {
	Page       *int64      `json:"page,omitempty"`
	PageSize   *int64      `json:"pageSize,omitempty"`
	OriginType *string     `json:"originType,omitempty"`
	Instrument *Instrument `json:"instrument,omitempty"`
}

// FuturesCheckGetHistoricalOrders3MonthsRequest — https://www.coinw.com/api-doc/en/futures-trading/check/get-historical-orders-3months
type FuturesCheckGetHistoricalOrders3MonthsRequest struct {
	Page       *int64      `json:"page,omitempty"`
	PageSize   *int64      `json:"pageSize,omitempty"`
	OriginType *string     `json:"originType,omitempty"`
	Instrument *Instrument `json:"instrument,omitempty"`
}

// FuturesCheckGetOrderInformationRequest — https://www.coinw.com/api-doc/en/futures-trading/check/get-order-information
type FuturesCheckGetOrderInformationRequest struct {
	SourceIds    *string     `json:"sourceIds,omitempty"`
	PositionType OrderType   `json:"positionType"`
	Instrument   *Instrument `json:"instrument,omitempty"`
}

// FuturesCheckGetPendingOrderCountRequest — https://www.coinw.com/api-doc/en/futures-trading/check/get-pending-order-count
type FuturesCheckGetPendingOrderCountRequest struct {
}

// FuturesCheckGetStopLossTakeProfitInformationRequest — https://www.coinw.com/api-doc/en/futures-trading/check/get-stop-loss-take-profit-information
type FuturesCheckGetStopLossTakeProfitInformationRequest struct {
	OrderID     *int64      `json:"orderId,omitempty"`
	OpenID      *int64      `json:"openId,omitempty"`
	PlanOrderID *int64      `json:"planOrderId,omitempty"`
	StopFrom    int64       `json:"stopFrom"`
	Instrument  *Instrument `json:"instrument,omitempty"`
}

// FuturesCheckGetTrailingStopLossTakeProfitInformationRequest — https://www.coinw.com/api-doc/en/futures-trading/check/get-trailing-stop-loss-take-profit-information
type FuturesCheckGetTrailingStopLossTakeProfitInformationRequest struct {
}
