package config

import (
	"encoding/json"

	"github.com/ShahoBashoki/kucoin/object"
)

type (
	// OtelConfigger is an interface.
	OtelConfigger interface {
		// GetExporterJaegerEndpoint is a function.
		GetExporterJaegerEndpoint() string
		// GetExporterJaegerPassword is a function.
		GetExporterJaegerPassword() string
		// GetExporterJaegerUsername is a function.
		GetExporterJaegerUsername() string
		// GetInstrumentationName is a function.
		GetInstrumentationName() string
		// GetServiceInstanceID is a function.
		GetServiceInstanceID() string
		// GetServiceName is a function.
		GetServiceName() string
		// GetServiceNamespace is a function.
		GetServiceNamespace() string
		// GetServiceVersion is a function.
		GetServiceVersion() string
	}

	// GetOtelConfigger is an interface.
	GetOtelConfigger interface {
		// GetOtelConfigger is a function.
		GetOtelConfigger() OtelConfigger
	}

	otelConfig struct {
		exporterJaegerEndpoint string
		exporterJaegerPassword string
		exporterJaegerUsername string
		instrumentationName    string
		serviceInstanceID      string
		serviceName            string
		serviceNamespace       string
		serviceVersion         string
	}

	otelConfigOptioner interface {
		apply(*otelConfig)
	}

	otelConfigOptionerFunc func(*otelConfig)
)

var (
	_ OtelConfigger  = (*otelConfig)(nil)
	_ json.Marshaler = (*otelConfig)(nil)
	_ object.GetMap  = (*otelConfig)(nil)
)

// NewOtelConfig is a function.
func NewOtelConfig(
	optioners ...otelConfigOptioner,
) *otelConfig {
	otelConfig := &otelConfig{
		exporterJaegerEndpoint: object.URIEmpty,
		exporterJaegerPassword: object.URIEmpty,
		exporterJaegerUsername: object.URIEmpty,
		instrumentationName:    object.URIEmpty,
		serviceInstanceID:      object.URIEmpty,
		serviceName:            object.URIEmpty,
		serviceNamespace:       object.URIEmpty,
		serviceVersion:         object.URIEmpty,
	}

	return otelConfig.WithOptioners(optioners...)
}

// WithOtelConfigExporterJaegerEndpoint is a function.
func WithOtelConfigExporterJaegerEndpoint(
	exporterJaegerEndpoint string,
) otelConfigOptioner {
	return otelConfigOptionerFunc(func(
		config *otelConfig,
	) {
		config.exporterJaegerEndpoint = exporterJaegerEndpoint
	})
}

// WithOtelConfigExporterJaegerPassword is a function.
func WithOtelConfigExporterJaegerPassword(
	exporterJaegerPassword string,
) otelConfigOptioner {
	return otelConfigOptionerFunc(func(
		config *otelConfig,
	) {
		config.exporterJaegerPassword = exporterJaegerPassword
	})
}

// WithOtelConfigExporterJaegerUsername is a function.
func WithOtelConfigExporterJaegerUsername(
	exporterJaegerUsername string,
) otelConfigOptioner {
	return otelConfigOptionerFunc(func(
		config *otelConfig,
	) {
		config.exporterJaegerUsername = exporterJaegerUsername
	})
}

// WithOtelConfigInstrumentationName is a function.
func WithOtelConfigInstrumentationName(
	instrumentationName string,
) otelConfigOptioner {
	return otelConfigOptionerFunc(func(
		config *otelConfig,
	) {
		config.instrumentationName = instrumentationName
	})
}

// WithOtelConfigServiceInstanceID is a function.
func WithOtelConfigServiceInstanceID(
	serviceInstanceID string,
) otelConfigOptioner {
	return otelConfigOptionerFunc(func(
		config *otelConfig,
	) {
		config.serviceInstanceID = serviceInstanceID
	})
}

// WithOtelConfigServiceName is a function.
func WithOtelConfigServiceName(
	serviceName string,
) otelConfigOptioner {
	return otelConfigOptionerFunc(func(
		config *otelConfig,
	) {
		config.serviceName = serviceName
	})
}

// WithOtelConfigServiceNamespace is a function.
func WithOtelConfigServiceNamespace(
	serviceNamespace string,
) otelConfigOptioner {
	return otelConfigOptionerFunc(func(
		config *otelConfig,
	) {
		config.serviceNamespace = serviceNamespace
	})
}

// WithOtelConfigServiceVersion is a function.
func WithOtelConfigServiceVersion(
	serviceVersion string,
) otelConfigOptioner {
	return otelConfigOptionerFunc(func(
		config *otelConfig,
	) {
		config.serviceVersion = serviceVersion
	})
}

// GetExporterJaegerEndpoint is a function.
func (config *otelConfig) GetExporterJaegerEndpoint() string {
	return config.exporterJaegerEndpoint
}

// GetExporterJaegerPassword is a function.
func (config *otelConfig) GetExporterJaegerPassword() string {
	return config.exporterJaegerPassword
}

// GetExporterJaegerUsername is a function.
func (config *otelConfig) GetExporterJaegerUsername() string {
	return config.exporterJaegerUsername
}

// GetInstrumentationName is a function.
func (config *otelConfig) GetInstrumentationName() string {
	return config.instrumentationName
}

// GetServiceInstanceID is a function.
func (config *otelConfig) GetServiceInstanceID() string {
	return config.serviceInstanceID
}

// GetServiceName is a function.
func (config *otelConfig) GetServiceName() string {
	return config.serviceName
}

// GetServiceNamespace is a function.
func (config *otelConfig) GetServiceNamespace() string {
	return config.serviceNamespace
}

// GetServiceVersion is a function.
func (config *otelConfig) GetServiceVersion() string {
	return config.serviceVersion
}

// GetMap is a function.
func (config *otelConfig) GetMap() map[string]any {
	return map[string]any{
		"exporter_jaeger_endpoint": config.GetExporterJaegerEndpoint(),
		"exporter_jaeger_password": config.GetExporterJaegerPassword(),
		"exporter_jaeger_username": config.GetExporterJaegerUsername(),
		"instrumentation_name":     config.GetInstrumentationName(),
		"service_instance_id":      config.GetServiceInstanceID(),
		"service_name":             config.GetServiceName(),
		"service_namespace":        config.GetServiceNamespace(),
		"service_version":          config.GetServiceVersion(),
	}
}

// MarshalJSON is a function.
// read more https://pkg.go.dev/encoding/json#Marshaler
func (config *otelConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(config.GetMap())
}

// WithOptioners is a function.
func (config *otelConfig) WithOptioners(
	optioners ...otelConfigOptioner,
) *otelConfig {
	newConfig := config.clone()
	for _, optioner := range optioners {
		optioner.apply(newConfig)
	}

	return newConfig
}

func (config *otelConfig) clone() *otelConfig {
	newConfig := config

	return newConfig
}

func (optionerFunc otelConfigOptionerFunc) apply(
	config *otelConfig,
) {
	optionerFunc(config)
}
