package dao

import (
	"encoding/json"

	"github.com/ShahoBashoki/kucoin/object"
	"gorm.io/gorm"
)

type (
	// Paginationer is an interface.
	Paginationer interface {
		// GetCursorer is a function.
		GetCursorer() Cursorer
		// GetLimit is a function.
		GetLimit() uint32
		// Pagination is a function.
		Pagination(
			table string,
		) func(*gorm.DB) *gorm.DB
	}

	pagination struct {
		cursorer Cursorer
		limit    uint32
	}
)

var (
	_ Paginationer   = (*pagination)(nil)
	_ json.Marshaler = (*pagination)(nil)
	_ object.GetMap  = (*pagination)(nil)
)

// NewPagination is a function.
func NewPagination(
	cursorer Cursorer,
	limit uint32,
) *pagination {
	return &pagination{
		cursorer: cursorer,
		limit:    limit,
	}
}

// GetCursorer is a function.
func (pagination *pagination) GetCursorer() Cursorer {
	return pagination.cursorer
}

// GetLimit is a function.
func (pagination *pagination) GetLimit() uint32 {
	return pagination.limit
}

// GetMap is a function.
func (pagination *pagination) GetMap() map[string]any {
	return map[string]any{
		"cursorer": pagination.GetCursorer(),
		"limit":    pagination.GetLimit(),
	}
}

// MarshalJSON is a function.
// read more https://pkg.go.dev/encoding/json#Marshaler
func (pagination *pagination) MarshalJSON() ([]byte, error) {
	return json.Marshal(pagination.GetMap())
}

// Pagination is a function.
func (pagination *pagination) Pagination(
	table string,
) func(*gorm.DB) *gorm.DB {
	return func(
		gormDB *gorm.DB,
	) *gorm.DB {
		if pagination.GetCursorer() != nil {
			gormDB.Scopes(pagination.GetCursorer().Query(table))
		}

		if pagination.GetLimit() != 0 {
			gormDB.Limit(int(pagination.GetLimit()) + 1)
		}

		return gormDB
	}
}
