package util

import (
	"context"

	"github.com/ShahoBashoki/kucoin/config"
	"github.com/ShahoBashoki/kucoin/log"
	"github.com/ShahoBashoki/kucoin/object"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	sdkResource "go.opentelemetry.io/otel/sdk/resource"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

type (
	// Tracer is an interface.
	Tracer interface {
		// GetTracer is a function.
		GetTracer() trace.Tracer
	}

	// GetTracer is an interface.
	GetTracer interface {
		// GetTracer is a function.
		GetTracer() trace.Tracer
	}
)

// NewTracer is a function.
func NewTracer(
	ctx context.Context,
	configConfigger config.Configger,
	logRuntimeLogger log.RuntimeLogger,
	utilUUIDer UUIDer,
) trace.Tracer {
	runtimeContext := NewRuntimeContext(ctx, utilUUIDer)
	spanContext := NewSpanContext(nil)
	fields := map[string]any{
		"name":   "NewOpenTelemetry",
		"rt_ctx": runtimeContext,
		"sp_ctx": spanContext,
		"config": configConfigger,
	}

	logRuntimeLogger.
		WithFields(fields).
		Info(object.URIEmpty)

	sdkResourceResource := sdkResource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceInstanceIDKey.String(
			configConfigger.GetOtelConfigger().GetServiceInstanceID(),
		),
		semconv.ServiceNameKey.String(configConfigger.GetOtelConfigger().GetServiceName()),
		semconv.ServiceNamespaceKey.String(
			configConfigger.GetOtelConfigger().GetServiceNamespace(),
		),
		semconv.ServiceVersionKey.String(configConfigger.GetOtelConfigger().GetServiceVersion()),
	)

	logRuntimeLogger.
		WithFields(fields).
		WithField(object.URIFieldSDKResourceResource, sdkResourceResource).
		Debug(object.URIEmpty)

	newSDKResourceResource, err := sdkResource.Merge(sdkResource.Default(), sdkResourceResource)
	if err != nil {
		logRuntimeLogger.
			WithFields(fields).
			WithField(object.URIFieldError, err).
			Error(object.ErrSDKResourceMerge.Error())
	}

	logRuntimeLogger.
		WithFields(fields).
		WithField(object.URIFieldSDKResourceResource, newSDKResourceResource).
		Debug(object.URIEmpty)

	sdkResourceResource3, err := sdkResource.New(
		ctx,
		sdkResource.WithContainer(),
		sdkResource.WithFromEnv(),
		sdkResource.WithOS(),
		sdkResource.WithProcess(),
	)
	if err != nil {
		logRuntimeLogger.
			WithFields(fields).
			WithField(object.URIFieldError, err).
			Error(object.ErrSDKResourceNew.Error())
	}

	logRuntimeLogger.
		WithFields(fields).
		WithField(object.URIFieldSDKResourceResource, sdkResourceResource3).
		Debug(object.URIEmpty)

	sdkResourceResource4, err := sdkResource.Merge(newSDKResourceResource, sdkResourceResource3)
	if err != nil {
		logRuntimeLogger.
			WithFields(fields).
			WithField(object.URIFieldError, err).
			Error(object.ErrSDKResourceMerge.Error())
	}

	logRuntimeLogger.
		WithFields(fields).
		WithField(object.URIFieldSDKResourceResource, sdkResourceResource4).
		Debug(object.URIEmpty)

	jaegerExporter, err := jaeger.New(
		jaeger.WithCollectorEndpoint(
			// Other options.
			// Jaeger.WithAgentEndpoint(),
			// jaeger.WithAgentHost(),
			// jaeger.WithAgentPort(),
			// jaeger.WithLogger(),
			// jaeger.WithDisableAttemptReconnecting(),
			// jaeger.WithAttemptReconnectingInterval(),
			// jaeger.WithMaxPacketSize(),
			// jaeger.WithCollectorEndpoint(),.
			// Jaeger.WithHTTPClient(http.DefaultClient),.
			jaeger.WithEndpoint(configConfigger.GetOtelConfigger().GetExporterJaegerEndpoint()),
			jaeger.WithUsername(configConfigger.GetOtelConfigger().GetExporterJaegerPassword()),
			jaeger.WithUsername(configConfigger.GetOtelConfigger().GetExporterJaegerUsername()),
		),
	)
	if err != nil {
		logRuntimeLogger.
			WithFields(fields).
			WithField(object.URIFieldError, err).
			Error(object.ErrJaegerNew.Error())
	}

	logRuntimeLogger.
		WithFields(fields).
		WithField(object.URIFieldJaegerExporter, jaegerExporter).
		Debug(object.URIEmpty)

	sdkTraceTracerProvider := sdkTrace.NewTracerProvider(
		// Other options.
		// SdkTrace.WithBatchTimeout(),
		// sdkTrace.WithBlocking(),
		// sdkTrace.WithExportTimeout(),
		// sdkTrace.WithMaxExportBatchSize(),
		// sdkTrace.WithMaxQueueSize(),
		// ##############################
		// sdkTrace.WithIDGenerator(),
		// sdkTrace.WithRawSpanLimits(),.
		// SdkTrace.WithSampler(),
		// sdkTrace.WithSpanLimits(),
		// sdkTrace.WithSpanProcessor(),
		// sdkTrace.WithSyncer(),.
		sdkTrace.WithBatcher(jaegerExporter),
		sdkTrace.WithResource(sdkResourceResource4),
	)

	logRuntimeLogger.
		WithFields(fields).
		WithField(object.URIFieldTracerProvider, sdkTraceTracerProvider).
		Debug(object.URIEmpty)

	otel.SetTracerProvider(sdkTraceTracerProvider)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(b3.New(b3.WithInjectEncoding(b3.B3SingleHeader))),
	)

	traceTracer := otel.Tracer(configConfigger.GetOtelConfigger().GetInstrumentationName())

	logRuntimeLogger.
		WithFields(fields).
		WithField(object.URIFieldTracer, traceTracer).
		Debug(object.URIEmpty)

	go func() {
		<-ctx.Done()

		logRuntimeLogger.
			WithFields(fields).
			Debug(`shutting down gracefully the tracing`)

		ctxWT, ctxWTCancelFunc := context.WithTimeout(ctx, object.NUMSystemGracefulShutdown)
		defer ctxWTCancelFunc()

		if err = sdkTraceTracerProvider.Shutdown(ctxWT); err != nil {
			logRuntimeLogger.
				WithFields(fields).
				WithField(object.URIFieldError, err).
				Error(object.ErrTracerProviderShutdown.Error())
		}
	}()

	return traceTracer
}
