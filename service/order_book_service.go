package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Kucoin/kucoin-go-sdk"
	"github.com/ShahoBashoki/kucoin/config"
	"github.com/ShahoBashoki/kucoin/log"
	"github.com/ShahoBashoki/kucoin/object"
	"github.com/ShahoBashoki/kucoin/util"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type (
	// OrderBookServicer is an interface.
	OrderBookServicer interface {
		// GetMarketRatioFromRemote is a function.
		GetMarketRatioFromRemote(
			context.Context,
			string, float64,
		) (bool, error)
	}

	// GetOrderBookServicer is an interface.
	GetOrderBookServicer interface {
		// GetOrderBookServicer is a function.
		GetOrderBookServicer() OrderBookServicer
	}

	orderBookService struct {
		configConfigger  config.Configger
		logRuntimeLogger log.RuntimeLogger
		servicer         Servicer
		traceTracer      trace.Tracer
		utilUUIDer       util.UUIDer
		kucoinAPIService *kucoin.ApiService
	}
)

var (
	_ GetServicer          = (*orderBookService)(nil)
	_ OrderBookServicer    = (*orderBookService)(nil)
	_ WithServicer         = (*orderBookService)(nil)
	_ config.GetConfigger  = (*orderBookService)(nil)
	_ log.GetRuntimeLogger = (*orderBookService)(nil)
	_ util.GetTracer       = (*orderBookService)(nil)
	_ util.GetUUIDer       = (*orderBookService)(nil)
)

// NewOrderBookServicer is a function.
func NewOrderBookServicer(
	configConfigger config.Configger,
	logRuntimeLogger log.RuntimeLogger,
	traceTracer trace.Tracer,
	utilUUIDer util.UUIDer,
	kucoinAPIService *kucoin.ApiService,
) OrderBookServicer {
	return &orderBookService{
		configConfigger:  configConfigger,
		logRuntimeLogger: logRuntimeLogger,
		servicer:         nil,
		traceTracer:      traceTracer,
		utilUUIDer:       utilUUIDer,
		kucoinAPIService: kucoinAPIService,
	}
}

// GetConfigger is a function.
func (service *orderBookService) GetConfigger() config.Configger {
	return service.configConfigger
}

// GetRuntimeLogger is a function.
func (service *orderBookService) GetRuntimeLogger() log.RuntimeLogger {
	return service.logRuntimeLogger
}

// GetServicer is a function.
func (service *orderBookService) GetServicer() Servicer {
	return service.servicer
}

// GetTracer is a function.
func (service *orderBookService) GetTracer() trace.Tracer {
	return service.traceTracer
}

// GetUUIDer is a function.
func (service *orderBookService) GetUUIDer() util.UUIDer {
	return service.utilUUIDer
}

// GetAPIService is a function.
func (service *orderBookService) GetAPIService() *kucoin.ApiService {
	return service.kucoinAPIService
}

// WithServicer is a function.
func (service *orderBookService) WithServicer(
	servicer Servicer,
) {
	service.servicer = servicer
}

// GetMarketRatioFromRemote is a function.
func (service *orderBookService) GetMarketRatioFromRemote(
	ctx context.Context,
	symbol string,
	ratio float64,
) (bool, error) {
	var traceSpan trace.Span

	ctx, traceSpan = service.GetTracer().Start(
		ctx,
		"GetMarketRatioFromRemote",
		trace.WithSpanKind(trace.SpanKindProducer),
	)
	defer traceSpan.End()

	utilRuntimeContext := util.NewRuntimeContext(ctx, service.GetUUIDer())
	utilSpanContext := util.NewSpanContext(traceSpan)
	fields := map[string]any{
		"name":   "GetMarketRatioFromRemote",
		"rt_ctx": utilRuntimeContext,
		"sp_ctx": utilSpanContext,
		"config": service.configConfigger,
		"symbol": symbol,
		"ratio":  ratio,
	}

	service.GetRuntimeLogger().
		WithFields(fields).
		Info(object.URIEmpty)

	response, err := service.kucoinAPIService.AggregatedFullOrderBookV3(symbol)
	if err != nil {
		service.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, err).
			Error(object.ErrOrderBookKucoinServiceGetList.Error())
		traceSpan.RecordError(err)
		traceSpan.SetStatus(codes.Error, object.ErrOrderBookKucoinServiceGetList.Error())

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
			Error(object.ErrOrderBookKucoinServiceGetList.Error())
		traceSpan.RecordError(object.ErrOrderBookKucoinServiceGetList)
		traceSpan.SetStatus(codes.Error, object.ErrOrderBookKucoinServiceGetList.Error())

		return false, object.ErrOrderBookKucoinServiceGetList
	}

	kucoinFullOrderBookModel := kucoin.FullOrderBookModel{
		Sequence: "",
		Time:     0,
		Bids:     [][]string{},
		Asks:     [][]string{},
	}

	if err = response.ReadData(&kucoinFullOrderBookModel); err != nil {
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
		WithField(object.URIFieldKucoinPaginationModel, kucoinFullOrderBookModel).
		Debug(object.URIEmpty)

	bidsValue := 0.0

	for key, value := range kucoinFullOrderBookModel.Bids {
		service.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldKey, key).
			WithField(object.URIFieldValue, value).
			Debug(object.URIEmpty)

		if key > len(kucoinFullOrderBookModel.Bids)/2 {
			service.GetRuntimeLogger().
				WithFields(fields).
				Debug(`key > len(kucoinFullOrderBookModel.Bids)/2`)

			break
		}

		firstValue, errParseFloat := strconv.ParseFloat(value[0], 64)
		if errParseFloat != nil {
			service.GetRuntimeLogger().
				WithFields(fields).
				WithField(object.URIFieldError, errParseFloat).
				Error(object.ErrSTRCONVParseFloat.Error())
			traceSpan.RecordError(errParseFloat)
			traceSpan.SetStatus(codes.Error, object.ErrSTRCONVParseFloat.Error())

			return false, errParseFloat
		}

		secondValue, errParseFloat := strconv.ParseFloat(value[1], 64)
		if errParseFloat != nil {
			service.GetRuntimeLogger().
				WithFields(fields).
				WithField(object.URIFieldError, errParseFloat).
				Error(object.ErrSTRCONVParseFloat.Error())
			traceSpan.RecordError(errParseFloat)
			traceSpan.SetStatus(codes.Error, object.ErrSTRCONVParseFloat.Error())

			return false, errParseFloat
		}

		bidsValue += firstValue + secondValue
	}

	service.GetRuntimeLogger().
		WithFields(fields).
		WithField(object.URIFieldBidsValue, bidsValue).
		Debug(object.URIEmpty)

	asksValue := 0.0

	for key, value := range kucoinFullOrderBookModel.Asks {
		service.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldKey, key).
			WithField(object.URIFieldValue, value).
			Debug(object.URIEmpty)

		if key > len(kucoinFullOrderBookModel.Asks)/2 {
			service.GetRuntimeLogger().
				WithFields(fields).
				Debug(`key > len(kucoinFullOrderBookModel.Asks)/2`)

			break
		}

		firstValue, errParseFloat := strconv.ParseFloat(value[0], 64)
		if errParseFloat != nil {
			service.GetRuntimeLogger().
				WithFields(fields).
				WithField(object.URIFieldError, errParseFloat).
				Error(object.ErrSTRCONVParseFloat.Error())
			traceSpan.RecordError(errParseFloat)
			traceSpan.SetStatus(codes.Error, object.ErrSTRCONVParseFloat.Error())

			return false, errParseFloat
		}

		secondValue, errParseFloat := strconv.ParseFloat(value[1], 64)
		if errParseFloat != nil {
			service.GetRuntimeLogger().
				WithFields(fields).
				WithField(object.URIFieldError, errParseFloat).
				Error(object.ErrSTRCONVParseFloat.Error())
			traceSpan.RecordError(errParseFloat)
			traceSpan.SetStatus(codes.Error, object.ErrSTRCONVParseFloat.Error())

			return false, errParseFloat
		}

		asksValue += firstValue + secondValue
	}

	service.GetRuntimeLogger().
		WithFields(fields).
		WithField(object.URIFieldAsksValue, asksValue).
		Debug(object.URIEmpty)

	if bidsValue/asksValue > ratio {
		service.GetRuntimeLogger().
			WithFields(fields).
			Debug(`bidsValue/asksValue > ratio`)

		return true, nil
	}

	return false, nil
}
