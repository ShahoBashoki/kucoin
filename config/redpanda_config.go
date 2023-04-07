package config

import (
	"encoding/json"
	"fmt"

	"github.com/ShahoBashoki/kucoin/object"
)

type (
	// RedpandaConfigger is an interface.
	RedpandaConfigger interface {
		// GetProxyURL is a function.
		GetProxyURL() string
		// GetTopic is a function.
		GetTopic() string
	}

	// GetRedpandaConfigger is an interface.
	GetRedpandaConfigger interface {
		// GetRedpandaConfigger is a function.
		GetRedpandaConfigger() RedpandaConfigger
	}

	redpandaConfig struct {
		proxyURL string
		topic    string
	}

	redpandaConfigOptioner interface {
		apply(*redpandaConfig)
	}

	redpandaConfigOptionerFunc func(*redpandaConfig)
)

var (
	_ RedpandaConfigger = (*redpandaConfig)(nil)
	_ json.Marshaler    = (*redpandaConfig)(nil)
	_ object.GetMap     = (*redpandaConfig)(nil)
)

// NewRedpandaConfig is a function.
func NewRedpandaConfig(
	optioners ...redpandaConfigOptioner,
) *redpandaConfig {
	redpandaConfig := &redpandaConfig{
		proxyURL: object.URIEmpty,
		topic:    object.URIEmpty,
	}

	return redpandaConfig.WithOptioners(optioners...)
}

// WithRedpandaConfigProxyURL is a function.
func WithRedpandaConfigProxyURL(
	proxyURL string,
) redpandaConfigOptioner {
	return redpandaConfigOptionerFunc(func(
		config *redpandaConfig,
	) {
		config.proxyURL = proxyURL
	})
}

// WithRedpandaConfigTopic is a function.
func WithRedpandaConfigTopic(
	topic string,
) redpandaConfigOptioner {
	return redpandaConfigOptionerFunc(func(
		config *redpandaConfig,
	) {
		config.topic = topic
	})
}

// GetProxyURL is a function.
func (config *redpandaConfig) GetProxyURL() string {
	return fmt.Sprintf(
		object.URIURLPath,
		config.proxyURL,
		fmt.Sprintf(object.URIRedpandaTopic, config.GetTopic()),
	)
}

// GetTopic is a function.
func (config *redpandaConfig) GetTopic() string {
	return config.topic
}

// GetMap is a function.
func (config *redpandaConfig) GetMap() map[string]any {
	return map[string]any{
		"proxy_url": config.GetProxyURL(),
		"topic":     config.GetTopic(),
	}
}

// MarshalJSON is a function.
// read more https://pkg.go.dev/encoding/json#Marshaler
func (config *redpandaConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(config.GetMap())
}

// WithOptioners is a function.
func (config *redpandaConfig) WithOptioners(
	optioners ...redpandaConfigOptioner,
) *redpandaConfig {
	newConfig := config.clone()
	for _, optioner := range optioners {
		optioner.apply(newConfig)
	}

	return newConfig
}

func (config *redpandaConfig) clone() *redpandaConfig {
	newConfig := config

	return newConfig
}

func (optionerFunc redpandaConfigOptionerFunc) apply(
	config *redpandaConfig,
) {
	optionerFunc(config)
}
