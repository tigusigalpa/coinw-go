package services

import (
	"context"
	"regexp"
	"strings"

	coinwerr "github.com/tigusigalpa/coinw-go/errors"
	"github.com/tigusigalpa/coinw-go/model"
	"github.com/tigusigalpa/coinw-go/transport"
)

// ---------------------------------------------------------------------------
// Common > Accounts (Spot private profile)
// ---------------------------------------------------------------------------

// AccountService — deposit/withdrawal/transfer operations.
type AccountService struct{ base }

func NewAccountService(ex *transport.Executor) *AccountService {
	return &AccountService{newBase(ex)}
}

// GetDepositWithdrawalAddress — command=returnDepositAddresses.
// symbolID is the published currency identifier; chain is required by its table.
func (s *AccountService) GetDepositWithdrawalAddress(ctx context.Context, symbolID, chain string) (*transport.Response[[]model.DepositAddress], error) {
	if _, err := requireNonEmpty(symbolID, "symbolId"); err != nil {
		return nil, err
	}
	return callTyped[[]model.DepositAddress](s.base, ctx, "common.account.getDepositWithdrawalAddress", nonNull(map[string]interface{}{
		"symbolId": symbolID,
		"chain":    emptyToNil(chain),
	}))
}

// GetDepositWithdrawalHistory — command=returnDepositsWithdrawals.
// symbol is required by the table; up to 20 comma-separated currencies.
func (s *AccountService) GetDepositWithdrawalHistory(ctx context.Context, symbol string) (*transport.Response[[]model.DepositWithdrawal], error) {
	return callTyped[[]model.DepositWithdrawal](s.base, ctx, "common.account.getDepositWithdrawalHistory", nonNull(map[string]interface{}{
		"symbol": emptyToNil(symbol),
	}))
}

// WithdrawalRequest — documented doWithdraw fields.
type WithdrawalRequest struct {
	Amount      model.Decimal // decimal string — required
	Currency    string        // e.g. "HC" — required
	Address     string        // destination — required
	Chain       string        // documented blockchain name — required
	Type        string        // e.g. "internal_transfer" — "" omits
	Memo        string        // "" omits
	InnerToType int           // <=0 omits
}

// InitiateWithdrawal — command=doWithdraw. HIGH-RISK MUTATION.
//
// Funds leave the account. The SDK never calls this implicitly and never
// retries it. memo/address-type fields are passed exactly as given — the
// documentation does not define their value space.
func (s *AccountService) InitiateWithdrawal(ctx context.Context, r WithdrawalRequest) (*transport.Response[model.WithdrawalAcknowledgement], error) {
	if _, err := decimalArg(string(r.Amount), "amount"); err != nil {
		return nil, err
	}
	if _, err := requireNonEmpty(r.Currency, "currency"); err != nil {
		return nil, err
	}
	if _, err := requireNonEmpty(r.Address, "address"); err != nil {
		return nil, err
	}
	params := map[string]interface{}{
		"amount":      string(r.Amount),
		"currency":    r.Currency,
		"address":     r.Address,
		"chain":       r.Chain,
		"memo":        r.Memo,
		"innerToType": r.InnerToType,
	}
	if r.Type != "" {
		params["type"] = r.Type
	}
	if r.Memo != "" {
		params["memo"] = r.Memo
	}
	if r.InnerToType > 0 {
		params["innerToType"] = r.InnerToType
	}
	return callTyped[model.WithdrawalAcknowledgement](s.base, ctx, "common.account.initiateWithdrawal", params)
}

// CancelWithdrawal — command=cancelWithdraw. MUTATION.
func (s *AccountService) CancelWithdrawal(ctx context.Context, id string) (*transport.Response[model.NoData], error) {
	if _, err := requireNonEmpty(id, "id"); err != nil {
		return nil, err
	}
	return callTyped[model.NoData](s.base, ctx, "common.account.cancelWithdrawal", map[string]interface{}{"id": id})
}

// TransferAssets — command=spotWealthTransfer. MUTATION.
// amount is a decimal string; the other fields are the documented
// account/biz codes (e.g. "SPOT", "WEALTH", "SPOT_TO_WEALTH", "BTC").
func (s *AccountService) TransferAssets(ctx context.Context, accountType, targetAccountType, bizType, coinCode, amount string) (*transport.Response[bool], error) {
	for name, v := range map[string]string{
		"accountType": accountType, "targetAccountType": targetAccountType,
		"bizType": bizType, "coinCode": coinCode,
	} {
		if _, err := requireNonEmpty(v, name); err != nil {
			return nil, err
		}
	}
	if _, err := decimalArg(amount, "amount"); err != nil {
		return nil, err
	}
	return callTyped[bool](s.base, ctx, "common.account.transferAssets", map[string]interface{}{
		"accountType":       accountType,
		"targetAccountType": targetAccountType,
		"bizType":           bizType,
		"coinCode":          coinCode,
		"amount":            amount,
	})
}

// ---------------------------------------------------------------------------
// Common > Sub-accounts (Spot private profile, master key only)
// ---------------------------------------------------------------------------

// Sub-account transfer direction and account-type constants.
const (
	SubSideOut = "OUT" // Master → Sub
	SubSideIn  = "IN"  // Sub → Master

	SubAccountAsset    = 0
	SubAccountSpot     = 1
	SubAccountContract = 2
	SubAccountSub      = 3
	SubAccountMaster   = 4
)

// SubAccountService — master-key sub-account operations.
type SubAccountService struct{ base }

func NewSubAccountService(ex *transport.Executor) *SubAccountService {
	return &SubAccountService{newBase(ex)}
}

// GetSubAccountList — command=subAccount_list.
func (s *SubAccountService) GetSubAccountList(ctx context.Context) (*transport.Response[[]model.SubAccount], error) {
	return callTyped[[]model.SubAccount](s.base, ctx, "common.subAccount.getSubAccountList", nil)
}

// GetSubAccountApiKeys — command=subAccount_apikey.
func (s *SubAccountService) GetSubAccountApiKeys(ctx context.Context) (*transport.Response[[]model.SubAccountAPIKey], error) {
	return callTyped[[]model.SubAccountAPIKey](s.base, ctx, "common.subAccount.getSubAccountApiKeys", nil)
}

// ResetSubAccountApiKey — command=reset_subAccount_apikey. MUTATION.
func (s *SubAccountService) ResetSubAccountApiKey(ctx context.Context, subAccountAPIKey string) (*transport.Response[model.NoData], error) {
	if _, err := requireNonEmpty(subAccountAPIKey, "subAccount_apiKey"); err != nil {
		return nil, err
	}
	return callTyped[model.NoData](s.base, ctx, "common.subAccount.resetSubAccountApiKey",
		map[string]interface{}{"subAccount_apiKey": subAccountAPIKey})
}

// GetMaxTransferableAmount — command=subAccount_transferAvailable.
//
// GATED: the endpoint body shows /api/v1/private while its hyperlink
// resolves to /api/v1/api/v1/private (doubled prefix). Returns
// AmbiguousDocumentationError until the route is verified.
func (s *SubAccountService) GetMaxTransferableAmount(ctx context.Context, coinCode, subUID string) (*transport.Response[model.WireDecimal], error) {
	if _, err := requireNonEmpty(coinCode, "coinCode"); err != nil {
		return nil, err
	}
	return nil, &coinwerr.AmbiguousDocumentationError{
		EndpointID: "common.subAccount.getMaxTransferableAmount",
		Note: "endpoint body shows /api/v1/private; endpoint hyperlink resolves " +
			"to /api/v1/api/v1/private (doubled prefix). Unverified route",
	}
}

// GetTransferHistory — command=returnTransferHistory.
// All five documented parameters are mandatory. from/to are documented
// account types 0..4 (Sub* constants).
func (s *SubAccountService) GetTransferHistory(ctx context.Context, currency string, from, to, pageNo, pageSize int) (*transport.Response[[]model.Transfer], error) {
	for _, t := range []int{from, to} {
		if t < 0 || t > 4 {
			return nil, validationErr("from/to must be a documented account type 0..4")
		}
	}
	if _, err := requireNonEmpty(currency, "currency"); err != nil {
		return nil, err
	}
	if pageNo < 1 || pageSize < 1 {
		return nil, validationErr("pageNo and pageSize must be >= 1")
	}
	return callTyped[[]model.Transfer](s.base, ctx, "common.subAccount.getTransferHistory", map[string]interface{}{
		"currency": currency,
		"from":     from,
		"to":       to,
		"pageNo":   pageNo,
		"pageSize": pageSize,
	})
}

// AccountTransfer — command=subAccount_transfer. MUTATION.
// side: SubSideOut (Master→Sub) or SubSideIn (Sub→Master); amount decimal string.
func (s *SubAccountService) AccountTransfer(ctx context.Context, side, subUID, currency, amount string) (*transport.Response[model.TransferAcknowledgement], error) {
	if side != SubSideOut && side != SubSideIn {
		return nil, validationErr("side must be OUT (Master→Sub) or IN (Sub→Master)")
	}
	if _, err := requireNonEmpty(currency, "currency"); err != nil {
		return nil, err
	}
	if _, err := decimalArg(amount, "amount"); err != nil {
		return nil, err
	}
	return callTyped[model.TransferAcknowledgement](s.base, ctx, "common.subAccount.accountTransfer", map[string]interface{}{
		"side":     side,
		"subUid":   subUID,
		"currency": currency,
		"amount":   amount,
	})
}

// ---------------------------------------------------------------------------
// Common > Affiliate (Futures private profile, /api/v3/patener/...)
// ---------------------------------------------------------------------------

var (
	affiliateDateTimePattern = regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$`)
	affiliateDatePattern     = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
)

func requireAffiliateDateTime(v, name string) error {
	if v != "" && !affiliateDateTimePattern.MatchString(v) {
		return validationErr("%s must match 'yyyy-MM-dd HH:mm:ss'", name)
	}
	return nil
}

func requireAffiliateDate(v, name string) error {
	if v != "" && !affiliateDatePattern.MatchString(v) {
		return validationErr("%s must match 'yyyy-MM-dd'", name)
	}
	return nil
}

// AffiliateService — partner endpoints under the documented
// /api/v3/patener/openapi/v1/partner/... base path (literal "patener").
type AffiliateService struct{ base }

func NewAffiliateService(ex *transport.Executor) *AffiliateService {
	return &AffiliateService{newBase(ex)}
}

// AffiliateListQuery — shared documented list filters.
// Time format: "yyyy-MM-dd HH:mm:ss". PageSize caps per docs.
type AffiliateListQuery struct {
	StartTime    string
	EndTime      string
	PageNo       int
	PageSize     int
	UID          int64
	ReferralCode string
	ShowDetail   string // "yes"/"no" — activity endpoint only
	Type         *int   // team endpoint only: nil omits; 0 Futures,1 Spot,2 ETF,3 Futures API
	PartnerUIDs  []int64
	Quote        string // "usdt"/"usdc" — assets endpoint only
	ActivityID   string // activity endpoint only — required there
}

func (q AffiliateListQuery) params(maxPageSize int) (map[string]interface{}, error) {
	if err := requireAffiliateDateTime(q.StartTime, "startTime"); err != nil {
		return nil, err
	}
	if err := requireAffiliateDateTime(q.EndTime, "endTime"); err != nil {
		return nil, err
	}
	params := map[string]interface{}{}
	if q.StartTime != "" {
		params["startTime"] = q.StartTime
	}
	if q.EndTime != "" {
		params["endTime"] = q.EndTime
	}
	if q.PageNo > 0 {
		params["pageNo"] = q.PageNo
	}
	if q.PageSize > 0 {
		if _, err := requirePageSize(q.PageSize, maxPageSize); err != nil {
			return nil, err
		}
		params["pageSize"] = q.PageSize
	}
	if q.UID > 0 {
		params["uid"] = q.UID
	}
	if q.ReferralCode != "" {
		params["referralCode"] = q.ReferralCode
	}
	return params, nil
}

// GetInvitedUserList — GET .../invitation/user/list.
func (s *AffiliateService) GetInvitedUserList(ctx context.Context, q AffiliateListQuery) (*transport.Response[model.ListPage[model.InvitedUser]], error) {
	p, err := q.params(1000)
	if err != nil {
		return nil, err
	}
	return callTyped[model.ListPage[model.InvitedUser]](s.base, ctx, "common.affiliate.getInvitedUserList", p)
}

// GetInvitedUsersKycResults — GET .../invitation/user/kyc/list.
func (s *AffiliateService) GetInvitedUsersKycResults(ctx context.Context, q AffiliateListQuery) (*transport.Response[model.ListPage[model.InvitedUserKYC]], error) {
	p, err := q.params(1000)
	if err != nil {
		return nil, err
	}
	return callTyped[model.ListPage[model.InvitedUserKYC]](s.base, ctx, "common.affiliate.getInvitedUsersKycResults", p)
}

// GetInvitedUsersTradingVolume — POST .../invitation/user/volume/list.
//
// Documented quirk: a request body is REQUIRED on every call, even an empty
// JSON object {} — omitting it may fail signature verification. With no
// filters the SDK sends "{}" explicitly. Date format: "yyyy-MM-dd".
func (s *AffiliateService) GetInvitedUsersTradingVolume(ctx context.Context, q AffiliateListQuery) (*transport.Response[model.ListPage[model.InvitedUserVolume]], error) {
	if err := requireAffiliateDate(q.StartTime, "startTime"); err != nil {
		return nil, err
	}
	if err := requireAffiliateDate(q.EndTime, "endTime"); err != nil {
		return nil, err
	}
	params := map[string]interface{}{}
	if q.StartTime != "" {
		params["startTime"] = q.StartTime
	}
	if q.EndTime != "" {
		params["endTime"] = q.EndTime
	}
	if q.PageNo > 0 {
		params["pageNo"] = q.PageNo
	}
	if q.PageSize > 0 {
		if _, err := requirePageSize(q.PageSize, 1000); err != nil {
			return nil, err
		}
		params["pageSize"] = q.PageSize
	}
	if q.UID > 0 {
		params["uid"] = q.UID
	}
	return callTyped[model.ListPage[model.InvitedUserVolume]](s.base, ctx, "common.affiliate.getInvitedUsersTradingVolume", params)
}

// GetAffiliateUserActivityData — POST .../invitation/user/activity/list.
// activityID is mandatory; uses "page" (not pageNo) per docs.
func (s *AffiliateService) GetAffiliateUserActivityData(ctx context.Context, activityID string, q AffiliateListQuery) (*transport.Response[model.ListPage[model.AffiliateActivity]], error) {
	if _, err := requireNonEmpty(activityID, "activityId"); err != nil {
		return nil, err
	}
	if q.ShowDetail != "" && q.ShowDetail != "yes" && q.ShowDetail != "no" {
		return nil, validationErr(`showDetail must be "yes" or "no"`)
	}
	p, err := q.params(1000)
	if err != nil {
		return nil, err
	}
	p["activityId"] = activityID
	if q.PageNo > 0 {
		p["page"] = q.PageNo // this endpoint documents "page", not pageNo
		delete(p, "pageNo")
	}
	if q.ShowDetail != "" {
		p["showDetail"] = q.ShowDetail
	}
	return callTyped[model.ListPage[model.AffiliateActivity]](s.base, ctx, "common.affiliate.getAffiliateUserActivityData", p)
}

// GetInvitedUserAssets — GET .../invitation/user/assets.
// quote: "usdt" (default) or "usdc".
func (s *AffiliateService) GetInvitedUserAssets(ctx context.Context, q AffiliateListQuery) (*transport.Response[model.ListPage[model.InvitedUserAssets]], error) {
	if q.Quote != "" {
		lq := strings.ToLower(q.Quote)
		if lq != "usdt" && lq != "usdc" {
			return nil, validationErr(`quote must be "usdt" or "usdc"`)
		}
	}
	p, err := q.params(100)
	if err != nil {
		return nil, err
	}
	if q.Quote != "" {
		p["quote"] = strings.ToLower(q.Quote)
	}
	return callTyped[model.ListPage[model.InvitedUserAssets]](s.base, ctx, "common.affiliate.getInvitedUserAssets", p)
}

// GetAffiliateTeamDirectUserRange — POST .../team/users/range.
// Documented constraints: partnerUids max 100; type ∈ {0,1,2,3} (0 Futures,
// 1 Spot, 2 ETF, 3 Futures API) — nil Type omits the filter.
func (s *AffiliateService) GetAffiliateTeamDirectUserRange(ctx context.Context, q AffiliateListQuery) (*transport.Response[model.AffiliateTeam], error) {
	if q.Type != nil && (*q.Type < 0 || *q.Type > 3) {
		return nil, validationErr("type must be 0 (Futures), 1 (Spot), 2 (ETF) or 3 (Futures API)")
	}
	if len(q.PartnerUIDs) > 100 {
		return nil, validationErr("partnerUids supports at most 100 UIDs")
	}
	p, err := q.params(1000)
	if err != nil {
		return nil, err
	}
	if q.Type != nil {
		p["type"] = *q.Type
	}
	if len(q.PartnerUIDs) > 0 {
		p["partnerUids"] = q.PartnerUIDs
	}
	return callTyped[model.AffiliateTeam](s.base, ctx, "common.affiliate.getAffiliateTeamDirectUserRange", p)
}
