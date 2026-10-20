package endpoints

import (
	"strings"
	"testing"
)

var validAuth = map[string]bool{AuthPublic: true, AuthSpot: true, AuthFutures: true}
var validGroups = map[string]bool{
	GroupSpot: true, GroupFuturesA: true, GroupFuturesB: true, GroupAffiliate: true,
}
var validMethods = map[string]bool{"GET": true, "POST": true, "PUT": true, "DELETE": true}

func TestRegistryCoverage(t *testing.T) {
	all := All()
	if len(all) < 50 {
		t.Fatalf("registry should cover the full documented catalogue, got %d", len(all))
	}
	for id, d := range all {
		if d.ID != id {
			t.Errorf("key/id mismatch: %s vs %s", id, d.ID)
		}
		if !validAuth[d.Auth] {
			t.Errorf("%s: invalid auth %q", id, d.Auth)
		}
		if !validGroups[d.RateGroup] {
			t.Errorf("%s: invalid rate group %q", id, d.RateGroup)
		}
		if !validMethods[d.Method] {
			t.Errorf("%s: invalid method %q", id, d.Method)
		}
		if d.Path == "" {
			t.Errorf("%s: empty path", id)
		}
		if d.Ambiguous && d.AmbiguityNote == "" {
			t.Errorf("%s: ambiguous without a note", id)
		}
	}
}

func TestMutationsAreNotReadOnly(t *testing.T) {
	mutations := []string{
		"spot.trade.placeOrder", "spot.trade.cancelOrder", "spot.trade.cancelAllOrders",
		"common.account.initiateWithdrawal", "common.account.cancelWithdrawal",
		"common.account.transferAssets", "common.subAccount.accountTransfer",
		"common.subAccount.resetSubAccountApiKey",
		"futures.trade.placeAnOrder", "futures.trade.placeBatchOrders",
		"futures.trade.cancelAnOrder", "futures.trade.cancelBatchOrders",
		"futures.trade.closeAPosition", "futures.trade.closeBatchPositions",
		"futures.trade.closePositionsAtMarketPrice", "futures.trade.modifyAnOrder",
		"futures.trade.reverseAPosition", "futures.trade.adjustMargin",
		"futures.trade.setStopLossTakeProfit", "futures.trade.setStopLossTakeProfitInBatches",
		"futures.trade.setTrailingStopLossTakeProfit",
		"futures.assets.setMarginMode", "futures.assets.coinConversion",
		"futures.assets.enableDisableMegaCoupon",
	}
	for _, id := range mutations {
		d, ok := Get(id)
		if !ok {
			t.Fatalf("endpoint not registered: %s", id)
		}
		if d.ReadOnly {
			t.Errorf("%s must be marked non-readOnly (mutation)", id)
		}
	}
}

func TestAmbiguousEndpoints(t *testing.T) {
	for _, id := range []string{
		"futures.trade.modifyStopLossTakeProfitInBatches",
		"common.subAccount.getMaxTransferableAmount",
		"spot.check.getBatchHistoricalOrders",
	} {
		d, ok := Get(id)
		if !ok || !d.Ambiguous {
			t.Errorf("%s must be registered as ambiguous", id)
		}
	}
}

func TestHistoricalPublicTradesUsesFuturesAuth(t *testing.T) {
	d, _ := Get("futures.market.getHistoricalPublicTrades")
	if d.Auth != AuthFutures {
		t.Fatalf("endpoint must use Futures HMAC profile per doc page, got %s", d.Auth)
	}
}

func TestDocumentedPathShapes(t *testing.T) {
	for id, d := range All() {
		if strings.HasPrefix(id, "spot.") || strings.HasPrefix(id, "common.account") || strings.HasPrefix(id, "common.subAccount") {
			if !strings.Contains(d.Path, "?command=") {
				t.Errorf("%s: spot-family endpoint must embed ?command=…", id)
			}
		}
		if strings.HasPrefix(id, "common.affiliate") &&
			!strings.HasPrefix(d.Path, "/api/v3/patener/openapi/v1/partner/") {
			t.Errorf("%s: affiliate path must use the documented patener base", id)
		}
	}
}

func TestDocumentedRateScopes(t *testing.T) {
	// Independent checks for the restrictive and asymmetric limits most easily
	// lost when converting the official Frequency Limit tables into metadata.
	for _, tc := range []struct {
		id      string
		ip, uid float64
		group   string
	}{
		{"futures.market.getBatchLastTradeSummary", 0.2, 0.2, GroupFuturesB},
		{"futures.assets.getMegaCouponBalance", 0.5, 0.5, GroupFuturesB},
		{"futures.position.getAllCurrentPositions", 50000, 30, GroupFuturesB},
		{"futures.position.getHistoricalPositionInformation", 100000, 15, GroupFuturesB},
		{"spot.account.getFullSpotAccountBalance", 0, 2.5, GroupSpot},
		{"spot.market.getRecentTrades", 10, 0, GroupSpot},
	} {
		ep, ok := Get(tc.id)
		if !ok || ep.PerIPRPS != tc.ip || ep.PerUIDRPS != tc.uid || ep.RateGroup != tc.group {
			t.Errorf("%s: incorrect limits: %+v", tc.id, ep)
		}
	}
}
