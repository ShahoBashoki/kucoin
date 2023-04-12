package dao

import (
	"database/sql"
	"time"

	"github.com/ShahoBashoki/kucoin/object"
	"github.com/google/uuid"
)

type (

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

	// DAOer is an interface.
	DAOer interface {
		DAOJoiner
		// GetID is a function.
		GetID() uuid.UUID
	}

	daoJoin struct {
		createdAt time.Time
		updatedAt time.Time
		deletedAt sql.NullTime
	}

	dao struct {
		daoJoin
		id uuid.UUID
	}
)

// DAOJoinerComparer is a function.
func DAOJoinerComparer(
	first DAOJoiner,
	second DAOJoiner,
) bool {
	return first.GetCreatedAt().Equal(second.GetCreatedAt()) &&
		first.GetUpdatedAt().Equal(second.GetUpdatedAt()) &&
		first.GetDeletedAt() == second.GetDeletedAt()
}

// DAOerComparer is a function.
func DAOerComparer(
	first DAOer,
	second DAOer,
) bool {
	return DAOJoinerComparer(first, second) &&
		first.GetID() == second.GetID()
}
