package dao

import (
	"database/sql"
	"time"

	"github.com/ShahoBashoki/kucoin/object"
	"github.com/google/uuid"
)

type (

	// DAOer is an interface.
	DAOer interface {
		DAOJoiner
		// GetID is a function.
		GetID() uuid.UUID
	}

	// DAOJoiner is an interface.
	DAOJoiner interface {
		object.GetMap
		// GetCreatedAt is a function.
		GetCreatedAt() time.Time
		// GetUpdatedAt is a function.
		GetUpdatedAt() time.Time
		// GetDeletedAt is a function.
		GetDeletedAt() sql.NullTime
	}

	dao struct {
		createdAt time.Time
		updatedAt time.Time
		deletedAt sql.NullTime
		id        uuid.UUID
	}

	daoJoin struct {
		createdAt time.Time
		updatedAt time.Time
		deletedAt sql.NullTime
	}
)
