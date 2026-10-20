// Package endpoints is the registry of every documented CoinW REST
// operation (catalogue snapshot 2026-10-20,
// https://www.coinw.com/api-doc/en/sitemap.xml).
//
// Auth profile legend: AuthPublic = no credentials; AuthSpot = Spot MD5
// signer; AuthFutures = Futures HMAC-SHA256 signer. Never mixed.
//
// Ambiguous operations are registered but gated: the executor returns
// errors.AmbiguousDocumentationError instead of sending an unverified route.
package endpoints

// Auth profiles.
const (
	AuthPublic  = "public"  // no credentials
	AuthSpot    = "spot"    // Spot MD5 scheme
	AuthFutures = "futures" // Futures HMAC scheme
)

// Rate-limit global groups per Common > Precautions.
const (
	GroupSpot      = "spot"      // 100 rps/IP + 300 rps/UID
	GroupFuturesA  = "futures_a" // 30 rps/IP (perpumPublic market data)
	GroupFuturesB  = "futures_b" // 100 rps/UID
	GroupAffiliate = "affiliate" // 5 rps/IP+UID per docs pages
)

// Definition statically describes one documented REST operation.
type Definition struct {
	ID              string  // stable operation id (coverage mapping key)
	Method          string  // HTTP verb
	Path            string  // documented endpoint path (spot paths may embed ?command=)
	Auth            string  // one of Auth*
	RateGroup       string  // one of Group*
	PerInterfaceRPS float64 // documented per-interface limit; <=0 = unspecified
	PerIPRPS        float64 // explicit per-IP interface rate; overrides the legacy shared rate
	PerUIDRPS       float64 // explicit per-UID interface rate; overrides the legacy shared rate
	ReadOnly        bool    // idempotent — eligible for auto-retry
	Ambiguous       bool    // gated by AmbiguousDocumentationError until verified
	AmbiguityNote   string  // description of the documented conflict
}

// IsMutation reports whether the endpoint mutates state (never retried).
func (d Definition) IsMutation() bool { return !d.ReadOnly }

// registry holds all definitions, built once at init.
var registry = build()

// Get returns the definition for id; ok=false when unregistered.
func Get(id string) (Definition, bool) {
	d, ok := registry[id]
	return d, ok
}

// All returns a copy of the full registry.
func All() map[string]Definition {
	out := make(map[string]Definition, len(registry))
	for k, v := range registry {
		out[k] = v
	}
	return out
}

func def(id, method, path, auth, group string, rps float64, readOnly bool) Definition {
	return Definition{ID: id, Method: method, Path: path, Auth: auth,
		RateGroup: group, PerInterfaceRPS: rps, ReadOnly: readOnly}
}

func ambiguous(d Definition, note string) Definition {
	d.Ambiguous = true
	d.AmbiguityNote = note
	return d
}

func build() map[string]Definition {
	d := map[string]Definition{}
	put := func(e Definition) { d[e.ID] = e }

	// ---- Common > Accounts (Spot private profile) ------------------------
	put(def("common.account.cancelWithdrawal", "POST", "/api/v1/private?command=cancelWithdraw",
		AuthSpot, GroupSpot, 0, false))
	put(def("common.account.getDepositWithdrawalAddress", "POST", "/api/v1/private?command=returnDepositAddresses",
		AuthSpot, GroupSpot, 0, true))
	put(def("common.account.getDepositWithdrawalHistory", "POST", "/api/v1/private?command=returnDepositsWithdrawals",
		AuthSpot, GroupSpot, 0, true))
	put(def("common.account.initiateWithdrawal", "POST", "/api/v1/private?command=doWithdraw",
		AuthSpot, GroupSpot, 0, false))
	put(def("common.account.transferAssets", "POST", "/api/v1/private?command=spotWealthTransfer",
		AuthSpot, GroupSpot, 0, false))

	// ---- Common > Sub-accounts (Spot private profile, master key only) ---
	put(def("common.subAccount.getSubAccountList", "POST", "/api/v1/private?command=subAccount_list",
		AuthSpot, GroupSpot, 3.0, true))
	put(def("common.subAccount.getSubAccountApiKeys", "POST", "/api/v1/private?command=subAccount_apikey",
		AuthSpot, GroupSpot, 3.0, true))
	put(def("common.subAccount.resetSubAccountApiKey", "POST", "/api/v1/private?command=reset_subAccount_apikey",
		AuthSpot, GroupSpot, 3.0, false))
	put(ambiguous(def("common.subAccount.getMaxTransferableAmount", "POST", "/api/v1/private?command=subAccount_transferAvailable",
		AuthSpot, GroupSpot, 3.0, true),
		"Endpoint body shows /api/v1/private but its hyperlink resolves to "+
			"/api/v1/api/v1/private (doubled prefix). Unverified route."))
	put(def("common.subAccount.getTransferHistory", "POST", "/api/v1/private?command=returnTransferHistory",
		AuthSpot, GroupSpot, 3.0, true))
	put(def("common.subAccount.accountTransfer", "POST", "/api/v1/private?command=subAccount_transfer",
		AuthSpot, GroupSpot, 3.0, false))

	// ---- Common > Affiliate (Futures private profile; literal "patener"
	// in the documented path — it is part of the real route) --------------
	affiliate := func(id, method, path string, readOnly bool) Definition {
		return def(id, method, path, AuthFutures, GroupAffiliate, 5.0, readOnly)
	}
	put(affiliate("common.affiliate.getInvitedUserList", "GET",
		"/api/v3/patener/openapi/v1/partner/invitation/user/list", true))
	put(affiliate("common.affiliate.getInvitedUsersKycResults", "GET",
		"/api/v3/patener/openapi/v1/partner/invitation/user/kyc/list", true))
	put(affiliate("common.affiliate.getInvitedUsersTradingVolume", "POST",
		"/api/v3/patener/openapi/v1/partner/invitation/user/volume/list", false))
	put(affiliate("common.affiliate.getAffiliateUserActivityData", "POST",
		"/api/v3/patener/openapi/v1/partner/invitation/user/activity/list", false))
	put(affiliate("common.affiliate.getInvitedUserAssets", "GET",
		"/api/v3/patener/openapi/v1/partner/invitation/user/assets", true))
	put(affiliate("common.affiliate.getAffiliateTeamDirectUserRange", "POST",
		"/api/v3/patener/openapi/v1/partner/team/users/range", false))

	// ---- Spot > Market (public command API on /api/v1/public) ------------
	spotPublic := func(id, command string) Definition {
		return def(id, "GET", "/api/v1/public?command="+command, AuthPublic, GroupSpot, 0, true)
	}
	put(spotPublic("spot.market.getInstrumentInformation", "returnSymbol"))
	put(spotPublic("spot.market.get24hTradeSummaryForAllInstruments", "returnTicker"))
	put(spotPublic("spot.market.getCoinDepositWithdrawalLimits", "returnCurrencies"))
	put(spotPublic("spot.market.getOrderBook", "returnOrderBook"))
	put(spotPublic("spot.market.getRecentTrades", "returnTradeHistory"))
	put(spotPublic("spot.market.getKLine", "returnChartData"))
	put(spotPublic("spot.market.get24hVolumeForPopularInstruments", "return24hVolume"))

	// ---- Spot > Account / Check / Trade (Spot private profile) -----------
	put(def("spot.account.getSpotAccountBalance", "POST", "/api/v1/private?command=returnBalances",
		AuthSpot, GroupSpot, 0, true))
	put(def("spot.account.getFullSpotAccountBalance", "POST", "/api/v1/private?command=returnCompleteBalances",
		AuthSpot, GroupSpot, 0, true))

	put(ambiguous(def("spot.check.getBatchHistoricalOrders", "POST", "/v1/private?command=getBatchHistoryOrders",
		AuthSpot, GroupSpot, 0.5, true),
		"Documented endpoint /v1/private omits the /api prefix used by every "+
			"other spot private route and conflicts on host/path. Unverified."))
	put(def("spot.check.getCurrentOrders", "POST", "/api/v1/private?command=returnOpenOrders",
		AuthSpot, GroupSpot, 0, true))
	put(def("spot.check.getHistoricalOrders", "POST", "/api/v1/private?command=getUserTrades",
		AuthSpot, GroupSpot, 0, true))
	put(def("spot.check.getOrderDetails", "POST", "/api/v1/private?command=returnOrderTrades",
		AuthSpot, GroupSpot, 0, true))
	put(def("spot.check.getOrderStatus", "POST", "/api/v1/private?command=returnOrderStatus",
		AuthSpot, GroupSpot, 0, true))
	put(def("spot.check.getTransactionHistory", "POST", "/api/v1/private?command=returnUTradeHistory",
		AuthSpot, GroupSpot, 0, true))

	put(def("spot.trade.placeOrder", "POST", "/api/v1/private?command=doTrade",
		AuthSpot, GroupSpot, 0, false))
	put(def("spot.trade.cancelOrder", "POST", "/api/v1/private?command=cancelOrder",
		AuthSpot, GroupSpot, 0, false))
	put(def("spot.trade.cancelAllOrders", "POST", "/api/v1/private?command=cancelAllOrder",
		AuthSpot, GroupSpot, 0, false))

	// ---- Futures > Market -----------------------------------------------
	futPublicA := func(id, path string) Definition {
		return def(id, "GET", path, AuthPublic, GroupFuturesA, 0, true)
	}
	futPublicB := func(id, path string) Definition {
		return def(id, "GET", path, AuthPublic, GroupFuturesB, 0, true)
	}
	futPrivate := func(id, method, path string, readOnly bool, rps float64) Definition {
		return def(id, method, path, AuthFutures, GroupFuturesB, rps, readOnly)
	}

	put(futPublicB("futures.market.getInstrumentInfo", "/v1/perpum/instruments"))
	put(futPublicB("futures.market.getBatchInstrumentsInfo", "/v1/perpum/instrumentList"))
	put(futPublicA("futures.market.getLastTradeSummaryOfAnInstrument", "/v1/perpumPublic/ticker"))
	put(futPublicB("futures.market.getBatchLastTradeSummary", "/v1/perpumPublic/ticker/list"))
	put(futPublicA("futures.market.getLastTradeSummaryOfAllInstruments", "/v1/perpumPublic/tickers"))
	put(futPublicA("futures.market.getKLineOfAnInstrument", "/v1/perpumPublic/klines"))
	put(futPublicB("futures.market.getLastSettlementFundingFeeRate", "/v1/perpum/fundingRate"))
	put(futPublicA("futures.market.getOrderBookOfAnInstrument", "/v1/perpumPublic/depth"))
	put(futPublicA("futures.market.getTradeDataOfAnInstrument", "/v1/perpumPublic/trades"))
	// Documented as public-market endpoints but requiring Futures private auth:
	put(futPrivate("futures.market.getMarginRequirementsOfAllInstruments", "GET", "/v1/perpum/ladders", true, 0))
	put(futPrivate("futures.market.getHistoricalPublicTrades", "GET", "/v1/perpum/orders/trades", true, 5.0))

	// ---- Futures > Account & Assets + Common -----------------------------
	put(futPrivate("futures.assets.coinConversion", "POST", "/v1/perpum/pieceConvert", false, 0))
	put(futPrivate("futures.assets.enableDisableMegaCoupon", "POST", "/v1/perpum/account/almightyGoldInfo", false, 0))
	put(futPrivate("futures.assets.getFuturesAccountAssets", "GET", "/v1/perpum/account/getUserAssets", true, 0))
	put(futPrivate("futures.assets.getFuturesAccountFees", "GET", "/v1/perpum/account/fees", true, 0))
	put(futPrivate("futures.assets.getMarginMode", "GET", "/v1/perpum/positions/type", true, 0))
	put(futPrivate("futures.assets.getMaxTransferableBalance", "GET", "/v1/perpum/account/available", true, 0))
	put(futPrivate("futures.assets.getMaxUsersContractSize", "GET", "/v1/perpum/orders/availSize", true, 0))
	put(futPrivate("futures.assets.getMegaCouponBalance", "GET", "/v1/perpum/account/almightyGoldInfo", true, 0))
	put(futPrivate("futures.assets.getTransactionDetails3Days", "GET", "/v1/perpum/orders/deals", true, 0))
	put(futPrivate("futures.assets.getTransactionDetails3Months", "GET", "/v1/perpum/orders/deals/history", true, 0))
	put(futPrivate("futures.assets.setMarginMode", "POST", "/v1/perpum/positions/type", false, 0))
	put(futPrivate("futures.common.getLeverageInformation", "GET", "/v1/perpum/positions/leverage", true, 0))

	// ---- Futures > Check ---------------------------------------------------
	put(futPrivate("futures.check.getBatchHistoricalOrders3Months", "GET", "/v1/perpum/orders/batchInsArchive", true, 1.0))
	put(futPrivate("futures.check.getBatchHistoricalOrders7Days", "GET", "/v1/perpum/orders/batchInsHistory", true, 1.0))
	put(futPrivate("futures.check.getCurrentOrders", "GET", "/v1/perpum/orders/open", true, 0))
	put(futPrivate("futures.check.getHistoricalOrders3Months", "GET", "/v1/perpum/orders/archive", true, 0))
	put(futPrivate("futures.check.getHistoricalOrders7Days", "GET", "/v1/perpum/orders/history", true, 0))
	put(futPrivate("futures.check.getOrderInformation", "GET", "/v1/perpum/order", true, 0))
	put(futPrivate("futures.check.getPendingOrderCount", "GET", "/v1/perpum/orders/openQuantity", true, 0))
	put(futPrivate("futures.check.getStopLossTakeProfitInformation", "GET", "/v1/perpum/TPSL", true, 0))
	put(futPrivate("futures.check.getTrailingStopLossTakeProfitInformation", "GET", "/v1/perpum/moveTPSL", true, 0))

	// ---- Futures > Position ------------------------------------------------
	put(futPrivate("futures.position.getCurrentPositionInformation", "GET", "/v1/perpum/positions", true, 0))
	put(futPrivate("futures.position.getAllCurrentPositions", "GET", "/v1/perpum/positions/all", true, 0))
	put(futPrivate("futures.position.getHistoricalPositionInformation", "GET", "/v1/perpum/positions/history", true, 0))
	put(futPrivate("futures.position.getMaxContractSize", "GET", "/v1/perpum/orders/maxSize", true, 0))
	put(futPrivate("futures.position.getPositionMarginRate", "GET", "/v1/perpum/positions/marginRate", true, 5.0))

	// ---- Futures > Trade (all mutations — never auto-retried) --------------
	put(futPrivate("futures.trade.adjustMargin", "POST", "/v1/perpum/positions/margin", false, 5.0))
	put(futPrivate("futures.trade.cancelAnOrder", "DELETE", "/v1/perpum/order", false, 0))
	put(futPrivate("futures.trade.cancelBatchOrders", "DELETE", "/v1/perpum/batchOrders", false, 0))
	put(futPrivate("futures.trade.closeAPosition", "DELETE", "/v1/perpum/positions", false, 0))
	put(futPrivate("futures.trade.closeBatchPositions", "DELETE", "/v1/perpum/batchClose", false, 0))
	put(futPrivate("futures.trade.closePositionsAtMarketPrice", "DELETE", "/v1/perpum/allpositions", false, 0))
	put(futPrivate("futures.trade.modifyAnOrder", "PUT", "/v1/perpum/order", false, 10.0))
	put(ambiguous(def("futures.trade.modifyStopLossTakeProfitInBatches", "POST", "/v1/perpum/updateTpsl",
		AuthFutures, GroupFuturesB, 5.0, false),
		"Endpoint field declares /v1/perpum/updateTpsl while the Futures "+
			"code-snippet routes this operation via /v1/perpum/TPSL."))
	put(futPrivate("futures.trade.placeAnOrder", "POST", "/v1/perpum/order", false, 0))
	put(futPrivate("futures.trade.placeBatchOrders", "POST", "/v1/perpum/batchOrders", false, 0))
	put(futPrivate("futures.trade.reverseAPosition", "POST", "/v1/perpum/positions/reverse", false, 0))
	put(futPrivate("futures.trade.setStopLossTakeProfit", "POST", "/v1/perpum/TPSL", false, 0))
	put(futPrivate("futures.trade.setStopLossTakeProfitInBatches", "POST", "/v1/perpum/addTpsl", false, 0))
	put(futPrivate("futures.trade.setTrailingStopLossTakeProfit", "POST", "/v1/perpum/moveTPSL", false, 0))

	applyDocumentedLimits(d)
	return d
}
