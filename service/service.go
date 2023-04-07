package service

import (
	"github.com/ShahoBashoki/kucoin/config"
	"github.com/ShahoBashoki/kucoin/log"
	"github.com/ShahoBashoki/kucoin/repository"
	"github.com/ShahoBashoki/kucoin/util"
	"go.opentelemetry.io/otel/trace"
)

type (
	// Servicer is an interface.
	Servicer any

	// GetServicer is an interface.
	GetServicer interface {
		// GetServicer is a function.
		GetServicer() Servicer
	}

	// WithServicer is an interface.
	WithServicer interface {
		// WithServicer is a function.
		WithServicer(
			Servicer,
		)
	}

	service struct{}
)

var _ Servicer = (*service)(nil)

// NewServicer is a function.
func NewServicer(
	configConfigger config.Configger,
	repositorier repository.Repositorier,
	logRuntimeLogger log.RuntimeLogger,
	traceTracer trace.Tracer,
	utilUUIDer util.UUIDer,
) Servicer {
	service := &service{}

	return service
}
