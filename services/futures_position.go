package services

import (
	"context"

	"github.com/tigusigalpa/coinw-go/transport"
)

// FuturesPositionService — futures position reads.
type FuturesPositionService struct{ base }

func NewFuturesPositionService(ex *transport.Executor) *FuturesPositionService {
	return &FuturesPositionService{newBase(ex)}
}

// GetCurrentPositionInformation — GET /v1/perpum/positions.
// openIDs: documented id list (CSV); empty slice omits.
func (s *FuturesPositionService) GetCurrentPositionInformation(ctx context.Context, instrument string, openIDs []string) (*transport.Result, error) {
	if _, err := requireNonEmpty(instrument, "instrument"); err != nil {
		return nil, err
	}
	params := map[string]interface{}{"instrument": instrument}
	if len(openIDs) > 0 {
		params["openIds"] = csv(openIDs)
	}
	return s.call(ctx, "futures.position.getCurrentPositionInformation", params)
}

// GetAllCurrentPositions — GET /v1/perpum/positions/all.
func (s *FuturesPositionService) GetAllCurrentPositions(ctx context.Context) (*transport.Result, error) {
	return s.call(ctx, "futures.position.getAllCurrentPositions", nil)
}

// GetHistoricalPositionInformation — GET /v1/perpum/positions/history.
// positionModel<1 omits the documented filter.
func (s *FuturesPositionService) GetHistoricalPositionInformation(ctx context.Context, instrument string, positionModel int) (*transport.Result, error) {
	if _, err := requireNonEmpty(instrument, "instrument"); err != nil {
		return nil, err
	}
	params := map[string]interface{}{"instrument": instrument}
	if positionModel > 0 {
		params["positionModel"] = positionModel
	}
	return s.call(ctx, "futures.position.getHistoricalPositionInformation", params)
}

// GetMaxContractSize — GET /v1/perpum/orders/maxSize.
func (s *FuturesPositionService) GetMaxContractSize(ctx context.Context, instrument string, leverage, positionModel int) (*transport.Result, error) {
	if _, err := requireNonEmpty(instrument, "instrument"); err != nil {
		return nil, err
	}
	return s.call(ctx, "futures.position.getMaxContractSize", map[string]interface{}{
		"instrument":    instrument,
		"leverage":      leverage,
		"positionModel": positionModel,
	})
}

// GetPositionMarginRate — GET /v1/perpum/positions/marginRate.
func (s *FuturesPositionService) GetPositionMarginRate(ctx context.Context, positionID string) (*transport.Result, error) {
	return s.call(ctx, "futures.position.getPositionMarginRate", nonNull(map[string]interface{}{
		"positionId": emptyToNil(positionID),
	}))
}
