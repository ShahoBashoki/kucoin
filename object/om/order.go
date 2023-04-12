package om

import (
	"encoding/json"

	"github.com/google/uuid"
)

type (
	// Orderer is an interface.
	Orderer interface {
		OMer
		// GetChannel is a function.
		GetChannel() string
		// GetClientOID is a function.
		GetClientOID() string
		// GetDealFunds is a function.
		GetDealFunds() string
		// GetDealSize is a function.
		GetDealSize() string
		// GetFee is a function.
		GetFee() string
		// GetFeeCurrency is a function.
		GetFeeCurrency() string
		// GetFunds is a function.
		GetFunds() string
		// GetKucoinID is a function.
		GetKucoinID() string
		// GetKucoinType is a function.
		GetKucoinType() string
		// GetOPType is a function.
		GetOPType() string
		// GetPrice is a function.
		GetPrice() string
		// GetRemark is a function.
		GetRemark() string
		// GetSide is a function.
		GetSide() string
		// GetSize is a function.
		GetSize() string
		// GetStop is a function.
		GetStop() string
		// GetStopPrice is a function.
		GetStopPrice() string
		// GetSTP is a function.
		GetSTP() string
		// GetSymbol is a function.
		GetSymbol() string
		// GetTags is a function.
		GetTags() string
		// GetTimeInForce is a function.
		GetTimeInForce() string
		// GetTradeType is a function.
		GetTradeType() string
		// GetVisibleSize is a function.
		GetVisibleSize() string
		// GetCancelAfter is a function.
		GetCancelAfter() uint32
		// GetKucoinCreatedAt is a function.
		GetKucoinCreatedAt() uint32
		// GetCancelExist is a function.
		GetCancelExist() bool
		// GetHidden is a function.
		GetHidden() bool
		// GetIceBerg is a function.
		GetIceBerg() bool
		// GetIsActive is a function.
		GetIsActive() bool
		// GetPostOnly is a function.
		GetPostOnly() bool
		// GetStopTriggered is a function.
		GetStopTriggered() bool
	}

	order struct {
		channel         string
		clientOID       string
		dealFunds       string
		dealSize        string
		fee             string
		feeCurrency     string
		funds           string
		kucoinID        string
		kucoinType      string
		opType          string
		price           string
		remark          string
		side            string
		size            string
		stop            string
		stopPrice       string
		stp             string
		symbol          string
		tags            string
		timeInForce     string
		tradeType       string
		visibleSize     string
		cancelAfter     uint32
		kucoinCreatedAt uint32
		cancelExist     bool
		hidden          bool
		iceBerg         bool
		isActive        bool
		postOnly        bool
		stopTriggered   bool
		id              uuid.UUID
	}
)

var _ Orderer = (*order)(nil)

// NewOrder is a function.
func NewOrder(
	channel string,
	clientOID string,
	dealFunds string,
	dealSize string,
	fee string,
	feeCurrency string,
	funds string,
	kucoinID string,
	kucoinType string,
	opType string,
	price string,
	remark string,
	side string,
	size string,
	stop string,
	stopPrice string,
	stp string,
	symbol string,
	tags string,
	timeInForce string,
	tradeType string,
	visibleSize string,
	cancelAfter uint32,
	kucoinCreatedAt uint32,
	cancelExist bool,
	hidden bool,
	iceBerg bool,
	isActive bool,
	postOnly bool,
	stopTriggered bool,
	id uuid.UUID,
) *order {
	return &order{
		channel:         channel,
		clientOID:       clientOID,
		dealFunds:       dealFunds,
		dealSize:        dealSize,
		fee:             fee,
		feeCurrency:     feeCurrency,
		funds:           funds,
		kucoinID:        kucoinID,
		kucoinType:      kucoinType,
		opType:          opType,
		price:           price,
		remark:          remark,
		side:            side,
		size:            size,
		stop:            stop,
		stopPrice:       stopPrice,
		stp:             stp,
		symbol:          symbol,
		tags:            tags,
		timeInForce:     timeInForce,
		tradeType:       tradeType,
		visibleSize:     visibleSize,
		cancelAfter:     cancelAfter,
		kucoinCreatedAt: kucoinCreatedAt,
		cancelExist:     cancelExist,
		hidden:          hidden,
		iceBerg:         iceBerg,
		isActive:        isActive,
		postOnly:        postOnly,
		stopTriggered:   stopTriggered,
		id:              id,
	}
}

// OrdererComparer is a function.
func OrdererComparer(
	first Orderer,
	second Orderer,
) bool {
	return OMerComparer(first, second) &&
		first.GetChannel() == second.GetChannel() &&
		first.GetClientOID() == second.GetClientOID() &&
		first.GetDealFunds() == second.GetDealFunds() &&
		first.GetDealSize() == second.GetDealSize() &&
		first.GetFee() == second.GetFee() &&
		first.GetFeeCurrency() == second.GetFeeCurrency() &&
		first.GetFunds() == second.GetFunds() &&
		first.GetKucoinID() == second.GetKucoinID() &&
		first.GetKucoinType() == second.GetKucoinType() &&
		first.GetOPType() == second.GetOPType() &&
		first.GetPrice() == second.GetPrice() &&
		first.GetRemark() == second.GetRemark() &&
		first.GetSide() == second.GetSide() &&
		first.GetSize() == second.GetSize() &&
		first.GetStop() == second.GetStop() &&
		first.GetStopPrice() == second.GetStopPrice() &&
		first.GetSTP() == second.GetSTP() &&
		first.GetSymbol() == second.GetSymbol() &&
		first.GetTags() == second.GetTags() &&
		first.GetTimeInForce() == second.GetTimeInForce() &&
		first.GetTradeType() == second.GetTradeType() &&
		first.GetVisibleSize() == second.GetVisibleSize() &&
		first.GetCancelAfter() == second.GetCancelAfter() &&
		first.GetKucoinCreatedAt() == second.GetKucoinCreatedAt() &&
		first.GetCancelExist() == second.GetCancelExist() &&
		first.GetHidden() == second.GetHidden() &&
		first.GetIceBerg() == second.GetIceBerg() &&
		first.GetIsActive() == second.GetIsActive() &&
		first.GetPostOnly() == second.GetPostOnly() &&
		first.GetStopTriggered() == second.GetStopTriggered()
}

// GetID is a function.
func (order *order) GetID() uuid.UUID {
	return order.id
}

// GetChannel is a function.
func (order *order) GetChannel() string {
	return order.channel
}

// GetClientOID is a function.
func (order *order) GetClientOID() string {
	return order.clientOID
}

// GetDealFunds is a function.
func (order *order) GetDealFunds() string {
	return order.dealFunds
}

// GetDealSize is a function.
func (order *order) GetDealSize() string {
	return order.dealSize
}

// GetFee is a function.
func (order *order) GetFee() string {
	return order.fee
}

// GetFeeCurrency is a function.
func (order *order) GetFeeCurrency() string {
	return order.feeCurrency
}

// GetFunds is a function.
func (order *order) GetFunds() string {
	return order.funds
}

// GetKucoinID is a function.
func (order *order) GetKucoinID() string {
	return order.kucoinID
}

// GetKucoinType is a function.
func (order *order) GetKucoinType() string {
	return order.kucoinType
}

// GetOPType is a function.
func (order *order) GetOPType() string {
	return order.opType
}

// GetPrice is a function.
func (order *order) GetPrice() string {
	return order.price
}

// GetRemark is a function.
func (order *order) GetRemark() string {
	return order.remark
}

// GetSide is a function.
func (order *order) GetSide() string {
	return order.side
}

// GetSize is a function.
func (order *order) GetSize() string {
	return order.size
}

// GetStop is a function.
func (order *order) GetStop() string {
	return order.stop
}

// GetStopPrice is a function.
func (order *order) GetStopPrice() string {
	return order.stopPrice
}

// GetSTP is a function.
func (order *order) GetSTP() string {
	return order.stp
}

// GetSymbol is a function.
func (order *order) GetSymbol() string {
	return order.symbol
}

// GetTags is a function.
func (order *order) GetTags() string {
	return order.tags
}

// GetTimeInForce is a function.
func (order *order) GetTimeInForce() string {
	return order.timeInForce
}

// GetTradeType is a function.
func (order *order) GetTradeType() string {
	return order.tradeType
}

// GetVisibleSize is a function.
func (order *order) GetVisibleSize() string {
	return order.visibleSize
}

// GetCancelAfter is a function.
func (order *order) GetCancelAfter() uint32 {
	return order.cancelAfter
}

// GetKucoinCreatedAt is a function.
func (order *order) GetKucoinCreatedAt() uint32 {
	return order.kucoinCreatedAt
}

// GetCancelExist is a function.
func (order *order) GetCancelExist() bool {
	return order.cancelExist
}

// GetHidden is a function.
func (order *order) GetHidden() bool {
	return order.hidden
}

// GetIceBerg is a function.
func (order *order) GetIceBerg() bool {
	return order.iceBerg
}

// GetIsActive is a function.
func (order *order) GetIsActive() bool {
	return order.isActive
}

// GetPostOnly is a function.
func (order *order) GetPostOnly() bool {
	return order.postOnly
}

// GetStopTriggered is a function.
func (order *order) GetStopTriggered() bool {
	return order.stopTriggered
}

// GetMap is a function.
func (order *order) GetMap() map[string]any {
	return map[string]any{
		"id":                order.GetID(),
		"channel":           order.GetChannel(),
		"client_oid":        order.GetClientOID(),
		"deal_funds":        order.GetDealFunds(),
		"deal_size":         order.GetDealSize(),
		"fee":               order.GetFee(),
		"fee_currency":      order.GetFeeCurrency(),
		"funds":             order.GetFunds(),
		"kucoin_id":         order.GetKucoinID(),
		"kucoin_type":       order.GetKucoinType(),
		"op_type":           order.GetOPType(),
		"price":             order.GetPrice(),
		"remark":            order.GetRemark(),
		"side":              order.GetSide(),
		"size":              order.GetSize(),
		"stop":              order.GetStop(),
		"stop_price":        order.GetStopPrice(),
		"stp":               order.GetSTP(),
		"symbol":            order.GetSymbol(),
		"tags":              order.GetTags(),
		"time_in_force":     order.GetTimeInForce(),
		"trade_type":        order.GetTradeType(),
		"visible_size":      order.GetVisibleSize(),
		"cancel_after":      order.GetCancelAfter(),
		"kucoin_created_at": order.GetKucoinCreatedAt(),
		"cancel_exist":      order.GetCancelExist(),
		"hidden":            order.GetHidden(),
		"ice_berg":          order.GetIceBerg(),
		"is_active":         order.GetIsActive(),
		"post_only":         order.GetPostOnly(),
		"stop_triggered":    order.GetStopTriggered(),
	}
}

// MarshalJSON is a function.
// read more https://pkg.go.dev/encoding/json#Marshaler
func (order *order) MarshalJSON() ([]byte, error) {
	return json.Marshal(order.GetMap())
}
