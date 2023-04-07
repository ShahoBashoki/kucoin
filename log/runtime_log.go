package log

import (
	"fmt"
	"net/http"

	"github.com/ShahoBashoki/kucoin/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	// RuntimeLogger exposes a logging framework to use in modules.
	// It exposes level-specific logging functions and a set of common functions for compatibility.
	RuntimeLogger interface {
		// Log a message with optional arguments at DEBUG level. Arguments are handled in the manner of fmt.Printf.
		Debug(
			string,
			...any,
		)
		// Log a message with optional arguments at INFO level. Arguments are handled in the manner of fmt.Printf.
		Info(
			string,
			...any,
		)
		// Log a message with optional arguments at WARN level. Arguments are handled in the manner of fmt.Printf.
		Warn(
			string,
			...any,
		)
		// Log a message with optional arguments at ERROR level. Arguments are handled in the manner of fmt.Printf.
		Error(
			string,
			...any,
		)
		// Log a message with optional arguments at FATAL level. Arguments are handled in the manner of fmt.Printf.
		Fatal(
			string,
			...any,
		)
		// Return a runtimeLog with the specified field set so that they are included in subsequent logging calls.
		WithField(
			string,
			any,
		) RuntimeLogger
		// Return a runtimeLog with the specified fields set so that they are included in subsequent logging calls.
		WithFields(
			map[string]any,
		) RuntimeLogger
		// Returns the fields set in this RuntimeLogger.
		Fields() map[string]any
	}

	// GetRuntimeLogger is an interface.
	GetRuntimeLogger interface {
		// GetRuntimeLogger is a function.
		GetRuntimeLogger() RuntimeLogger
	}

	runtimeLog struct {
		configConfigger config.Configger
		fields          map[string]any
		zapLogger       *zap.Logger
	}
)

var (
	_ GetLogger            = (*runtimeLog)(nil)
	_ RuntimeLogger        = (*runtimeLog)(nil)
	_ config.GetConfigger  = (*runtimeLog)(nil)
	_ zapcore.LevelEnabler = (*runtimeLog)(nil)
)

// NewRuntimeLog is function.
func NewRuntimeLog(
	configConfigger config.Configger,
	fields map[string]any,
	zapLogger *zap.Logger,
) *runtimeLog {
	return &runtimeLog{
		configConfigger: configConfigger,
		fields:          fields,
		zapLogger:       zapLogger,
	}
}

// GetLogger is a function.
func (runtimeLog *runtimeLog) GetLogger() *zap.Logger {
	return runtimeLog.zapLogger
}

// GetConfigger is a function.
func (runtimeLog *runtimeLog) GetConfigger() config.Configger {
	return runtimeLog.configConfigger
}

// Debug implements RuntimeLogger.Debug.
func (runtimeLog *runtimeLog) Debug(
	format string,
	args ...any,
) {
	if runtimeLog.GetLogger().Core().Enabled(zap.DebugLevel) {
		msg := fmt.Sprintf(format, args...)
		runtimeLog.GetLogger().Debug(msg)
	}
}

// Info implements RuntimeLogger.Info.
func (runtimeLog *runtimeLog) Info(
	format string,
	args ...any,
) {
	if runtimeLog.GetLogger().Core().Enabled(zap.InfoLevel) {
		msg := fmt.Sprintf(format, args...)
		runtimeLog.GetLogger().Info(msg, GetFileLine())
	}
}

// Warn implements RuntimeLogger.Warn.
func (runtimeLog *runtimeLog) Warn(
	format string,
	args ...any,
) {
	if runtimeLog.GetLogger().Core().Enabled(zap.WarnLevel) {
		msg := fmt.Sprintf(format, args...)
		runtimeLog.GetLogger().Warn(msg, GetFileLine())
	}
}

// Error implements RuntimeLogger.Error.
func (runtimeLog *runtimeLog) Error(
	format string,
	args ...any,
) {
	if runtimeLog.GetLogger().Core().Enabled(zap.ErrorLevel) {
		msg := fmt.Sprintf(format, args...)
		runtimeLog.GetLogger().Error(msg, GetFileLine())
	}
}

// Fatal implements RuntimeLogger.Error.
func (runtimeLog *runtimeLog) Fatal(
	format string,
	args ...any,
) {
	if runtimeLog.GetLogger().Core().Enabled(zap.FatalLevel) {
		msg := fmt.Sprintf(format, args...)
		runtimeLog.GetLogger().Fatal(msg, GetFileLine())
	}
}

// WithField implements RuntimeLogger.WithFieldV1.
func (runtimeLog *runtimeLog) WithField(
	key string,
	value any,
) RuntimeLogger {
	return runtimeLog.WithFields(map[string]any{
		key: value,
	})
}

// WithFields implements RuntimeLogger.WithFieldsV1.
func (runtimeLog *runtimeLog) WithFields(
	fields map[string]any,
) RuntimeLogger {
	newFields := make(map[string]any, len(fields)+len(runtimeLog.fields))
	zapcoreFields := make([]zap.Field, 0, len(fields)+len(runtimeLog.fields))

	for key, value := range runtimeLog.fields {
		newFields[key] = value
	}

	for key, value := range fields {
		if value, ok := value.(http.Request); ok && key == "req" {
			newFields["req_body"] = value.Body
			newFields["req_cookies"] = value.Cookies()
			newFields["req_header"] = value.Header
			zapcoreFields = append(zapcoreFields, zap.Any(key, value))

			continue
		}

		newFields[key] = value
		zapcoreFields = append(zapcoreFields, zap.Any(key, value))
	}

	return NewRuntimeLog(
		runtimeLog.GetConfigger(),
		newFields,
		runtimeLog.GetLogger().With(zapcoreFields...),
	)
}

// Fields implements RuntimeLogger.Fields.
func (runtimeLog *runtimeLog) Fields() map[string]any {
	return runtimeLog.fields
}

// Enabled is a function.
func (runtimeLog *runtimeLog) Enabled(
	lvl zapcore.Level,
) bool {
	zapcoreLevel, err := zapcore.ParseLevel(runtimeLog.GetConfigger().GetLogConfigger().GetLevel())
	if err != nil {
		zapcoreLevel = zapcore.InfoLevel
	}

	return zapcoreLevel <= lvl
}
