package main

import (
	"context"
	"time"

	"github.com/ShahoBashoki/kucoin/config"
	"github.com/ShahoBashoki/kucoin/log"
	"github.com/ShahoBashoki/kucoin/object"
	"github.com/ShahoBashoki/kucoin/repository"
	"github.com/ShahoBashoki/kucoin/server"
	"github.com/ShahoBashoki/kucoin/service"
	"github.com/ShahoBashoki/kucoin/util"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	viper.AutomaticEnv()
	viper.SetDefault("DATABASE_DSN", "postgresql://root@127.0.0.1:26257/defaultdb?sslmode=disable")
	viper.SetDefault("LOG_FILE", "file.log")
	viper.SetDefault("LOG_FORMAT", "json")
	viper.SetDefault("LOG_LEVEL", "debug")
	viper.SetDefault("LOG_SQL_SLOW_THRESHOLD", 1*time.Second)
	viper.SetDefault("LOG_MAX_AGE", 0)
	viper.SetDefault("LOG_MAX_BACKUPS", 0)
	viper.SetDefault("LOG_MAX_SIZE", object.NUMLogConfigDefaultLogMaxSize)
	viper.SetDefault("LOG_COMPRESS", false)
	viper.SetDefault("LOG_LOCAL_TIME", false)
	viper.SetDefault("LOG_ROTATION", false)
	viper.SetDefault("LOG_STDOUT", true)
	viper.SetDefault("OTEL_EXPORTER_JAEGER_ENDPOINT", "http://otelcol:14268/api/traces")
	viper.SetDefault("OTEL_EXPORTER_JAEGER_PASSWORD", object.URIEmpty)
	viper.SetDefault("OTEL_EXPORTER_JAEGER_USER", object.URIEmpty)
	viper.SetDefault("OTEL_INSTRUMENTATION_NAME", "kucoin")
	viper.SetDefault("OTEL_SERVICE_INSTANCE_ID", "kucoin")
	viper.SetDefault("OTEL_SERVICE_NAME", "kucoin")
	viper.SetDefault("OTEL_SERVICE_NAMESPACE", "kucoin")
	viper.SetDefault("OTEL_SERVICE_VERSION", "v0.1.0")
	viper.SetDefault("REDPANDA_PROXY_URL", "http://redpanda:8082")
	viper.SetDefault("REDPANDA_TOPIC", "kucoin")
	viper.SetDefault("RUNTIME_VALIDATE_MAP_RULES", `{"rules":[{"version":"1"}]}`)
	viper.SetDefault("RUNTIME_NODE", "kucoin")
	viper.SetDefault("SERVER_ENDPOINT_ADDR", ":8080")
	viper.SetDefault("SERVER_ENDPOINT_NETWORK", "tcp")

	configConfig := config.NewConfig(
		config.WithDatabaseConfigger(
			config.WithDatabaseConfigDSN(viper.GetString("DATABASE_DSN")),
		),
		config.WithLogConfigger(
			config.WithLogConfigFile(viper.GetString("LOG_FILE")),
			config.WithLogConfigFormat(viper.GetString("LOG_FORMAT")),
			config.WithLogConfigLevel(viper.GetString("LOG_LEVEL")),
			config.WithLogConfigSQLSlowThreshold(viper.GetDuration("LOG_SQL_SLOW_THRESHOLD")),
			config.WithLogConfigMaxAge(viper.GetInt("LOG_MAX_AGE")),
			config.WithLogConfigMaxBackups(viper.GetInt("LOG_MAX_BACKUPS")),
			config.WithLogConfigMaxSize(viper.GetInt("LOG_MAX_SIZE")),
			config.WithLogConfigCompress(viper.GetBool("LOG_COMPRESS")),
			config.WithLogConfigLocalTime(viper.GetBool("LOG_LOCAL_TIME")),
			config.WithLogConfigRotation(viper.GetBool("LOG_ROTATION")),
			config.WithLogConfigStdout(viper.GetBool("LOG_STDOUT")),
		),
		config.WithOtelConfigger(
			config.WithOtelConfigExporterJaegerEndpoint(
				viper.GetString("OTEL_EXPORTER_JAEGER_ENDPOINT"),
			),
			config.WithOtelConfigExporterJaegerPassword(
				viper.GetString("OTEL_EXPORTER_JAEGER_PASSWORD"),
			),
			config.WithOtelConfigExporterJaegerUsername(
				viper.GetString("OTEL_EXPORTER_JAEGER_USER"),
			),
			config.WithOtelConfigInstrumentationName(viper.GetString("OTEL_INSTRUMENTATION_NAME")),
			config.WithOtelConfigServiceInstanceID(viper.GetString("OTEL_SERVICE_INSTANCE_ID")),
			config.WithOtelConfigServiceName(viper.GetString("OTEL_SERVICE_NAME")),
			config.WithOtelConfigServiceNamespace(viper.GetString("OTEL_SERVICE_NAMESPACE")),
			config.WithOtelConfigServiceVersion(viper.GetString("OTEL_SERVICE_VERSION")),
		),
		config.WithRedpandaConfigger(
			config.WithRedpandaConfigProxyURL(viper.GetString("REDPANDA_PROXY_URL")),
			config.WithRedpandaConfigTopic(viper.GetString("REDPANDA_TOPIC")),
		),
		config.WithRuntimeConfigger(
			config.WithRuntimeConfigValidateMapRules(
				util.Cast(viper.GetStringMap("RUNTIME_VALIDATE_MAP_RULES")),
			),
			config.WithRuntimeConfigNode(viper.GetString("RUNTIME_NODE")),
			config.WithRuntimeConfigClientPaginationRequestSizeMax(
				viper.GetInt("RUNTIME_CLIENT_PAGINATION_REQUEST_SIZE_MAX"),
			),
		),
		config.WithServerConfigger(
			config.WithServerConfigEndpointConfigger(
				config.NewEndpointConfig(
					config.WithEndpointConfigAddr(viper.GetString("SERVER_ENDPOINT_ADDR")),
					config.WithEndpointConfigNetwork(viper.GetString("SERVER_ENDPOINT_NETWORK")),
				),
			),
		),
	)

	logZapLogger := log.NewZapLogger(configConfig)
	// ObjectTime := object.NewTime()
	// logGormLog := log.NewGormLog(configConfig, map[string]any{}, objectTime, logZapLogger).
	logRuntimeLog := log.NewRuntimeLog(
		configConfig,
		map[string]any{},
		logZapLogger.WithOptions(zap.AddCallerSkip(1)),
	)
	utilUUID := util.NewUUID()
	traceTracer := util.NewTracer(ctx, configConfig, logRuntimeLog, utilUUID)

	var traceSpan trace.Span
	ctx, traceSpan = traceTracer.Start(
		ctx,
		"main",
		trace.WithSpanKind(trace.SpanKindServer))

	defer traceSpan.End()

	utilRuntimeContext := util.NewRuntimeContext(ctx, utilUUID)
	utilSpanContext := util.NewSpanContext(traceSpan)
	fields := map[string]any{
		"name":   "main",
		"rt_ctx": utilRuntimeContext,
		"sp_ctx": utilSpanContext,
		"config": configConfig,
	}

	logRuntimeLog.
		WithFields(fields).
		Info(object.URIEmpty)

	// GormDB, err := gorm.Open(
	// 	postgres.Open(configConfig.GetDatabaseConfigger().GetDSN()),
	// 	&gorm.Config{
	// 		SkipDefaultTransaction: true,
	// 		NamingStrategy:         nil,
	// 		FullSaveAssociations:   false,
	// 		Logger:                 logGormLog,
	// 		NowFunc: func() time.Time {
	// 			return object.NewTime().NowUTC()
	// 		},
	// 		DryRun:                                   false,
	// 		PrepareStmt:                              true,
	// 		DisableAutomaticPing:                     false,
	// 		DisableForeignKeyConstraintWhenMigrating: true,
	// 		IgnoreRelationshipsWhenMigrating:         false,
	// 		DisableNestedTransaction:                 true,
	// 		AllowGlobalUpdate:                        false,
	// 		QueryFields:                              true,
	// 		CreateBatchSize:                          0,
	// 		ClauseBuilders:                           map[string]clause.ClauseBuilder{},
	// 		ConnPool:                                 nil,
	// 		Dialector:                                nil,
	// 		Plugins:                                  map[string]gorm.Plugin{},
	// 	},
	// )
	// if err != nil {
	// 	logRuntimeLog.
	// 		WithFields(fields).
	// 		WithField(object.URIFieldError, err).
	// 		Error(object.ErrGormOpen.Error())
	// 	traceSpan.RecordError(err)
	// 	traceSpan.SetStatus(codes.Error, object.ErrGormOpen.Error())
	// }.

	repositoryRepository := repository.NewRepository()
	servicer := service.NewServicer(
		configConfig,
		repositoryRepository,
		logRuntimeLog,
		traceTracer,
		utilUUID,
	)
	serverer := server.NewServerrer(
		configConfig,
		ctx,
		logRuntimeLog,
		servicer,
		traceTracer,
		utilUUID,
	)

	if err := serverer.Run(ctx); err != nil {
		logRuntimeLog.
			WithFields(fields).
			WithField(object.URIFieldError, err).
			Error(object.ErrServerRun.Error())
		traceSpan.RecordError(err)
		traceSpan.SetStatus(codes.Error, object.ErrServerRun.Error())
	}
}
