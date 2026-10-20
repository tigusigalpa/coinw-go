package services

import (
	"bytes"
	"context"
	"encoding/json"
	"regexp"

	"github.com/tigusigalpa/coinw-go/model"
	"github.com/tigusigalpa/coinw-go/transport"
)

var clientOrderIDPattern = regexp.MustCompile(`^[A-Za-z0-9_-]{1,50}$`)

func requestMap(v any) (map[string]interface{}, error) {
	raw, err := json.Marshal(v)
	if err != nil {
		return nil, validationErr("request contains an invalid value")
	}
	var result map[string]interface{}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	if err = dec.Decode(&result); err != nil {
		return nil, validationErr("request must be an object")
	}
	for k, v := range result {
		if n, ok := v.(json.Number); ok {
			i, err := n.Int64()
			if err != nil || int64(int(i)) != i {
				return nil, validationErr("parameter %s must be a supported integer", k)
			}
			result[k] = int(i)
		}
	}
	return result, nil
}

func orderParams(r model.FuturesOrderRequest, extended bool) (map[string]interface{}, error) {
	params, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	if !validDirection(string(r.Direction)) || !validQuantityUnit(int(r.QuantityUnit)) || !validPositionModel(int(r.PositionModel)) {
		return nil, validationErr("invalid direction, quantity unit or position model")
	}
	if !validPositionType(string(r.PositionType)) || (!extended && r.PositionType != model.LimitOrder && r.PositionType != model.TriggerOrder && r.PositionType != model.MarketOrder) {
		return nil, validationErr("unsupported order type for this operation")
	}
	if r.PositionType == model.LimitOrder && r.OpenPrice == nil {
		return nil, validationErr("openPrice is mandatory for plan orders")
	}
	if r.TriggerType != nil {
		if *r.TriggerType != 0 && *r.TriggerType != 1 {
			return nil, validationErr("triggerType must be 0 or 1")
		}
		if *r.TriggerType == 0 && r.OpenPrice == nil {
			return nil, validationErr("openPrice is mandatory for a limit trigger")
		}
	}
	if r.PositionType == model.TriggerOrder && r.TriggerPrice == nil {
		return nil, validationErr("triggerPrice is mandatory for planTrigger orders")
	}
	if r.ThirdOrderID != "" && !clientOrderIDPattern.MatchString(r.ThirdOrderID) {
		return nil, validationErr("thirdOrderId must contain 1..50 Latin letters, digits, hyphens or underscores")
	}
	return params, nil
}

// PlaceAnOrder explicitly submits one order. The acknowledgement is an order
// identifier, not proof of a fill. The transport never retries this mutation.
func (s *FuturesTradeService) PlaceAnOrder(ctx context.Context, r model.FuturesOrderRequest) (*transport.Response[model.ValueOrScalar[model.StringOrInteger]], error) {
	p, err := orderParams(r, true)
	if err != nil {
		return nil, err
	}
	return callTyped[model.ValueOrScalar[model.StringOrInteger]](s.base, ctx, "futures.trade.placeAnOrder", p)
}

// ModifyAnOrder explicitly modifies an existing order without changing margin
// mode implicitly. The caller must supply the existing positionModel.
func (s *FuturesTradeService) ModifyAnOrder(ctx context.Context, r model.ModifyOrderRequest) (*transport.Response[model.ModifiedOrderAcknowledgement], error) {
	if _, err := requireNonEmpty(r.ID, "id"); err != nil {
		return nil, err
	}
	p, err := orderParams(r.FuturesOrderRequest, false)
	if err != nil {
		return nil, err
	}
	extra, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	for k, v := range extra {
		p[k] = v
	}
	if r.FeeType != nil && *r.FeeType != 0 && *r.FeeType != 1 {
		return nil, validationErr("feeType must be 0 or 1")
	}
	return callTyped[model.ModifiedOrderAcknowledgement](s.base, ctx, "futures.trade.modifyAnOrder", p)
}

func (s *FuturesTradeService) PlaceBatchOrders(ctx context.Context, orders []model.FuturesOrderRequest) (*transport.Response[[]model.BatchOrderAcknowledgement], error) {
	if len(orders) == 0 {
		return nil, validationErr("batch orders require a non-empty list")
	}
	items := make([]map[string]interface{}, len(orders))
	for i, r := range orders {
		p, err := orderParams(r, false)
		if err != nil {
			return nil, err
		}
		items[i] = p
	}
	return callTyped[[]model.BatchOrderAcknowledgement](s.base, ctx, "futures.trade.placeBatchOrders", items)
}

func (s *FuturesTradeService) CloseBatchPositions(ctx context.Context, items []model.CloseBatchPositionRequest) (*transport.Response[[]model.BatchOrderAcknowledgement], error) {
	if len(items) == 0 {
		return nil, validationErr("batch close requires a non-empty list")
	}
	for _, r := range items {
		if !clientOrderIDPattern.MatchString(r.ThirdOrderID) {
			return nil, validationErr("invalid thirdOrderId")
		}
	}
	return callTyped[[]model.BatchOrderAcknowledgement](s.base, ctx, "futures.trade.closeBatchPositions", items)
}

func (s *FuturesTradeService) SetStopLossTakeProfit(ctx context.Context, id, instrument string, r model.SLTPPrices) (*transport.Response[model.NoData], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return s.setStopLossTakeProfit(ctx, id, instrument, p)
}
func (s *FuturesTradeService) SetStopLossTakeProfitInBatches(ctx context.Context, id, instrument string, r model.BatchSLTPRequest) (*transport.Response[model.NoData], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return s.setStopLossTakeProfitInBatches(ctx, id, instrument, p)
}
func (s *FuturesTradeService) ModifyStopLossTakeProfitInBatches(ctx context.Context, r model.ModifyBatchSLTPRequest) (*transport.Response[model.NoData], error) {
	p, err := requestMap(r)
	if err != nil {
		return nil, err
	}
	return s.modifyStopLossTakeProfitInBatches(ctx, p)
}
