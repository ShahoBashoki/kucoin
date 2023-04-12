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
	// TickerRepositorier is a interface.
	TickerRepositorier interface {
		DAORepositorier[dao.Tickerer, dao.TickerFilterer]
	}

	// GetTickerRepositorier is an interface.
	GetTickerRepositorier interface {
		// GetTickerRepositorier is a function.
		GetTickerRepositorier() TickerRepositorier
	}

	tickerRepository struct {
		configConfigger  config.Configger
		gormDB           *gorm.DB
		logRuntimeLogger log.RuntimeLogger
		objectTimer      object.Timer
		traceTracer      trace.Tracer
		utilUUIDer       util.UUIDer
	}

	tickerRepositoryOptioner interface {
		apply(*tickerRepository)
	}

	tickerRepositoryOptionerFunc func(*tickerRepository)
)

var (
	_ TickerRepositorier   = (*tickerRepository)(nil)
	_ GetDB                = (*tickerRepository)(nil)
	_ config.GetConfigger  = (*tickerRepository)(nil)
	_ log.GetRuntimeLogger = (*tickerRepository)(nil)
	_ object.GetTimer      = (*tickerRepository)(nil)
	_ util.GetTracer       = (*tickerRepository)(nil)
	_ util.GetUUIDer       = (*tickerRepository)(nil)
)

// NewTickerRepository is a function.
func NewTickerRepository(
	configConfigger config.Configger,
	logRuntimeLogger log.RuntimeLogger,
	traceTracer trace.Tracer,
	utilUUIDer util.UUIDer,
	optioners ...tickerRepositoryOptioner,
) *tickerRepository {
	tickerRepository := &tickerRepository{
		configConfigger:  configConfigger,
		gormDB:           nil,
		logRuntimeLogger: logRuntimeLogger,
		objectTimer:      nil,
		traceTracer:      traceTracer,
		utilUUIDer:       utilUUIDer,
	}

	return tickerRepository.WithOptioners(optioners...)
}

// WithTickerRepositoryTimer is a function.
func WithTickerRepositoryTimer(
	objectTimer object.Timer,
) tickerRepositoryOptioner {
	return tickerRepositoryOptionerFunc(func(
		config *tickerRepository,
	) {
		config.objectTimer = objectTimer
	})
}

// WithTickerRepositoryDB is a function.
func WithTickerRepositoryDB(
	gormDB *gorm.DB,
) tickerRepositoryOptioner {
	return tickerRepositoryOptionerFunc(func(
		config *tickerRepository,
	) {
		config.gormDB = gormDB.
			Table(object.URITableTicker).
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
func (repository *tickerRepository) GetDB() *gorm.DB {
	return repository.gormDB
}

// GetConfigger is a function.
func (repository *tickerRepository) GetConfigger() config.Configger {
	return repository.configConfigger
}

// GetRuntimeLogger is a function.
func (repository *tickerRepository) GetRuntimeLogger() log.RuntimeLogger {
	return repository.logRuntimeLogger
}

// GetTimer is a function.
func (repository *tickerRepository) GetTimer() object.Timer {
	return repository.objectTimer
}

// GetTracer is a function.
func (repository *tickerRepository) GetTracer() trace.Tracer {
	return repository.traceTracer
}

// GetUUIDer is a function.
func (repository *tickerRepository) GetUUIDer() util.UUIDer {
	return repository.utilUUIDer
}

// Create is a function.
func (repository *tickerRepository) Create(
	ctx context.Context,
	daoTickerer dao.Tickerer,
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
		"name":         "Create",
		"rt_ctx":       utilRuntimeContext,
		"sp_ctx":       utilSpanContext,
		"config":       repository.GetConfigger(),
		"dao_tickerer": daoTickerer,
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

	daoTicker := dao.NewTicker(
		nowUTC,
		nowUTC,
		sql.NullTime{
			Time:  time.Time{},
			Valid: false,
		},
		id,
		daoTickerer.GetAveragePrice(),
		daoTickerer.GetBuy(),
		daoTickerer.GetChangePrice(),
		daoTickerer.GetChangeRate(),
		daoTickerer.GetHigh(),
		daoTickerer.GetLast(),
		daoTickerer.GetLow(),
		daoTickerer.GetMakerCoefficient(),
		daoTickerer.GetMakerFeeRate(),
		daoTickerer.GetSell(),
		daoTickerer.GetSymbol(),
		daoTickerer.GetSymbolName(),
		daoTickerer.GetTakerCoefficient(),
		daoTickerer.GetTakerFeeRate(),
		daoTickerer.GetVol(),
		daoTickerer.GetVolValue(),
	)

	repository.GetRuntimeLogger().
		WithFields(fields).
		WithField(object.URIFieldDAOTicker, daoTicker).
		Debug(object.URIEmpty)

	gormDB := repository.GetDB().
		WithContext(ctx).
		Create(daoTicker.GetMap())
	if err = gormDB.Error; err != nil {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, err).
			Error(object.ErrTickerRepositoryCreate.Error())
		traceSpan.RecordError(err)
		traceSpan.SetStatus(codes.Error, object.ErrTickerRepositoryCreate.Error())

		return uuid.Nil, err
	}

	return daoTicker.GetID(), nil
}

// Delete is a function.
func (repository *tickerRepository) Delete(
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
			Error(object.ErrTickerRepositoryDelete.Error())
		traceSpan.RecordError(err)
		traceSpan.SetStatus(codes.Error, object.ErrTickerRepositoryDelete.Error())

		return time.Time{}, err
	}

	if gormDB.RowsAffected == 0 {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, object.ErrTickerRepositoryDelete).
			Error(object.ErrTickerRepositoryDelete.Error())
		traceSpan.RecordError(object.ErrTickerRepositoryDelete)
		traceSpan.SetStatus(codes.Error, object.ErrTickerRepositoryDelete.Error())

		return time.Time{}, object.ErrTickerRepositoryDelete
	}

	return nowUTC, nil
}

// DeleteAll is a function.
func (repository *tickerRepository) DeleteAll(
	ctx context.Context,
) (time.Time, error) {
	var traceSpan trace.Span

	ctx, traceSpan = repository.GetTracer().Start(
		ctx,
		"DeleteAll",
		trace.WithSpanKind(trace.SpanKindProducer),
	)
	defer traceSpan.End()

	utilRuntimeContext := util.NewRuntimeContext(ctx, repository.GetUUIDer())
	utilSpanContext := util.NewSpanContext(traceSpan)
	fields := map[string]any{
		"name":   "DeleteAll",
		"rt_ctx": utilRuntimeContext,
		"sp_ctx": utilSpanContext,
		"config": repository.GetConfigger(),
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
		Exec(fmt.Sprintf("DELETE FROM %s", object.URITableTicker))
	if err := gormDB.Error; err != nil {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, err).
			Error(object.ErrTickerRepositoryDelete.Error())
		traceSpan.RecordError(err)
		traceSpan.SetStatus(codes.Error, object.ErrTickerRepositoryDelete.Error())

		return time.Time{}, err
	}

	if gormDB.RowsAffected == 0 {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, object.ErrTickerRepositoryDelete).
			Error(object.ErrTickerRepositoryDelete.Error())
		traceSpan.RecordError(object.ErrTickerRepositoryDelete)
		traceSpan.SetStatus(codes.Error, object.ErrTickerRepositoryDelete.Error())

		return time.Time{}, object.ErrTickerRepositoryDelete
	}

	return nowUTC, nil
}

// Read is a function.
func (repository *tickerRepository) Read(
	ctx context.Context,
	id uuid.UUID,
) (dao.Tickerer, error) {
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
		Select(fmt.Sprintf("%s.*", object.URITableTicker)).
		Find(result)
	if err := gormDB.Error; err != nil {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, err).
			Error(object.ErrTickerRepositoryRead.Error())
		traceSpan.RecordError(err)
		traceSpan.SetStatus(codes.Error, object.ErrTickerRepositoryRead.Error())

		return nil, err
	}

	if gormDB.RowsAffected == 0 {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, object.ErrTickerRepositoryRead).
			Error(object.ErrTickerRepositoryRead.Error())
		traceSpan.RecordError(object.ErrTickerRepositoryRead)
		traceSpan.SetStatus(codes.Error, object.ErrTickerRepositoryRead.Error())

		return nil, object.ErrTickerRepositoryRead
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

	averagePrice, ok := result["average_price"].(string)
	if !ok {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, object.ErrTypeAssertion).
			Error(object.ErrTypeAssertion.Error())
		traceSpan.RecordError(object.ErrTypeAssertion)
		traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

		return nil, object.ErrTypeAssertion
	}

	buy, ok := result["buy"].(string)
	if !ok {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, object.ErrTypeAssertion).
			Error(object.ErrTypeAssertion.Error())
		traceSpan.RecordError(object.ErrTypeAssertion)
		traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

		return nil, object.ErrTypeAssertion
	}

	changePrice, ok := result["change_price"].(string)
	if !ok {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, object.ErrTypeAssertion).
			Error(object.ErrTypeAssertion.Error())
		traceSpan.RecordError(object.ErrTypeAssertion)
		traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

		return nil, object.ErrTypeAssertion
	}

	changeRate, ok := result["change_rate"].(string)
	if !ok {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, object.ErrTypeAssertion).
			Error(object.ErrTypeAssertion.Error())
		traceSpan.RecordError(object.ErrTypeAssertion)
		traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

		return nil, object.ErrTypeAssertion
	}

	high, ok := result["high"].(string)
	if !ok {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, object.ErrTypeAssertion).
			Error(object.ErrTypeAssertion.Error())
		traceSpan.RecordError(object.ErrTypeAssertion)
		traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

		return nil, object.ErrTypeAssertion
	}

	last, ok := result["last"].(string)
	if !ok {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, object.ErrTypeAssertion).
			Error(object.ErrTypeAssertion.Error())
		traceSpan.RecordError(object.ErrTypeAssertion)
		traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

		return nil, object.ErrTypeAssertion
	}

	low, ok := result["low"].(string)
	if !ok {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, object.ErrTypeAssertion).
			Error(object.ErrTypeAssertion.Error())
		traceSpan.RecordError(object.ErrTypeAssertion)
		traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

		return nil, object.ErrTypeAssertion
	}

	makerCoefficient, ok := result["maker_coefficient"].(string)
	if !ok {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, object.ErrTypeAssertion).
			Error(object.ErrTypeAssertion.Error())
		traceSpan.RecordError(object.ErrTypeAssertion)
		traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

		return nil, object.ErrTypeAssertion
	}

	makerFeeRate, ok := result["maker_fee_rate"].(string)
	if !ok {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, object.ErrTypeAssertion).
			Error(object.ErrTypeAssertion.Error())
		traceSpan.RecordError(object.ErrTypeAssertion)
		traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

		return nil, object.ErrTypeAssertion
	}

	sell, ok := result["sell"].(string)
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

	symbolName, ok := result["symbol_name"].(string)
	if !ok {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, object.ErrTypeAssertion).
			Error(object.ErrTypeAssertion.Error())
		traceSpan.RecordError(object.ErrTypeAssertion)
		traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

		return nil, object.ErrTypeAssertion
	}

	takerCoefficient, ok := result["taker_coefficient"].(string)
	if !ok {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, object.ErrTypeAssertion).
			Error(object.ErrTypeAssertion.Error())
		traceSpan.RecordError(object.ErrTypeAssertion)
		traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

		return nil, object.ErrTypeAssertion
	}

	takerFeeRate, ok := result["taker_fee_rate"].(string)
	if !ok {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, object.ErrTypeAssertion).
			Error(object.ErrTypeAssertion.Error())
		traceSpan.RecordError(object.ErrTypeAssertion)
		traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

		return nil, object.ErrTypeAssertion
	}

	vol, ok := result["vol"].(string)
	if !ok {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, object.ErrTypeAssertion).
			Error(object.ErrTypeAssertion.Error())
		traceSpan.RecordError(object.ErrTypeAssertion)
		traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

		return nil, object.ErrTypeAssertion
	}

	volValue, ok := result["vol_value"].(string)
	if !ok {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, object.ErrTypeAssertion).
			Error(object.ErrTypeAssertion.Error())
		traceSpan.RecordError(object.ErrTypeAssertion)
		traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

		return nil, object.ErrTypeAssertion
	}

	daoTicker := dao.NewTicker(
		createdAT,
		updatedAT,
		sql.NullTime{
			Time:  time.Time{},
			Valid: false,
		},
		id,
		averagePrice,
		buy,
		changePrice,
		changeRate,
		high,
		last,
		low,
		makerCoefficient,
		makerFeeRate,
		sell,
		symbol,
		symbolName,
		takerCoefficient,
		takerFeeRate,
		vol,
		volValue,
	)

	repository.GetRuntimeLogger().
		WithFields(fields).
		WithField(object.URIFieldDAOTicker, daoTicker).
		Debug(object.URIEmpty)

	return daoTicker, nil
}

// ReadList is a function.
func (repository *tickerRepository) ReadList(
	ctx context.Context,
	daoPaginationer dao.Paginationer,
	daoTickerFilterer dao.TickerFilterer,
) ([]dao.Tickerer, dao.Cursorer, error) {
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
		"name":                "ReadList",
		"rt_ctx":              utilRuntimeContext,
		"sp_ctx":              utilSpanContext,
		"config":              repository.GetConfigger(),
		"dao_paginationer":    daoPaginationer,
		"dao_ticker_filterer": daoTickerFilterer,
	}

	repository.GetRuntimeLogger().
		WithFields(fields).
		Info(object.URIEmpty)

	result := make([]map[string]any, 0, daoPaginationer.GetLimit()+1)

	gormDB := repository.GetDB().
		WithContext(ctx).
		Scopes(
			daoTickerFilterer.Filter,
			daoPaginationer.Pagination(object.URITableTicker),
		).
		Where(map[string]any{
			"deleted_at": nil,
		}).
		Select(fmt.Sprintf("%s.*", object.URITableTicker)).
		Find(&result)
	if err := gormDB.Error; err != nil {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, err).
			Error(object.ErrTickerRepositoryReadList.Error())
		traceSpan.RecordError(err)
		traceSpan.SetStatus(codes.Error, object.ErrTickerRepositoryReadList.Error())

		return nil, nil, err
	}

	repository.GetRuntimeLogger().
		WithFields(fields).
		WithField(object.URIFieldResult, result).
		Debug(object.URIEmpty)

	daoTickerers := make([]dao.Tickerer, 0, daoPaginationer.GetLimit())

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

		averagePrice, ok := value["average_price"].(string)
		if !ok {
			repository.GetRuntimeLogger().
				WithFields(fields).
				WithField(object.URIFieldError, object.ErrTypeAssertion).
				Error(object.ErrTypeAssertion.Error())
			traceSpan.RecordError(object.ErrTypeAssertion)
			traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

			return nil, nil, object.ErrTypeAssertion
		}

		buy, ok := value["buy"].(string)
		if !ok {
			repository.GetRuntimeLogger().
				WithFields(fields).
				WithField(object.URIFieldError, object.ErrTypeAssertion).
				Error(object.ErrTypeAssertion.Error())
			traceSpan.RecordError(object.ErrTypeAssertion)
			traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

			return nil, nil, object.ErrTypeAssertion
		}

		changePrice, ok := value["change_price"].(string)
		if !ok {
			repository.GetRuntimeLogger().
				WithFields(fields).
				WithField(object.URIFieldError, object.ErrTypeAssertion).
				Error(object.ErrTypeAssertion.Error())
			traceSpan.RecordError(object.ErrTypeAssertion)
			traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

			return nil, nil, object.ErrTypeAssertion
		}

		changeRate, ok := value["change_rate"].(string)
		if !ok {
			repository.GetRuntimeLogger().
				WithFields(fields).
				WithField(object.URIFieldError, object.ErrTypeAssertion).
				Error(object.ErrTypeAssertion.Error())
			traceSpan.RecordError(object.ErrTypeAssertion)
			traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

			return nil, nil, object.ErrTypeAssertion
		}

		high, ok := value["high"].(string)
		if !ok {
			repository.GetRuntimeLogger().
				WithFields(fields).
				WithField(object.URIFieldError, object.ErrTypeAssertion).
				Error(object.ErrTypeAssertion.Error())
			traceSpan.RecordError(object.ErrTypeAssertion)
			traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

			return nil, nil, object.ErrTypeAssertion
		}

		last, ok := value["last"].(string)
		if !ok {
			repository.GetRuntimeLogger().
				WithFields(fields).
				WithField(object.URIFieldError, object.ErrTypeAssertion).
				Error(object.ErrTypeAssertion.Error())
			traceSpan.RecordError(object.ErrTypeAssertion)
			traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

			return nil, nil, object.ErrTypeAssertion
		}

		low, ok := value["low"].(string)
		if !ok {
			repository.GetRuntimeLogger().
				WithFields(fields).
				WithField(object.URIFieldError, object.ErrTypeAssertion).
				Error(object.ErrTypeAssertion.Error())
			traceSpan.RecordError(object.ErrTypeAssertion)
			traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

			return nil, nil, object.ErrTypeAssertion
		}

		makerCoefficient, ok := value["maker_coefficient"].(string)
		if !ok {
			repository.GetRuntimeLogger().
				WithFields(fields).
				WithField(object.URIFieldError, object.ErrTypeAssertion).
				Error(object.ErrTypeAssertion.Error())
			traceSpan.RecordError(object.ErrTypeAssertion)
			traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

			return nil, nil, object.ErrTypeAssertion
		}

		makerFeeRate, ok := value["maker_fee_rate"].(string)
		if !ok {
			repository.GetRuntimeLogger().
				WithFields(fields).
				WithField(object.URIFieldError, object.ErrTypeAssertion).
				Error(object.ErrTypeAssertion.Error())
			traceSpan.RecordError(object.ErrTypeAssertion)
			traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

			return nil, nil, object.ErrTypeAssertion
		}

		sell, ok := value["sell"].(string)
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

		symbolName, ok := value["symbol_name"].(string)
		if !ok {
			repository.GetRuntimeLogger().
				WithFields(fields).
				WithField(object.URIFieldError, object.ErrTypeAssertion).
				Error(object.ErrTypeAssertion.Error())
			traceSpan.RecordError(object.ErrTypeAssertion)
			traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

			return nil, nil, object.ErrTypeAssertion
		}

		takerCoefficient, ok := value["taker_coefficient"].(string)
		if !ok {
			repository.GetRuntimeLogger().
				WithFields(fields).
				WithField(object.URIFieldError, object.ErrTypeAssertion).
				Error(object.ErrTypeAssertion.Error())
			traceSpan.RecordError(object.ErrTypeAssertion)
			traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

			return nil, nil, object.ErrTypeAssertion
		}

		takerFeeRate, ok := value["taker_fee_rate"].(string)
		if !ok {
			repository.GetRuntimeLogger().
				WithFields(fields).
				WithField(object.URIFieldError, object.ErrTypeAssertion).
				Error(object.ErrTypeAssertion.Error())
			traceSpan.RecordError(object.ErrTypeAssertion)
			traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

			return nil, nil, object.ErrTypeAssertion
		}

		vol, ok := value["vol"].(string)
		if !ok {
			repository.GetRuntimeLogger().
				WithFields(fields).
				WithField(object.URIFieldError, object.ErrTypeAssertion).
				Error(object.ErrTypeAssertion.Error())
			traceSpan.RecordError(object.ErrTypeAssertion)
			traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

			return nil, nil, object.ErrTypeAssertion
		}

		volValue, ok := value["vol_value"].(string)
		if !ok {
			repository.GetRuntimeLogger().
				WithFields(fields).
				WithField(object.URIFieldError, object.ErrTypeAssertion).
				Error(object.ErrTypeAssertion.Error())
			traceSpan.RecordError(object.ErrTypeAssertion)
			traceSpan.SetStatus(codes.Error, object.ErrTypeAssertion.Error())

			return nil, nil, object.ErrTypeAssertion
		}

		daoTickerers = append(daoTickerers, dao.NewTicker(
			createdAT,
			updatedAT,
			sql.NullTime{
				Time:  time.Time{},
				Valid: false,
			},
			id,
			averagePrice,
			buy,
			changePrice,
			changeRate,
			high,
			last,
			low,
			makerCoefficient,
			makerFeeRate,
			sell,
			symbol,
			symbolName,
			takerCoefficient,
			takerFeeRate,
			vol,
			volValue,
		))
	}

	repository.GetRuntimeLogger().
		WithFields(fields).
		WithField(object.URIFieldDAOTickerers, daoTickerers).
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
		WithField(object.URIFieldDAOCursor, daoCursorer).
		Debug(object.URIEmpty)

	return daoTickerers, daoCursorer, nil
}

// Update is a function.
func (repository *tickerRepository) Update(
	ctx context.Context,
	daoTickerer dao.Tickerer,
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
		"name":         "Update",
		"rt_ctx":       utilRuntimeContext,
		"sp_ctx":       utilSpanContext,
		"config":       repository.GetConfigger(),
		"dao_tickerer": daoTickerer,
	}

	repository.GetRuntimeLogger().
		WithFields(fields).
		Info(object.URIEmpty)

	nowUTC := repository.GetTimer().NowUTC()

	repository.GetRuntimeLogger().
		WithFields(fields).
		WithField(object.URIFieldNowUTC, nowUTC).
		Debug(object.URIEmpty)

	daoTicker := dao.NewTicker(
		daoTickerer.GetCreatedAt(),
		nowUTC,
		sql.NullTime{
			Time:  time.Time{},
			Valid: false,
		},
		daoTickerer.GetID(),
		daoTickerer.GetAveragePrice(),
		daoTickerer.GetBuy(),
		daoTickerer.GetChangePrice(),
		daoTickerer.GetChangeRate(),
		daoTickerer.GetHigh(),
		daoTickerer.GetLast(),
		daoTickerer.GetLow(),
		daoTickerer.GetMakerCoefficient(),
		daoTickerer.GetMakerFeeRate(),
		daoTickerer.GetSell(),
		daoTickerer.GetSymbol(),
		daoTickerer.GetSymbolName(),
		daoTickerer.GetTakerCoefficient(),
		daoTickerer.GetTakerFeeRate(),
		daoTickerer.GetVol(),
		daoTickerer.GetVolValue(),
	)

	repository.GetRuntimeLogger().
		WithFields(fields).
		WithField(object.URIFieldDAOTicker, daoTicker).
		Debug(object.URIEmpty)

	gormDB := repository.GetDB().
		WithContext(ctx).
		Where(map[string]any{
			"id": daoTickerer.GetID(),
		}).
		Updates(daoTicker.GetMap())
	if err := gormDB.Error; err != nil {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, err).
			Error(object.ErrTickerRepositoryUpdate.Error())
		traceSpan.RecordError(err)
		traceSpan.SetStatus(codes.Error, object.ErrTickerRepositoryUpdate.Error())

		return time.Time{}, err
	}

	if gormDB.RowsAffected == 0 {
		repository.GetRuntimeLogger().
			WithFields(fields).
			WithField(object.URIFieldError, object.ErrTickerRepositoryUpdate).
			Error(object.ErrTickerRepositoryUpdate.Error())
		traceSpan.RecordError(object.ErrTickerRepositoryUpdate)
		traceSpan.SetStatus(codes.Error, object.ErrTickerRepositoryUpdate.Error())

		return time.Time{}, object.ErrTickerRepositoryUpdate
	}

	return daoTicker.GetUpdatedAt(), nil
}

// WithOptioners is a function.
func (repository *tickerRepository) WithOptioners(
	optioners ...tickerRepositoryOptioner,
) *tickerRepository {
	newRepository := repository.clone()
	for _, optioner := range optioners {
		optioner.apply(newRepository)
	}

	return newRepository
}

func (repository *tickerRepository) clone() *tickerRepository {
	newRepository := repository

	return newRepository
}

func (optionerFunc tickerRepositoryOptionerFunc) apply(
	repository *tickerRepository,
) {
	optionerFunc(repository)
}
