package dao

import (
	"encoding/json"

	"github.com/ShahoBashoki/kucoin/object"
	"gorm.io/gorm"
)

type (

	// OrderFilterer is an interface.
	OrderFilterer interface {
		Filterer
		// GetIsActive is a function.
		GetIsActive() bool
	}

	orderFilter struct {
		isActive bool
	}
)

var (
	_ OrderFilterer  = (*orderFilter)(nil)
	_ json.Marshaler = (*orderFilter)(nil)
	_ object.GetMap  = (*orderFilter)(nil)
)

// NewOrderFilter is a function.
func NewOrderFilter(
	isActive bool,
) *orderFilter {
	return &orderFilter{
		isActive: isActive,
	}
}

// GetIsActive is a function.
func (filter *orderFilter) GetIsActive() bool {
	return filter.isActive
}

// GetMap is a function.
func (filter *orderFilter) GetMap() map[string]any {
	return map[string]any{
		"is_active": filter.GetIsActive(),
	}
}

// MarshalJSON is a function.
// read more https://pkg.go.dev/encoding/json#Marshaler
func (filter *orderFilter) MarshalJSON() ([]byte, error) {
	return json.Marshal(filter.GetMap())
}

// Filter is a function.
func (filter *orderFilter) Filter(
	gormDB *gorm.DB,
) *gorm.DB {
	if filter.GetIsActive() {
		gormDB.Where("is_active = ?", filter.GetIsActive())
	}

	return gormDB
}
