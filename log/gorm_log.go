package log

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ShahoBashoki/kucoin/config"
	"github.com/ShahoBashoki/kucoin/object"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type gormLog struct {
	configConfigger config.Configger
	fields          map[string]any
	objectTimer     object.Timer
	zapLogger       *zap.Logger
	logger.Config
}

var (
	_ GetLogger            = (*gormLog)(nil)
	_ config.GetConfigger  = (*gormLog)(nil)
	_ logger.Interface     = (*gormLog)(nil)
	_ object.GetTimer      = (*gormLog)(nil)
	_ zapcore.LevelEnabler = (*gormLog)(nil)
)

// NewGormLog is function.
func NewGormLog(
	configConfigger config.Configger,
	fields map[string]any,
	objectTimer object.Timer,
	zapLogger *zap.Logger,
) *gormLog {
	return &gormLog{
		configConfigger: configConfigger,
		objectTimer:     objectTimer,
		fields:          fields,
		zapLogger:       zapLogger,
		Config: logger.Config{
			SlowThreshold:             configConfigger.GetLogConfigger().GetSQLSlowThreshold(),
			Colorful:                  false,
			IgnoreRecordNotFoundError: false,
			ParameterizedQueries:      false,
			LogLevel:                  logger.Info,
		},
	}
}

// GetLogger is a function.
func (gormLog *gormLog) GetLogger() *zap.Logger {
	return gormLog.zapLogger
}

// GetConfigger is a function.
func (gormLog *gormLog) GetConfigger() config.Configger {
	return gormLog.configConfigger
}

// LogMode is a function.
func (gormLog *gormLog) LogMode(
	level logger.LogLevel,
) logger.Interface {
	var zapCoreLevel zapcore.Level

	switch level {
	case logger.Silent:
		zapCoreLevel = zapcore.DPanicLevel

	case logger.Error:
		zapCoreLevel = zapcore.ErrorLevel

	case logger.Warn:
		zapCoreLevel = zapcore.WarnLevel

	case logger.Info:
		zapCoreLevel = zapcore.DebugLevel

	default:
		zapCoreLevel = zapcore.DebugLevel
	}

	gormLog.Enabled(zapCoreLevel)

	return gormLog
}

// Info is a function.
func (gormLog *gormLog) Info(
	_ context.Context,
	format string,
	args ...any,
) {
	if gormLog.GetLogger().Core().Enabled(zap.InfoLevel) {
		msg := fmt.Sprintf(format, args...)
		gormLog.GetLogger().Info(msg, GetFileLine())
	}
}

// Warn is a function.
func (gormLog *gormLog) Warn(
	_ context.Context,
	format string,
	args ...any,
) {
	if gormLog.GetLogger().Core().Enabled(zap.InfoLevel) {
		msg := fmt.Sprintf(format, args...)
		gormLog.GetLogger().Warn(msg, GetFileLine())
	}
}

// Error is a function.
func (gormLog *gormLog) Error(
	_ context.Context,
	format string,
	args ...any,
) {
	if gormLog.GetLogger().Core().Enabled(zap.InfoLevel) {
		msg := fmt.Sprintf(format, args...)
		gormLog.GetLogger().Error(msg, GetFileLine())
	}
}

// Trace is a function.
func (gormLog *gormLog) Trace(
	ctx context.Context,
	begin time.Time,
	fc func() (sql string, rowsAffected int64),
	err error,
) {
	elapsed := gormLog.GetTimer().Since(begin)
	sql, rows := fc()
	fields := map[string]any{
		"elapsed": float64(elapsed.Nanoseconds()) / 1e6,
		"rows":    rows,
		"sql":     sql,
	}

	if rows == -1 {
		fields[object.URIFieldRows] = -1
	}

	switch {
	case err != nil && (!errors.Is(err, gorm.ErrRecordNotFound) || !gormLog.IgnoreRecordNotFoundError):
		fields[object.URIFieldError] = err
		gormLog.
			WithFields(fields).
			Error(ctx, "%v", object.ErrSQL.Error())

	case gormLog.SlowThreshold < elapsed && gormLog.SlowThreshold != 0:
		fields["sql_slow_threshold"] = true
		gormLog.
			WithFields(fields).
			Warn(ctx, "%v", object.URIEmpty)

	default:
		gormLog.
			WithFields(fields).
			Info(ctx, "%v", object.URIEmpty)
	}
}

// WithFields implements RuntimeLogger.WithFieldsV1.
func (gormLog *gormLog) WithFields(
	fields map[string]any,
) logger.Interface {
	newFields := make(map[string]any, len(fields))
	zapcoreFields := make([]zap.Field, 0, len(fields))

	for key, value := range fields {
		newFields[key] = value
		zapcoreFields = append(zapcoreFields, zap.Any(key, value))
	}

	return NewGormLog(
		gormLog.GetConfigger(),
		newFields,
		gormLog.GetTimer(),
		gormLog.GetLogger().With(zapcoreFields...),
	)
}

// GetTimer is a function.
func (gormLog *gormLog) GetTimer() object.Timer {
	return gormLog.objectTimer
}

// Enabled is a function.
func (gormLog *gormLog) Enabled(
	lvl zapcore.Level,
) bool {
	zapCoreLevel, err := zapcore.ParseLevel(gormLog.GetConfigger().GetLogConfigger().GetLevel())
	if err != nil {
		zapCoreLevel = zapcore.InfoLevel
	}

	return zapCoreLevel <= lvl
}
