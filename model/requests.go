package model

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Direction string

const (
	Long  Direction = "long"
	Short Direction = "short"
)

type QuantityUnit int

const (
	QuoteQuantity QuantityUnit = iota
	ContractQuantity
	BaseQuantity
)

type PositionModel int

const (
	IsolatedMargin PositionModel = iota
	CrossMargin
)

type OrderType string

const (
	LimitOrder               OrderType = "plan"
	TriggerOrder             OrderType = "planTrigger"
	MarketOrder              OrderType = "execute"
	BBOOrder                 OrderType = "bbo"
	IOCOrder                 OrderType = "IOC"
	FOKOrder                 OrderType = "FOK"
	PostOnlyOrder            OrderType = "PostOnly"
	TrailingOrderQuery       OrderType = "moveStop"
	SLTPOrderQuery           OrderType = "stopProfitLoss"
	TrailingSLTPCurrentQuery OrderType = "moveStopProfitLoss"
)

// Instrument keeps USDT base identifiers and USDC pair identifiers distinct.
// Constructors validate the published naming distinction, without converting
// quantities or looking up an assumed list of listed instruments.
type Instrument struct {
	name       string
	settlement SettlementCurrency
}
type SettlementCurrency string

const (
	USDT SettlementCurrency = "USDT"
	USDC SettlementCurrency = "USDC"
)

func USDTInstrument(base string) (Instrument, error) {
	if strings.TrimSpace(base) == "" || strings.Contains(base, "_") || base != strings.TrimSpace(base) {
		return Instrument{}, fmt.Errorf("USDT instrument requires a base identifier")
	}
	return Instrument{name: base, settlement: USDT}, nil
}
func USDCInstrument(pair string) (Instrument, error) {
	if !strings.HasSuffix(strings.ToUpper(pair), "_USDC") || len(pair) <= 5 || strings.Count(pair, "_") != 1 || pair != strings.TrimSpace(pair) {
		return Instrument{}, fmt.Errorf("USDC instrument requires BASE_USDC")
	}
	return Instrument{name: pair, settlement: USDC}, nil
}
func (i Instrument) String() string                 { return i.name }
func (i Instrument) Settlement() SettlementCurrency { return i.settlement }
func (i Instrument) MarshalJSON() ([]byte, error) {
	if i.name == "" || (i.settlement != USDT && i.settlement != USDC) {
		return nil, fmt.Errorf("uninitialized instrument")
	}
	return json.Marshal(i.name)
}

// FuturesOrderRequest is the published order contract. Optional decimals are
// pointers so omission is distinct from an explicit zero. All money is strict
// Decimal, including quantities measured in contracts.
type FuturesOrderRequest struct {
	Instrument      Instrument    `json:"instrument"`
	Direction       Direction     `json:"direction"`
	Leverage        int           `json:"leverage"`
	QuantityUnit    QuantityUnit  `json:"quantityUnit"`
	Quantity        Decimal       `json:"quantity"`
	PositionModel   PositionModel `json:"positionModel"`
	PositionType    OrderType     `json:"positionType"`
	OpenPrice       *Decimal      `json:"openPrice,omitempty"`
	StopLossPrice   *Decimal      `json:"stopLossPrice,omitempty"`
	StopProfitPrice *Decimal      `json:"stopProfitPrice,omitempty"`
	TriggerPrice    *Decimal      `json:"triggerPrice,omitempty"`
	TriggerType     *int          `json:"triggerType,omitempty"`
	GoldID          *int64        `json:"goldId,omitempty"`
	ThirdOrderID    string        `json:"thirdOrderId,omitempty"`
	UseAlmightyGold *bool         `json:"useAlmightyGold,omitempty"`
}

type ModifyOrderRequest struct {
	FuturesOrderRequest
	ID          string   `json:"id"`
	FeeType     *int     `json:"feeType,omitempty"`
	OrderPrice  *Decimal `json:"orderPrice,omitempty"`
	OrderBookID string   `json:"orderBookId,omitempty"`
}

type CloseBatchPositionRequest struct {
	ThirdOrderID string `json:"thirdOrderId"`
}

// SLTPPrices contains only the fields published on the single SL/TP page.
type SLTPPrices struct {
	StopLossOrderPrice   *Decimal `json:"stopLossOrderPrice,omitempty"`
	StopProfitOrderPrice *Decimal `json:"stopProfitOrderPrice,omitempty"`
	StopLossPrice        *Decimal `json:"stopLossPrice,omitempty"`
	StopProfitPrice      *Decimal `json:"stopProfitPrice,omitempty"`
	StopLossRate         *Decimal `json:"stopLossRate,omitempty"`
	StopProfitRate       *Decimal `json:"stopProfitRate,omitempty"`
}

type BatchSLTPRequest struct {
	SLTPPrices
	StopFrom   int      `json:"stopFrom"`
	PriceType  *int     `json:"priceType,omitempty"`
	StopType   *int     `json:"stopType,omitempty"`
	ClosePiece *Decimal `json:"closePiece,omitempty"`
}

type ModifyBatchSLTPRequest struct {
	ID                   string      `json:"id"`
	Instrument           *Instrument `json:"instrument,omitempty"`
	StopProfitOrderPrice *Decimal    `json:"stopProfitOrderPrice,omitempty"`
	StopLossPrice        *Decimal    `json:"stopLossPrice,omitempty"`
	StopProfitPrice      *Decimal    `json:"stopProfitPrice,omitempty"`
	StopLossRate         *Decimal    `json:"stopLossRate,omitempty"`
	StopProfitRate       *Decimal    `json:"stopProfitRate,omitempty"`
	StopFrom             int         `json:"stopFrom"`
	PriceType            *int        `json:"priceType,omitempty"`
	StopType             *int        `json:"stopType,omitempty"`
	ClosePiece           *Decimal    `json:"closePiece,omitempty"`
}
