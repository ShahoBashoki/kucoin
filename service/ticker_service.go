package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Kucoin/kucoin-go-sdk"
	"github.com/ShahoBashoki/kucoin/config"
	"github.com/ShahoBashoki/kucoin/log"
	"github.com/ShahoBashoki/kucoin/object"
	"github.com/ShahoBashoki/kucoin/object/dao"
	"github.com/ShahoBashoki/kucoin/object/om"
	"github.com/ShahoBashoki/kucoin/repository"
	"github.com/ShahoBashoki/kucoin/util"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type (
	// TickerServicer is an interface.
	TickerServicer interface {
		// Create is a function.
		Create(
			context.Context,
			om.Tickerer,
		) (uuid.UUID, error)
		// DeleteAll is a function.
		DeleteAll(
			context.Context,
		) (time.Time, error)
		// Get is a function.
		Get(
			context.Context,
			uuid.UUID,
		) (om.Tickerer, error)
		// GetListFromRemote is a function.
		GetListFromRemote(
			context.Context,
		) error
		// GetListFromRepository is a function.
		GetListFromRepository(
			context.Context,
			dao.Paginationer,
			dao.TickerFilterer,
		) ([]om.Tickerer, dao.Cursorer, error)
	}

	// GetTickerServicer is an interface.
	GetTickerServicer interface {
		// GetTickerServicer is a function.
		GetTickerServicer() TickerServicer
	}

	tickerService struct {
		configConfigger  config.Configger
		repositorier     repository.TickerRepositorier
		logRuntimeLogger log.RuntimeLogger
		servicer         Servicer
		traceTracer      trace.Tracer
		utilUUIDer       util.UUIDer
		kucoinAPIService *kucoin.ApiService
	}
)

var (
	_ GetServicer                      = (*tickerService)(nil)
	_ TickerServicer                   = (*tickerService)(nil)
	_ WithServicer                     = (*tickerService)(nil)
	_ config.GetConfigger              = (*tickerService)(nil)
	_ log.GetRuntimeLogger             = (*tickerService)(nil)
	_ repository.GetTickerRepositorier = (*tickerService)(nil)
	_ util.GetTracer                   = (*tickerService)(nil)
	_ util.GetUUIDer                   = (*tickerService)(nil)
)

// NewTickerServicer is a function.
func NewTickerServicer(
	configConfigger config.Configger,
	repositorier repository.TickerRepositorier,
	logRuntimeLogger log.RuntimeLogger,
	traceTracer trace.Tracer,
	utilUUIDer util.UUIDer,
	kucoinAPIService *kucoin.ApiService,
) TickerServicer {
	return &tickerService{
		configConfigger:  configConfigger,
		repositorier:     repositorier,
		logRuntimeLogger: logRuntimeLogger,
		servicer:         nil,
		traceTracer:      traceTracer,
		utilUUIDer:       utilUUIDer,
		kucoinAPIService: kucoinAPIService,
	}
}

// GetConfigger is a function.
func (service *tickerService) GetConfigger() config.Configger {
	return service.configConfigger
}

// GetRuntimeLogger is a function.
func (service *tickerService) GetRuntimeLogger() log.RuntimeLogger {
	return service.logRuntimeLogger
}

// GetServicer is a function.
func (service *tickerService) GetServicer() Servicer {
	return service.servicer
}

// GetTickerRepositorier is a function.
func (service *tickerService) GetTickerRepositorier() repository.TickerRepositorier {
	return service.repositorier
}

// GetTracer is a function.
func (service *tickerService) GetTracer() trace.Tracer {
	return service.traceTracer
}

// GetUUIDer is a function.
func (service *tickerService) GetUUIDer() util.UUIDer {
	return service.utilUUIDer
}

// GetAPIService is a function.
func (service *tickerService) GetAPIService() *kucoin.ApiService {
	return service.kucoinAPIService
}

// WithServicer is a function.
func (service *tickerService) WithServicer(
	servicer Servicer,
) {
	service.servicer = servicer
}

// Create is a function.
func (service *tickerService) Create(
	ctx context.Context,
	omTickerer om.Tickerer,
) (uuid.UUID, error) {
	var traceSpan trace.Span

	ctx, traceSpan = service.GetTracer().Start(
		ctx,
		"Create",
		trace.WithSpanKind(trace.SpanKindProducer),
	)
	defer traceSpan.End()

	utilRuntimeContext := util.NewRuntimeContext(ctx, service.GetUUIDer())
	utilSpanContext := util.NewSpanContext(traceSpan)
	fields := map[string]any{
		"name":        "Create",
		"rt_ctx":      utilRuntimeContext,
		"sp_ctx":      utilSpanContext,
		"config":      service.configConfigger,
		"om_tickerer": omTickerer,
	}

	service.GetRuntimeLogger().
		WithFields(fields).
		Info(object.URIEmpty)

	daoTicker := dao.NewTicker(
		time.Time{},
		time.Time{},
		sql.NullTime{
			Time:  time.Time{},
			Valid: false,
		},
		uuid.Nil,
		omTickerer.GetAveragePrice(),
		omTickerer.GetBuy(),
		omTickerer.GetChangePrice(),
		omTickerer.GetChangeRate(),
		omTickerer.GetHigh(),
		omTickerer.GetLast(),
		omTickerer.GetLow(),
		omTickerer.GetMakerCoefficient(),
		omTickerer.GetMakerFeeRate(),
		omTickerer.GetSell(),
		omTickerer.GetSymbol(),
		omTickerer.GetSymbolName(),
		omTickerer.GetTakerCoefficient(),
		omTickerer.GetTakerFeeRate(),
		omTickerer.GetVol(),
		omTickerer.GetVolValue(),
	)

	service.GetRuntimeLogger().
		WithFields(fields).
		WithField(object.URIFieldDAOTicker, daoTicker).
		Debug(object.URIEmpty)

	tickerID, err := service.GetTickerRepositorier().Create(ctx, daoTicker)
	if err != nil {
		service.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, err).
			Error(object.ErrTickerRepositoryCreate.Error())
		traceSpan.RecordError(err)
		traceSpan.SetStatus(codes.Error, object.ErrTickerRepositoryCreate.Error())

		return uuid.Nil, err
	}

	service.GetRuntimeLogger().
		WithFields(fields).
		WithField(object.URIFieldTickerID, tickerID).
		Debug(object.URIEmpty)

	return tickerID, nil
}

// DeleteAll is a function.
func (service *tickerService) DeleteAll(
	ctx context.Context,
) (time.Time, error) {
	var traceSpan trace.Span

	ctx, traceSpan = service.GetTracer().Start(
		ctx,
		"DeleteAll",
		trace.WithSpanKind(trace.SpanKindProducer),
	)
	defer traceSpan.End()

	utilRuntimeContext := util.NewRuntimeContext(ctx, service.GetUUIDer())
	utilSpanContext := util.NewSpanContext(traceSpan)
	fields := map[string]any{
		"name":   "DeleteAll",
		"rt_ctx": utilRuntimeContext,
		"sp_ctx": utilSpanContext,
		"config": service.configConfigger,
	}

	service.GetRuntimeLogger().
		WithFields(fields).
		Info(object.URIEmpty)

	deletedAt, err := service.GetTickerRepositorier().DeleteAll(ctx)
	if err != nil {
		service.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, err).
			Error(object.ErrTickerRepositoryDeleteAll.Error())
		traceSpan.RecordError(err)
		traceSpan.SetStatus(codes.Error, object.ErrTickerRepositoryDeleteAll.Error())

		return time.Time{}, err
	}

	service.GetRuntimeLogger().
		WithFields(fields).
		WithField(object.URIFieldDeletedAt, deletedAt).
		Debug(object.URIEmpty)

	return deletedAt, nil
}

// Get is a function.
func (service *tickerService) Get(
	ctx context.Context,
	id uuid.UUID,
) (om.Tickerer, error) {
	var traceSpan trace.Span

	ctx, traceSpan = service.GetTracer().Start(
		ctx,
		"Get",
		trace.WithSpanKind(trace.SpanKindProducer),
	)
	defer traceSpan.End()

	utilRuntimeContext := util.NewRuntimeContext(ctx, service.GetUUIDer())
	utilSpanContext := util.NewSpanContext(traceSpan)
	fields := map[string]any{
		"name":   "Get",
		"rt_ctx": utilRuntimeContext,
		"sp_ctx": utilSpanContext,
		"config": service.configConfigger,
		"id":     id,
	}

	service.GetRuntimeLogger().
		WithFields(fields).
		Info(object.URIEmpty)

	daoTicker, err := service.GetTickerRepositorier().Read(ctx, id)
	if err != nil {
		service.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, err).
			Error(object.ErrTickerRepositoryRead.Error())
		traceSpan.RecordError(err)
		traceSpan.SetStatus(codes.Error, object.ErrTickerRepositoryRead.Error())

		return nil, err
	}

	service.GetRuntimeLogger().
		WithFields(fields).
		WithField(object.URIFieldDAOTicker, daoTicker).
		Debug(object.URIEmpty)

	omTicker := om.NewTicker(
		daoTicker.GetAveragePrice(),
		daoTicker.GetBuy(),
		daoTicker.GetChangePrice(),
		daoTicker.GetChangeRate(),
		daoTicker.GetHigh(),
		daoTicker.GetLast(),
		daoTicker.GetLow(),
		daoTicker.GetMakerCoefficient(),
		daoTicker.GetMakerFeeRate(),
		daoTicker.GetSell(),
		daoTicker.GetSymbol(),
		daoTicker.GetSymbolName(),
		daoTicker.GetTakerCoefficient(),
		daoTicker.GetTakerFeeRate(),
		daoTicker.GetVol(),
		daoTicker.GetVolValue(),
		daoTicker.GetID(),
	)

	service.GetRuntimeLogger().
		WithFields(fields).
		WithField(object.URIFieldOMTicker, omTicker).
		Debug(object.URIEmpty)

	return omTicker, nil
}

// GetListFromRemote is a function.
func (service *tickerService) GetListFromRemote(
	ctx context.Context,
) error {
	var traceSpan trace.Span

	ctx, traceSpan = service.GetTracer().Start(
		ctx,
		"GetListFromRemote",
		trace.WithSpanKind(trace.SpanKindProducer),
	)
	defer traceSpan.End()

	utilRuntimeContext := util.NewRuntimeContext(ctx, service.GetUUIDer())
	utilSpanContext := util.NewSpanContext(traceSpan)
	fields := map[string]any{
		"name":   "GetListFromRemote",
		"rt_ctx": utilRuntimeContext,
		"sp_ctx": utilSpanContext,
		"config": service.configConfigger,
	}

	service.GetRuntimeLogger().
		WithFields(fields).
		Info(object.URIEmpty)

	response, err := service.kucoinAPIService.Tickers()
	if err != nil {
		service.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, err).
			Error(object.ErrTickerKucoinServiceGetList.Error())
		traceSpan.RecordError(err)
		traceSpan.SetStatus(codes.Error, object.ErrTickerKucoinServiceGetList.Error())

		return fmt.Errorf("%w", err)
	}

	service.GetRuntimeLogger().
		WithFields(fields).
		WithField(object.URIFieldResponse, response).
		Debug(object.URIEmpty)

	if response.Code != "200000" {
		service.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, response.Message).
			Error(object.ErrTickerKucoinServiceGetList.Error())
		traceSpan.RecordError(object.ErrTickerKucoinServiceGetList)
		traceSpan.SetStatus(codes.Error, object.ErrTickerKucoinServiceGetList.Error())

		return object.ErrTickerKucoinServiceGetList
	}

	kucoinTickersModel := kucoin.TickersResponseModel{
		Time:    0,
		Tickers: []*kucoin.TickerModel{},
	}

	if err = response.ReadData(&kucoinTickersModel); err != nil {
		service.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, err).
			Error(object.ErrKucoinServiceReadPaginationData.Error())
		traceSpan.RecordError(err)
		traceSpan.SetStatus(codes.Error, object.ErrKucoinServiceReadPaginationData.Error())

		return fmt.Errorf("%w", err)
	}

	service.GetRuntimeLogger().
		WithFields(fields).
		WithField(object.URIFieldKucoinPaginationModel, kucoinTickersModel).
		Debug(object.URIEmpty)

	for key, value := range kucoinTickersModel.Tickers {
		service.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldKey, key).
			WithField(object.URIFieldValue, value).
			Debug(object.URIEmpty)

		omTicker := om.NewTicker(
			value.AveragePrice,
			value.Buy,
			value.ChangePrice,
			value.ChangeRate,
			value.High,
			value.Last,
			value.Low,
			value.MakerCoefficient,
			value.MakerFeeRate,
			value.Sell,
			value.Symbol,
			value.SymbolName,
			value.TakerCoefficient,
			value.TakerFeeRate,
			value.Vol,
			value.VolValue,
			uuid.Nil,
		)

		tickerID, errTickerCreate := service.GetServicer().GetTickerServicer().Create(ctx, omTicker)
		if errTickerCreate != nil {
			service.GetRuntimeLogger().
				WithFields(fields).
				WithField(object.URIFieldError, errTickerCreate).
				Error(object.ErrTickerServiceCreate.Error())
			traceSpan.RecordError(errTickerCreate)
			traceSpan.SetStatus(codes.Error, object.ErrTickerServiceCreate.Error())

			return errTickerCreate
		}

		service.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldTickerID, tickerID).
			Debug(object.URIEmpty)
	}

	return nil
}

// GetListFromRepository is a function.
func (service *tickerService) GetListFromRepository(
	ctx context.Context,
	daoPaginator dao.Paginationer,
	daoTickerFilterer dao.TickerFilterer,
) ([]om.Tickerer, dao.Cursorer, error) {
	var traceSpan trace.Span

	ctx, traceSpan = service.GetTracer().Start(
		ctx,
		"GetListFromRepository",
		trace.WithSpanKind(trace.SpanKindProducer),
	)
	defer traceSpan.End()

	utilRuntimeContext := util.NewRuntimeContext(ctx, service.GetUUIDer())
	utilSpanContext := util.NewSpanContext(traceSpan)
	fields := map[string]any{
		"name":                "GetListFromRepository",
		"rt_ctx":              utilRuntimeContext,
		"sp_ctx":              utilSpanContext,
		"config":              service.configConfigger,
		"dao_paginator":       daoPaginator,
		"dao_ticker_filterer": daoTickerFilterer,
	}

	service.GetRuntimeLogger().
		WithFields(fields).
		Info(object.URIEmpty)

	daoTickers, daoCursorer, err := service.GetTickerRepositorier().
		ReadList(ctx, daoPaginator, daoTickerFilterer)
	if err != nil {
		service.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, err).
			Error(object.ErrTickerRepositoryReadList.Error())
		traceSpan.RecordError(err)
		traceSpan.SetStatus(codes.Error, object.ErrTickerRepositoryReadList.Error())

		return nil, nil, err
	}

	service.GetRuntimeLogger().
		WithFields(fields).
		WithField(object.URIFieldDAOTickers, daoTickers).
		WithField(object.URIFieldDAOCursor, daoCursorer).
		Debug(object.URIEmpty)

	omTickers := make([]om.Tickerer, 0, len(daoTickers))

	for key, daoTicker := range daoTickers {
		service.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldKey, key).
			WithField(object.URIFieldDAOTicker, daoTicker).
			Debug(object.URIEmpty)

		omTickers = append(omTickers, om.NewTicker(
			daoTicker.GetAveragePrice(),
			daoTicker.GetBuy(),
			daoTicker.GetChangePrice(),
			daoTicker.GetChangeRate(),
			daoTicker.GetHigh(),
			daoTicker.GetLast(),
			daoTicker.GetLow(),
			daoTicker.GetMakerCoefficient(),
			daoTicker.GetMakerFeeRate(),
			daoTicker.GetSell(),
			daoTicker.GetSymbol(),
			daoTicker.GetSymbolName(),
			daoTicker.GetTakerCoefficient(),
			daoTicker.GetTakerFeeRate(),
			daoTicker.GetVol(),
			daoTicker.GetVolValue(),
			daoTicker.GetID(),
		))
	}

	service.GetRuntimeLogger().
		WithFields(fields).
		WithField(object.URIFieldOMTickers, omTickers).
		Debug(object.URIEmpty)

	return omTickers, daoCursorer, nil
}
