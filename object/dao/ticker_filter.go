package dao

import (
	"encoding/json"

	"github.com/ShahoBashoki/kucoin/object"
	"gorm.io/gorm"
)

type (

	// TickerFilterer is an interface.
	TickerFilterer interface {
		Filterer
		// GetSymbol is a function.
		GetSymbol() string
		// GetSortChangeRate is a function.
		GetSortChangeRate() bool
	}

	tickerFilter struct {
		symbol         string
		sortChangeRate bool
	}
)

var (
	_ TickerFilterer = (*tickerFilter)(nil)
	_ json.Marshaler = (*tickerFilter)(nil)
	_ object.GetMap  = (*tickerFilter)(nil)
)

// NewTickerFilter is a function.
func NewTickerFilter(
	symbol string,
	sortChangeRate bool,
) *tickerFilter {
	return &tickerFilter{
		symbol:         symbol,
		sortChangeRate: sortChangeRate,
	}
}

// GetSymbol is a function.
func (filter *tickerFilter) GetSymbol() string {
	return filter.symbol
}

// GetSortChangeRate is a function.
func (filter *tickerFilter) GetSortChangeRate() bool {
	return filter.sortChangeRate
}

// GetMap is a function.
func (filter *tickerFilter) GetMap() map[string]any {
	return map[string]any{
		"symbol":           filter.GetSymbol(),
		"sort_change_rate": filter.GetSortChangeRate(),
	}
}

// MarshalJSON is a function.
// read more https://pkg.go.dev/encoding/json#Marshaler
func (filter *tickerFilter) MarshalJSON() ([]byte, error) {
	return json.Marshal(filter.GetMap())
}

// Filter is a function.
func (filter *tickerFilter) Filter(
	gormDB *gorm.DB,
) *gorm.DB {
	if filter.GetSymbol() != object.URIEmpty {
		gormDB.Where("symbol = ?", filter.GetSymbol())
	}

	if filter.GetSortChangeRate() {
		gormDB.Where("symbol LIKE '%-USDT' ORDER BY change_rate DESC")
	}

	return gormDB
}
