package server

import (
	"context"
	"os"

	"github.com/ShahoBashoki/kucoin/config"
	"github.com/ShahoBashoki/kucoin/log"
	"github.com/ShahoBashoki/kucoin/object"
	"github.com/ShahoBashoki/kucoin/service"
	"github.com/ShahoBashoki/kucoin/util"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type (
	// Fn is a function.
	Fn func(
		ctx context.Context,
	) (context.Context, error)

	// Serverer is an interface.
	Serverer interface {
		config.GetConfigger
		log.GetRuntimeLogger
		service.GetServicer
		util.GetTracer
		// Run is a function.
		Run(
			context.Context,
		) error
	}

	// GetServerer is an interface.
	GetServerer interface {
		// GetServerer is a function.
		GetServerer() Serverer
	}

	// WithServerer is an interface.
	WithServerer interface {
		// WithServerer is a function.
		WithServerer(
			Serverer,
		)
	}

	server struct {
		configConfigger  config.Configger
		logRuntimeLogger log.RuntimeLogger
		servicer         service.Servicer
		traceTracer      trace.Tracer
		utilUUIDer       util.UUIDer
	}
)

var _ Serverer = (*server)(nil)

// NewServerrer is a function.
func NewServerrer(
	configConfigger config.Configger,
	ctx context.Context,
	logRuntimeLogger log.RuntimeLogger,
	servicer service.Servicer,
	traceTracer trace.Tracer,
	utilUUIDer util.UUIDer,
) Serverer {
	server := &server{
		configConfigger:  configConfigger,
		logRuntimeLogger: logRuntimeLogger,
		servicer:         servicer,
		traceTracer:      traceTracer,
		utilUUIDer:       utilUUIDer,
	}

	return server
}

// GetConfigger is a function.
func (server *server) GetConfigger() config.Configger {
	return server.configConfigger
}

// GetRuntimeLogger is a function.
func (server *server) GetRuntimeLogger() log.RuntimeLogger {
	return server.logRuntimeLogger
}

// GetTracer is a function.
func (server *server) GetServicer() service.Servicer {
	return server.servicer
}

// GetTracer is a function.
func (server *server) GetTracer() trace.Tracer {
	return server.traceTracer
}

// GetUUIDer is a function.
func (server *server) GetUUIDer() util.UUIDer {
	return server.utilUUIDer
}

// Run is a function.
func (server *server) Run(
	ctx context.Context,
) error {
	var traceSpan trace.Span

	ctx, traceSpan = server.GetTracer().Start(
		ctx,
		"Run",
		trace.WithSpanKind(trace.SpanKindProducer),
	)
	defer traceSpan.End()

	utilRuntimeContext := util.NewRuntimeContext(ctx, server.GetUUIDer())
	utilSpanContext := util.NewSpanContext(traceSpan)
	fields := map[string]any{
		"name":   "Run",
		"rt_ctx": utilRuntimeContext,
		"sp_ctx": utilSpanContext,
		"config": server.configConfigger,
	}

	server.GetRuntimeLogger().
		WithFields(fields).
		Info(object.URIEmpty)

	gin.SetMode(os.Getenv("GIN_MODE"))
	router := gin.Default()

	errRouterRun := router.Run()
	if errRouterRun != nil {
		server.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, errRouterRun).
			Error(object.ErrRouterRun.Error())
		traceSpan.RecordError(errRouterRun)
		traceSpan.SetStatus(codes.Error, object.ErrRouterRun.Error())

		return errRouterRun
	}

	return nil
}
