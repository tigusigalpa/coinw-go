package services

import (
	"context"

	coinwerr "github.com/tigusigalpa/coinw-go/errors"
	"github.com/tigusigalpa/coinw-go/transport"
)

// Futures direction / unit / model enums (documented).
const (
	DirectionLong  = "long"
	DirectionShort = "short"

	QuantityUnitQuote     = 0 // USDT / USDC
	QuantityUnitContracts = 1
	QuantityUnitBase      = 2

	PositionModelIsolated = 0
	PositionModelCross    = 1

	PositionTypePlan        = "plan"
	PositionTypePlanTrigger = "planTrigger"
	PositionTypeExecute     = "execute"
)

func validDirection(d string) bool  { return d == DirectionLong || d == DirectionShort }
func validQuantityUnit(u int) bool  { return u >= 0 && u <= 2 }
func validPositionModel(m int) bool { return m == 0 || m == 1 }

func validPositionType(t string) bool {
	return t == PositionTypePlan || t == PositionTypePlanTrigger || t == PositionTypeExecute
}

// validateTPSLFields enforces the documented SL/TP invariants:
//   - stopType ∈ {1 batch/partial, 2 entire}
//   - stopFrom ∈ {1 limit, 2 position, 3 plan}
//   - priceType ∈ {1 index, 2 latest, 3 mark}
//   - closePiece mandatory when stopType=1
func validateTPSLFields(fields map[string]interface{}) error {
	if v, ok := fields["stopType"].(int); ok {
		if v != 1 && v != 2 {
			return validationErr("stopType must be 1 (batch/partial) or 2 (entire position)")
		}
		if v == 1 {
			if _, hasPiece := fields["closePiece"]; !hasPiece {
				return validationErr("closePiece is mandatory when stopType=1 (take profit/stop loss in batches)")
			}
		}
	}
	if v, ok := fields["stopFrom"].(int); ok && (v < 1 || v > 3) {
		return validationErr("stopFrom must be 1 (limit), 2 (position) or 3 (plan)")
	}
	if v, ok := fields["priceType"].(int); ok && (v < 1 || v > 3) {
		return validationErr("priceType must be 1 (index), 2 (latest) or 3 (mark)")
	}
	return nil
}

// FuturesTradeService — ALL MUTATING.
//
// Safety contract:
//   - never invoked implicitly (no background side effects)
//   - never retried automatically
//   - after a transport timeout, verify state via FuturesCheckService /
//     FuturesPositionService before resubmitting
//   - CoinW supports simultaneous long and short positions (hedging): an
//     opposite-direction order does NOT close a position. Use the dedicated
//     close APIs.
type FuturesTradeService struct{ base }

func NewFuturesTradeService(ex *transport.Executor) *FuturesTradeService {
	return &FuturesTradeService{newBase(ex)}
}

// AdjustMargin — POST /v1/perpum/positions/margin. MUTATION.
// Exactly one of addMargin/reduceMargin must be non-empty (decimal strings).
// typ e.g. "all" — "" omits.
func (s *FuturesTradeService) AdjustMargin(ctx context.Context, id string, addMargin, reduceMargin, typ string) (*transport.Result, error) {
	if (addMargin == "") == (reduceMargin == "") {
		return nil, validationErr("provide exactly one of addMargin or reduceMargin")
	}
	params := map[string]interface{}{"id": id}
	if addMargin != "" {
		if _, err := decimalArg(addMargin, "addMargin"); err != nil {
			return nil, err
		}
		params["addMargin"] = addMargin
	}
	if reduceMargin != "" {
		if _, err := decimalArg(reduceMargin, "reduceMargin"); err != nil {
			return nil, err
		}
		params["reduceMargin"] = reduceMargin
	}
	if typ != "" {
		params["type"] = typ
	}
	return s.call(ctx, "futures.trade.adjustMargin", params)
}

// CancelAnOrder — DELETE /v1/perpum/order. MUTATION.
func (s *FuturesTradeService) CancelAnOrder(ctx context.Context, id string) (*transport.Result, error) {
	return s.call(ctx, "futures.trade.cancelAnOrder", map[string]interface{}{"id": id})
}

// CancelBatchOrders — DELETE /v1/perpum/batchOrders. MUTATION.
// posType="" omits the documented filter.
func (s *FuturesTradeService) CancelBatchOrders(ctx context.Context, sourceIDs []string, posType string) (*transport.Result, error) {
	if len(sourceIDs) == 0 {
		return nil, validationErr("sourceIds must be a non-empty list")
	}
	params := map[string]interface{}{"sourceIds": sourceIDs}
	if posType != "" {
		params["posType"] = posType
	}
	return s.call(ctx, "futures.trade.cancelBatchOrders", params)
}

// CloseAPosition — DELETE /v1/perpum/positions. MUTATION.
// Exactly one of closeNum/closeRate per documented invariant; orderPrice=""
// omits. All decimal strings.
func (s *FuturesTradeService) CloseAPosition(ctx context.Context, id, positionType, closeNum, closeRate, orderPrice string) (*transport.Result, error) {
	if (closeNum == "") == (closeRate == "") {
		return nil, validationErr("provide exactly one of closeNum or closeRate")
	}
	if _, err := requireNonEmpty(positionType, "positionType"); err != nil {
		return nil, err
	}
	params := map[string]interface{}{"id": id, "positionType": positionType}
	if closeNum != "" {
		if _, err := decimalArg(closeNum, "closeNum"); err != nil {
			return nil, err
		}
		params["closeNum"] = closeNum
	}
	if closeRate != "" {
		if _, err := decimalArg(closeRate, "closeRate"); err != nil {
			return nil, err
		}
		params["closeRate"] = closeRate
	}
	if orderPrice != "" {
		if _, err := decimalArg(orderPrice, "orderPrice"); err != nil {
			return nil, err
		}
		params["orderPrice"] = orderPrice
	}
	return s.call(ctx, "futures.trade.closeAPosition", params)
}

// CloseBatchPositions — DELETE /v1/perpum/batchClose. MUTATION.
// items is the documented JSON array of close instructions (e.g. thirdOrderId).
func (s *FuturesTradeService) CloseBatchPositions(ctx context.Context, items []map[string]interface{}) (*transport.Result, error) {
	if len(items) == 0 {
		return nil, validationErr("batch close requires a non-empty list of items")
	}
	return s.call(ctx, "futures.trade.closeBatchPositions", items)
}

// ClosePositionsAtMarketPrice — DELETE /v1/perpum/allpositions. MUTATION.
func (s *FuturesTradeService) ClosePositionsAtMarketPrice(ctx context.Context, instrument string) (*transport.Result, error) {
	if _, err := requireNonEmpty(instrument, "instrument"); err != nil {
		return nil, err
	}
	return s.call(ctx, "futures.trade.closePositionsAtMarketPrice",
		map[string]interface{}{"instrument": instrument})
}

// ModifyAnOrder — PUT /v1/perpum/order. MUTATION.
// Margin mode cannot be changed for existing orders (server code 9050).
// extra carries additional documented fields (e.g. openPrice).
func (s *FuturesTradeService) ModifyAnOrder(ctx context.Context, id, instrument, direction string, leverage, quantityUnit int, quantity string, positionModel int, extra map[string]interface{}) (*transport.Result, error) {
	if !validDirection(direction) {
		return nil, validationErr(`direction must be "long" or "short"`)
	}
	if !validQuantityUnit(quantityUnit) {
		return nil, validationErr("quantityUnit must be 0 (quote), 1 (contracts) or 2 (base)")
	}
	if !validPositionModel(positionModel) {
		return nil, validationErr("positionModel must be 0 (isolated) or 1 (cross)")
	}
	if _, err := requireNonEmpty(instrument, "instrument"); err != nil {
		return nil, err
	}
	if _, err := decimalArg(quantity, "quantity"); err != nil {
		return nil, err
	}
	params := map[string]interface{}{
		"id":            id,
		"instrument":    instrument,
		"direction":     direction,
		"leverage":      leverage,
		"quantityUnit":  quantityUnit,
		"quantity":      quantity,
		"positionModel": positionModel,
	}
	for k, v := range extra {
		params[k] = v
	}
	return s.call(ctx, "futures.trade.modifyAnOrder", params)
}

// ModifyStopLossTakeProfitInBatches — POST (route per Endpoint field).
//
// GATED: the Endpoint field declares /v1/perpum/updateTpsl while the
// Futures code-snippet routes the same operation through /v1/perpum/TPSL.
// Returns AmbiguousDocumentationError until the route is verified.
func (s *FuturesTradeService) ModifyStopLossTakeProfitInBatches(ctx context.Context, fields map[string]interface{}) (*transport.Result, error) {
	if err := validateTPSLFields(fields); err != nil {
		return nil, err
	}
	return nil, &coinwerr.AmbiguousDocumentationError{
		EndpointID: "futures.trade.modifyStopLossTakeProfitInBatches",
		Note: "endpoint field declares /v1/perpum/updateTpsl; code-snippet page " +
			"routes the same operation via /v1/perpum/TPSL",
	}
}

// PlaceAnOrder — POST /v1/perpum/order. MUTATION.
// openPrice is mandatory for plan (limit) orders per docs; "" otherwise.
// thirdOrderID is the client reference; "" omits. extra merges other
// documented fields.
func (s *FuturesTradeService) PlaceAnOrder(ctx context.Context, instrument, direction string, leverage, quantityUnit int, quantity string, positionModel int, positionType, openPrice, thirdOrderID string, extra map[string]interface{}) (*transport.Result, error) {
	if !validPositionType(positionType) {
		return nil, validationErr("positionType must be plan, planTrigger or execute")
	}
	if positionType == PositionTypePlan && openPrice == "" {
		return nil, validationErr("openPrice is mandatory for plan (limit) orders")
	}
	if !validDirection(direction) {
		return nil, validationErr(`direction must be "long" or "short"`)
	}
	if !validQuantityUnit(quantityUnit) {
		return nil, validationErr("quantityUnit must be 0 (quote), 1 (contracts) or 2 (base)")
	}
	if !validPositionModel(positionModel) {
		return nil, validationErr("positionModel must be 0 (isolated) or 1 (cross)")
	}
	if _, err := requireNonEmpty(instrument, "instrument"); err != nil {
		return nil, err
	}
	if _, err := decimalArg(quantity, "quantity"); err != nil {
		return nil, err
	}
	params := map[string]interface{}{
		"instrument":    instrument,
		"direction":     direction,
		"leverage":      leverage,
		"quantityUnit":  quantityUnit,
		"quantity":      quantity,
		"positionModel": positionModel,
		"positionType":  positionType,
	}
	if openPrice != "" {
		if _, err := decimalArg(openPrice, "openPrice"); err != nil {
			return nil, err
		}
		params["openPrice"] = openPrice
	}
	if thirdOrderID != "" {
		params["thirdOrderId"] = thirdOrderID
	}
	for k, v := range extra {
		params[k] = v
	}
	return s.call(ctx, "futures.trade.placeAnOrder", params)
}

// PlaceBatchOrders — POST /v1/perpum/batchOrders. MUTATION.
// orders is the documented JSON array of order objects (placeAnOrder fields).
func (s *FuturesTradeService) PlaceBatchOrders(ctx context.Context, orders []map[string]interface{}) (*transport.Result, error) {
	if len(orders) == 0 {
		return nil, validationErr("batch orders require a non-empty list")
	}
	return s.call(ctx, "futures.trade.placeBatchOrders", orders)
}

// ReverseAPosition — POST /v1/perpum/positions/reverse. MUTATION.
func (s *FuturesTradeService) ReverseAPosition(ctx context.Context, id string) (*transport.Result, error) {
	return s.call(ctx, "futures.trade.reverseAPosition", map[string]interface{}{"id": id})
}

// SetStopLossTakeProfit — POST /v1/perpum/TPSL. MUTATION.
// fields carries the documented SL/TP params (stopLossOrderPrice,
// stopProfitOrderPrice, stopLossPrice/Rates, priceType, stopFrom, stopType,
// closePiece, …) with documented invariants enforced.
func (s *FuturesTradeService) SetStopLossTakeProfit(ctx context.Context, id, instrument string, fields map[string]interface{}) (*transport.Result, error) {
	if _, err := requireNonEmpty(instrument, "instrument"); err != nil {
		return nil, err
	}
	if err := validateTPSLFields(fields); err != nil {
		return nil, err
	}
	params := map[string]interface{}{"id": id, "instrument": instrument}
	for k, v := range fields {
		params[k] = v
	}
	return s.call(ctx, "futures.trade.setStopLossTakeProfit", params)
}

// SetStopLossTakeProfitInBatches — POST /v1/perpum/addTpsl. MUTATION.
func (s *FuturesTradeService) SetStopLossTakeProfitInBatches(ctx context.Context, id, instrument string, fields map[string]interface{}) (*transport.Result, error) {
	if _, err := requireNonEmpty(instrument, "instrument"); err != nil {
		return nil, err
	}
	if err := validateTPSLFields(fields); err != nil {
		return nil, err
	}
	params := map[string]interface{}{"id": id, "instrument": instrument}
	for k, v := range fields {
		params[k] = v
	}
	return s.call(ctx, "futures.trade.setStopLossTakeProfitInBatches", params)
}

// SetTrailingStopLossTakeProfit — POST /v1/perpum/moveTPSL. MUTATION.
// callbackRate/quantity are decimal strings; quantityUnit is 0 (USDT) or
// 1 (contracts) per docs; triggerPrice="" omits.
func (s *FuturesTradeService) SetTrailingStopLossTakeProfit(ctx context.Context, openID, callbackRate, quantity string, quantityUnit int, triggerPrice string) (*transport.Result, error) {
	if quantityUnit != 0 && quantityUnit != 1 {
		return nil, validationErr("quantityUnit must be 0 (USDT) or 1 (contracts)")
	}
	if _, err := decimalArg(callbackRate, "callbackRate"); err != nil {
		return nil, err
	}
	if _, err := decimalArg(quantity, "quantity"); err != nil {
		return nil, err
	}
	params := map[string]interface{}{
		"openId":       openID,
		"callbackRate": callbackRate,
		"quantity":     quantity,
		"quantityUnit": quantityUnit,
	}
	if triggerPrice != "" {
		if _, err := decimalArg(triggerPrice, "triggerPrice"); err != nil {
			return nil, err
		}
		params["triggerPrice"] = triggerPrice
	}
	return s.call(ctx, "futures.trade.setTrailingStopLossTakeProfit", params)
}
