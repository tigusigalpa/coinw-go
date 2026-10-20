package model

// DepositAddress is a documented CoinW response record.
type DepositAddress struct {
	ChainName         string          `json:"chainName" required:"true"`
	Address           StringOrInteger `json:"address" required:"true"`
	Memo              *string         `json:"memo,omitempty"`
	MinRechargeAmount WireDecimal     `json:"minRechargeAmount" required:"true"`
}

// DepositWithdrawal is a documented CoinW response record.
type DepositWithdrawal struct {
	Amount        WireDecimal      `json:"amount" required:"true"`
	Chain         *string          `json:"chain,omitempty"`
	Side          *int             `json:"side,omitempty"`
	DepositNumber *StringOrInteger `json:"depositNumber,omitempty"`
	Address       *string          `json:"address,omitempty"`
	Txid          *string          `json:"txid,omitempty"`
	Memo          *string          `json:"memo,omitempty"`
	Currency      string           `json:"currency" required:"true"`
	Time          StringOrInteger  `json:"time" required:"true"`
	Confirmations *StringOrInteger `json:"confirmations,omitempty"`
	Status        StringOrInteger  `json:"status" required:"true"`
	Dest          *string          `json:"dest,omitempty"`
	FromAddress   *string          `json:"fromAddress,omitempty"`
	ToAddress     *string          `json:"toAddress,omitempty"`
}

// WithdrawalAcknowledgement is a documented CoinW response record.
type WithdrawalAcknowledgement struct {
	DepositNumber StringOrInteger `json:"depositNumber" required:"true"`
}

// SubAccount is a documented CoinW response record.
type SubAccount struct {
	SpotUID           StringOrInteger `json:"spotUid" required:"true"`
	UserCode          string          `json:"userCode" required:"true"`
	NickName          *string         `json:"nickName,omitempty"`
	FrozenFlag        *int            `json:"frozenFlag,omitempty"`
	SpotTradeFlag     *int            `json:"spotTradeFlag,omitempty"`
	ContractTradeFlag *int            `json:"contractTradeFlag,omitempty"`
}

// SubAccountAPIKey is a documented CoinW response record.
type SubAccountAPIKey struct {
	UserCode          string           `json:"userCode" required:"true"`
	NickName          *string          `json:"nickName,omitempty"`
	APIRemark         *string          `json:"apiRemark,omitempty"`
	APIKey            string           `json:"apiKey" required:"true"`
	CreatedAt         *StringOrInteger `json:"createdAt,omitempty"`
	SpotTradeFlag     *int             `json:"spotTradeFlag,omitempty"`
	ContractTradeFlag *int             `json:"contractTradeFlag,omitempty"`
	ReadonlyFlag      *int             `json:"readonlyFlag,omitempty"`
	APIWhiteListFlag  *int             `json:"apiWhiteListFlag,omitempty"`
	APIWhiteList      *string          `json:"apiWhiteList,omitempty"`
}

// Transfer is a documented CoinW response record.
type Transfer struct {
	Currency   string           `json:"currency" required:"true"`
	CreateTime *string          `json:"createTime,omitempty"`
	Amount     WireDecimal      `json:"amount" required:"true"`
	From       *int             `json:"from,omitempty"`
	To         *int             `json:"to,omitempty"`
	Status     *StringOrInteger `json:"status,omitempty"`
	LedgerID   StringOrInteger  `json:"ledgerId" required:"true"`
}

// TransferAcknowledgement is a documented CoinW response record.
type TransferAcknowledgement struct {
	BizNo StringOrInteger `json:"bizNo" required:"true"`
}

// InvitedUser is a documented CoinW response record.
type InvitedUser struct {
	UID          StringOrInteger `json:"uid" required:"true"`
	RegisterTime string          `json:"registerTime" required:"true"`
}

// InvitedUserKYC is a documented CoinW response record.
type InvitedUserKYC struct {
	UID       StringOrInteger `json:"uid" required:"true"`
	KycResult string          `json:"kycResult" required:"true"`
}

// InvitedUserVolume is a documented CoinW response record.
type InvitedUserVolume struct {
	UID    StringOrInteger `json:"uid" required:"true"`
	Volume WireDecimal     `json:"volume" required:"true"`
	Time   string          `json:"time" required:"true"`
}

// ActivityDay is a documented CoinW response record.
type ActivityDay struct {
	Date             string       `json:"date" required:"true"`
	ContractTransIn  *WireDecimal `json:"contractTransIn,omitempty"`
	ContractTransOut *WireDecimal `json:"contractTransOut,omitempty"`
	ContractVolume   *WireDecimal `json:"contractVolume,omitempty"`
}

// AffiliateActivity is a documented CoinW response record.
type AffiliateActivity struct {
	UID                StringOrInteger `json:"uid" required:"true"`
	ContractTransIn    *WireDecimal    `json:"contractTransIn,omitempty"`
	ContractTransOut   *WireDecimal    `json:"contractTransOut,omitempty"`
	ContractNetTransIn *WireDecimal    `json:"contractNetTransIn,omitempty"`
	ContractVolume     *WireDecimal    `json:"contractVolume,omitempty"`
	ContractProfit     *WireDecimal    `json:"contractProfit,omitempty"`
	ContractRoi        *WireDecimal    `json:"contractRoi,omitempty"`
	InitialNetValue    *WireDecimal    `json:"initialNetValue,omitempty"`
	Detail             []ActivityDay   `json:"detail" required:"true"`
}

// InvitedUserAssets is a documented CoinW response record.
type InvitedUserAssets struct {
	UID                StringOrInteger `json:"uid" required:"true"`
	TotalDeposit       *WireDecimal    `json:"totalDeposit,omitempty"`
	TotalWithdrawal    *WireDecimal    `json:"totalWithdrawal,omitempty"`
	CurrentBalance     WireDecimal     `json:"currentBalance" required:"true"`
	FirstDepositAmount *WireDecimal    `json:"firstDepositAmount,omitempty"`
	UpdateTime         *string         `json:"updateTime,omitempty"`
	FirstDepositTime   *string         `json:"firstDepositTime,omitempty"`
}

// AffiliateTeamUser is a documented CoinW response record.
type AffiliateTeamUser struct {
	UID                             StringOrInteger `json:"uid" required:"true"`
	RegisterTime                    *string         `json:"registerTime,omitempty"`
	ReferralCode                    *string         `json:"referralCode,omitempty"`
	Deposit                         *string         `json:"deposit,omitempty"`
	TradeVolume                     *WireDecimal    `json:"tradeVolume,omitempty"`
	TradingFee                      *WireDecimal    `json:"tradingFee,omitempty"`
	Rebate                          *WireDecimal    `json:"rebate,omitempty"`
	TotalDepositAmount              *WireDecimal    `json:"totalDepositAmount,omitempty"`
	TotalDepositGreaterThanHundred  *string         `json:"totalDepositGreaterThanHundred,omitempty"`
	FirstDepositGreaterThan3Hundred *string         `json:"firstDepositGreaterThan3Hundred,omitempty"`
	FirstWithdrawLessThan2Hundred   *string         `json:"firstWithdrawLessThan2Hundred,omitempty"`
	FirstWithdrawDate               *string         `json:"firstWithdrawDate,omitempty"`
	FirstTradeDate                  *string         `json:"firstTradeDate,omitempty"`
	ContractBalance                 *MaskedDecimal  `json:"contractBalance,omitempty"`
}

// AffiliateTeam is a documented CoinW response record.
type AffiliateTeam struct {
	Total                    StringOrInteger     `json:"total" required:"true"`
	TotalUsers               *StringOrInteger    `json:"totalUsers,omitempty"`
	TotalUsersInTeam         *StringOrInteger    `json:"totalUsersInTeam,omitempty"`
	TotalTradingVolumeInTeam *WireDecimal        `json:"totalTradingVolumeInTeam,omitempty"`
	List                     []AffiliateTeamUser `json:"list" required:"true"`
}

// SpotInstrument is a documented CoinW response record.
type SpotInstrument struct {
	CurrencyPair   string       `json:"currencyPair" required:"true"`
	CurrencyBase   string       `json:"currencyBase" required:"true"`
	CurrencyQuote  string       `json:"currencyQuote" required:"true"`
	MaxBuyCount    *WireDecimal `json:"maxBuyCount,omitempty"`
	MinBuyCount    *WireDecimal `json:"minBuyCount,omitempty"`
	MinBuyAmount   *WireDecimal `json:"minBuyAmount,omitempty"`
	MaxBuyAmount   *WireDecimal `json:"maxBuyAmount,omitempty"`
	MinBuyPrice    *WireDecimal `json:"minBuyPrice,omitempty"`
	MaxBuyPrice    *WireDecimal `json:"maxBuyPrice,omitempty"`
	PricePrecision int          `json:"pricePrecision" required:"true"`
	CountPrecision int          `json:"countPrecision" required:"true"`
	State          *int         `json:"state,omitempty"`
}

// SpotCurrency is a documented CoinW response record.
type SpotCurrency struct {
	SymbolID StringOrInteger  `json:"symbolId" required:"true"`
	Symbol   string           `json:"symbol" required:"true"`
	WithDraw *StringOrInteger `json:"withDraw,omitempty"`
	Recharge *StringOrInteger `json:"recharge,omitempty"`
	MaxQty   *WireDecimal     `json:"maxQty,omitempty"`
	MinQty   *WireDecimal     `json:"minQty,omitempty"`
	Chain    *string          `json:"chain,omitempty"`
}

// SpotTicker is a documented CoinW response record.
type SpotTicker struct {
	ID            StringOrInteger  `json:"id" required:"true"`
	IsFrozen      *StringOrInteger `json:"isFrozen,omitempty"`
	Last          WireDecimal      `json:"last" required:"true"`
	LowestAsk     *WireDecimal     `json:"lowestAsk,omitempty"`
	HighestBid    *WireDecimal     `json:"highestBid,omitempty"`
	PercentChange *WireDecimal     `json:"percentChange,omitempty"`
	High24hr      *WireDecimal     `json:"high24hr,omitempty"`
	Low24hr       *WireDecimal     `json:"low24hr,omitempty"`
	BaseVolume    *WireDecimal     `json:"baseVolume,omitempty"`
}

// SpotTrade is a documented CoinW response record.
type SpotTrade struct {
	ID     StringOrInteger `json:"id" required:"true"`
	Type   *string         `json:"type,omitempty"`
	Price  WireDecimal     `json:"price" required:"true"`
	Amount WireDecimal     `json:"amount" required:"true"`
	Total  *WireDecimal    `json:"total,omitempty"`
	Time   StringOrInteger `json:"time" required:"true"`
	Pair   *string         `json:"pair,omitempty"`
}

// SpotOrderBook is a documented CoinW response record.
type SpotOrderBook struct {
	Asks []PriceQuantity `json:"asks" required:"true"`
	Bids []PriceQuantity `json:"bids" required:"true"`
	Pair string          `json:"pair" required:"true"`
}

// SpotFullBalance is a documented CoinW response record.
type SpotFullBalance struct {
	Available WireDecimal `json:"available" required:"true"`
	OnOrders  WireDecimal `json:"onOrders" required:"true"`
}

// SpotCandle is a documented CoinW response record.
type SpotCandle struct {
	Date   StringOrInteger `json:"date" required:"true"`
	High   WireDecimal     `json:"high" required:"true"`
	Low    WireDecimal     `json:"low" required:"true"`
	Open   WireDecimal     `json:"open" required:"true"`
	Close  WireDecimal     `json:"close" required:"true"`
	Volume WireDecimal     `json:"volume" required:"true"`
	Pair   string          `json:"pair" required:"true"`
}

// SpotCurrentOrder is a documented CoinW response record.
type SpotCurrentOrder struct {
	OrderNumber    StringOrInteger `json:"orderNumber" required:"true"`
	Date           StringOrInteger `json:"date" required:"true"`
	StartingAmount *WireDecimal    `json:"startingAmount,omitempty"`
	Total          *WireDecimal    `json:"total,omitempty"`
	Prize          *WireDecimal    `json:"prize,omitempty"`
	SuccessCount   *WireDecimal    `json:"success_count,omitempty"`
	SuccessAmount  *WireDecimal    `json:"success_amount,omitempty"`
	Type           string          `json:"type" required:"true"`
	Status         StringOrInteger `json:"status" required:"true"`
}

// SpotBatchHistoricalOrder is a documented CoinW response record.
type SpotBatchHistoricalOrder struct {
	OrderID      StringOrInteger  `json:"orderId" required:"true"`
	Date         *StringOrInteger `json:"date,omitempty"`
	Side         string           `json:"side" required:"true"`
	Type         *string          `json:"type,omitempty"`
	DealSize     *WireDecimal     `json:"dealSize,omitempty"`
	DealFunds    *WireDecimal     `json:"dealFunds,omitempty"`
	DealAvgPrice *WireDecimal     `json:"dealAvgPrice,omitempty"`
	Fee          *WireDecimal     `json:"fee,omitempty"`
	Size         *WireDecimal     `json:"size,omitempty"`
	Price        *WireDecimal     `json:"price,omitempty"`
	Status       StringOrInteger  `json:"status" required:"true"`
	Symbol       *string          `json:"symbol,omitempty"`
}

// SpotHistoricalOrder is a documented CoinW response record.
type SpotHistoricalOrder struct {
	TradeID   StringOrInteger `json:"tradeId" required:"true"`
	OrderID   StringOrInteger `json:"orderId" required:"true"`
	Price     WireDecimal     `json:"price" required:"true"`
	Size      WireDecimal     `json:"size" required:"true"`
	Fee       *WireDecimal    `json:"fee,omitempty"`
	Side      *string         `json:"side,omitempty"`
	OrderType *string         `json:"orderType,omitempty"`
	Time      StringOrInteger `json:"time" required:"true"`
}

// SpotOrderStatus is a documented CoinW response record.
type SpotOrderStatus struct {
	CurrencyPair   string           `json:"currencyPair" required:"true"`
	Type           string           `json:"type" required:"true"`
	Total          *WireDecimal     `json:"total,omitempty"`
	StartingAmount *WireDecimal     `json:"startingAmount,omitempty"`
	Status         StringOrInteger  `json:"status" required:"true"`
	Date           *StringOrInteger `json:"date,omitempty"`
}

// SpotOrderDetails is a documented CoinW response record.
type SpotOrderDetails struct {
	TradeID       StringOrInteger  `json:"tradeID" required:"true"`
	CurrencyPair  string           `json:"currencyPair" required:"true"`
	Type          *string          `json:"type,omitempty"`
	Amount        *WireDecimal     `json:"amount,omitempty"`
	SuccessAmount *WireDecimal     `json:"success_amount,omitempty"`
	Total         *WireDecimal     `json:"total,omitempty"`
	SuccessTotal  *WireDecimal     `json:"success_total,omitempty"`
	Fee           *WireDecimal     `json:"fee,omitempty"`
	Date          *StringOrInteger `json:"date,omitempty"`
	Status        StringOrInteger  `json:"status" required:"true"`
}

// SpotTransaction is a documented CoinW response record.
type SpotTransaction struct {
	TradeID       StringOrInteger  `json:"tradeID" required:"true"`
	Type          string           `json:"type" required:"true"`
	Amount        *WireDecimal     `json:"amount,omitempty"`
	SuccessAmount *WireDecimal     `json:"success_amount,omitempty"`
	Total         *WireDecimal     `json:"total,omitempty"`
	SuccessCount  *WireDecimal     `json:"success_count,omitempty"`
	Fee           *WireDecimal     `json:"fee,omitempty"`
	Prize         *WireDecimal     `json:"prize,omitempty"`
	Date          *StringOrInteger `json:"date,omitempty"`
	Status        StringOrInteger  `json:"status" required:"true"`
	Pair          *string          `json:"pair,omitempty"`
	OutTradeNo    *StringOrInteger `json:"out_trade_no,omitempty"`
}

// SpotOrderAcknowledgement is a documented CoinW response record.
type SpotOrderAcknowledgement struct {
	OrderNumber StringOrInteger `json:"orderNumber" required:"true"`
}

// SpotCancelAcknowledgement is a documented CoinW response record.
type SpotCancelAcknowledgement struct {
	ClientOrderID StringOrInteger `json:"clientOrderId" required:"true"`
}

// FuturesInstrumentConfig is a documented CoinW response record.
type FuturesInstrumentConfig struct {
	Margins          map[string]WireDecimal `json:"margins,omitempty"`
	SimulatedMargins map[string]WireDecimal `json:"simulatedMargins,omitempty"`
}

// FuturesInstrument is a documented CoinW response record.
type FuturesInstrument struct {
	Base                  string                                `json:"base" required:"true"`
	DefaultLeverage       *int                                  `json:"defaultLeverage,omitempty"`
	DefaultStopLossRate   *WireDecimal                          `json:"defaultStopLossRate,omitempty"`
	DefaultStopProfitRate *WireDecimal                          `json:"defaultStopProfitRate,omitempty"`
	IndexID               *StringOrInteger                      `json:"indexId,omitempty"`
	Leverage              *LeverageOptions                      `json:"leverage,omitempty"`
	MakerFee              *WireDecimal                          `json:"makerFee,omitempty"`
	MaxLeverage           *int                                  `json:"maxLeverage,omitempty"`
	MaxPosition           *WireDecimal                          `json:"maxPosition,omitempty"`
	MinLeverage           *int                                  `json:"minLeverage,omitempty"`
	MinSize               *WireDecimal                          `json:"minSize,omitempty"`
	Name                  string                                `json:"name" required:"true"`
	OneLotMargin          *WireDecimal                          `json:"oneLotMargin,omitempty"`
	OneMaxPosition        *WireDecimal                          `json:"oneMaxPosition,omitempty"`
	PricePrecision        *int                                  `json:"pricePrecision,omitempty"`
	Quote                 string                                `json:"quote" required:"true"`
	SettledAt             *StringOrInteger                      `json:"settledAt,omitempty"`
	SettlementRate        *WireDecimal                          `json:"settlementRate,omitempty"`
	Status                *string                               `json:"status,omitempty"`
	TakerFee              *WireDecimal                          `json:"takerFee,omitempty"`
	StopCrossPositionRate *WireDecimal                          `json:"stopCrossPositionRate,omitempty"`
	StopSurplusRate       *WireDecimal                          `json:"stopSurplusRate,omitempty"`
	UpdatedDate           *StringOrInteger                      `json:"updatedDate,omitempty"`
	Sort                  *int                                  `json:"sort,omitempty"`
	Selected              *int                                  `json:"selected,omitempty"`
	OpenSpread            *WireDecimal                          `json:"openSpread,omitempty"`
	OneLotSize            *WireDecimal                          `json:"oneLotSize,omitempty"`
	CloseSpread           *WireDecimal                          `json:"closeSpread,omitempty"`
	CommissionRate        *WireDecimal                          `json:"commissionRate,omitempty"`
	PartitionIds          *string                               `json:"partitionIds,omitempty"`
	DepthPrecision        *string                               `json:"depthPrecision,omitempty"`
	IconUrl               *string                               `json:"iconUrl,omitempty"`
	ID                    *StringOrInteger                      `json:"id,omitempty"`
	CreatedDate           *StringOrInteger                      `json:"createdDate,omitempty"`
	SettledPeriod         *int                                  `json:"settledPeriod,omitempty"`
	ConfigBo              *JSONEncoded[FuturesInstrumentConfig] `json:"configBo,omitempty"`
	Margins               map[string]WireDecimal                `json:"margins,omitempty"`
	SimulatedMargins      map[string]WireDecimal                `json:"simulatedMargins,omitempty"`
}

// FuturesTicker is a documented CoinW response record.
type FuturesTicker struct {
	ContractID   *StringOrInteger `json:"contract_id,omitempty"`
	Name         string           `json:"name" required:"true"`
	BaseCoin     *string          `json:"base_coin,omitempty"`
	QuoteCoin    *string          `json:"quote_coin,omitempty"`
	PriceCoin    *string          `json:"price_coin,omitempty"`
	MaxLeverage  *int             `json:"max_leverage,omitempty"`
	ContractSize *WireDecimal     `json:"contract_size,omitempty"`
	LastPrice    WireDecimal      `json:"last_price" required:"true"`
	High         *WireDecimal     `json:"high,omitempty"`
	Low          *WireDecimal     `json:"low,omitempty"`
	RiseFallRate *WireDecimal     `json:"rise_fall_rate,omitempty"`
	TotalVolume  *WireDecimal     `json:"total_volume,omitempty"`
	FairPrice    *WireDecimal     `json:"fair_price,omitempty"`
	TS           *StringOrInteger `json:"ts,omitempty"`
}

// FuturesDepthLevel is a documented CoinW response record.
type FuturesDepthLevel struct {
	M WireDecimal `json:"m" required:"true"`
	P WireDecimal `json:"p" required:"true"`
}

// FuturesOrderBook is a documented CoinW response record.
type FuturesOrderBook struct {
	Asks []FuturesDepthLevel `json:"asks" required:"true"`
	Bids []FuturesDepthLevel `json:"bids" required:"true"`
	N    string              `json:"n" required:"true"`
	T    *OpaqueNull         `json:"t,omitempty"`
	TS   StringOrInteger     `json:"ts" required:"true"`
}

// FuturesTrade is a documented CoinW response record.
type FuturesTrade struct {
	CreatedDate StringOrInteger `json:"createdDate" required:"true"`
	Piece       *WireDecimal    `json:"piece,omitempty"`
	Price       WireDecimal     `json:"price" required:"true"`
	Quantity    WireDecimal     `json:"quantity" required:"true"`
	Direction   *string         `json:"direction,omitempty"`
	ID          StringOrInteger `json:"id" required:"true"`
}

// DecimalAtTime is a documented CoinW response record.
type DecimalAtTime struct {
	Value WireDecimal      `json:"value" required:"true"`
	TS    *StringOrInteger `json:"ts,omitempty"`
}

// FuturesLadder is a documented CoinW response record.
type FuturesLadder struct {
	ID              *StringOrInteger `json:"id,omitempty"`
	Instrument      string           `json:"instrument" required:"true"`
	Ladder          int              `json:"ladder" required:"true"`
	LastLadder      *bool            `json:"lastLadder,omitempty"`
	MarginKeepRate  *WireDecimal     `json:"marginKeepRate,omitempty"`
	MarginStartRate *WireDecimal     `json:"marginStartRate,omitempty"`
	MaxLeverage     int              `json:"maxLeverage" required:"true"`
	MaxPiece        *WireDecimal     `json:"maxPiece,omitempty"`
}

// FuturesLadderGroup is a documented CoinW response record.
type FuturesLadderGroup struct {
	Name       string          `json:"name" required:"true"`
	LadderList []FuturesLadder `json:"ladderList" required:"true"`
}

// FuturesLadders is a documented CoinW response record.
type FuturesLadders struct {
	LadderConfig []FuturesLadderGroup `json:"ladderConfig" required:"true"`
}

// HistoricalPublicTrade is a documented CoinW response record.
type HistoricalPublicTrade struct {
	ClosedPiece *WireDecimal    `json:"closedPiece,omitempty"`
	CreatedDate StringOrInteger `json:"createdDate" required:"true"`
	DealPrice   WireDecimal     `json:"dealPrice" required:"true"`
	Direction   *string         `json:"direction,omitempty"`
	ID          StringOrInteger `json:"id" required:"true"`
}

// DecimalValue is a documented CoinW response record.
type DecimalValue struct {
	Value WireDecimal `json:"value" required:"true"`
}

// FuturesAccountAssets is a documented CoinW response record.
type FuturesAccountAssets struct {
	AvailableMargin WireDecimal      `json:"availableMargin" required:"true"`
	AlmightyGold    *WireDecimal     `json:"almightyGold,omitempty"`
	AvailableUSDT   *WireDecimal     `json:"availableUsdt,omitempty"`
	AlMargin        *WireDecimal     `json:"alMargin,omitempty"`
	AlFreeze        *WireDecimal     `json:"alFreeze,omitempty"`
	UserID          StringOrInteger  `json:"userId" required:"true"`
	Time            *StringOrInteger `json:"time,omitempty"`
	TS              *StringOrInteger `json:"ts,omitempty"`
}

// FuturesAccountFees is a documented CoinW response record.
type FuturesAccountFees struct {
	MakerFee WireDecimal      `json:"makerFee" required:"true"`
	TakerFee WireDecimal      `json:"takerFee" required:"true"`
	UserID   *StringOrInteger `json:"userId,omitempty"`
}

// MarginMode is a documented CoinW response record.
type MarginMode struct {
	Layout        StringOrInteger  `json:"layout" required:"true"`
	PositionModel StringOrInteger  `json:"positionModel" required:"true"`
	UserID        *StringOrInteger `json:"userId,omitempty"`
}

// AvailableContractSize is a documented CoinW response record.
type AvailableContractSize struct {
	AvailBuy  WireDecimal `json:"availBuy" required:"true"`
	AvailSell WireDecimal `json:"availSell" required:"true"`
}

// MegaCoupon is a documented CoinW response record.
type MegaCoupon struct {
	ID             *StringOrInteger `json:"id,omitempty"`
	AssetsOut      *int             `json:"assetsOut,omitempty"`
	AgRecordID     *StringOrInteger `json:"agRecordId,omitempty"`
	Kyc            *int             `json:"kyc,omitempty"`
	CurrentAmount  WireDecimal      `json:"currentAmount" required:"true"`
	TotalAmount    WireDecimal      `json:"totalAmount" required:"true"`
	Type           *int             `json:"type,omitempty"`
	StartTime      *StringOrInteger `json:"startTime,omitempty"`
	EndTime        *StringOrInteger `json:"endTime,omitempty"`
	CreatedDate    *StringOrInteger `json:"createdDate,omitempty"`
	UpdateDate     *StringOrInteger `json:"updateDate,omitempty"`
	Remark         *string          `json:"remark,omitempty"`
	UserID         *StringOrInteger `json:"userId,omitempty"`
	BackAmount     *WireDecimal     `json:"backAmount,omitempty"`
	BackMark       *int             `json:"backMark,omitempty"`
	Name           *string          `json:"name,omitempty"`
	ProcessStatus  *int             `json:"processStatus,omitempty"`
	TransferExpire *int             `json:"transferExpire,omitempty"`
	Status         *int             `json:"status,omitempty"`
}

// FuturesTransaction is a documented CoinW response record.
type FuturesTransaction struct {
	ID                   StringOrInteger  `json:"id" required:"true"`
	UserID               *StringOrInteger `json:"userId,omitempty"`
	ClosePrice           *WireDecimal     `json:"closePrice,omitempty"`
	ClosedPiece          *WireDecimal     `json:"closedPiece,omitempty"`
	CreatedDate          *StringOrInteger `json:"createdDate,omitempty"`
	CurrentPiece         *WireDecimal     `json:"currentPiece,omitempty"`
	Direction            *string          `json:"direction,omitempty"`
	Fee                  *WireDecimal     `json:"fee,omitempty"`
	FeeRate              *WireDecimal     `json:"feeRate,omitempty"`
	FundingSettle        *WireDecimal     `json:"fundingSettle,omitempty"`
	HedgeID              *StringOrInteger `json:"hedgeId,omitempty"`
	IndexPrice           *WireDecimal     `json:"indexPrice,omitempty"`
	Instrument           string           `json:"instrument" required:"true"`
	Leverage             *WireDecimal     `json:"leverage,omitempty"`
	LiquidateBy          *string          `json:"liquidateBy,omitempty"`
	Margin               *WireDecimal     `json:"margin,omitempty"`
	NetProfit            *WireDecimal     `json:"netProfit,omitempty"`
	OpenID               *StringOrInteger `json:"openId,omitempty"`
	OpenPrice            *WireDecimal     `json:"openPrice,omitempty"`
	OrderID              *StringOrInteger `json:"orderId,omitempty"`
	OrderPrice           *WireDecimal     `json:"orderPrice,omitempty"`
	OrderStatus          *string          `json:"orderStatus,omitempty"`
	OriginalType         *string          `json:"originalType,omitempty"`
	PositionMargin       *WireDecimal     `json:"positionMargin,omitempty"`
	PositionModel        *StringOrInteger `json:"positionModel,omitempty"`
	Quantity             *WireDecimal     `json:"quantity,omitempty"`
	QuantityUnit         *StringOrInteger `json:"quantityUnit,omitempty"`
	SettlementID         *StringOrInteger `json:"settlementId,omitempty"`
	Status               *StringOrInteger `json:"status,omitempty"`
	TakerMaker           *StringOrInteger `json:"takerMaker,omitempty"`
	TotalPiece           *WireDecimal     `json:"totalPiece,omitempty"`
	UpdatedDate          *StringOrInteger `json:"updatedDate,omitempty"`
	ThirdOrderID         *StringOrInteger `json:"thirdOrderId,omitempty"`
	ContractType         *StringOrInteger `json:"contractType,omitempty"`
	BaseSize             *WireDecimal     `json:"baseSize,omitempty"`
	IsProfession         *int             `json:"isProfession,omitempty"`
	ProcessStatus        *int             `json:"processStatus,omitempty"`
	Source               *string          `json:"source,omitempty"`
	FundingFee           *WireDecimal     `json:"fundingFee,omitempty"`
	TriggerPrice         *WireDecimal     `json:"triggerPrice,omitempty"`
	TriggerType          *StringOrInteger `json:"triggerType,omitempty"`
	PosType              *string          `json:"posType,omitempty"`
	FrozenFee            *WireDecimal     `json:"frozenFee,omitempty"`
	MakerFee             *WireDecimal     `json:"makerFee,omitempty"`
	TakerFee             *WireDecimal     `json:"takerFee,omitempty"`
	CompleteUSDT         *WireDecimal     `json:"completeUsdt,omitempty"`
	EntrustUSDT          *WireDecimal     `json:"entrustUsdt,omitempty"`
	HavShortfall         *bool            `json:"havShortfall,omitempty"`
	StepMarginVersion    *int             `json:"stepMarginVersion,omitempty"`
	CancelPiece          *WireDecimal     `json:"cancelPiece,omitempty"`
	TradePiece           *WireDecimal     `json:"tradePiece,omitempty"`
	AvgPrice             *WireDecimal     `json:"avgPrice,omitempty"`
	Base                 *string          `json:"base,omitempty"`
	Quote                *string          `json:"quote,omitempty"`
	FinalOrderStatus     *string          `json:"finalOrderStatus,omitempty"`
	PrimaryID            *StringOrInteger `json:"primaryId,omitempty"`
	StopLossPrice        *WireDecimal     `json:"stopLossPrice,omitempty"`
	StopLossRate         *WireDecimal     `json:"stopLossRate,omitempty"`
	StopProfitPrice      *WireDecimal     `json:"stopProfitPrice,omitempty"`
	StopProfitRate       *WireDecimal     `json:"stopProfitRate,omitempty"`
	StopLossOrderPrice   *WireDecimal     `json:"stopLossOrderPrice,omitempty"`
	StopProfitOrderPrice *WireDecimal     `json:"stopProfitOrderPrice,omitempty"`
	StopProfitType       *StringOrInteger `json:"stopProfitType,omitempty"`
	StopLossType         *StringOrInteger `json:"stopLossType,omitempty"`
}

// FuturesOrder is a documented CoinW response record.
type FuturesOrder struct {
	ID                   StringOrInteger  `json:"id" required:"true"`
	UserID               *StringOrInteger `json:"userId,omitempty"`
	ClosePrice           *WireDecimal     `json:"closePrice,omitempty"`
	ClosedPiece          *WireDecimal     `json:"closedPiece,omitempty"`
	CreatedDate          *StringOrInteger `json:"createdDate,omitempty"`
	CurrentPiece         *WireDecimal     `json:"currentPiece,omitempty"`
	Direction            *string          `json:"direction,omitempty"`
	Fee                  *WireDecimal     `json:"fee,omitempty"`
	FeeRate              *WireDecimal     `json:"feeRate,omitempty"`
	FundingSettle        *WireDecimal     `json:"fundingSettle,omitempty"`
	HedgeID              *StringOrInteger `json:"hedgeId,omitempty"`
	IndexPrice           *WireDecimal     `json:"indexPrice,omitempty"`
	Instrument           string           `json:"instrument" required:"true"`
	Leverage             *WireDecimal     `json:"leverage,omitempty"`
	LiquidateBy          *string          `json:"liquidateBy,omitempty"`
	Margin               *WireDecimal     `json:"margin,omitempty"`
	NetProfit            *WireDecimal     `json:"netProfit,omitempty"`
	OpenID               *StringOrInteger `json:"openId,omitempty"`
	OpenPrice            *WireDecimal     `json:"openPrice,omitempty"`
	OrderID              *StringOrInteger `json:"orderId,omitempty"`
	OrderPrice           *WireDecimal     `json:"orderPrice,omitempty"`
	OrderStatus          *string          `json:"orderStatus,omitempty"`
	OriginalType         *string          `json:"originalType,omitempty"`
	PositionMargin       *WireDecimal     `json:"positionMargin,omitempty"`
	PositionModel        *StringOrInteger `json:"positionModel,omitempty"`
	Quantity             *WireDecimal     `json:"quantity,omitempty"`
	QuantityUnit         *StringOrInteger `json:"quantityUnit,omitempty"`
	SettlementID         *StringOrInteger `json:"settlementId,omitempty"`
	Status               *StringOrInteger `json:"status,omitempty"`
	TakerMaker           *StringOrInteger `json:"takerMaker,omitempty"`
	TotalPiece           *WireDecimal     `json:"totalPiece,omitempty"`
	UpdatedDate          *StringOrInteger `json:"updatedDate,omitempty"`
	ThirdOrderID         *StringOrInteger `json:"thirdOrderId,omitempty"`
	ContractType         *StringOrInteger `json:"contractType,omitempty"`
	BaseSize             *WireDecimal     `json:"baseSize,omitempty"`
	IsProfession         *int             `json:"isProfession,omitempty"`
	ProcessStatus        *int             `json:"processStatus,omitempty"`
	Source               *string          `json:"source,omitempty"`
	FundingFee           *WireDecimal     `json:"fundingFee,omitempty"`
	TriggerPrice         *WireDecimal     `json:"triggerPrice,omitempty"`
	TriggerType          *StringOrInteger `json:"triggerType,omitempty"`
	PosType              *string          `json:"posType,omitempty"`
	FrozenFee            *WireDecimal     `json:"frozenFee,omitempty"`
	MakerFee             *WireDecimal     `json:"makerFee,omitempty"`
	TakerFee             *WireDecimal     `json:"takerFee,omitempty"`
	CompleteUSDT         *WireDecimal     `json:"completeUsdt,omitempty"`
	EntrustUSDT          *WireDecimal     `json:"entrustUsdt,omitempty"`
	HavShortfall         *bool            `json:"havShortfall,omitempty"`
	StepMarginVersion    *int             `json:"stepMarginVersion,omitempty"`
	CancelPiece          *WireDecimal     `json:"cancelPiece,omitempty"`
	TradePiece           *WireDecimal     `json:"tradePiece,omitempty"`
	AvgPrice             *WireDecimal     `json:"avgPrice,omitempty"`
	Base                 *string          `json:"base,omitempty"`
	Quote                *string          `json:"quote,omitempty"`
	FinalOrderStatus     *string          `json:"finalOrderStatus,omitempty"`
	PrimaryID            *StringOrInteger `json:"primaryId,omitempty"`
	StopLossPrice        *WireDecimal     `json:"stopLossPrice,omitempty"`
	StopLossRate         *WireDecimal     `json:"stopLossRate,omitempty"`
	StopProfitPrice      *WireDecimal     `json:"stopProfitPrice,omitempty"`
	StopProfitRate       *WireDecimal     `json:"stopProfitRate,omitempty"`
	StopLossOrderPrice   *WireDecimal     `json:"stopLossOrderPrice,omitempty"`
	StopProfitOrderPrice *WireDecimal     `json:"stopProfitOrderPrice,omitempty"`
	StopProfitType       *StringOrInteger `json:"stopProfitType,omitempty"`
	StopLossType         *StringOrInteger `json:"stopLossType,omitempty"`
}

// FuturesPosition is a documented CoinW response record.
type FuturesPosition struct {
	ID                    StringOrInteger  `json:"id" required:"true"`
	UserID                *StringOrInteger `json:"userId,omitempty"`
	ClosePrice            *WireDecimal     `json:"closePrice,omitempty"`
	ClosedPiece           *WireDecimal     `json:"closedPiece,omitempty"`
	CreatedDate           *StringOrInteger `json:"createdDate,omitempty"`
	CurrentPiece          WireDecimal      `json:"currentPiece" required:"true"`
	Direction             *string          `json:"direction,omitempty"`
	Fee                   *WireDecimal     `json:"fee,omitempty"`
	FeeRate               *WireDecimal     `json:"feeRate,omitempty"`
	FundingSettle         *WireDecimal     `json:"fundingSettle,omitempty"`
	HedgeID               *StringOrInteger `json:"hedgeId,omitempty"`
	IndexPrice            *WireDecimal     `json:"indexPrice,omitempty"`
	Instrument            string           `json:"instrument" required:"true"`
	Leverage              *WireDecimal     `json:"leverage,omitempty"`
	LiquidateBy           *string          `json:"liquidateBy,omitempty"`
	Margin                *WireDecimal     `json:"margin,omitempty"`
	NetProfit             *WireDecimal     `json:"netProfit,omitempty"`
	OpenID                *StringOrInteger `json:"openId,omitempty"`
	OpenPrice             *WireDecimal     `json:"openPrice,omitempty"`
	OrderID               *StringOrInteger `json:"orderId,omitempty"`
	OrderPrice            *WireDecimal     `json:"orderPrice,omitempty"`
	OrderStatus           *string          `json:"orderStatus,omitempty"`
	OriginalType          *string          `json:"originalType,omitempty"`
	PositionMargin        *WireDecimal     `json:"positionMargin,omitempty"`
	PositionModel         *StringOrInteger `json:"positionModel,omitempty"`
	Quantity              *WireDecimal     `json:"quantity,omitempty"`
	QuantityUnit          *StringOrInteger `json:"quantityUnit,omitempty"`
	SettlementID          *StringOrInteger `json:"settlementId,omitempty"`
	Status                *StringOrInteger `json:"status,omitempty"`
	TakerMaker            *StringOrInteger `json:"takerMaker,omitempty"`
	TotalPiece            *WireDecimal     `json:"totalPiece,omitempty"`
	UpdatedDate           *StringOrInteger `json:"updatedDate,omitempty"`
	ThirdOrderID          *StringOrInteger `json:"thirdOrderId,omitempty"`
	ContractType          *StringOrInteger `json:"contractType,omitempty"`
	BaseSize              *WireDecimal     `json:"baseSize,omitempty"`
	IsProfession          *int             `json:"isProfession,omitempty"`
	ProcessStatus         *int             `json:"processStatus,omitempty"`
	Source                *string          `json:"source,omitempty"`
	FundingFee            *WireDecimal     `json:"fundingFee,omitempty"`
	TriggerPrice          *WireDecimal     `json:"triggerPrice,omitempty"`
	TriggerType           *StringOrInteger `json:"triggerType,omitempty"`
	PosType               *string          `json:"posType,omitempty"`
	FrozenFee             *WireDecimal     `json:"frozenFee,omitempty"`
	MakerFee              *WireDecimal     `json:"makerFee,omitempty"`
	TakerFee              *WireDecimal     `json:"takerFee,omitempty"`
	CompleteUSDT          *WireDecimal     `json:"completeUsdt,omitempty"`
	EntrustUSDT           *WireDecimal     `json:"entrustUsdt,omitempty"`
	HavShortfall          *bool            `json:"havShortfall,omitempty"`
	StepMarginVersion     *int             `json:"stepMarginVersion,omitempty"`
	CancelPiece           *WireDecimal     `json:"cancelPiece,omitempty"`
	TradePiece            *WireDecimal     `json:"tradePiece,omitempty"`
	AvgPrice              *WireDecimal     `json:"avgPrice,omitempty"`
	Base                  *string          `json:"base,omitempty"`
	Quote                 *string          `json:"quote,omitempty"`
	FinalOrderStatus      *string          `json:"finalOrderStatus,omitempty"`
	PrimaryID             *StringOrInteger `json:"primaryId,omitempty"`
	StopLossPrice         *WireDecimal     `json:"stopLossPrice,omitempty"`
	StopLossRate          *WireDecimal     `json:"stopLossRate,omitempty"`
	StopProfitPrice       *WireDecimal     `json:"stopProfitPrice,omitempty"`
	StopProfitRate        *WireDecimal     `json:"stopProfitRate,omitempty"`
	StopLossOrderPrice    *WireDecimal     `json:"stopLossOrderPrice,omitempty"`
	StopProfitOrderPrice  *WireDecimal     `json:"stopProfitOrderPrice,omitempty"`
	StopProfitType        *StringOrInteger `json:"stopProfitType,omitempty"`
	StopLossType          *StringOrInteger `json:"stopLossType,omitempty"`
	AutoDeleveragingScore *StringOrInteger `json:"autoDeleveragingScore,omitempty"`
	ProfitUnreal          *WireDecimal     `json:"profitUnreal,omitempty"`
	LiquidationPrice      *WireDecimal     `json:"liquidationPrice,omitempty"`
	ProfitReal            *WireDecimal     `json:"profitReal,omitempty"`
	PositionVersion       *StringOrInteger `json:"positionVersion,omitempty"`
	ClosedQuantity        *WireDecimal     `json:"closedQuantity,omitempty"`
	DistID                *StringOrInteger `json:"distId,omitempty"`
	LeaderID              *StringOrInteger `json:"leaderId,omitempty"`
	ParentID              *StringOrInteger `json:"parentId,omitempty"`
	PartnerID             *StringOrInteger `json:"partnerId,omitempty"`
	RemainCurrentPiece    *WireDecimal     `json:"remainCurrentPiece,omitempty"`
	SalesID               *StringOrInteger `json:"salesId,omitempty"`
}

// FuturesPositionChange is a documented CoinW response record.
type FuturesPositionChange struct {
	ID                   StringOrInteger  `json:"id" required:"true"`
	UserID               *StringOrInteger `json:"userId,omitempty"`
	ClosePrice           *WireDecimal     `json:"closePrice,omitempty"`
	ClosedPiece          *WireDecimal     `json:"closedPiece,omitempty"`
	CreatedDate          *StringOrInteger `json:"createdDate,omitempty"`
	CurrentPiece         *WireDecimal     `json:"currentPiece,omitempty"`
	Direction            *string          `json:"direction,omitempty"`
	Fee                  *WireDecimal     `json:"fee,omitempty"`
	FeeRate              *WireDecimal     `json:"feeRate,omitempty"`
	FundingSettle        *WireDecimal     `json:"fundingSettle,omitempty"`
	HedgeID              *StringOrInteger `json:"hedgeId,omitempty"`
	IndexPrice           *WireDecimal     `json:"indexPrice,omitempty"`
	Instrument           string           `json:"instrument" required:"true"`
	Leverage             *WireDecimal     `json:"leverage,omitempty"`
	LiquidateBy          *string          `json:"liquidateBy,omitempty"`
	Margin               *WireDecimal     `json:"margin,omitempty"`
	NetProfit            *WireDecimal     `json:"netProfit,omitempty"`
	OpenID               *StringOrInteger `json:"openId,omitempty"`
	OpenPrice            *WireDecimal     `json:"openPrice,omitempty"`
	OrderID              *StringOrInteger `json:"orderId,omitempty"`
	OrderPrice           *WireDecimal     `json:"orderPrice,omitempty"`
	OrderStatus          *string          `json:"orderStatus,omitempty"`
	OriginalType         *string          `json:"originalType,omitempty"`
	PositionMargin       *WireDecimal     `json:"positionMargin,omitempty"`
	PositionModel        *StringOrInteger `json:"positionModel,omitempty"`
	Quantity             *WireDecimal     `json:"quantity,omitempty"`
	QuantityUnit         *StringOrInteger `json:"quantityUnit,omitempty"`
	SettlementID         *StringOrInteger `json:"settlementId,omitempty"`
	Status               *StringOrInteger `json:"status,omitempty"`
	TakerMaker           *StringOrInteger `json:"takerMaker,omitempty"`
	TotalPiece           *WireDecimal     `json:"totalPiece,omitempty"`
	UpdatedDate          *StringOrInteger `json:"updatedDate,omitempty"`
	ThirdOrderID         *StringOrInteger `json:"thirdOrderId,omitempty"`
	ContractType         *StringOrInteger `json:"contractType,omitempty"`
	BaseSize             *WireDecimal     `json:"baseSize,omitempty"`
	IsProfession         *int             `json:"isProfession,omitempty"`
	ProcessStatus        *int             `json:"processStatus,omitempty"`
	Source               *string          `json:"source,omitempty"`
	FundingFee           *WireDecimal     `json:"fundingFee,omitempty"`
	TriggerPrice         *WireDecimal     `json:"triggerPrice,omitempty"`
	TriggerType          *StringOrInteger `json:"triggerType,omitempty"`
	PosType              *string          `json:"posType,omitempty"`
	FrozenFee            *WireDecimal     `json:"frozenFee,omitempty"`
	MakerFee             *WireDecimal     `json:"makerFee,omitempty"`
	TakerFee             *WireDecimal     `json:"takerFee,omitempty"`
	CompleteUSDT         *WireDecimal     `json:"completeUsdt,omitempty"`
	EntrustUSDT          *WireDecimal     `json:"entrustUsdt,omitempty"`
	HavShortfall         *bool            `json:"havShortfall,omitempty"`
	StepMarginVersion    *int             `json:"stepMarginVersion,omitempty"`
	CancelPiece          *WireDecimal     `json:"cancelPiece,omitempty"`
	TradePiece           *WireDecimal     `json:"tradePiece,omitempty"`
	AvgPrice             *WireDecimal     `json:"avgPrice,omitempty"`
	Base                 *string          `json:"base,omitempty"`
	Quote                *string          `json:"quote,omitempty"`
	FinalOrderStatus     *string          `json:"finalOrderStatus,omitempty"`
	PrimaryID            *StringOrInteger `json:"primaryId,omitempty"`
	StopLossPrice        *WireDecimal     `json:"stopLossPrice,omitempty"`
	StopLossRate         *WireDecimal     `json:"stopLossRate,omitempty"`
	StopProfitPrice      *WireDecimal     `json:"stopProfitPrice,omitempty"`
	StopProfitRate       *WireDecimal     `json:"stopProfitRate,omitempty"`
	StopLossOrderPrice   *WireDecimal     `json:"stopLossOrderPrice,omitempty"`
	StopProfitOrderPrice *WireDecimal     `json:"stopProfitOrderPrice,omitempty"`
	StopProfitType       *StringOrInteger `json:"stopProfitType,omitempty"`
	StopLossType         *StringOrInteger `json:"stopLossType,omitempty"`
	RealPrice            *WireDecimal     `json:"realPrice,omitempty"`
}

// HistoricalPosition is a documented CoinW response record.
type HistoricalPosition struct {
	AvgOpenPrice    *WireDecimal     `json:"avgOpenPrice,omitempty"`
	CompleteUSDT    *WireDecimal     `json:"completeUsdt,omitempty"`
	EntrustUSDT     *WireDecimal     `json:"entrustUsdt,omitempty"`
	Fee             *WireDecimal     `json:"fee,omitempty"`
	IndexPrice      *WireDecimal     `json:"indexPrice,omitempty"`
	Leverage        *WireDecimal     `json:"leverage,omitempty"`
	Margin          *WireDecimal     `json:"margin,omitempty"`
	TotalPiece      *WireDecimal     `json:"totalPiece,omitempty"`
	TradePiece      *WireDecimal     `json:"tradePiece,omitempty"`
	NetProfit       *WireDecimal     `json:"netProfit,omitempty"`
	Direction       *string          `json:"direction,omitempty"`
	HavShortfall    *bool            `json:"havShortfall,omitempty"`
	Instrument      string           `json:"instrument" required:"true"`
	LiquidateBy     *string          `json:"liquidateBy,omitempty"`
	OpenID          *StringOrInteger `json:"openId,omitempty"`
	OrderID         StringOrInteger  `json:"orderId" required:"true"`
	UserID          *StringOrInteger `json:"userId,omitempty"`
	OrderStatus     *string          `json:"orderStatus,omitempty"`
	OriginalType    *string          `json:"originalType,omitempty"`
	Status          *StringOrInteger `json:"status,omitempty"`
	PositionModel   *StringOrInteger `json:"positionModel,omitempty"`
	TradeStartDate  *StringOrInteger `json:"tradeStartDate,omitempty"`
	AvgClosePrice   *WireDecimal     `json:"avgClosePrice,omitempty"`
	AvgPrice        *WireDecimal     `json:"avgPrice,omitempty"`
	Base            *string          `json:"base,omitempty"`
	Quote           *string          `json:"quote,omitempty"`
	IsMultiCurrency *bool            `json:"isMultiCurrency,omitempty"`
	ProfitRate      *PercentageText  `json:"profitRate,omitempty"`
}

// MaximumContractSize is a documented CoinW response record.
type MaximumContractSize struct {
	MaxBuy  WireDecimal      `json:"maxBuy" required:"true"`
	MaxSell WireDecimal      `json:"maxSell" required:"true"`
	TS      *StringOrInteger `json:"ts,omitempty"`
}

// StopLossTakeProfit is a documented CoinW response record.
type StopLossTakeProfit struct {
	CurrentPiece         *WireDecimal     `json:"currentPiece,omitempty"`
	ClosePiece           *WireDecimal     `json:"closePiece,omitempty"`
	Direction            *string          `json:"direction,omitempty"`
	CreatedDate          *StringOrInteger `json:"createdDate,omitempty"`
	UpdatedDate          *StringOrInteger `json:"updatedDate,omitempty"`
	ID                   StringOrInteger  `json:"id" required:"true"`
	IndexPrice           *WireDecimal     `json:"indexPrice,omitempty"`
	Instrument           string           `json:"instrument" required:"true"`
	Leverage             *WireDecimal     `json:"leverage,omitempty"`
	OpenID               *StringOrInteger `json:"openId,omitempty"`
	PriceType            *StringOrInteger `json:"priceType,omitempty"`
	StopLossPrice        *WireDecimal     `json:"stopLossPrice,omitempty"`
	StopLossRate         *WireDecimal     `json:"stopLossRate,omitempty"`
	StopProfitPrice      *WireDecimal     `json:"stopProfitPrice,omitempty"`
	StopProfitRate       *WireDecimal     `json:"stopProfitRate,omitempty"`
	StopLossOrderPrice   *WireDecimal     `json:"stopLossOrderPrice,omitempty"`
	StopProfitOrderPrice *WireDecimal     `json:"stopProfitOrderPrice,omitempty"`
	StopProfitType       *StringOrInteger `json:"stopProfitType,omitempty"`
	StopLossType         *StringOrInteger `json:"stopLossType,omitempty"`
	TriggerStatus        *StringOrInteger `json:"triggerStatus,omitempty"`
	UserID               *StringOrInteger `json:"userId,omitempty"`
	PositionModel        *StringOrInteger `json:"positionModel,omitempty"`
	StopType             *StringOrInteger `json:"stopType,omitempty"`
	OrderID              *StringOrInteger `json:"orderId,omitempty"`
}

// TrailingStopLossTakeProfit is a documented CoinW response record.
type TrailingStopLossTakeProfit struct {
	BaseSize       *WireDecimal     `json:"baseSize,omitempty"`
	CallbackRate   WireDecimal      `json:"callbackRate" required:"true"`
	CreatedDate    *StringOrInteger `json:"createdDate,omitempty"`
	CurrentPiece   *WireDecimal     `json:"currentPiece,omitempty"`
	Direction      *string          `json:"direction,omitempty"`
	Fee            *WireDecimal     `json:"fee,omitempty"`
	Instrument     string           `json:"instrument" required:"true"`
	Leverage       *WireDecimal     `json:"leverage,omitempty"`
	Margin         *WireDecimal     `json:"margin,omitempty"`
	OpenID         *StringOrInteger `json:"openId,omitempty"`
	OpenPrice      *WireDecimal     `json:"openPrice,omitempty"`
	OrderPrice     *WireDecimal     `json:"orderPrice,omitempty"`
	OriginalType   *string          `json:"originalType,omitempty"`
	PosType        *string          `json:"posType,omitempty"`
	PositionModel  *StringOrInteger `json:"positionModel,omitempty"`
	PositionMargin *WireDecimal     `json:"positionMargin,omitempty"`
	Quantity       *WireDecimal     `json:"quantity,omitempty"`
	QuantityUnit   *StringOrInteger `json:"quantityUnit,omitempty"`
	Status         *StringOrInteger `json:"status,omitempty"`
	TriggerPrice   *WireDecimal     `json:"triggerPrice,omitempty"`
	TriggerType    *StringOrInteger `json:"triggerType,omitempty"`
	UpdatedDate    *StringOrInteger `json:"updatedDate,omitempty"`
	ContractType   *StringOrInteger `json:"contractType,omitempty"`
	FinishStatus   *int             `json:"finishStatus,omitempty"`
	FundingFee     *WireDecimal     `json:"fundingFee,omitempty"`
	ID             StringOrInteger  `json:"id" required:"true"`
	UserID         *StringOrInteger `json:"userId,omitempty"`
	ProcessStatus  *int             `json:"processStatus,omitempty"`
	Source         *string          `json:"source,omitempty"`
	Version        *StringOrInteger `json:"version,omitempty"`
}

// BatchOrderAcknowledgement is a documented CoinW response record.
type BatchOrderAcknowledgement struct {
	MsgCode      StringOrInteger  `json:"msgCode" required:"true"`
	OpenID       *StringOrInteger `json:"openId,omitempty"`
	ThirdOrderID *StringOrInteger `json:"thirdOrderId,omitempty"`
	TS           *StringOrInteger `json:"ts,omitempty"`
}

// ModifiedOrderAcknowledgement is a documented CoinW response record.
type ModifiedOrderAcknowledgement struct {
	EditID   StringOrInteger  `json:"editId" required:"true"`
	OriginID StringOrInteger  `json:"originId" required:"true"`
	TS       *StringOrInteger `json:"ts,omitempty"`
}

// SpotAssetEvent is a documented CoinW response record.
type SpotAssetEvent struct {
	Available WireDecimal      `json:"available" required:"true"`
	Currency  string           `json:"currency" required:"true"`
	Time      StringOrInteger  `json:"time" required:"true"`
	Type      *string          `json:"type,omitempty"`
	LedgerID  *StringOrInteger `json:"ledger_id,omitempty"`
	Account   *string          `json:"account,omitempty"`
	Hold      *WireDecimal     `json:"hold,omitempty"`
}

// SpotOrderEvent is a documented CoinW response record.
type SpotOrderEvent struct {
	Fee           *WireDecimal     `json:"fee,omitempty"`
	DealFunds     *WireDecimal     `json:"dealFunds,omitempty"`
	Type          string           `json:"type" required:"true"`
	Time          StringOrInteger  `json:"time" required:"true"`
	ProductID     string           `json:"product_id" required:"true"`
	OrderID       StringOrInteger  `json:"order_id" required:"true"`
	ClientID      *StringOrInteger `json:"client_id,omitempty"`
	Size          *WireDecimal     `json:"size,omitempty"`
	RemainingSize *WireDecimal     `json:"remaining_size,omitempty"`
	Price         *WireDecimal     `json:"price,omitempty"`
	Side          *string          `json:"side,omitempty"`
	OrderType     *string          `json:"order_type,omitempty"`
	Reason        *string          `json:"reason,omitempty"`
	DealAvgPrice  *WireDecimal     `json:"dealAvgPrice,omitempty"`
}

// SpotTickerEvent is a documented CoinW response record.
type SpotTickerEvent struct {
	Buy         *WireDecimal    `json:"buy,omitempty"`
	ChangePrice *WireDecimal    `json:"changePrice,omitempty"`
	ChangeRate  *WireDecimal    `json:"changeRate,omitempty"`
	High        *WireDecimal    `json:"high,omitempty"`
	Last        WireDecimal     `json:"last" required:"true"`
	Low         *WireDecimal    `json:"low,omitempty"`
	Open        *WireDecimal    `json:"open,omitempty"`
	Sell        *WireDecimal    `json:"sell,omitempty"`
	Symbol      StringOrInteger `json:"symbol" required:"true"`
	Vol         *WireDecimal    `json:"vol,omitempty"`
	VolValue    *WireDecimal    `json:"volValue,omitempty"`
}

// SpotAllTickerEvent is a documented CoinW response record.
type SpotAllTickerEvent struct {
	ActivityState  *int            `json:"activityState,omitempty"`
	CurrencyVol    *WireDecimal    `json:"currencyVol,omitempty"`
	FPartitionIds  *string         `json:"fPartitionIds,omitempty"`
	Fav            *bool           `json:"fav,omitempty"`
	HotCoinSort    *int            `json:"hotCoinSort,omitempty"`
	LeftCoinName   *string         `json:"leftCoinName,omitempty"`
	LeftCoinUrl    *string         `json:"leftCoinUrl,omitempty"`
	NewCoinSort    *int            `json:"newCoinSort,omitempty"`
	OneDayHighest  *WireDecimal    `json:"oneDayHighest,omitempty"`
	OneDayLowest   *WireDecimal    `json:"oneDayLowest,omitempty"`
	OneDayTotal    *WireDecimal    `json:"oneDayTotal,omitempty"`
	Price          WireDecimal     `json:"price" required:"true"`
	RightCoinName  *string         `json:"rightCoinName,omitempty"`
	Rose           *WireDecimal    `json:"rose,omitempty"`
	Selective      *bool           `json:"selective,omitempty"`
	Symbol         *string         `json:"symbol,omitempty"`
	TmID           StringOrInteger `json:"tmId" required:"true"`
	TransferPrice  *WireDecimal    `json:"transferPrice,omitempty"`
	TransferSymbol *string         `json:"transferSymbol,omitempty"`
}

// SpotDepthSnapshot is a documented CoinW response record.
type SpotDepthSnapshot struct {
	Asks []PriceQuantity `json:"asks" required:"true"`
	Bids []PriceQuantity `json:"bids" required:"true"`
	Time StringOrInteger `json:"time" required:"true"`
	Seq  StringOrInteger `json:"seq" required:"true"`
}

// SpotDepthIncrement is a documented CoinW response record.
type SpotDepthIncrement struct {
	StartSeq StringOrInteger          `json:"startSeq" required:"true"`
	EndSeq   StringOrInteger          `json:"endSeq" required:"true"`
	Asks     []SequencedPriceQuantity `json:"asks" required:"true"`
	Bids     []SequencedPriceQuantity `json:"bids" required:"true"`
}

// SpotStreamTrade is a documented CoinW response record.
type SpotStreamTrade struct {
	Price   WireDecimal      `json:"price" required:"true"`
	Side    *string          `json:"side,omitempty"`
	Size    *WireDecimal     `json:"size,omitempty"`
	Time    *StringOrInteger `json:"time,omitempty"`
	TradeID *StringOrInteger `json:"tradeId,omitempty"`
	Seq     *StringOrInteger `json:"seq,omitempty"`
	Symbol  *StringOrInteger `json:"symbol,omitempty"`
}

// FuturesTickerEvent is a documented CoinW response record.
type FuturesTickerEvent struct {
	High         *WireDecimal `json:"high,omitempty"`
	Vol          *WireDecimal `json:"vol,omitempty"`
	VolUSDT      *WireDecimal `json:"volUsdt,omitempty"`
	Last         WireDecimal  `json:"last" required:"true"`
	Low          *WireDecimal `json:"low,omitempty"`
	ChangeRate   *WireDecimal `json:"changeRate,omitempty"`
	CurrencyCode string       `json:"currencyCode" required:"true"`
	Open         *WireDecimal `json:"open,omitempty"`
}

// FuturesDepthEvent is a documented CoinW response record.
type FuturesDepthEvent struct {
	Ask  []FuturesDepthLevel `json:"ask,omitempty"`
	Asks []FuturesDepthLevel `json:"asks,omitempty"`
	Bids []FuturesDepthLevel `json:"bids" required:"true"`
	N    *string             `json:"n,omitempty"`
	TS   *StringOrInteger    `json:"ts,omitempty"`
}

// FuturesPriceEvent is a documented CoinW response record.
type FuturesPriceEvent struct {
	P  WireDecimal     `json:"p" required:"true"`
	N  string          `json:"n" required:"true"`
	TS StringOrInteger `json:"ts" required:"true"`
}

// FuturesFundingEvent is a documented CoinW response record.
type FuturesFundingEvent struct {
	R  WireDecimal     `json:"r" required:"true"`
	Nt StringOrInteger `json:"nt" required:"true"`
	TS StringOrInteger `json:"ts" required:"true"`
	N  string          `json:"n" required:"true"`
}

// FuturesAssetEvent is a documented CoinW response record.
type FuturesAssetEvent struct {
	Margin            *WireDecimal    `json:"margin,omitempty"`
	ProfitUnreal      *WireDecimal    `json:"profitUnreal,omitempty"`
	Freeze            *WireDecimal    `json:"freeze,omitempty"`
	Size              *WireDecimal    `json:"size,omitempty"`
	Available         WireDecimal     `json:"available" required:"true"`
	Currency          string          `json:"currency" required:"true"`
	Type              *int            `json:"type,omitempty"`
	AlmightyGold      *WireDecimal    `json:"almightyGold,omitempty"`
	TransferAvailable *WireDecimal    `json:"transferAvailable,omitempty"`
	UserID            StringOrInteger `json:"userId" required:"true"`
	AvailableMargin   *WireDecimal    `json:"availableMargin,omitempty"`
	Hold              *WireDecimal    `json:"hold,omitempty"`
}
