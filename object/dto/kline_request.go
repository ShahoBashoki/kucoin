package dto

import (
	"encoding/json"

	"github.com/ShahoBashoki/kucoin/object"
)

type (
	// KlineRequester is an interface.
	KlineRequester interface {
		// GetKlineType is a function.
		GetKlineType() object.KlineTypeType
		// GetSymbol is a function.
		GetSymbol() string
		// GetEndAt is a function.
		GetEndAt() int64
		// GetStartAt is a function.
		GetStartAt() int64
	}

	klineRequest struct {
		klineType object.KlineTypeType
		symbol    string
		endAt     int64
		startAt   int64
	}
)

var (
	_ KlineRequester = (*klineRequest)(nil)
	_ json.Marshaler = (*klineRequest)(nil)
	_ object.GetMap  = (*klineRequest)(nil)
)

// NewKlineRequest is a function.
func NewKlineRequest(
	klineType object.KlineTypeType,
	symbol string,
	endAt int64,
	startAt int64,
) *klineRequest {
	return &klineRequest{
		klineType: klineType,
		symbol:    symbol,
		endAt:     endAt,
		startAt:   startAt,
	}
}

// KlineRequesterComparer is a function.
func KlineRequesterComparer(
	first KlineRequester,
	second KlineRequester,
) bool {
	return first.GetKlineType() == second.GetKlineType() &&
		first.GetSymbol() == second.GetSymbol() &&
		first.GetEndAt() == second.GetEndAt() &&
		first.GetStartAt() == second.GetStartAt()
}

// GetKlineType is a function.
func (klineRequest *klineRequest) GetKlineType() object.KlineTypeType {
	return klineRequest.klineType
}

// GetSymbol is a function.
func (klineRequest *klineRequest) GetSymbol() string {
	return klineRequest.symbol
}

// GetEndAt is a function.
func (klineRequest *klineRequest) GetEndAt() int64 {
	return klineRequest.endAt
}

// GetStartAt is a function.
func (klineRequest *klineRequest) GetStartAt() int64 {
	return klineRequest.startAt
}

// GetMap is a function.
func (klineRequest *klineRequest) GetMap() map[string]any {
	return map[string]any{
		"type":    string(klineRequest.GetKlineType()),
		"symbol":  klineRequest.GetSymbol(),
		"endAt":   klineRequest.GetEndAt(),
		"startAt": klineRequest.GetStartAt(),
	}
}

// MarshalJSON is a function.
// read more https://pkg.go.dev/encoding/json#Marshaler
func (klineRequest *klineRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(klineRequest.GetMap())
}
