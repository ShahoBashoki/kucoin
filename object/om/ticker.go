package om

import (
	"encoding/json"

	"github.com/google/uuid"
)

type (
	// Tickerer is an interface.
	Tickerer interface {
		OMer
		// GetAveragePrice is a function.
		GetAveragePrice() string
		// GetBuy is a function.
		GetBuy() string
		// GetChangePrice is a function.
		GetChangePrice() string
		// GetChangeRate is a function.
		GetChangeRate() string
		// GetHigh is a function.
		GetHigh() string
		// GetLast is a function.
		GetLast() string
		// GetLow is a function.
		GetLow() string
		// GetMakerCoefficient is a function.
		GetMakerCoefficient() string
		// GetMakerFeeRate is a function.
		GetMakerFeeRate() string
		// GetSell is a function.
		GetSell() string
		// GetSymbol is a function.
		GetSymbol() string
		// GetSymbolName is a function.
		GetSymbolName() string
		// GetTakerCoefficient is a function.
		GetTakerCoefficient() string
		// GetTakerFeeRate is a function.
		GetTakerFeeRate() string
		// GetVol is a function.
		GetVol() string
		// GetVolValue is a function.
		GetVolValue() string
	}

	ticker struct {
		averagePrice     string
		buy              string
		changePrice      string
		changeRate       string
		high             string
		last             string
		low              string
		makerCoefficient string
		makerFeeRate     string
		sell             string
		symbol           string
		symbolName       string
		takerCoefficient string
		takerFeeRate     string
		vol              string
		volValue         string
		id               uuid.UUID
	}
)

var _ Tickerer = (*ticker)(nil)

// NewTicker is a function.
func NewTicker(
	averagePrice string,
	buy string,
	changePrice string,
	changeRate string,
	high string,
	last string,
	low string,
	makerCoefficient string,
	makerFeeRate string,
	sell string,
	symbol string,
	symbolName string,
	takerCoefficient string,
	takerFeeRate string,
	vol string,
	volValue string,
	id uuid.UUID,
) *ticker {
	return &ticker{
		averagePrice:     averagePrice,
		buy:              buy,
		changePrice:      changePrice,
		changeRate:       changeRate,
		high:             high,
		last:             last,
		low:              low,
		makerCoefficient: makerCoefficient,
		makerFeeRate:     makerFeeRate,
		sell:             sell,
		symbol:           symbol,
		symbolName:       symbolName,
		takerCoefficient: takerCoefficient,
		takerFeeRate:     takerFeeRate,
		vol:              vol,
		volValue:         volValue,
		id:               id,
	}
}

// TickererComparer is a function.
func TickererComparer(
	first Tickerer,
	second Tickerer,
) bool {
	return OMerComparer(first, second) &&
		first.GetAveragePrice() == second.GetAveragePrice() &&
		first.GetBuy() == second.GetBuy() &&
		first.GetChangePrice() == second.GetChangePrice() &&
		first.GetChangeRate() == second.GetChangeRate() &&
		first.GetHigh() == second.GetHigh() &&
		first.GetLast() == second.GetLast() &&
		first.GetLow() == second.GetLow() &&
		first.GetMakerCoefficient() == second.GetMakerCoefficient() &&
		first.GetMakerFeeRate() == second.GetMakerFeeRate() &&
		first.GetSell() == second.GetSell() &&
		first.GetSymbol() == second.GetSymbol() &&
		first.GetSymbolName() == second.GetSymbolName() &&
		first.GetTakerCoefficient() == second.GetTakerCoefficient() &&
		first.GetTakerFeeRate() == second.GetTakerFeeRate() &&
		first.GetVol() == second.GetVol() &&
		first.GetVolValue() == second.GetVolValue()
}

// GetID is a function.
func (ticker *ticker) GetID() uuid.UUID {
	return ticker.id
}

// GetAveragePrice is a function.
func (ticker *ticker) GetAveragePrice() string {
	return ticker.averagePrice
}

// GetBuy is a function.
func (ticker *ticker) GetBuy() string {
	return ticker.buy
}

// GetChangePrice is a function.
func (ticker *ticker) GetChangePrice() string {
	return ticker.changePrice
}

// GetChangeRate is a function.
func (ticker *ticker) GetChangeRate() string {
	return ticker.changeRate
}

// GetHigh is a function.
func (ticker *ticker) GetHigh() string {
	return ticker.high
}

// GetLast is a function.
func (ticker *ticker) GetLast() string {
	return ticker.last
}

// GetLow is a function.
func (ticker *ticker) GetLow() string {
	return ticker.low
}

// GetMakerCoefficient is a function.
func (ticker *ticker) GetMakerCoefficient() string {
	return ticker.makerCoefficient
}

// GetMakerFeeRate is a function.
func (ticker *ticker) GetMakerFeeRate() string {
	return ticker.makerFeeRate
}

// GetSell is a function.
func (ticker *ticker) GetSell() string {
	return ticker.sell
}

// GetSymbol is a function.
func (ticker *ticker) GetSymbol() string {
	return ticker.symbol
}

// GetSymbolName is a function.
func (ticker *ticker) GetSymbolName() string {
	return ticker.symbolName
}

// GetTakerCoefficient is a function.
func (ticker *ticker) GetTakerCoefficient() string {
	return ticker.takerCoefficient
}

// GetTakerFeeRate is a function.
func (ticker *ticker) GetTakerFeeRate() string {
	return ticker.takerFeeRate
}

// GetVol is a function.
func (ticker *ticker) GetVol() string {
	return ticker.vol
}

// GetVolValue is a function.
func (ticker *ticker) GetVolValue() string {
	return ticker.volValue
}

// GetMap is a function.
func (ticker *ticker) GetMap() map[string]any {
	return map[string]any{
		"id":                ticker.GetID(),
		"average_price":     ticker.GetAveragePrice(),
		"buy":               ticker.GetBuy(),
		"change_price":      ticker.GetChangePrice(),
		"change_rate":       ticker.GetChangeRate(),
		"high":              ticker.GetHigh(),
		"last":              ticker.GetLast(),
		"low":               ticker.GetLow(),
		"maker_coefficient": ticker.GetMakerCoefficient(),
		"maker_fee_rate":    ticker.GetMakerFeeRate(),
		"sell":              ticker.GetSell(),
		"symbol":            ticker.GetSymbol(),
		"symbol_name":       ticker.GetSymbolName(),
		"taker_coefficient": ticker.GetTakerCoefficient(),
		"taker_fee_rate":    ticker.GetTakerFeeRate(),
		"vol":               ticker.GetVol(),
		"vol_value":         ticker.GetVolValue(),
	}
}

// MarshalJSON is a function.
// read more https://pkg.go.dev/encoding/json#Marshaler
func (ticker *ticker) MarshalJSON() ([]byte, error) {
	return json.Marshal(ticker.GetMap())
}
