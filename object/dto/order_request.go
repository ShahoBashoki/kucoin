package dto

import (
	"encoding/json"

	"github.com/ShahoBashoki/kucoin/object"
)

type (
	// OrderRequester is an interface.
	OrderRequester interface {
		// GetEndAt is a function.
		GetEndAt() string
		// GetMap is a function.
		GetMap() map[string]any
		// GetOrderType is a function.
		GetOrderType() object.OrderTypeType
		// GetSide is a function.
		GetSide() object.OrderSideType
		// GetStartAt is a function.
		GetStartAt() string
		// GetStatus is a function.
		GetStatus() object.OrderStateType
		// GetSymbol is a function.
		GetSymbol() string
		// GetTradeType is a function.
		GetTradeType() object.OrderTypeType
	}

	orderRequest struct {
		endAt     string
		orderType object.OrderTypeType
		side      object.OrderSideType
		startAt   string
		status    object.OrderStateType
		symbol    string
		tradeType object.OrderTypeType
	}
)

var (
	_ OrderRequester = (*orderRequest)(nil)
	_ json.Marshaler = (*orderRequest)(nil)
	_ object.GetMap  = (*orderRequest)(nil)
)

// NewOrderRequest is a function.
func NewOrderRequest(
	endAt string,
	orderType object.OrderTypeType,
	side object.OrderSideType,
	startAt string,
	status object.OrderStateType,
	symbol string,
	tradeType object.OrderTypeType,
) *orderRequest {
	return &orderRequest{
		endAt:     endAt,
		orderType: orderType,
		side:      side,
		startAt:   startAt,
		status:    status,
		symbol:    symbol,
		tradeType: tradeType,
	}
}

// OrderRequesterComparer is a function.
func OrderRequesterComparer(
	first OrderRequester,
	second OrderRequester,
) bool {
	return first.GetEndAt() == second.GetEndAt() &&
		first.GetOrderType() == second.GetOrderType() &&
		first.GetSide() == second.GetSide() &&
		first.GetStartAt() == second.GetStartAt() &&
		first.GetStatus() == second.GetStatus() &&
		first.GetSymbol() == second.GetSymbol() &&
		first.GetTradeType() == second.GetTradeType()
}

// GetEndAt is a function.
func (orderRequest *orderRequest) GetEndAt() string {
	return orderRequest.endAt
}

// GetOrderType is a function.
func (orderRequest *orderRequest) GetOrderType() object.OrderTypeType {
	return orderRequest.orderType
}

// GetSide is a function.
func (orderRequest *orderRequest) GetSide() object.OrderSideType {
	return orderRequest.side
}

// GetStartAt is a function.
func (orderRequest *orderRequest) GetStartAt() string {
	return orderRequest.startAt
}

// GetStatus is a function.
func (orderRequest *orderRequest) GetStatus() object.OrderStateType {
	return orderRequest.status
}

// GetSymbol is a function.
func (orderRequest *orderRequest) GetSymbol() string {
	return orderRequest.symbol
}

// GetTradeType is a function.
func (orderRequest *orderRequest) GetTradeType() object.OrderTypeType {
	return orderRequest.tradeType
}

// GetMap is a function.
func (orderRequest *orderRequest) GetMap() map[string]any {
	return map[string]any{
		"endAt":     orderRequest.GetEndAt(),
		"type":      string(orderRequest.GetOrderType()),
		"side":      string(orderRequest.GetSide()),
		"startAt":   orderRequest.GetStartAt(),
		"status":    string(orderRequest.GetStatus()),
		"symbol":    orderRequest.GetSymbol(),
		"tradeType": string(orderRequest.GetTradeType()),
	}
}

// MarshalJSON is a function.
// read more https://pkg.go.dev/encoding/json#Marshaler
func (orderRequest *orderRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(orderRequest.GetMap())
}
