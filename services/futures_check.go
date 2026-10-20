package services

import (
	"context"

	"github.com/tigusigalpa/coinw-go/transport"
)

// FuturesOriginTypes are the documented originType values.
var FuturesOriginTypes = []string{"plan", "planTrigger", "execute"}

func validOriginType(v string) bool {
	for _, t := range FuturesOriginTypes {
		if t == v {
			return true
		}
	}
	return false
}

// FuturesCheckService — futures order/SL-TP state reads.
type FuturesCheckService struct{ base }

func NewFuturesCheckService(ex *transport.Executor) *FuturesCheckService {
	return &FuturesCheckService{newBase(ex)}
}

func (s *FuturesCheckService) batchParams(symbols []string, originType string, page, pageSize int) (map[string]interface{}, error) {
	if _, err := requireMaxInstruments(symbols, 20); err != nil {
		return nil, err
	}
	if !validOriginType(originType) {
		return nil, validationErr("originType must be one of: %v", FuturesOriginTypes)
	}
	params := map[string]interface{}{"symbols": csv(symbols), "originType": originType}
	if page > 0 {
		params["page"] = page
	}
	if pageSize > 0 {
		if _, err := requirePageSize(pageSize, 500); err != nil {
			return nil, err
		}
		params["pageSize"] = pageSize
	}
	return params, nil
}

// GetBatchHistoricalOrders7Days — GET /v1/perpum/orders/batchInsHistory.
func (s *FuturesCheckService) GetBatchHistoricalOrders7Days(ctx context.Context, symbols []string, originType string, page, pageSize int) (*transport.Result, error) {
	p, err := s.batchParams(symbols, originType, page, pageSize)
	if err != nil {
		return nil, err
	}
	return s.call(ctx, "futures.check.getBatchHistoricalOrders7Days", p)
}

// GetBatchHistoricalOrders3Months — GET /v1/perpum/orders/batchInsArchive.
func (s *FuturesCheckService) GetBatchHistoricalOrders3Months(ctx context.Context, symbols []string, originType string, page, pageSize int) (*transport.Result, error) {
	p, err := s.batchParams(symbols, originType, page, pageSize)
	if err != nil {
		return nil, err
	}
	return s.call(ctx, "futures.check.getBatchHistoricalOrders3Months", p)
}

// GetCurrentOrders — GET /v1/perpum/orders/open.
func (s *FuturesCheckService) GetCurrentOrders(ctx context.Context, instrument, positionType string) (*transport.Result, error) {
	if _, err := requireNonEmpty(instrument, "instrument"); err != nil {
		return nil, err
	}
	if _, err := requireNonEmpty(positionType, "positionType"); err != nil {
		return nil, err
	}
	return s.call(ctx, "futures.check.getCurrentOrders", map[string]interface{}{
		"instrument":   instrument,
		"positionType": positionType,
	})
}

// GetHistoricalOrders7Days — GET /v1/perpum/orders/history.
// page/pageSize<1 omit; empty originType omits.
func (s *FuturesCheckService) GetHistoricalOrders7Days(ctx context.Context, instrument string, page, pageSize int, originType string) (*transport.Result, error) {
	if _, err := requireNonEmpty(instrument, "instrument"); err != nil {
		return nil, err
	}
	if originType != "" && !validOriginType(originType) {
		return nil, validationErr("originType must be one of: %v", FuturesOriginTypes)
	}
	params := map[string]interface{}{"instrument": instrument}
	if page > 0 {
		params["page"] = page
	}
	if pageSize > 0 {
		params["pageSize"] = pageSize
	}
	if originType != "" {
		params["originType"] = originType
	}
	return s.call(ctx, "futures.check.getHistoricalOrders7Days", params)
}

// GetHistoricalOrders3Months — GET /v1/perpum/orders/archive.
func (s *FuturesCheckService) GetHistoricalOrders3Months(ctx context.Context, instrument string, page, pageSize int, originType string) (*transport.Result, error) {
	if _, err := requireNonEmpty(instrument, "instrument"); err != nil {
		return nil, err
	}
	if originType != "" && !validOriginType(originType) {
		return nil, validationErr("originType must be one of: %v", FuturesOriginTypes)
	}
	params := map[string]interface{}{"instrument": instrument}
	if page > 0 {
		params["page"] = page
	}
	if pageSize > 0 {
		params["pageSize"] = pageSize
	}
	if originType != "" {
		params["originType"] = originType
	}
	return s.call(ctx, "futures.check.getHistoricalOrders3Months", params)
}

// GetOrderInformation — GET /v1/perpum/order.
// sourceIds: documented list of order ids (sent as JSON array in query? —
// the documented field is a comma list; pass []string, encoded CSV).
func (s *FuturesCheckService) GetOrderInformation(ctx context.Context, positionType string, sourceIDs []string, instrument string) (*transport.Result, error) {
	if _, err := requireNonEmpty(positionType, "positionType"); err != nil {
		return nil, err
	}
	params := map[string]interface{}{"positionType": positionType}
	if len(sourceIDs) > 0 {
		params["sourceIds"] = csv(sourceIDs)
	}
	if instrument != "" {
		params["instrument"] = instrument
	}
	return s.call(ctx, "futures.check.getOrderInformation", params)
}

// GetPendingOrderCount — GET /v1/perpum/orders/openQuantity.
func (s *FuturesCheckService) GetPendingOrderCount(ctx context.Context) (*transport.Result, error) {
	return s.call(ctx, "futures.check.getPendingOrderCount", nil)
}

// SLTPQuery carries the documented optional SL/TP lookup fields.
type SLTPQuery struct {
	StopFrom    int    // 1 limit, 2 position, 3 plan — 0 omits
	OrderID     string // "" omits
	OpenID      string
	PlanOrderID string
	Instrument  string
}

// GetStopLossTakeProfitInformation — GET /v1/perpum/TPSL.
func (s *FuturesCheckService) GetStopLossTakeProfitInformation(ctx context.Context, q SLTPQuery) (*transport.Result, error) {
	params := map[string]interface{}{}
	if q.StopFrom != 0 {
		params["stopFrom"] = q.StopFrom
	}
	if q.OrderID != "" {
		params["orderId"] = q.OrderID
	}
	if q.OpenID != "" {
		params["openId"] = q.OpenID
	}
	if q.PlanOrderID != "" {
		params["planOrderId"] = q.PlanOrderID
	}
	if q.Instrument != "" {
		params["instrument"] = q.Instrument
	}
	return s.call(ctx, "futures.check.getStopLossTakeProfitInformation", params)
}

// GetTrailingStopLossTakeProfitInformation — GET /v1/perpum/moveTPSL.
func (s *FuturesCheckService) GetTrailingStopLossTakeProfitInformation(ctx context.Context) (*transport.Result, error) {
	return s.call(ctx, "futures.check.getTrailingStopLossTakeProfitInformation", nil)
}
