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
		GetKlineServicer
		GetOrderBookServicer
		GetOrderServicer
		GetTickerServicer
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
		klineServicer     KlineServicer
		orderBookServicer OrderBookServicer
		orderServicer     OrderServicer
		tickerServicer    TickerServicer
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
	klineServicer := NewKlineServicer(
		configConfigger,
		logRuntimeLogger,
		traceTracer,
		utilUUIDer,
		kucoinAPIService,
	)

	orderBookServicer := NewOrderBookServicer(
		configConfigger,
		logRuntimeLogger,
		traceTracer,
		utilUUIDer,
		kucoinAPIService,
	)

	orderServicer := NewOrderServicer(
		configConfigger,
		repositorier.GetOrderRepositorier(),
		logRuntimeLogger,
		traceTracer,
		utilUUIDer,
		kucoinAPIService,
	)

	tickerServicer := NewTickerServicer(
		configConfigger,
		repositorier.GetTickerRepositorier(),
		logRuntimeLogger,
		traceTracer,
		utilUUIDer,
		kucoinAPIService,
	)

	service := &service{
		klineServicer:     klineServicer,
		orderBookServicer: orderBookServicer,
		orderServicer:     orderServicer,
		tickerServicer:    tickerServicer,
	}

	klineServicerWithTypeCheck, ok := klineServicer.(WithServicer)
	if ok {
		klineServicerWithTypeCheck.WithServicer(service)
	}

	orderBookServicerWithTypeCheck, ok := orderBookServicer.(WithServicer)
	if ok {
		orderBookServicerWithTypeCheck.WithServicer(service)
	}

	orderServicerWithTypeCheck, ok := orderServicer.(WithServicer)
	if ok {
		orderServicerWithTypeCheck.WithServicer(service)
	}

	tickerServicerWithTypeCheck, ok := tickerServicer.(WithServicer)
	if ok {
		tickerServicerWithTypeCheck.WithServicer(service)
	}

	return service
}

// GetKlineServicer is a function.
func (service *service) GetKlineServicer() KlineServicer {
	return service.klineServicer
}

// GetOrderBookServicer is a function.
func (service *service) GetOrderBookServicer() OrderBookServicer {
	return service.orderBookServicer
}

// GetOrderServicer is a function.
func (service *service) GetOrderServicer() OrderServicer {
	return service.orderServicer
}

// GetTickerServicer is a function.
func (service *service) GetTickerServicer() TickerServicer {
	return service.tickerServicer
}
