package services

import (
	"context"

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
func (s *FuturesAssetsService) CoinConversion(ctx context.Context, convertType int, dealPiece, faceValue, baseSize string) (*transport.Result, error) {
	for name, v := range map[string]string{"dealPiece": dealPiece, "faceValue": faceValue, "baseSize": baseSize} {
		if _, err := decimalArg(v, name); err != nil {
			return nil, err
		}
	}
	return s.call(ctx, "futures.assets.coinConversion", map[string]interface{}{
		"convertType": convertType,
		"dealPiece":   dealPiece,
		"faceValue":   faceValue,
		"baseSize":    baseSize,
	})
}

// EnableDisableMegaCoupon — POST /v1/perpum/account/almightyGoldInfo. MUTATION.
func (s *FuturesAssetsService) EnableDisableMegaCoupon(ctx context.Context, status int) (*transport.Result, error) {
	return s.call(ctx, "futures.assets.enableDisableMegaCoupon", map[string]interface{}{"status": status})
}

// GetFuturesAccountAssets — GET /v1/perpum/account/getUserAssets.
func (s *FuturesAssetsService) GetFuturesAccountAssets(ctx context.Context) (*transport.Result, error) {
	return s.call(ctx, "futures.assets.getFuturesAccountAssets", nil)
}

// GetFuturesAccountFees — GET /v1/perpum/account/fees.
func (s *FuturesAssetsService) GetFuturesAccountFees(ctx context.Context) (*transport.Result, error) {
	return s.call(ctx, "futures.assets.getFuturesAccountFees", nil)
}

// GetMarginMode — GET /v1/perpum/positions/type.
func (s *FuturesAssetsService) GetMarginMode(ctx context.Context) (*transport.Result, error) {
	return s.call(ctx, "futures.assets.getMarginMode", nil)
}

// GetMaxTransferableBalance — GET /v1/perpum/account/available.
func (s *FuturesAssetsService) GetMaxTransferableBalance(ctx context.Context) (*transport.Result, error) {
	return s.call(ctx, "futures.assets.getMaxTransferableBalance", nil)
}

// GetMaxUsersContractSize — GET /v1/perpum/orders/availSize.
func (s *FuturesAssetsService) GetMaxUsersContractSize(ctx context.Context, instrument string) (*transport.Result, error) {
	if _, err := requireNonEmpty(instrument, "instrument"); err != nil {
		return nil, err
	}
	return s.call(ctx, "futures.assets.getMaxUsersContractSize",
		map[string]interface{}{"instrument": instrument})
}

// GetMegaCouponBalance — GET /v1/perpum/account/almightyGoldInfo.
// typ<0 omits the documented type filter.
func (s *FuturesAssetsService) GetMegaCouponBalance(ctx context.Context, typ int) (*transport.Result, error) {
	params := map[string]interface{}{}
	if typ >= 0 {
		params["type"] = typ
	}
	return s.call(ctx, "futures.assets.getMegaCouponBalance", params)
}

// GetTransactionDetails3Days — GET /v1/perpum/orders/deals.
// page/pageSize/positionModel<1 omit the optional documented fields.
func (s *FuturesAssetsService) GetTransactionDetails3Days(ctx context.Context, instrument string, page, pageSize, positionModel int) (*transport.Result, error) {
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
	if positionModel > 0 {
		params["positionModel"] = positionModel
	}
	return s.call(ctx, "futures.assets.getTransactionDetails3Days", params)
}

// GetTransactionDetails3Months — GET /v1/perpum/orders/deals/history.
func (s *FuturesAssetsService) GetTransactionDetails3Months(ctx context.Context, instrument string, page, pageSize, positionModel int) (*transport.Result, error) {
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
	if positionModel > 0 {
		params["positionModel"] = positionModel
	}
	return s.call(ctx, "futures.assets.getTransactionDetails3Months", params)
}

// SetMarginMode — POST /v1/perpum/positions/type. MUTATION.
func (s *FuturesAssetsService) SetMarginMode(ctx context.Context, layout, positionModel int) (*transport.Result, error) {
	return s.call(ctx, "futures.assets.setMarginMode", map[string]interface{}{
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
func (s *FuturesCommonService) GetLeverageInformation(ctx context.Context, orderID, positionID string) (*transport.Result, error) {
	return s.call(ctx, "futures.common.getLeverageInformation", nonNull(map[string]interface{}{
		"orderId":    emptyToNil(orderID),
		"positionId": emptyToNil(positionID),
	}))
}
