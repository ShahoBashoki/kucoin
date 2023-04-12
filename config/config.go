package config

import (
	"encoding/json"

	"github.com/ShahoBashoki/kucoin/object"
)

type (
	// Configger interface is the core configuration.
	Configger interface {
		GetDatabaseConfigger
		GetKucoinConfigger
		GetLogConfigger
		GetOtelConfigger
		GetRedpandaConfigger
		GetRuntimeConfigger
		GetServerConfigger
	}

	// GetConfigger is an interface.
	GetConfigger interface {
		// GetConfigger is a function.
		GetConfigger() Configger
	}

	config struct {
		databaseConfigger DatabaseConfigger
		kucoinConfigger   KucoinConfigger
		logConfigger      LogConfigger
		otelConfigger     OtelConfigger
		redpandaConfigger RedpandaConfigger
		runtimeConfigger  RuntimeConfigger
		serverConfigger   ServerConfigger
	}

	configOptioner interface {
		apply(*config)
	}

	configOptionerFunc func(*config)
)

var (
	_ Configger            = (*config)(nil)
	_ GetDatabaseConfigger = (*config)(nil)
	_ GetKucoinConfigger   = (*config)(nil)
	_ GetLogConfigger      = (*config)(nil)
	_ GetOtelConfigger     = (*config)(nil)
	_ GetRedpandaConfigger = (*config)(nil)
	_ GetRuntimeConfigger  = (*config)(nil)
	_ GetServerConfigger   = (*config)(nil)
	_ json.Marshaler       = (*config)(nil)
	_ object.GetMap        = (*config)(nil)
)

// NewConfig constructs a Config struct which represents server settings,
// and populates it with default values.
func NewConfig(
	optioners ...configOptioner,
) *config {
	config := &config{
		databaseConfigger: nil,
		kucoinConfigger:   nil,
		logConfigger:      nil,
		otelConfigger:     nil,
		redpandaConfigger: nil,
		runtimeConfigger:  nil,
		serverConfigger:   nil,
	}

	return config.WithOptioners(optioners...)
}

// WithDatabaseConfigger is a function.
func WithDatabaseConfigger(
	optioners ...databaseConfigOptioner,
) configOptioner {
	return configOptionerFunc(func(
		config *config,
	) {
		config.databaseConfigger = NewDatabaseConfig(optioners...)
	})
}

// WithKucoinConfigger is a function.
func WithKucoinConfigger(
	optioners ...kucoinConfigOptioner,
) configOptioner {
	return configOptionerFunc(func(
		config *config,
	) {
		config.kucoinConfigger = NewKucoinConfig(optioners...)
	})
}

// WithLogConfigger is a function.
func WithLogConfigger(
	optioners ...logConfigOptioner,
) configOptioner {
	return configOptionerFunc(func(
		config *config,
	) {
		config.logConfigger = NewLogConfig(optioners...)
	})
}

// WithOtelConfigger is a function.
func WithOtelConfigger(
	optioners ...otelConfigOptioner,
) configOptioner {
	return configOptionerFunc(func(
		config *config,
	) {
		config.otelConfigger = NewOtelConfig(optioners...)
	})
}

// WithRedpandaConfigger is a function.
func WithRedpandaConfigger(
	optioners ...redpandaConfigOptioner,
) configOptioner {
	return configOptionerFunc(func(
		config *config,
	) {
		config.redpandaConfigger = NewRedpandaConfig(optioners...)
	})
}

// WithRuntimeConfigger is a function.
func WithRuntimeConfigger(
	optioners ...runtimeConfigOptioner,
) configOptioner {
	return configOptionerFunc(func(
		config *config,
	) {
		config.runtimeConfigger = NewRuntimeConfig(optioners...)
	})
}

// WithServerConfigger is a function.
func WithServerConfigger(
	optioners ...serverConfigOptioner,
) configOptioner {
	return configOptionerFunc(func(
		config *config,
	) {
		config.serverConfigger = NewServerConfig(optioners...)
	})
}

// GetDatabaseConfigger is a function.
func (config *config) GetDatabaseConfigger() DatabaseConfigger {
	return config.databaseConfigger
}

// GetKucoinConfigger is a function.
func (config *config) GetKucoinConfigger() KucoinConfigger {
	return config.kucoinConfigger
}

// GetLogConfigger is a function.
func (config *config) GetLogConfigger() LogConfigger {
	return config.logConfigger
}

// GetOtelConfigger is a function.
func (config *config) GetOtelConfigger() OtelConfigger {
	return config.otelConfigger
}

// GetRedpandaConfigger is a function.
func (config *config) GetRedpandaConfigger() RedpandaConfigger {
	return config.redpandaConfigger
}

// GetRuntimeConfigger is a function.
func (config *config) GetRuntimeConfigger() RuntimeConfigger {
	return config.runtimeConfigger
}

// GetServerConfigger is a function.
func (config *config) GetServerConfigger() ServerConfigger {
	return config.serverConfigger
}

// GetMap is a function.
func (config *config) GetMap() map[string]any {
	return map[string]any{
		"database_configger": config.GetDatabaseConfigger(),
		"kucoin_configger":   config.GetKucoinConfigger(),
		"logger_configger":   config.GetLogConfigger(),
		"otel_configger":     config.GetOtelConfigger(),
		"redpanda_configger": config.GetRedpandaConfigger(),
		"runtime_configger":  config.GetRuntimeConfigger(),
		"server_configger":   config.GetServerConfigger(),
	}
}

// MarshalJSON is a function.
// read more https://pkg.go.dev/encoding/json#Marshaler
func (config *config) MarshalJSON() ([]byte, error) {
	return json.Marshal(config.GetMap())
}

// WithOptioners is a function.
func (config *config) WithOptioners(
	optioners ...configOptioner,
) *config {
	newConfig := config.clone()
	for _, optioner := range optioners {
		optioner.apply(newConfig)
	}

	return newConfig
}

func (config *config) clone() *config {
	newConfig := config

	return newConfig
}

func (optionerFunc configOptionerFunc) apply(
	config *config,
) {
	optionerFunc(config)
}
