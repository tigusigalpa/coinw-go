package services

import (
	"fmt"
	"math/big"
	"strings"

	coinwerr "github.com/tigusigalpa/coinw-go/errors"
)

func validateContractRequest(id string, value any) error {
	p, ok := value.(map[string]interface{})
	if !ok {
		return nil
	}
	contract, exists := requestContracts[id]
	if !exists {
		return nil
	}
	for key := range p {
		found := false
		for _, allowed := range contract.allowed {
			if key == allowed {
				found = true
				break
			}
		}
		if !found {
			return validationErr("undocumented parameter %s for %s", key, id)
		}
	}
	for _, key := range contract.required {
		v, ok := p[key]
		if !ok || v == nil {
			return validationErr("parameter %s is required", key)
		}
		if s, ok := v.(string); ok && strings.TrimSpace(s) == "" && key != "memo" {
			return validationErr("parameter %s must not be empty", key)
		}
	}
	// These sets are intentionally operation-specific: query-order categories
	// are different from the placement order types on CoinW's trade pages.
	if v, ok := p["originType"].(string); ok && !validOriginType(v) {
		return validationErr("originType must be plan, planTrigger or execute")
	}
	if id == "futures.check.getOrderInformation" || id == "futures.trade.cancelBatchOrders" {
		key := "positionType"
		if id == "futures.trade.cancelBatchOrders" {
			key = "posType"
		}
		if v, ok := p[key].(string); ok {
			switch v {
			case "execute", "plan", "moveStop", "stopProfitLoss":
			default:
				return validationErr("unsupported %s for this operation", key)
			}
		}
	}
	if id == "futures.check.getCurrentOrders" {
		switch p["positionType"] {
		case "plan", "execute", "planTrigger", "moveStopProfitLoss":
		default:
			return validationErr("unsupported current-order positionType")
		}
	}
	if id == "futures.assets.enableDisableMegaCoupon" {
		switch fmt.Sprint(p["status"]) {
		case "0", "1":
		default:
			return validationErr("coupon status must be 0 or 1")
		}
	}
	if id == "common.subAccount.accountTransfer" {
		switch p["side"] {
		case "IN", "OUT":
		default:
			return validationErr("transfer side must be IN or OUT")
		}
	}
	if id == "common.account.transferAssets" {
		for _, key := range []string{"accountType", "targetAccountType"} {
			if p[key] != "SPOT" && p[key] != "WEALTH" {
				return validationErr("%s must be SPOT or WEALTH", key)
			}
		}
		if p["bizType"] != "SPOT_TO_WEALTH" && p["bizType"] != "WEALTH_TO_SPOT" {
			return validationErr("unsupported transfer bizType")
		}
	}
	if id == "spot.trade.cancelAllOrders" {
		pair, ok := p["currencyPair"].(string)
		if !ok || strings.TrimSpace(pair) == "" {
			return &coinwerr.AmbiguousDocumentationError{EndpointID: id, Note: "unfiltered cancellation scope is not defined consistently; specify currencyPair"}
		}
	}
	if v, ok := p["symbols"].(string); ok {
		if _, err := requireMaxInstruments(strings.Split(v, ","), 20); err != nil {
			return err
		}
	}
	for _, key := range []string{"openIds", "sourceIds"} {
		if v, ok := p[key].(string); ok {
			if _, err := requireMaxInstruments(strings.Split(v, ","), 20); err != nil {
				return err
			}
		}
	}
	if strings.HasPrefix(id, "spot.market.") || id == "common.account.getDepositWithdrawalHistory" || id == "spot.check.getTransactionHistory" {
		for _, key := range []string{"symbol", "currencyPair"} {
			if v, ok := p[key].(string); ok {
				if _, err := requireMaxInstruments(strings.Split(v, ","), 20); err != nil {
					return err
				}
			}
		}
	}
	if id == "spot.check.getBatchHistoricalOrders" {
		if a, ok := p["orderIds"].([]interface{}); ok && len(a) > 50 {
			return validationErr("orderIds supports at most 50 entries")
		}
	}
	if strings.HasPrefix(id, "common.affiliate.") {
		for _, key := range []string{"startTime", "endTime"} {
			if v, ok := p[key].(string); ok {
				var err error
				if id == "common.affiliate.getInvitedUsersTradingVolume" {
					err = requireAffiliateDate(v, key)
				} else {
					err = requireAffiliateDateTime(v, key)
				}
				if err != nil {
					return err
				}
			}
		}
		if start, ok := p["startTime"].(string); ok {
			if end, ok := p["endTime"].(string); ok && start > end {
				return validationErr("startTime must not exceed endTime")
			}
		}
		if v, ok := p["pageSize"]; ok {
			n, ok := integerParam(v)
			max := 1000
			if id == "common.affiliate.getInvitedUserAssets" {
				max = 100
			}
			if !ok || n < 1 || n > max {
				return validationErr("pageSize exceeds documented bounds")
			}
		}
		if a, ok := p["partnerUids"].([]interface{}); ok && len(a) > 100 {
			return validationErr("partnerUids supports at most 100 entries")
		}
	}
	if id == "futures.market.getHistoricalPublicTrades" || id == "futures.check.getBatchHistoricalOrders7Days" || id == "futures.check.getBatchHistoricalOrders3Months" {
		if v, ok := p["pageSize"]; ok {
			n, ok := integerParam(v)
			if !ok || n < 1 || n > 500 {
				return validationErr("pageSize must be 1..500")
			}
		}
	}
	if id == "futures.assets.coinConversion" {
		n, ok := integerParam(p["convertType"])
		if !ok || (n != 1 && n != 2) {
			return validationErr("convertType must be 1 or 2")
		}
		key := "dealPiece"
		if n == 2 {
			key = "baseSize"
		}
		if p[key] == nil {
			return validationErr("%s is mandatory for this convertType", key)
		}
	}
	if id == "spot.market.getKLine" {
		if v, ok := p["currencyPair"].(string); ok {
			if _, err := requireMaxInstruments(strings.Split(v, ","), 20); err != nil {
				return err
			}
		}
		switch fmt.Sprint(p["period"]) {
		case "60", "180", "300", "900", "1800", "7200", "14400":
		default:
			return validationErr("unsupported Spot K-line period")
		}
	}
	if id == "futures.market.getKLineOfAnInstrument" {
		switch fmt.Sprint(p["granuality"]) {
		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
		default:
			return validationErr("granuality must be 0..9")
		}
		if v, ok := p["limit"]; ok {
			n, ok := integerParam(v)
			if !ok || n < 1 || n > 1500 {
				return validationErr("limit must be 1..1500")
			}
		}
		if v, ok := p["klineType"]; ok && v != "0" && v != "1" {
			return validationErr("klineType must be 0 or 1")
		}
	}
	if v, ok := p["positionModel"]; ok {
		n, ok := integerParam(v)
		if !ok || !validPositionModel(n) {
			return validationErr("invalid positionModel")
		}
	}
	if v, ok := p["quantityUnit"]; ok {
		n, ok := integerParam(v)
		if !ok || !validQuantityUnit(n) {
			return validationErr("invalid quantityUnit")
		}
	}
	if v, ok := p["stopFrom"]; ok {
		n, ok := integerParam(v)
		if !ok || n < 1 || n > 3 {
			return validationErr("stopFrom must be 1..3")
		}
	}
	if id == "futures.assets.setMarginMode" {
		n, ok := integerParam(p["layout"])
		if !ok || n < 0 || n > 1 {
			return validationErr("layout must be 0 or 1")
		}
	}
	if id == "futures.check.getStopLossTakeProfitInformation" {
		n, _ := integerParam(p["stopFrom"])
		key := map[int]string{1: "orderId", 2: "openId", 3: "planOrderId"}[n]
		if p[key] == nil {
			return validationErr("%s is mandatory for stopFrom", key)
		}
	}
	if id == "spot.trade.placeOrder" {
		if v, ok := p["out_trade_no"].(string); ok && !clientOrderIDPattern.MatchString(v) {
			return validationErr("invalid out_trade_no")
		}
		market, ok := p["isMarket"].(bool)
		if !ok {
			return validationErr("isMarket must be a boolean")
		}
		if !market && (p["amount"] == nil || p["rate"] == nil) {
			return validationErr("amount and rate are required for a limit order")
		}
	}
	if id == "futures.trade.cancelBatchOrders" {
		if a, ok := p["sourceIds"].([]string); ok {
			if _, err := requireMaxInstruments(a, 20); err != nil {
				return err
			}
		}
	}
	if strings.HasPrefix(id, "futures.trade.") {
		if id == "futures.trade.closeAPosition" {
			positionType, _ := p["positionType"].(string)
			if positionType != "" && positionType != "plan" && positionType != "execute" {
				return validationErr("closing positionType must be plan or execute")
			}
			if p["orderPrice"] != nil && positionType != "plan" {
				return validationErr("orderPrice requires positionType=plan")
			}
		}
		for _, key := range []string{"closeRate", "callbackRate"} {
			if v, ok := p[key]; ok {
				n, ok := new(big.Rat).SetString(fmt.Sprint(v))
				if !ok || n.Sign() < 0 || n.Cmp(big.NewRat(1, 1)) > 0 {
					return validationErr("%s must be between 0 and 1", key)
				}
			}
		}
		for _, key := range []string{"addMargin", "reduceMargin"} {
			if v, ok := p[key]; ok {
				n, ok := new(big.Rat).SetString(fmt.Sprint(v))
				if !ok || n.Sign() <= 0 {
					return validationErr("%s must be greater than zero", key)
				}
			}
		}
	}
	return nil
}
func integerParam(v any) (int, bool) {
	switch n := v.(type) {
	case int:
		return n, true
	case int64:
		return int(n), true
	case string:
		if n == "0" {
			return 0, true
		}
		if n == "1" {
			return 1, true
		}
	}
	return 0, false
}
