package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ShahoBashoki/kucoin/config"
	"github.com/ShahoBashoki/kucoin/log"
	"github.com/ShahoBashoki/kucoin/object"
	"github.com/ShahoBashoki/kucoin/object/dao"
	"github.com/ShahoBashoki/kucoin/util"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

type (
	// OrderRepositorier is a interface.
	OrderRepositorier interface {
		DAORepositorier[dao.Orderer, dao.OrderFilterer]
	}

	// GetOrderRepositorier is an interface.
	GetOrderRepositorier interface {
		// GetOrderRepositorier is a function.
		GetOrderRepositorier() OrderRepositorier
	}

	orderRepository struct {
		configConfigger  config.Configger
		gormDB           *gorm.DB
		logRuntimeLogger log.RuntimeLogger
		objectTimer      object.Timer
		traceTracer      trace.Tracer
		utilUUIDer       util.UUIDer
	}

	orderRepositoryOptioner interface {
		apply(*orderRepository)
	}

	orderRepositoryOptionerFunc func(*orderRepository)
)

var (
	_ OrderRepositorier    = (*orderRepository)(nil)
	_ GetDB                = (*orderRepository)(nil)
	_ config.GetConfigger  = (*orderRepository)(nil)
	_ log.GetRuntimeLogger = (*orderRepository)(nil)
	_ object.GetTimer      = (*orderRepository)(nil)
	_ util.GetTracer       = (*orderRepository)(nil)
	_ util.GetUUIDer       = (*orderRepository)(nil)
)

// NewOrderRepository is a function.
func NewOrderRepository(
	configConfigger config.Configger,
	logRuntimeLogger log.RuntimeLogger,
	traceTracer trace.Tracer,
	utilUUIDer util.UUIDer,
	optioners ...orderRepositoryOptioner,
) *orderRepository {
	orderRepository := &orderRepository{
		configConfigger:  configConfigger,
		gormDB:           nil,
		logRuntimeLogger: logRuntimeLogger,
		objectTimer:      nil,
		traceTracer:      traceTracer,
		utilUUIDer:       utilUUIDer,
	}

	return orderRepository.WithOptioners(optioners...)
}

// WithOrderRepositoryTimer is a function.
func WithOrderRepositoryTimer(
	objectTimer object.Timer,
) orderRepositoryOptioner {
	return orderRepositoryOptionerFunc(func(
		config *orderRepository,
	) {
		config.objectTimer = objectTimer
	})
}

// WithOrderRepositoryDB is a function.
func WithOrderRepositoryDB(
	gormDB *gorm.DB,
) orderRepositoryOptioner {
	return orderRepositoryOptionerFunc(func(
		config *orderRepository,
	) {
		config.gormDB = gormDB.
			Table(object.URITableKucoinOrder).
			Session(&gorm.Session{
				DryRun:                   false,
				PrepareStmt:              true,
				NewDB:                    true,
				Initialized:              false,
				SkipHooks:                true,
				SkipDefaultTransaction:   true,
				DisableNestedTransaction: true,
				AllowGlobalUpdate:        false,
				FullSaveAssociations:     false,
				QueryFields:              true,
				Context:                  nil,
				Logger:                   nil,
				NowFunc:                  nil,
				CreateBatchSize:          0,
			})
	})
}

// GetDB is a function.
func (repository *orderRepository) GetDB() *gorm.DB {
	return repository.gormDB
}

// GetConfigger is a function.
func (repository *orderRepository) GetConfigger() config.Configger {
	return repository.configConfigger
}

// GetRuntimeLogger is a function.
func (repository *orderRepository) GetRuntimeLogger() log.RuntimeLogger {
	return repository.logRuntimeLogger
}

// GetTimer is a function.
func (repository *orderRepository) GetTimer() object.Timer {
	return repository.objectTimer
}

// GetTracer is a function.
func (repository *orderRepository) GetTracer() trace.Tracer {
	return repository.traceTracer
}

// GetUUIDer is a function.
func (repository *orderRepository) GetUUIDer() util.UUIDer {
	return repository.utilUUIDer
}

// Create is a function.
func (repository *orderRepository) Create(
	ctx context.Context,
	daoOrderer dao.Orderer,
) (uuid.UUID, error) {
	var traceSpan trace.Span

	ctx, traceSpan = repository.GetTracer().Start(
		ctx,
		"Create",
		trace.WithSpanKind(trace.SpanKindProducer),
	)
	defer traceSpan.End()

	utilRuntimeContext := util.NewRuntimeContext(ctx, repository.GetUUIDer())
	utilSpanContext := util.NewSpanContext(traceSpan)
	fields := map[string]any{
		"name":        "Create",
		"rt_ctx":      utilRuntimeContext,
		"sp_ctx":      utilSpanContext,
		"config":      repository.GetConfigger(),
		"dao_orderer": daoOrderer,
	}

	repository.GetRuntimeLogger().
		WithFields(fields).
		Info(object.URIEmpty)

	id, err := repository.GetUUIDer().NewRandom()
	if err != nil {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, err).
			Error(object.ErrUUIDerNewRandom.Error())
		traceSpan.RecordError(err)
		traceSpan.SetStatus(codes.Error, object.ErrUUIDerNewRandom.Error())

		return uuid.Nil, err
	}

	repository.GetRuntimeLogger().
		WithFields(fields).
		WithField(object.URIFieldID, id).
		Debug(object.URIEmpty)

	nowUTC := repository.GetTimer().NowUTC()

	repository.GetRuntimeLogger().
		WithFields(fields).
		WithField(object.URIFieldNowUTC, nowUTC).
		Debug(object.URIEmpty)

	daoOrder := dao.NewOrder(
		nowUTC,
		nowUTC,
		sql.NullTime{
			Time:  time.Time{},
			Valid: false,
		},
		id,
		daoOrderer.GetChannel(),
		daoOrderer.GetClientOID(),
		daoOrderer.GetDealFunds(),
		daoOrderer.GetDealSize(),
		daoOrderer.GetFee(),
		daoOrderer.GetFeeCurrency(),
		daoOrderer.GetFunds(),
		daoOrderer.GetKucoinID(),
		daoOrderer.GetKucoinType(),
		daoOrderer.GetOPType(),
		daoOrderer.GetPrice(),
		daoOrderer.GetRemark(),
		daoOrderer.GetSide(),
		daoOrderer.GetSize(),
		daoOrderer.GetStop(),
		daoOrderer.GetStopPrice(),
		daoOrderer.GetSTP(),
		daoOrderer.GetSymbol(),
		daoOrderer.GetTags(),
		daoOrderer.GetTimeInForce(),
		daoOrderer.GetTradeType(),
		daoOrderer.GetVisibleSize(),
		daoOrderer.GetCancelAfter(),
		daoOrderer.GetKucoinCreatedAt(),
		daoOrderer.GetCancelExist(),
		daoOrderer.GetHidden(),
		daoOrderer.GetIceBerg(),
		daoOrderer.GetIsActive(),
		daoOrderer.GetPostOnly(),
		daoOrderer.GetStopTriggered(),
	)

	repository.GetRuntimeLogger().
		WithFields(fields).
		WithField(object.URIFieldDAOOrder, daoOrder).
		Debug(object.URIEmpty)

	gormDB := repository.GetDB().
		WithContext(ctx).
		Create(daoOrder.GetMap())
	if err = gormDB.Error; err != nil {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, err).
			Error(object.ErrOrderRepositoryCreate.Error())
		traceSpan.RecordError(err)
		traceSpan.SetStatus(codes.Error, object.ErrOrderRepositoryCreate.Error())

		return uuid.Nil, err
	}

	return daoOrder.GetID(), nil
}

// Delete is a function.
func (repository *orderRepository) Delete(
	ctx context.Context,
	id uuid.UUID,
) (time.Time, error) {
	var traceSpan trace.Span

	ctx, traceSpan = repository.GetTracer().Start(
		ctx,
		"Delete",
		trace.WithSpanKind(trace.SpanKindProducer),
	)
	defer traceSpan.End()

	utilRuntimeContext := util.NewRuntimeContext(ctx, repository.GetUUIDer())
	utilSpanContext := util.NewSpanContext(traceSpan)
	fields := map[string]any{
		"name":   "Delete",
		"rt_ctx": utilRuntimeContext,
		"sp_ctx": utilSpanContext,
		"config": repository.GetConfigger(),
		"id":     id,
	}

	repository.GetRuntimeLogger().
		WithFields(fields).
		Info(object.URIEmpty)

	nowUTC := repository.GetTimer().NowUTC()

	repository.GetRuntimeLogger().
		WithFields(fields).
		WithField(object.URIFieldNowUTC, nowUTC).
		Debug(object.URIEmpty)

	gormDB := repository.GetDB().
		WithContext(ctx).
		Where(map[string]any{
			"id": id,
		}).
		Updates(map[string]any{
			"deleted_at": sql.NullTime{
				Time:  nowUTC,
				Valid: true,
			},
		})
	if err := gormDB.Error; err != nil {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, err).
			Error(object.ErrOrderRepositoryDelete.Error())
		traceSpan.RecordError(err)
		traceSpan.SetStatus(codes.Error, object.ErrOrderRepositoryDelete.Error())

		return time.Time{}, err
	}

	if gormDB.RowsAffected == 0 {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, object.ErrOrderRepositoryDelete).
			Error(object.ErrOrderRepositoryDelete.Error())
		traceSpan.RecordError(object.ErrOrderRepositoryDelete)
		traceSpan.SetStatus(codes.Error, object.ErrOrderRepositoryDelete.Error())

		return time.Time{}, object.ErrOrderRepositoryDelete
	}

	return nowUTC, nil
}

// Read is a function.
func (repository *orderRepository) Read(
	ctx context.Context,
	id uuid.UUID,
) (dao.Orderer, error) {
	var traceSpan trace.Span

	ctx, traceSpan = repository.GetTracer().Start(
		ctx,
		"Read",
		trace.WithSpanKind(trace.SpanKindProducer),
	)
	defer traceSpan.End()

	utilRuntimeContext := util.NewRuntimeContext(ctx, repository.GetUUIDer())
	utilSpanContext := util.NewSpanContext(traceSpan)
	fields := map[string]any{
		"name":   "Read",
		"rt_ctx": utilRuntimeContext,
		"sp_ctx": utilSpanContext,
		"config": repository.GetConfigger(),
		"id":     id,
	}

	repository.GetRuntimeLogger().
		WithFields(fields).
		Info(object.URIEmpty)

	result := map[string]any{}

	gormDB := repository.GetDB().
		WithContext(ctx).
		Where(map[string]any{
			"id":         id,
			"deleted_at": nil,
		}).
		Select(fmt.Sprintf("%s.*", object.URITableKucoinOrder)).
		Find(result)
	if err := gormDB.Error; err != nil {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, err).
			Error(object.ErrOrderRepositoryRead.Error())
		traceSpan.RecordError(err)
		traceSpan.SetStatus(codes.Error, object.ErrOrderRepositoryRead.Error())

		return nil, err
	}

	if gormDB.RowsAffected == 0 {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, object.ErrOrderRepositoryRead).
			Error(object.ErrOrderRepositoryRead.Error())
		traceSpan.RecordError(object.ErrOrderRepositoryRead)
		traceSpan.SetStatus(codes.Error, object.ErrOrderRepositoryRead.Error())

		return nil, object.ErrOrderRepositoryRead
	}

	repository.GetRuntimeLogger().
		WithFields(fields).
		WithField(object.URIFieldResult, result).
		Debug(object.URIEmpty)

	createdAT, ok := result["created_at"].(time.Time)
	if !ok {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, object.ErrTypeAssertion).
			Error(object.ErrTypeAssertion.Error())
		traceSpan.RecordError(object.ErrTypeAssertion)
		traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

		return nil, object.ErrTypeAssertion
	}

	updatedAT, ok := result["updated_at"].(time.Time)
	if !ok {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, object.ErrTypeAssertion).
			Error(object.ErrTypeAssertion.Error())
		traceSpan.RecordError(object.ErrTypeAssertion)
		traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

		return nil, object.ErrTypeAssertion
	}

	channel, ok := result["channel"].(string)
	if !ok {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, object.ErrTypeAssertion).
			Error(object.ErrTypeAssertion.Error())
		traceSpan.RecordError(object.ErrTypeAssertion)
		traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

		return nil, object.ErrTypeAssertion
	}

	clientOID, ok := result["client_oid"].(string)
	if !ok {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, object.ErrTypeAssertion).
			Error(object.ErrTypeAssertion.Error())
		traceSpan.RecordError(object.ErrTypeAssertion)
		traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

		return nil, object.ErrTypeAssertion
	}

	dealFunds, ok := result["deal_funds"].(string)
	if !ok {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, object.ErrTypeAssertion).
			Error(object.ErrTypeAssertion.Error())
		traceSpan.RecordError(object.ErrTypeAssertion)
		traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

		return nil, object.ErrTypeAssertion
	}

	dealSize, ok := result["deal_size"].(string)
	if !ok {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, object.ErrTypeAssertion).
			Error(object.ErrTypeAssertion.Error())
		traceSpan.RecordError(object.ErrTypeAssertion)
		traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

		return nil, object.ErrTypeAssertion
	}

	fee, ok := result["fee"].(string)
	if !ok {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, object.ErrTypeAssertion).
			Error(object.ErrTypeAssertion.Error())
		traceSpan.RecordError(object.ErrTypeAssertion)
		traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

		return nil, object.ErrTypeAssertion
	}

	feeCurrency, ok := result["fee_currency"].(string)
	if !ok {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, object.ErrTypeAssertion).
			Error(object.ErrTypeAssertion.Error())
		traceSpan.RecordError(object.ErrTypeAssertion)
		traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

		return nil, object.ErrTypeAssertion
	}

	funds, ok := result["funds"].(string)
	if !ok {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, object.ErrTypeAssertion).
			Error(object.ErrTypeAssertion.Error())
		traceSpan.RecordError(object.ErrTypeAssertion)
		traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

		return nil, object.ErrTypeAssertion
	}

	kucoinID, ok := result["kucoin_id"].(string)
	if !ok {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, object.ErrTypeAssertion).
			Error(object.ErrTypeAssertion.Error())
		traceSpan.RecordError(object.ErrTypeAssertion)
		traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

		return nil, object.ErrTypeAssertion
	}

	kucoinType, ok := result["kucoin_type"].(string)
	if !ok {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, object.ErrTypeAssertion).
			Error(object.ErrTypeAssertion.Error())
		traceSpan.RecordError(object.ErrTypeAssertion)
		traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

		return nil, object.ErrTypeAssertion
	}

	opType, ok := result["op_type"].(string)
	if !ok {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, object.ErrTypeAssertion).
			Error(object.ErrTypeAssertion.Error())
		traceSpan.RecordError(object.ErrTypeAssertion)
		traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

		return nil, object.ErrTypeAssertion
	}

	price, ok := result["price"].(string)
	if !ok {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, object.ErrTypeAssertion).
			Error(object.ErrTypeAssertion.Error())
		traceSpan.RecordError(object.ErrTypeAssertion)
		traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

		return nil, object.ErrTypeAssertion
	}

	remark, ok := result["remark"].(string)
	if !ok {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, object.ErrTypeAssertion).
			Error(object.ErrTypeAssertion.Error())
		traceSpan.RecordError(object.ErrTypeAssertion)
		traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

		return nil, object.ErrTypeAssertion
	}

	side, ok := result["side"].(string)
	if !ok {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, object.ErrTypeAssertion).
			Error(object.ErrTypeAssertion.Error())
		traceSpan.RecordError(object.ErrTypeAssertion)
		traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

		return nil, object.ErrTypeAssertion
	}

	size, ok := result["size"].(string)
	if !ok {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, object.ErrTypeAssertion).
			Error(object.ErrTypeAssertion.Error())
		traceSpan.RecordError(object.ErrTypeAssertion)
		traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

		return nil, object.ErrTypeAssertion
	}

	stop, ok := result["stop"].(string)
	if !ok {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, object.ErrTypeAssertion).
			Error(object.ErrTypeAssertion.Error())
		traceSpan.RecordError(object.ErrTypeAssertion)
		traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

		return nil, object.ErrTypeAssertion
	}

	stopPrice, ok := result["stop_price"].(string)
	if !ok {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, object.ErrTypeAssertion).
			Error(object.ErrTypeAssertion.Error())
		traceSpan.RecordError(object.ErrTypeAssertion)
		traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

		return nil, object.ErrTypeAssertion
	}

	stp, ok := result["stp"].(string)
	if !ok {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, object.ErrTypeAssertion).
			Error(object.ErrTypeAssertion.Error())
		traceSpan.RecordError(object.ErrTypeAssertion)
		traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

		return nil, object.ErrTypeAssertion
	}

	symbol, ok := result["symbol"].(string)
	if !ok {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, object.ErrTypeAssertion).
			Error(object.ErrTypeAssertion.Error())
		traceSpan.RecordError(object.ErrTypeAssertion)
		traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

		return nil, object.ErrTypeAssertion
	}

	tags, ok := result["tags"].(string)
	if !ok {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, object.ErrTypeAssertion).
			Error(object.ErrTypeAssertion.Error())
		traceSpan.RecordError(object.ErrTypeAssertion)
		traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

		return nil, object.ErrTypeAssertion
	}

	timeInForce, ok := result["time_in_force"].(string)
	if !ok {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, object.ErrTypeAssertion).
			Error(object.ErrTypeAssertion.Error())
		traceSpan.RecordError(object.ErrTypeAssertion)
		traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

		return nil, object.ErrTypeAssertion
	}

	tradeType, ok := result["trade_type"].(string)
	if !ok {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, object.ErrTypeAssertion).
			Error(object.ErrTypeAssertion.Error())
		traceSpan.RecordError(object.ErrTypeAssertion)
		traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

		return nil, object.ErrTypeAssertion
	}

	visibleSize, ok := result["visible_size"].(string)
	if !ok {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, object.ErrTypeAssertion).
			Error(object.ErrTypeAssertion.Error())
		traceSpan.RecordError(object.ErrTypeAssertion)
		traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

		return nil, object.ErrTypeAssertion
	}

	cancelAfter, ok := result["cancel_after"].(int64)
	if !ok {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, object.ErrTypeAssertion).
			Error(object.ErrTypeAssertion.Error())
		traceSpan.RecordError(object.ErrTypeAssertion)
		traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

		return nil, object.ErrTypeAssertion
	}

	kucoinCreatedAt, ok := result["kucoin_created_at"].(int64)
	if !ok {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, object.ErrTypeAssertion).
			Error(object.ErrTypeAssertion.Error())
		traceSpan.RecordError(object.ErrTypeAssertion)
		traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

		return nil, object.ErrTypeAssertion
	}

	cancelExist, ok := result["cancel_exist"].(bool)
	if !ok {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, object.ErrTypeAssertion).
			Error(object.ErrTypeAssertion.Error())
		traceSpan.RecordError(object.ErrTypeAssertion)
		traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

		return nil, object.ErrTypeAssertion
	}

	hidden, ok := result["hidden"].(bool)
	if !ok {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, object.ErrTypeAssertion).
			Error(object.ErrTypeAssertion.Error())
		traceSpan.RecordError(object.ErrTypeAssertion)
		traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

		return nil, object.ErrTypeAssertion
	}

	iceBerg, ok := result["ice_berg"].(bool)
	if !ok {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, object.ErrTypeAssertion).
			Error(object.ErrTypeAssertion.Error())
		traceSpan.RecordError(object.ErrTypeAssertion)
		traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

		return nil, object.ErrTypeAssertion
	}

	isActive, ok := result["is_active"].(bool)
	if !ok {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, object.ErrTypeAssertion).
			Error(object.ErrTypeAssertion.Error())
		traceSpan.RecordError(object.ErrTypeAssertion)
		traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

		return nil, object.ErrTypeAssertion
	}

	postOnly, ok := result["post_only"].(bool)
	if !ok {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, object.ErrTypeAssertion).
			Error(object.ErrTypeAssertion.Error())
		traceSpan.RecordError(object.ErrTypeAssertion)
		traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

		return nil, object.ErrTypeAssertion
	}

	stopTriggered, ok := result["stop_triggered"].(bool)
	if !ok {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, object.ErrTypeAssertion).
			Error(object.ErrTypeAssertion.Error())
		traceSpan.RecordError(object.ErrTypeAssertion)
		traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

		return nil, object.ErrTypeAssertion
	}

	daoOrder := dao.NewOrder(
		createdAT,
		updatedAT,
		sql.NullTime{
			Time:  time.Time{},
			Valid: false,
		},
		id,
		channel,
		clientOID,
		dealFunds,
		dealSize,
		fee,
		feeCurrency,
		funds,
		kucoinID,
		kucoinType,
		opType,
		price,
		remark,
		side,
		size,
		stop,
		stopPrice,
		stp,
		symbol,
		tags,
		timeInForce,
		tradeType,
		visibleSize,
		uint32(cancelAfter),
		uint32(kucoinCreatedAt),
		cancelExist,
		hidden,
		iceBerg,
		isActive,
		postOnly,
		stopTriggered,
	)

	repository.GetRuntimeLogger().
		WithFields(fields).
		WithField(object.URIFieldDAOOrder, daoOrder).
		Debug(object.URIEmpty)

	return daoOrder, nil
}

// ReadList is a function.
func (repository *orderRepository) ReadList(
	ctx context.Context,
	daoPaginationer dao.Paginationer,
	daoOrderFilterer dao.OrderFilterer,
) ([]dao.Orderer, dao.Cursorer, error) {
	var traceSpan trace.Span

	ctx, traceSpan = repository.GetTracer().Start(
		ctx,
		"ReadList",
		trace.WithSpanKind(trace.SpanKindProducer),
	)
	defer traceSpan.End()

	utilRuntimeContext := util.NewRuntimeContext(ctx, repository.GetUUIDer())
	utilSpanContext := util.NewSpanContext(traceSpan)
	fields := map[string]any{
		"name":               "ReadList",
		"rt_ctx":             utilRuntimeContext,
		"sp_ctx":             utilSpanContext,
		"config":             repository.GetConfigger(),
		"dao_paginationer":   daoPaginationer,
		"dao_order_filterer": daoOrderFilterer,
	}

	repository.GetRuntimeLogger().
		WithFields(fields).
		Info(object.URIEmpty)

	result := make([]map[string]any, 0, daoPaginationer.GetLimit()+1)

	gormDB := repository.GetDB().
		WithContext(ctx).
		Scopes(
			daoOrderFilterer.Filter,
			daoPaginationer.Pagination(object.URITableKucoinOrder),
		).
		Where(map[string]any{
			"deleted_at": nil,
		}).
		Select(fmt.Sprintf("%s.*", object.URITableKucoinOrder)).
		Find(&result)
	if err := gormDB.Error; err != nil {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, err).
			Error(object.ErrOrderRepositoryReadList.Error())
		traceSpan.RecordError(err)
		traceSpan.SetStatus(codes.Error, object.ErrOrderRepositoryReadList.Error())

		return nil, nil, err
	}

	repository.GetRuntimeLogger().
		WithFields(fields).
		WithField(object.URIFieldResult, result).
		Debug(object.URIEmpty)

	daoOrderers := make([]dao.Orderer, 0, daoPaginationer.GetLimit())

	for key, value := range result {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldKey, key).
			WithField(object.URIFieldValue, value).
			Debug(object.URIEmpty)

		if uint32(key) == daoPaginationer.GetLimit() {
			repository.GetRuntimeLogger().
				WithFields(fields).
				Debug(`uint32(key) == daoPaginationer.GetLimit()`)

			break
		}

		id, err := repository.GetUUIDer().Parse(value["id"].(string))
		if err != nil {
			repository.logRuntimeLogger.
				WithFields(fields).
				WithField(object.URIFieldError, err).
				Error(object.ErrUUIDerParse.Error())
			traceSpan.RecordError(err)
			traceSpan.SetStatus(codes.Error, object.ErrUUIDerParse.Error())

			return nil, nil, err
		}

		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldID, id).
			Debug(object.URIEmpty)

		createdAT, ok := value["created_at"].(time.Time)
		if !ok {
			repository.GetRuntimeLogger().
				WithFields(fields).
				WithField(object.URIFieldError, object.ErrTypeAssertion).
				Error(object.ErrTypeAssertion.Error())
			traceSpan.RecordError(object.ErrTypeAssertion)
			traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

			return nil, nil, object.ErrTypeAssertion
		}

		updatedAT, ok := value["updated_at"].(time.Time)
		if !ok {
			repository.GetRuntimeLogger().
				WithFields(fields).
				WithField(object.URIFieldError, object.ErrTypeAssertion).
				Error(object.ErrTypeAssertion.Error())
			traceSpan.RecordError(object.ErrTypeAssertion)
			traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

			return nil, nil, object.ErrTypeAssertion
		}

		channel, ok := value["channel"].(string)
		if !ok {
			repository.GetRuntimeLogger().
				WithFields(fields).
				WithField(object.URIFieldError, object.ErrTypeAssertion).
				Error(object.ErrTypeAssertion.Error())
			traceSpan.RecordError(object.ErrTypeAssertion)
			traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

			return nil, nil, object.ErrTypeAssertion
		}

		clientOID, ok := value["client_oid"].(string)
		if !ok {
			repository.GetRuntimeLogger().
				WithFields(fields).
				WithField(object.URIFieldError, object.ErrTypeAssertion).
				Error(object.ErrTypeAssertion.Error())
			traceSpan.RecordError(object.ErrTypeAssertion)
			traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

			return nil, nil, object.ErrTypeAssertion
		}

		dealFunds, ok := value["deal_funds"].(string)
		if !ok {
			repository.GetRuntimeLogger().
				WithFields(fields).
				WithField(object.URIFieldError, object.ErrTypeAssertion).
				Error(object.ErrTypeAssertion.Error())
			traceSpan.RecordError(object.ErrTypeAssertion)
			traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

			return nil, nil, object.ErrTypeAssertion
		}

		dealSize, ok := value["deal_size"].(string)
		if !ok {
			repository.GetRuntimeLogger().
				WithFields(fields).
				WithField(object.URIFieldError, object.ErrTypeAssertion).
				Error(object.ErrTypeAssertion.Error())
			traceSpan.RecordError(object.ErrTypeAssertion)
			traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

			return nil, nil, object.ErrTypeAssertion
		}

		fee, ok := value["fee"].(string)
		if !ok {
			repository.GetRuntimeLogger().
				WithFields(fields).
				WithField(object.URIFieldError, object.ErrTypeAssertion).
				Error(object.ErrTypeAssertion.Error())
			traceSpan.RecordError(object.ErrTypeAssertion)
			traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

			return nil, nil, object.ErrTypeAssertion
		}

		feeCurrency, ok := value["fee_currency"].(string)
		if !ok {
			repository.GetRuntimeLogger().
				WithFields(fields).
				WithField(object.URIFieldError, object.ErrTypeAssertion).
				Error(object.ErrTypeAssertion.Error())
			traceSpan.RecordError(object.ErrTypeAssertion)
			traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

			return nil, nil, object.ErrTypeAssertion
		}

		funds, ok := value["funds"].(string)
		if !ok {
			repository.GetRuntimeLogger().
				WithFields(fields).
				WithField(object.URIFieldError, object.ErrTypeAssertion).
				Error(object.ErrTypeAssertion.Error())
			traceSpan.RecordError(object.ErrTypeAssertion)
			traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

			return nil, nil, object.ErrTypeAssertion
		}

		kucoinID, ok := value["kucoin_id"].(string)
		if !ok {
			repository.GetRuntimeLogger().
				WithFields(fields).
				WithField(object.URIFieldError, object.ErrTypeAssertion).
				Error(object.ErrTypeAssertion.Error())
			traceSpan.RecordError(object.ErrTypeAssertion)
			traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

			return nil, nil, object.ErrTypeAssertion
		}

		kucoinType, ok := value["kucoin_type"].(string)
		if !ok {
			repository.GetRuntimeLogger().
				WithFields(fields).
				WithField(object.URIFieldError, object.ErrTypeAssertion).
				Error(object.ErrTypeAssertion.Error())
			traceSpan.RecordError(object.ErrTypeAssertion)
			traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

			return nil, nil, object.ErrTypeAssertion
		}

		opType, ok := value["op_type"].(string)
		if !ok {
			repository.GetRuntimeLogger().
				WithFields(fields).
				WithField(object.URIFieldError, object.ErrTypeAssertion).
				Error(object.ErrTypeAssertion.Error())
			traceSpan.RecordError(object.ErrTypeAssertion)
			traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

			return nil, nil, object.ErrTypeAssertion
		}

		price, ok := value["price"].(string)
		if !ok {
			repository.GetRuntimeLogger().
				WithFields(fields).
				WithField(object.URIFieldError, object.ErrTypeAssertion).
				Error(object.ErrTypeAssertion.Error())
			traceSpan.RecordError(object.ErrTypeAssertion)
			traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

			return nil, nil, object.ErrTypeAssertion
		}

		remark, ok := value["remark"].(string)
		if !ok {
			repository.GetRuntimeLogger().
				WithFields(fields).
				WithField(object.URIFieldError, object.ErrTypeAssertion).
				Error(object.ErrTypeAssertion.Error())
			traceSpan.RecordError(object.ErrTypeAssertion)
			traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

			return nil, nil, object.ErrTypeAssertion
		}

		side, ok := value["side"].(string)
		if !ok {
			repository.GetRuntimeLogger().
				WithFields(fields).
				WithField(object.URIFieldError, object.ErrTypeAssertion).
				Error(object.ErrTypeAssertion.Error())
			traceSpan.RecordError(object.ErrTypeAssertion)
			traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

			return nil, nil, object.ErrTypeAssertion
		}

		size, ok := value["size"].(string)
		if !ok {
			repository.GetRuntimeLogger().
				WithFields(fields).
				WithField(object.URIFieldError, object.ErrTypeAssertion).
				Error(object.ErrTypeAssertion.Error())
			traceSpan.RecordError(object.ErrTypeAssertion)
			traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

			return nil, nil, object.ErrTypeAssertion
		}

		stop, ok := value["stop"].(string)
		if !ok {
			repository.GetRuntimeLogger().
				WithFields(fields).
				WithField(object.URIFieldError, object.ErrTypeAssertion).
				Error(object.ErrTypeAssertion.Error())
			traceSpan.RecordError(object.ErrTypeAssertion)
			traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

			return nil, nil, object.ErrTypeAssertion
		}

		stopPrice, ok := value["stop_price"].(string)
		if !ok {
			repository.GetRuntimeLogger().
				WithFields(fields).
				WithField(object.URIFieldError, object.ErrTypeAssertion).
				Error(object.ErrTypeAssertion.Error())
			traceSpan.RecordError(object.ErrTypeAssertion)
			traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

			return nil, nil, object.ErrTypeAssertion
		}

		stp, ok := value["stp"].(string)
		if !ok {
			repository.GetRuntimeLogger().
				WithFields(fields).
				WithField(object.URIFieldError, object.ErrTypeAssertion).
				Error(object.ErrTypeAssertion.Error())
			traceSpan.RecordError(object.ErrTypeAssertion)
			traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

			return nil, nil, object.ErrTypeAssertion
		}

		symbol, ok := value["symbol"].(string)
		if !ok {
			repository.GetRuntimeLogger().
				WithFields(fields).
				WithField(object.URIFieldError, object.ErrTypeAssertion).
				Error(object.ErrTypeAssertion.Error())
			traceSpan.RecordError(object.ErrTypeAssertion)
			traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

			return nil, nil, object.ErrTypeAssertion
		}

		tags, ok := value["tags"].(string)
		if !ok {
			repository.GetRuntimeLogger().
				WithFields(fields).
				WithField(object.URIFieldError, object.ErrTypeAssertion).
				Error(object.ErrTypeAssertion.Error())
			traceSpan.RecordError(object.ErrTypeAssertion)
			traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

			return nil, nil, object.ErrTypeAssertion
		}

		timeInForce, ok := value["time_in_force"].(string)
		if !ok {
			repository.GetRuntimeLogger().
				WithFields(fields).
				WithField(object.URIFieldError, object.ErrTypeAssertion).
				Error(object.ErrTypeAssertion.Error())
			traceSpan.RecordError(object.ErrTypeAssertion)
			traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

			return nil, nil, object.ErrTypeAssertion
		}

		tradeType, ok := value["trade_type"].(string)
		if !ok {
			repository.GetRuntimeLogger().
				WithFields(fields).
				WithField(object.URIFieldError, object.ErrTypeAssertion).
				Error(object.ErrTypeAssertion.Error())
			traceSpan.RecordError(object.ErrTypeAssertion)
			traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

			return nil, nil, object.ErrTypeAssertion
		}

		visibleSize, ok := value["visible_size"].(string)
		if !ok {
			repository.GetRuntimeLogger().
				WithFields(fields).
				WithField(object.URIFieldError, object.ErrTypeAssertion).
				Error(object.ErrTypeAssertion.Error())
			traceSpan.RecordError(object.ErrTypeAssertion)
			traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

			return nil, nil, object.ErrTypeAssertion
		}

		cancelAfter, ok := value["cancel_after"].(int64)
		if !ok {
			repository.GetRuntimeLogger().
				WithFields(fields).
				WithField(object.URIFieldError, object.ErrTypeAssertion).
				Error(object.ErrTypeAssertion.Error())
			traceSpan.RecordError(object.ErrTypeAssertion)
			traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

			return nil, nil, object.ErrTypeAssertion
		}

		kucoinCreatedAt, ok := value["kucoin_created_at"].(int64)
		if !ok {
			repository.GetRuntimeLogger().
				WithFields(fields).
				WithField(object.URIFieldError, object.ErrTypeAssertion).
				Error(object.ErrTypeAssertion.Error())
			traceSpan.RecordError(object.ErrTypeAssertion)
			traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

			return nil, nil, object.ErrTypeAssertion
		}

		cancelExist, ok := value["cancel_exist"].(bool)
		if !ok {
			repository.GetRuntimeLogger().
				WithFields(fields).
				WithField(object.URIFieldError, object.ErrTypeAssertion).
				Error(object.ErrTypeAssertion.Error())
			traceSpan.RecordError(object.ErrTypeAssertion)
			traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

			return nil, nil, object.ErrTypeAssertion
		}

		hidden, ok := value["hidden"].(bool)
		if !ok {
			repository.GetRuntimeLogger().
				WithFields(fields).
				WithField(object.URIFieldError, object.ErrTypeAssertion).
				Error(object.ErrTypeAssertion.Error())
			traceSpan.RecordError(object.ErrTypeAssertion)
			traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

			return nil, nil, object.ErrTypeAssertion
		}

		iceBerg, ok := value["ice_berg"].(bool)
		if !ok {
			repository.GetRuntimeLogger().
				WithFields(fields).
				WithField(object.URIFieldError, object.ErrTypeAssertion).
				Error(object.ErrTypeAssertion.Error())
			traceSpan.RecordError(object.ErrTypeAssertion)
			traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

			return nil, nil, object.ErrTypeAssertion
		}

		isActive, ok := value["is_active"].(bool)
		if !ok {
			repository.GetRuntimeLogger().
				WithFields(fields).
				WithField(object.URIFieldError, object.ErrTypeAssertion).
				Error(object.ErrTypeAssertion.Error())
			traceSpan.RecordError(object.ErrTypeAssertion)
			traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

			return nil, nil, object.ErrTypeAssertion
		}

		postOnly, ok := value["post_only"].(bool)
		if !ok {
			repository.GetRuntimeLogger().
				WithFields(fields).
				WithField(object.URIFieldError, object.ErrTypeAssertion).
				Error(object.ErrTypeAssertion.Error())
			traceSpan.RecordError(object.ErrTypeAssertion)
			traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

			return nil, nil, object.ErrTypeAssertion
		}

		stopTriggered, ok := value["stop_triggered"].(bool)
		if !ok {
			repository.GetRuntimeLogger().
				WithFields(fields).
				WithField(object.URIFieldError, object.ErrTypeAssertion).
				Error(object.ErrTypeAssertion.Error())
			traceSpan.RecordError(object.ErrTypeAssertion)
			traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

			return nil, nil, object.ErrTypeAssertion
		}

		daoOrderers = append(daoOrderers, dao.NewOrder(
			createdAT,
			updatedAT,
			sql.NullTime{
				Time:  time.Time{},
				Valid: false,
			},
			id,
			channel,
			clientOID,
			dealFunds,
			dealSize,
			fee,
			feeCurrency,
			funds,
			kucoinID,
			kucoinType,
			opType,
			price,
			remark,
			side,
			size,
			stop,
			stopPrice,
			stp,
			symbol,
			tags,
			timeInForce,
			tradeType,
			visibleSize,
			uint32(cancelAfter),
			uint32(kucoinCreatedAt),
			cancelExist,
			hidden,
			iceBerg,
			isActive,
			postOnly,
			stopTriggered,
		))
	}

	repository.GetRuntimeLogger().
		WithFields(fields).
		WithField(object.URIFieldDAOOrderers, daoOrderers).
		Debug(object.URIEmpty)

	var daoCursorer dao.Cursorer

	if daoPaginationer.GetLimit() < uint32(len(result)) {
		repository.GetRuntimeLogger().
			WithFields(fields).
			Debug(`daoPaginationer.GetLimit() < uint32(len(result))`)

		daoCursorer = dao.NewCursor(
			daoPaginationer.GetCursorer().GetOffset() + daoPaginationer.GetLimit(),
		)
	}

	repository.GetRuntimeLogger().
		WithFields(fields).
		WithField(object.URIFieldDAOCursorer, daoCursorer).
		Debug(object.URIEmpty)

	return daoOrderers, daoCursorer, nil
}

// Update is a function.
func (repository *orderRepository) Update(
	ctx context.Context,
	daoOrderer dao.Orderer,
) (time.Time, error) {
	var traceSpan trace.Span

	ctx, traceSpan = repository.GetTracer().Start(
		ctx,
		"Update",
		trace.WithSpanKind(trace.SpanKindProducer),
	)
	defer traceSpan.End()

	utilRuntimeContext := util.NewRuntimeContext(ctx, repository.GetUUIDer())
	utilSpanContext := util.NewSpanContext(traceSpan)
	fields := map[string]any{
		"name":        "Update",
		"rt_ctx":      utilRuntimeContext,
		"sp_ctx":      utilSpanContext,
		"config":      repository.GetConfigger(),
		"dao_orderer": daoOrderer,
	}

	repository.GetRuntimeLogger().
		WithFields(fields).
		Info(object.URIEmpty)

	nowUTC := repository.GetTimer().NowUTC()

	repository.GetRuntimeLogger().
		WithFields(fields).
		WithField(object.URIFieldNowUTC, nowUTC).
		Debug(object.URIEmpty)

	daoOrder := dao.NewOrder(
		daoOrderer.GetCreatedAt(),
		nowUTC,
		sql.NullTime{
			Time:  time.Time{},
			Valid: false,
		},
		daoOrderer.GetID(),
		daoOrderer.GetChannel(),
		daoOrderer.GetClientOID(),
		daoOrderer.GetDealFunds(),
		daoOrderer.GetDealSize(),
		daoOrderer.GetFee(),
		daoOrderer.GetFeeCurrency(),
		daoOrderer.GetFunds(),
		daoOrderer.GetKucoinID(),
		daoOrderer.GetKucoinType(),
		daoOrderer.GetOPType(),
		daoOrderer.GetPrice(),
		daoOrderer.GetRemark(),
		daoOrderer.GetSide(),
		daoOrderer.GetSize(),
		daoOrderer.GetStop(),
		daoOrderer.GetStopPrice(),
		daoOrderer.GetSTP(),
		daoOrderer.GetSymbol(),
		daoOrderer.GetTags(),
		daoOrderer.GetTimeInForce(),
		daoOrderer.GetTradeType(),
		daoOrderer.GetVisibleSize(),
		daoOrderer.GetCancelAfter(),
		daoOrderer.GetKucoinCreatedAt(),
		daoOrderer.GetCancelExist(),
		daoOrderer.GetHidden(),
		daoOrderer.GetIceBerg(),
		daoOrderer.GetIsActive(),
		daoOrderer.GetPostOnly(),
		daoOrderer.GetStopTriggered(),
	)

	repository.GetRuntimeLogger().
		WithFields(fields).
		WithField(object.URIFieldDAOOrder, daoOrder).
		Debug(object.URIEmpty)

	gormDB := repository.GetDB().
		WithContext(ctx).
		Where(map[string]any{
			"id": daoOrderer.GetID(),
		}).
		Updates(daoOrder.GetMap())
	if err := gormDB.Error; err != nil {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, err).
			Error(object.ErrOrderRepositoryUpdate.Error())
		traceSpan.RecordError(err)
		traceSpan.SetStatus(codes.Error, object.ErrOrderRepositoryUpdate.Error())

		return time.Time{}, err
	}

	if gormDB.RowsAffected == 0 {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, object.ErrOrderRepositoryUpdate).
			Error(object.ErrOrderRepositoryUpdate.Error())
		traceSpan.RecordError(object.ErrOrderRepositoryUpdate)
		traceSpan.SetStatus(codes.Error, object.ErrOrderRepositoryUpdate.Error())

		return time.Time{}, object.ErrOrderRepositoryUpdate
	}

	return daoOrder.GetUpdatedAt(), nil
}

// WithOptioners is a function.
func (repository *orderRepository) WithOptioners(
	optioners ...orderRepositoryOptioner,
) *orderRepository {
	newRepository := repository.clone()
	for _, optioner := range optioners {
		optioner.apply(newRepository)
	}

	return newRepository
}

func (repository *orderRepository) clone() *orderRepository {
	newRepository := repository

	return newRepository
}

func (optionerFunc orderRepositoryOptionerFunc) apply(
	repository *orderRepository,
) {
	optionerFunc(repository)
}
