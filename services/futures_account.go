package services

import (
	"context"

	"github.com/tigusigalpa/coinw-go/model"
	"github.com/tigusigalpa/coinw-go/transport"
)

// ---------------------------------------------------------------------------
// Futures > Account & Assets
// ---------------------------------------------------------------------------

// FuturesAssetsService — futures account/asset operations.
type FuturesAssetsService struct{ base }

func NewFuturesAssetsService(ex *transport.Executor) *FuturesAssetsService {
	return &FuturesAssetsService{newBase(ex)}
}

// CoinConversion — POST /v1/perpum/pieceConvert. MUTATION.
// All three sizes are decimal strings.
func (s *FuturesAssetsService) CoinConversion(ctx context.Context, convertType int, dealPiece, faceValue, baseSize string) (*transport.Response[model.DecimalValue], error) {
	for name, v := range map[string]string{"dealPiece": dealPiece, "faceValue": faceValue, "baseSize": baseSize} {
		if _, err := decimalArg(v, name); err != nil {
			return nil, err
		}
	}
	return callTyped[model.DecimalValue](s.base, ctx, "futures.assets.coinConversion", map[string]interface{}{
		"convertType": convertType,
		"dealPiece":   dealPiece,
		"faceValue":   faceValue,
		"baseSize":    baseSize,
	})
}

// EnableDisableMegaCoupon — POST /v1/perpum/account/almightyGoldInfo. MUTATION.
func (s *FuturesAssetsService) EnableDisableMegaCoupon(ctx context.Context, status int) (*transport.Response[model.NoData], error) {
	return callTyped[model.NoData](s.base, ctx, "futures.assets.enableDisableMegaCoupon", map[string]interface{}{"status": status})
}

// GetFuturesAccountAssets — GET /v1/perpum/account/getUserAssets.
func (s *FuturesAssetsService) GetFuturesAccountAssets(ctx context.Context) (*transport.Response[model.FuturesAccountAssets], error) {
	return callTyped[model.FuturesAccountAssets](s.base, ctx, "futures.assets.getFuturesAccountAssets", nil)
}

// GetFuturesAccountFees — GET /v1/perpum/account/fees.
func (s *FuturesAssetsService) GetFuturesAccountFees(ctx context.Context) (*transport.Response[model.FuturesAccountFees], error) {
	return callTyped[model.FuturesAccountFees](s.base, ctx, "futures.assets.getFuturesAccountFees", nil)
}

// GetMarginMode — GET /v1/perpum/positions/type.
func (s *FuturesAssetsService) GetMarginMode(ctx context.Context) (*transport.Response[model.MarginMode], error) {
	return callTyped[model.MarginMode](s.base, ctx, "futures.assets.getMarginMode", nil)
}

// GetMaxTransferableBalance — GET /v1/perpum/account/available.
func (s *FuturesAssetsService) GetMaxTransferableBalance(ctx context.Context) (*transport.Response[model.ValueOrScalar[model.WireDecimal]], error) {
	return callTyped[model.ValueOrScalar[model.WireDecimal]](s.base, ctx, "futures.assets.getMaxTransferableBalance", nil)
}

// GetMaxUsersContractSize — GET /v1/perpum/orders/availSize.
func (s *FuturesAssetsService) GetMaxUsersContractSize(ctx context.Context, instrument string) (*transport.Response[model.AvailableContractSize], error) {
	if _, err := requireNonEmpty(instrument, "instrument"); err != nil {
		return nil, err
	}
	return callTyped[model.AvailableContractSize](s.base, ctx, "futures.assets.getMaxUsersContractSize",
		map[string]interface{}{"instrument": instrument})
}

// GetMegaCouponBalance — GET /v1/perpum/account/almightyGoldInfo.
// typ<0 omits the documented type filter.
func (s *FuturesAssetsService) GetMegaCouponBalance(ctx context.Context, typ int) (*transport.Response[[]model.MegaCoupon], error) {
	params := map[string]interface{}{}
	if typ >= 0 {
		params["type"] = typ
	}
	return callTyped[[]model.MegaCoupon](s.base, ctx, "futures.assets.getMegaCouponBalance", params)
}

// GetTransactionDetails3Days — GET /v1/perpum/orders/deals.
// page/pageSize<1 omit; positionModel<0 omits (0 means isolated margin).
func (s *FuturesAssetsService) GetTransactionDetails3Days(ctx context.Context, instrument string, page, pageSize, positionModel int) (*transport.Response[model.RowsPage[model.FuturesTransaction]], error) {
	if _, err := requireNonEmpty(instrument, "instrument"); err != nil {
		return nil, err
	}
	params := map[string]interface{}{"instrument": instrument}
	if page > 0 {
		params["page"] = page
	}
	if pageSize > 0 {
		params["pageSize"] = pageSize
	}
	if positionModel >= 0 {
		if !validPositionModel(positionModel) {
			return nil, validationErr("positionModel must be 0 (isolated) or 1 (cross), or negative to omit")
		}
		params["positionModel"] = positionModel
	}
	return callTyped[model.RowsPage[model.FuturesTransaction]](s.base, ctx, "futures.assets.getTransactionDetails3Days", params)
}

// GetTransactionDetails3Months — GET /v1/perpum/orders/deals/history.
func (s *FuturesAssetsService) GetTransactionDetails3Months(ctx context.Context, instrument string, page, pageSize, positionModel int) (*transport.Response[model.RowsPage[model.FuturesTransaction]], error) {
	if _, err := requireNonEmpty(instrument, "instrument"); err != nil {
		return nil, err
	}
	params := map[string]interface{}{"instrument": instrument}
	if page > 0 {
		params["page"] = page
	}
	if pageSize > 0 {
		params["pageSize"] = pageSize
	}
	if positionModel >= 0 {
		if !validPositionModel(positionModel) {
			return nil, validationErr("positionModel must be 0 (isolated) or 1 (cross), or negative to omit")
		}
		params["positionModel"] = positionModel
	}
	return callTyped[model.RowsPage[model.FuturesTransaction]](s.base, ctx, "futures.assets.getTransactionDetails3Months", params)
}

// SetMarginMode — POST /v1/perpum/positions/type. MUTATION.
func (s *FuturesAssetsService) SetMarginMode(ctx context.Context, layout, positionModel int) (*transport.Response[string], error) {
	return callTyped[string](s.base, ctx, "futures.assets.setMarginMode", map[string]interface{}{
		"layout":        layout,
		"positionModel": positionModel,
	})
}

// ---------------------------------------------------------------------------
// Futures > Common
// ---------------------------------------------------------------------------

// FuturesCommonService — futures common reads.
type FuturesCommonService struct{ base }

func NewFuturesCommonService(ex *transport.Executor) *FuturesCommonService {
	return &FuturesCommonService{newBase(ex)}
}

// GetLeverageInformation — GET /v1/perpum/positions/leverage.
// orderID/positionID: pass "" to omit; at least one is meaningful per docs.
func (s *FuturesCommonService) GetLeverageInformation(ctx context.Context, orderID, positionID string) (*transport.Response[model.ValueOrScalar[model.WireDecimal]], error) {
	return callTyped[model.ValueOrScalar[model.WireDecimal]](s.base, ctx, "futures.common.getLeverageInformation", nonNull(map[string]interface{}{
		"orderId":    emptyToNil(orderID),
		"positionId": emptyToNil(positionID),
	}))
}
