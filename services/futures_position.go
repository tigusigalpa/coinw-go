package services

import (
	"context"

	"github.com/tigusigalpa/coinw-go/model"
	"github.com/tigusigalpa/coinw-go/transport"
)

// FuturesPositionService — futures position reads.
type FuturesPositionService struct{ base }

func NewFuturesPositionService(ex *transport.Executor) *FuturesPositionService {
	return &FuturesPositionService{newBase(ex)}
}

// GetCurrentPositionInformation — GET /v1/perpum/positions.
// openIDs: documented id list (CSV); empty slice omits.
func (s *FuturesPositionService) GetCurrentPositionInformation(ctx context.Context, instrument string, openIDs []string) (*transport.Response[[]model.FuturesPosition], error) {
	if _, err := requireNonEmpty(instrument, "instrument"); err != nil {
		return nil, err
	}
	params := map[string]interface{}{"instrument": instrument}
	if len(openIDs) > 0 {
		params["openIds"] = csv(openIDs)
	}
	return callTyped[[]model.FuturesPosition](s.base, ctx, "futures.position.getCurrentPositionInformation", params)
}

// GetAllCurrentPositions — GET /v1/perpum/positions/all.
func (s *FuturesPositionService) GetAllCurrentPositions(ctx context.Context) (*transport.Response[[]model.FuturesPosition], error) {
	return callTyped[[]model.FuturesPosition](s.base, ctx, "futures.position.getAllCurrentPositions", nil)
}

// GetHistoricalPositionInformation — GET /v1/perpum/positions/history.
// positionModel<0 omits the documented filter; 0 selects isolated margin.
func (s *FuturesPositionService) GetHistoricalPositionInformation(ctx context.Context, instrument string, positionModel int) (*transport.Response[model.RowsPage[model.HistoricalPosition]], error) {
	params := nonNull(map[string]interface{}{"instrument": emptyToNil(instrument)})
	if positionModel >= 0 {
		if !validPositionModel(positionModel) {
			return nil, validationErr("positionModel must be 0 (isolated) or 1 (cross), or negative to omit")
		}
		params["positionModel"] = positionModel
	}
	return callTyped[model.RowsPage[model.HistoricalPosition]](s.base, ctx, "futures.position.getHistoricalPositionInformation", params)
}

// GetMaxContractSize — GET /v1/perpum/orders/maxSize.
func (s *FuturesPositionService) GetMaxContractSize(ctx context.Context, instrument string, leverage, positionModel int) (*transport.Response[model.MaximumContractSize], error) {
	if _, err := requireNonEmpty(instrument, "instrument"); err != nil {
		return nil, err
	}
	return callTyped[model.MaximumContractSize](s.base, ctx, "futures.position.getMaxContractSize", map[string]interface{}{
		"instrument":    instrument,
		"leverage":      leverage,
		"positionModel": positionModel,
	})
}

// GetPositionMarginRate — GET /v1/perpum/positions/marginRate.
func (s *FuturesPositionService) GetPositionMarginRate(ctx context.Context, positionID string) (*transport.Response[model.ValueOrScalar[model.WireDecimal]], error) {
	return callTyped[model.ValueOrScalar[model.WireDecimal]](s.base, ctx, "futures.position.getPositionMarginRate", nonNull(map[string]interface{}{
		"positionId": emptyToNil(positionID),
	}))
}
