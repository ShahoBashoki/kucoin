package service

import (
	"context"
	"fmt"

	"github.com/Kucoin/kucoin-go-sdk"
	"github.com/ShahoBashoki/kucoin/config"
	"github.com/ShahoBashoki/kucoin/log"
	"github.com/ShahoBashoki/kucoin/object"
	"github.com/ShahoBashoki/kucoin/object/dto"
	"github.com/ShahoBashoki/kucoin/util"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type (
	// KlineServicer is an interface.
	KlineServicer interface {
		// GetOpenLowEqualityFromRemote is a function.
		GetOpenLowEqualityFromRemote(
			context.Context,
			dto.KlineRequester,
		) (bool, error)
	}

	// GetKlineServicer is an interface.
	GetKlineServicer interface {
		// GetKlineServicer is a function.
		GetKlineServicer() KlineServicer
	}

	klineService struct {
		configConfigger  config.Configger
		logRuntimeLogger log.RuntimeLogger
		servicer         Servicer
		traceTracer      trace.Tracer
		utilUUIDer       util.UUIDer
		kucoinAPIService *kucoin.ApiService
	}
)

var (
	_ GetServicer          = (*klineService)(nil)
	_ KlineServicer        = (*klineService)(nil)
	_ WithServicer         = (*klineService)(nil)
	_ config.GetConfigger  = (*klineService)(nil)
	_ log.GetRuntimeLogger = (*klineService)(nil)
	_ util.GetTracer       = (*klineService)(nil)
	_ util.GetUUIDer       = (*klineService)(nil)
)

// NewKlineServicer is a function.
func NewKlineServicer(
	configConfigger config.Configger,
	logRuntimeLogger log.RuntimeLogger,
	traceTracer trace.Tracer,
	utilUUIDer util.UUIDer,
	kucoinAPIService *kucoin.ApiService,
) KlineServicer {
	return &klineService{
		configConfigger:  configConfigger,
		logRuntimeLogger: logRuntimeLogger,
		servicer:         nil,
		traceTracer:      traceTracer,
		utilUUIDer:       utilUUIDer,
		kucoinAPIService: kucoinAPIService,
	}
}

// GetConfigger is a function.
func (service *klineService) GetConfigger() config.Configger {
	return service.configConfigger
}

// GetRuntimeLogger is a function.
func (service *klineService) GetRuntimeLogger() log.RuntimeLogger {
	return service.logRuntimeLogger
}

// GetServicer is a function.
func (service *klineService) GetServicer() Servicer {
	return service.servicer
}

// GetTracer is a function.
func (service *klineService) GetTracer() trace.Tracer {
	return service.traceTracer
}

// GetUUIDer is a function.
func (service *klineService) GetUUIDer() util.UUIDer {
	return service.utilUUIDer
}

// GetAPIService is a function.
func (service *klineService) GetAPIService() *kucoin.ApiService {
	return service.kucoinAPIService
}

// WithServicer is a function.
func (service *klineService) WithServicer(
	servicer Servicer,
) {
	service.servicer = servicer
}

// GetOpenLowEqualityFromRemote is a function.
func (service *klineService) GetOpenLowEqualityFromRemote(
	ctx context.Context,
	dtoKlineRequester dto.KlineRequester,
) (bool, error) {
	var traceSpan trace.Span

	ctx, traceSpan = service.GetTracer().Start(
		ctx,
		"GetOpenLowEqualityFromRemote",
		trace.WithSpanKind(trace.SpanKindProducer),
	)
	defer traceSpan.End()

	utilRuntimeContext := util.NewRuntimeContext(ctx, service.GetUUIDer())
	utilSpanContext := util.NewSpanContext(traceSpan)
	fields := map[string]any{
		"name":                "GetOpenLowEqualityFromRemote",
		"rt_ctx":              utilRuntimeContext,
		"sp_ctx":              utilSpanContext,
		"config":              service.configConfigger,
		"dto_kline_requester": dtoKlineRequester,
	}

	service.GetRuntimeLogger().
		WithFields(fields).
		Info(object.URIEmpty)

	response, err := service.kucoinAPIService.KLines(
		dtoKlineRequester.GetSymbol(),
		string(dtoKlineRequester.GetKlineType()),
		dtoKlineRequester.GetStartAt(),
		dtoKlineRequester.GetEndAt(),
	)
	if err != nil {
		service.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, err).
			Error(object.ErrKlineKucoinServiceGetList.Error())
		traceSpan.RecordError(err)
		traceSpan.SetStatus(codes.Error, object.ErrKlineKucoinServiceGetList.Error())

		return false, fmt.Errorf("%w", err)
	}

	service.GetRuntimeLogger().
		WithFields(fields).
		WithField(object.URIFieldResponse, response).
		Debug(object.URIEmpty)

	if response.Code != "200000" {
		service.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, response.Message).
			Error(object.ErrKlineKucoinServiceGetList.Error())
		traceSpan.RecordError(object.ErrKlineKucoinServiceGetList)
		traceSpan.SetStatus(codes.Error, object.ErrKlineKucoinServiceGetList.Error())

		return false, object.ErrKlineKucoinServiceGetList
	}

	var kucoinKLinesModel [][]string

	if err = response.ReadData(&kucoinKLinesModel); err != nil {
		service.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, err).
			Error(object.ErrKucoinServiceReadPaginationData.Error())
		traceSpan.RecordError(err)
		traceSpan.SetStatus(codes.Error, object.ErrKucoinServiceReadPaginationData.Error())

		return false, fmt.Errorf("%w", err)
	}

	service.GetRuntimeLogger().
		WithFields(fields).
		WithField(object.URIFieldKucoinKLinesModel, kucoinKLinesModel).
		Debug(object.URIEmpty)

	for key, value := range kucoinKLinesModel {
		service.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldKey, key).
			WithField(object.URIFieldValue, value).
			Debug(object.URIEmpty)

		if value[1] == value[4] {
			service.GetRuntimeLogger().
				WithFields(fields).
				Debug(`value[1] == value[4]`)

			return true, nil
		}
	}

	return false, nil
}
