package repository

import (
	"context"
	"time"

	"github.com/ShahoBashoki/kucoin/object/dao"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	// DAORepositorier is an interface.
	DAORepositorier[
		daoDAOer dao.DAOer,
		daoFilterer dao.Filterer,
	] interface {
		// Create is a function.
		Create(
			context.Context,
			daoDAOer,
		) (uuid.UUID, error)
		// Read is a function.
		Read(
			context.Context,
			uuid.UUID,
		) (daoDAOer, error)
		// ReadList is a function.
		ReadList(
			context.Context,
			dao.Paginationer,
			daoFilterer,
		) ([]daoDAOer, dao.Cursorer, error)
		// Update is a function.
		Update(
			context.Context,
			daoDAOer,
		) (time.Time, error)
		// Delete is a function.
		Delete(
			context.Context,
			uuid.UUID,
		) (time.Time, error)
	}

	// DAOJoinRepositorier is an interface.
	DAOJoinRepositorier[
		daoDAOJoiner dao.DAOJoiner,
		daoFilterer dao.Filterer,
	] interface {
		// Create is a function.
		Create(
			context.Context,
			daoDAOJoiner,
		) error
		// Read is a function.
		Read(
			context.Context,
			uuid.UUID,
			uuid.UUID,
		) (daoDAOJoiner, error)
		// ReadList is a function.
		ReadList(
			context.Context,
			dao.Paginationer,
			daoFilterer,
		) ([]daoDAOJoiner, dao.Cursorer, error)
		// Update is a function.
		Update(
			context.Context,
			daoDAOJoiner,
		) (time.Time, error)
		// Delete is a function.
		Delete(
			context.Context,
			uuid.UUID,
			uuid.UUID,
		) (time.Time, error)
	}

	// GetDB is an interface.
	GetDB interface {
		GetDB() *gorm.DB
	}

	// Repositorier is an interface.
	Repositorier any

	// GetRepositorier is an interface.
	GetRepositorier interface {
		// GetRepositorier is a function.
		GetRepositorier() Repositorier
	}

	repository struct{}

	optionRepositorier interface {
		apply(*repository)
	}

	optionRepositorierFunc func(*repository)
)

var _ Repositorier = (*repository)(nil)

// NewRepository is a function.
func NewRepository(
	optioners ...optionRepositorier,
) *repository {
	repository := &repository{}

	return repository.WithOptioners(optioners...)
}

// WithOptioners is a function.
func (repository *repository) WithOptioners(
	optioners ...optionRepositorier,
) *repository {
	newRepository := repository.clone()
	for _, optioner := range optioners {
		optioner.apply(newRepository)
	}

	return newRepository
}

func (repository *repository) clone() *repository {
	newRepository := repository

	return newRepository
}

func (optionerFunc optionRepositorierFunc) apply(
	repository *repository,
) {
	optionerFunc(repository)
}
