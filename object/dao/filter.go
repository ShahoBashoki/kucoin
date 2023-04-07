package dao

import "gorm.io/gorm"

// Filterer is an interface.
type Filterer interface {
	// Filter is a function.
	Filter(
		*gorm.DB,
	) *gorm.DB
}
