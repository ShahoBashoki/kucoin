package service

import (
	"github.com/Kucoin/kucoin-go-sdk"
	"github.com/ShahoBashoki/kucoin/config"
	"github.com/ShahoBashoki/kucoin/log"
	"github.com/ShahoBashoki/kucoin/repository"
	"github.com/ShahoBashoki/kucoin/util"
	"go.opentelemetry.io/otel/trace"
)

type (
	// Servicer is an interface.
	Servicer interface {
		GetOrderServicer
	}

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

	service struct {
		orderServicer OrderServicer
	}
)

var _ Servicer = (*service)(nil)

// NewServicer is a function.
func NewServicer(
	configConfigger config.Configger,
	repositorier repository.Repositorier,
	logRuntimeLogger log.RuntimeLogger,
	traceTracer trace.Tracer,
	utilUUIDer util.UUIDer,
	kucoinAPIService *kucoin.ApiService,
) Servicer {
	orderServicer := NewOrderServicer(
		configConfigger,
		repositorier.GetOrderRepositorier(),
		logRuntimeLogger,
		traceTracer,
		utilUUIDer,
		kucoinAPIService,
	)

	service := &service{
		orderServicer: orderServicer,
	}

	orderServicerWithTypeCheck, ok := orderServicer.(WithServicer)
	if ok {
		orderServicerWithTypeCheck.WithServicer(service)
	}

	return service
}

// GetOrderServicer is a function.
func (service *service) GetOrderServicer() OrderServicer {
	return service.orderServicer
}
