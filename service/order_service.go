package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Kucoin/kucoin-go-sdk"
	"github.com/ShahoBashoki/kucoin/config"
	"github.com/ShahoBashoki/kucoin/log"
	"github.com/ShahoBashoki/kucoin/object"
	"github.com/ShahoBashoki/kucoin/object/dao"
	"github.com/ShahoBashoki/kucoin/object/dto"
	"github.com/ShahoBashoki/kucoin/object/om"
	"github.com/ShahoBashoki/kucoin/repository"
	"github.com/ShahoBashoki/kucoin/util"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type (
	// OrderServicer is an interface.
	OrderServicer interface {
		// Create is a function.
		Create(
			context.Context,
			om.Orderer,
		) (uuid.UUID, error)
		// DeleteAll is a function.
		DeleteAll(
			context.Context,
		) (time.Time, error)
		// Get is a function.
		Get(
			context.Context,
			uuid.UUID,
		) (om.Orderer, error)
		// GetListFromRepository is a function.
		GetListFromRepository(
			context.Context,
			dao.Paginationer,
			dao.OrderFilterer,
		) ([]om.Orderer, dao.Cursorer, error)
		// GetListFromRemote is a function.
		GetListFromRemote(
			context.Context,
			dto.OrderRequester,
			int64,
		) error
	}

	// GetOrderServicer is an interface.
	GetOrderServicer interface {
		// GetOrderServicer is a function.
		GetOrderServicer() OrderServicer
	}

	orderService struct {
		configConfigger  config.Configger
		repositorier     repository.OrderRepositorier
		logRuntimeLogger log.RuntimeLogger
		servicer         Servicer
		traceTracer      trace.Tracer
		utilUUIDer       util.UUIDer
		kucoinAPIService *kucoin.ApiService
	}
)

var (
	_ GetServicer                     = (*orderService)(nil)
	_ OrderServicer                   = (*orderService)(nil)
	_ WithServicer                    = (*orderService)(nil)
	_ config.GetConfigger             = (*orderService)(nil)
	_ log.GetRuntimeLogger            = (*orderService)(nil)
	_ repository.GetOrderRepositorier = (*orderService)(nil)
	_ util.GetTracer                  = (*orderService)(nil)
	_ util.GetUUIDer                  = (*orderService)(nil)
)

// NewOrderServicer is a function.
func NewOrderServicer(
	configConfigger config.Configger,
	repositorier repository.OrderRepositorier,
	logRuntimeLogger log.RuntimeLogger,
	traceTracer trace.Tracer,
	utilUUIDer util.UUIDer,
	kucoinAPIService *kucoin.ApiService,
) OrderServicer {
	return &orderService{
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
func (service *orderService) GetConfigger() config.Configger {
	return service.configConfigger
}

// GetRuntimeLogger is a function.
func (service *orderService) GetRuntimeLogger() log.RuntimeLogger {
	return service.logRuntimeLogger
}

// GetServicer is a function.
func (service *orderService) GetServicer() Servicer {
	return service.servicer
}

// GetOrderRepositorier is a function.
func (service *orderService) GetOrderRepositorier() repository.OrderRepositorier {
	return service.repositorier
}

// GetTracer is a function.
func (service *orderService) GetTracer() trace.Tracer {
	return service.traceTracer
}

// GetUUIDer is a function.
func (service *orderService) GetUUIDer() util.UUIDer {
	return service.utilUUIDer
}

// GetAPIService is a function.
func (service *orderService) GetAPIService() *kucoin.ApiService {
	return service.kucoinAPIService
}

// WithServicer is a function.
func (service *orderService) WithServicer(
	servicer Servicer,
) {
	service.servicer = servicer
}

// Create is a function.
func (service *orderService) Create(
	ctx context.Context,
	omOrderer om.Orderer,
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
		"name":       "Create",
		"rt_ctx":     utilRuntimeContext,
		"sp_ctx":     utilSpanContext,
		"config":     service.configConfigger,
		"om_orderer": omOrderer,
	}

	service.GetRuntimeLogger().
		WithFields(fields).
		Info(object.URIEmpty)

	daoOrder := dao.NewOrder(
		time.Time{},
		time.Time{},
		sql.NullTime{
			Time:  time.Time{},
			Valid: false,
		},
		uuid.Nil,
		omOrderer.GetChannel(),
		omOrderer.GetClientOID(),
		omOrderer.GetDealFunds(),
		omOrderer.GetDealSize(),
		omOrderer.GetFee(),
		omOrderer.GetFeeCurrency(),
		omOrderer.GetFunds(),
		omOrderer.GetKucoinID(),
		omOrderer.GetKucoinType(),
		omOrderer.GetOPType(),
		omOrderer.GetPrice(),
		omOrderer.GetRemark(),
		omOrderer.GetSide(),
		omOrderer.GetSize(),
		omOrderer.GetStop(),
		omOrderer.GetStopPrice(),
		omOrderer.GetSTP(),
		omOrderer.GetSymbol(),
		omOrderer.GetTags(),
		omOrderer.GetTimeInForce(),
		omOrderer.GetTradeType(),
		omOrderer.GetVisibleSize(),
		omOrderer.GetCancelAfter(),
		omOrderer.GetKucoinCreatedAt(),
		omOrderer.GetCancelExist(),
		omOrderer.GetHidden(),
		omOrderer.GetIceBerg(),
		omOrderer.GetIsActive(),
		omOrderer.GetPostOnly(),
		omOrderer.GetStopTriggered(),
	)

	service.GetRuntimeLogger().
		WithFields(fields).
		WithField(object.URIFieldDAOOrder, daoOrder).
		Debug(object.URIEmpty)

	orderID, err := service.GetOrderRepositorier().Create(ctx, daoOrder)
	if err != nil {
		service.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, err).
			Error(object.ErrOrderRepositoryCreate.Error())
		traceSpan.RecordError(err)
		traceSpan.SetStatus(codes.Error, object.ErrOrderRepositoryCreate.Error())

		return uuid.Nil, err
	}

	service.GetRuntimeLogger().
		WithFields(fields).
		WithField(object.URIFieldOrderID, orderID).
		Debug(object.URIEmpty)

	return orderID, nil
}

// DeleteAll is a function.
func (service *orderService) DeleteAll(
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

	deletedAt, err := service.GetOrderRepositorier().DeleteAll(ctx)
	if err != nil {
		service.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, err).
			Error(object.ErrOrderRepositoryDeleteAll.Error())
		traceSpan.RecordError(err)
		traceSpan.SetStatus(codes.Error, object.ErrOrderRepositoryDeleteAll.Error())

		return time.Time{}, err
	}

	service.GetRuntimeLogger().
		WithFields(fields).
		WithField(object.URIFieldDeletedAt, deletedAt).
		Debug(object.URIEmpty)

	return deletedAt, nil
}

// Get is a function.
func (service *orderService) Get(
	ctx context.Context,
	id uuid.UUID,
) (om.Orderer, error) {
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

	daoOrder, err := service.GetOrderRepositorier().Read(ctx, id)
	if err != nil {
		service.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, err).
			Error(object.ErrOrderRepositoryRead.Error())
		traceSpan.RecordError(err)
		traceSpan.SetStatus(codes.Error, object.ErrOrderRepositoryRead.Error())

		return nil, err
	}

	service.GetRuntimeLogger().
		WithFields(fields).
		WithField(object.URIFieldDAOOrder, daoOrder).
		Debug(object.URIEmpty)

	omOrder := om.NewOrder(
		daoOrder.GetChannel(),
		daoOrder.GetClientOID(),
		daoOrder.GetDealFunds(),
		daoOrder.GetDealSize(),
		daoOrder.GetFee(),
		daoOrder.GetFeeCurrency(),
		daoOrder.GetFunds(),
		daoOrder.GetKucoinID(),
		daoOrder.GetKucoinType(),
		daoOrder.GetOPType(),
		daoOrder.GetPrice(),
		daoOrder.GetRemark(),
		daoOrder.GetSide(),
		daoOrder.GetSize(),
		daoOrder.GetStop(),
		daoOrder.GetStopPrice(),
		daoOrder.GetSTP(),
		daoOrder.GetSymbol(),
		daoOrder.GetTags(),
		daoOrder.GetTimeInForce(),
		daoOrder.GetTradeType(),
		daoOrder.GetVisibleSize(),
		daoOrder.GetCancelAfter(),
		daoOrder.GetKucoinCreatedAt(),
		daoOrder.GetCancelExist(),
		daoOrder.GetHidden(),
		daoOrder.GetIceBerg(),
		daoOrder.GetIsActive(),
		daoOrder.GetPostOnly(),
		daoOrder.GetStopTriggered(),
		daoOrder.GetID(),
	)

	service.GetRuntimeLogger().
		WithFields(fields).
		WithField(object.URIFieldOMOrder, omOrder).
		Debug(object.URIEmpty)

	return omOrder, nil
}

// GetListFromRepository is a function.
func (service *orderService) GetListFromRepository(
	ctx context.Context,
	daoPaginator dao.Paginationer,
	daoOrderFilterer dao.OrderFilterer,
) ([]om.Orderer, dao.Cursorer, error) {
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
		"name":               "GetListFromRepository",
		"rt_ctx":             utilRuntimeContext,
		"sp_ctx":             utilSpanContext,
		"config":             service.configConfigger,
		"dao_paginator":      daoPaginator,
		"dao_order_filterer": daoOrderFilterer,
	}

	service.GetRuntimeLogger().
		WithFields(fields).
		Info(object.URIEmpty)

	daoOrders, daoCursorer, err := service.GetOrderRepositorier().
		ReadList(ctx, daoPaginator, daoOrderFilterer)
	if err != nil {
		service.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, err).
			Error(object.ErrOrderRepositoryReadList.Error())
		traceSpan.RecordError(err)
		traceSpan.SetStatus(codes.Error, object.ErrOrderRepositoryReadList.Error())

		return nil, nil, err
	}

	service.GetRuntimeLogger().
		WithFields(fields).
		WithField(object.URIFieldDAOOrders, daoOrders).
		WithField(object.URIFieldDAOCursor, daoCursorer).
		Debug(object.URIEmpty)

	omOrders := make([]om.Orderer, 0, len(daoOrders))

	for key, daoOrder := range daoOrders {
		service.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldKey, key).
			WithField(object.URIFieldDAOOrder, daoOrder).
			Debug(object.URIEmpty)

		omOrders = append(omOrders, om.NewOrder(
			daoOrder.GetChannel(),
			daoOrder.GetClientOID(),
			daoOrder.GetDealFunds(),
			daoOrder.GetDealSize(),
			daoOrder.GetFee(),
			daoOrder.GetFeeCurrency(),
			daoOrder.GetFunds(),
			daoOrder.GetKucoinID(),
			daoOrder.GetKucoinType(),
			daoOrder.GetOPType(),
			daoOrder.GetPrice(),
			daoOrder.GetRemark(),
			daoOrder.GetSide(),
			daoOrder.GetSize(),
			daoOrder.GetStop(),
			daoOrder.GetStopPrice(),
			daoOrder.GetSTP(),
			daoOrder.GetSymbol(),
			daoOrder.GetTags(),
			daoOrder.GetTimeInForce(),
			daoOrder.GetTradeType(),
			daoOrder.GetVisibleSize(),
			daoOrder.GetCancelAfter(),
			daoOrder.GetKucoinCreatedAt(),
			daoOrder.GetCancelExist(),
			daoOrder.GetHidden(),
			daoOrder.GetIceBerg(),
			daoOrder.GetIsActive(),
			daoOrder.GetPostOnly(),
			daoOrder.GetStopTriggered(),
			daoOrder.GetID(),
		))
	}

	service.GetRuntimeLogger().
		WithFields(fields).
		WithField(object.URIFieldOMOrders, omOrders).
		Debug(object.URIEmpty)

	return omOrders, daoCursorer, nil
}

// GetListFromRemote is a function.
func (service *orderService) GetListFromRemote(
	ctx context.Context,
	dtoOrderRequester dto.OrderRequester,
	currentPage int64,
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
		"name":                "GetListFromRemote",
		"rt_ctx":              utilRuntimeContext,
		"sp_ctx":              utilSpanContext,
		"config":              service.configConfigger,
		"dto_order_requester": dtoOrderRequester,
		"current_page":        currentPage,
	}

	service.GetRuntimeLogger().
		WithFields(fields).
		Info(object.URIEmpty)

	params := map[string]string{}

	for key, value := range dtoOrderRequester.GetMap() {
		service.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldKey, key).
			WithField(object.URIFieldValue, value).
			Debug(object.URIEmpty)

		if newValue, ok := value.(string); ok && value != object.URIEmpty {
			service.GetRuntimeLogger().
				WithFields(fields).
				Debug(`newValue, ok := value.(string); ok && value != object.URIEmpty`)

			params[key] = newValue
		}
	}

	service.GetRuntimeLogger().
		WithFields(fields).
		WithField(object.URIFieldParams, params).
		Debug(object.URIEmpty)

	kucoinPaginationParam := &kucoin.PaginationParam{
		CurrentPage: currentPage,
		PageSize:    service.GetConfigger().GetRuntimeConfigger().GetKucoinPaginationRequestSize(),
	}

	service.GetRuntimeLogger().
		WithFields(fields).
		WithField(object.URIFieldKucoinPaginationParam, kucoinPaginationParam).
		Debug(object.URIEmpty)

	response, err := service.kucoinAPIService.Orders(params, kucoinPaginationParam)
	if err != nil {
		service.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, err).
			Error(object.ErrOrderKucoinServiceGetList.Error())
		traceSpan.RecordError(err)
		traceSpan.SetStatus(codes.Error, object.ErrOrderKucoinServiceGetList.Error())

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
			Error(object.ErrOrderKucoinServiceGetList.Error())
		traceSpan.RecordError(object.ErrOrderKucoinServiceGetList)
		traceSpan.SetStatus(codes.Error, object.ErrOrderKucoinServiceGetList.Error())

		return object.ErrOrderKucoinServiceGetList
	}

	kucoinPaginationModel, err := response.ReadPaginationData(&kucoin.OrdersModel{})
	if err != nil {
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
		WithField(object.URIFieldKucoinPaginationModel, kucoinPaginationModel).
		Debug(object.URIEmpty)

	kucoinOrdersModel := &[]kucoin.OrderModel{}

	if err = json.Unmarshal(
		kucoinPaginationModel.RawItems,
		kucoinOrdersModel,
	); err != nil {
		service.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, err).
			Error(object.ErrUnmarshalJSON.Error())
		traceSpan.RecordError(err)
		traceSpan.SetStatus(codes.Error, object.ErrUnmarshalJSON.Error())

		return err
	}

	service.GetRuntimeLogger().
		WithFields(fields).
		WithField(object.URIFieldKucoinOrdersModel, kucoinOrdersModel).
		Debug(object.URIEmpty)

	for key, value := range *kucoinOrdersModel {
		service.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldKey, key).
			WithField(object.URIFieldValue, value).
			Debug(object.URIEmpty)

		omOrder := om.NewOrder(
			value.Channel,
			value.ClientOid,
			value.DealFunds,
			value.DealSize,
			value.Fee,
			value.FeeCurrency,
			value.Funds,
			value.Id,
			value.Type,
			value.OpType,
			value.Price,
			value.Remark,
			value.Side,
			value.Size,
			value.Stop,
			value.StopPrice,
			value.Stp,
			value.Symbol,
			value.Tags,
			value.TimeInForce,
			value.TradeType,
			value.VisibleSize,
			uint32(value.CancelAfter),
			uint32(value.CreatedAt),
			value.CancelExist,
			value.Hidden,
			value.IceBerg,
			value.IsActive,
			value.PostOnly,
			value.StopTriggered,
			uuid.Nil,
		)

		orderID, errOrderCreate := service.GetServicer().GetOrderServicer().Create(ctx, omOrder)
		if errOrderCreate != nil {
			service.GetRuntimeLogger().
				WithFields(fields).
				WithField(object.URIFieldError, errOrderCreate).
				Error(object.ErrOrderServiceCreate.Error())
			traceSpan.RecordError(errOrderCreate)
			traceSpan.SetStatus(codes.Error, object.ErrOrderServiceCreate.Error())

			return errOrderCreate
		}

		service.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldOrderID, orderID).
			Debug(object.URIEmpty)
	}

	if kucoinPaginationModel.CurrentPage < kucoinPaginationModel.TotalPage {
		service.GetRuntimeLogger().
			WithFields(fields).
			Debug(`kucoinPaginationModel.CurrentPage < kucoinPaginationModel.TotalPage`)

		if err = service.GetServicer().
			GetOrderServicer().
			GetListFromRemote(ctx, dtoOrderRequester, currentPage+1); err != nil {
			service.GetRuntimeLogger().
				WithFields(fields).
				WithField(object.URIFieldError, err).
				Error(object.ErrOrderServiceGetListFromRemote.Error())
			traceSpan.RecordError(err)
			traceSpan.SetStatus(codes.Error, object.ErrOrderServiceGetListFromRemote.Error())

			return err
		}
	}

	return nil
}
