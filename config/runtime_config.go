package config

import (
	"encoding/json"

	"github.com/ShahoBashoki/kucoin/object"
)

type (
	// RuntimeConfigger is an interface.
	RuntimeConfigger interface {
		// GetValidateMapRules is a function.
		GetValidateMapRules() map[string][]map[string]any
		// GetNode is a function.
		GetNode() string
		// GetKucoinPaginationRequestSize is a function.
		GetKucoinPaginationRequestSize() int64
	}

	// GetRuntimeConfigger is an interface.
	GetRuntimeConfigger interface {
		// GetRuntimeConfigger is a function.
		GetRuntimeConfigger() RuntimeConfigger
	}

	runtimeConfig struct {
		validateMapRules            map[string][]map[string]any
		node                        string
		kucoinPaginationRequestSize int64
	}

	runtimeConfigOptioner interface {
		apply(*runtimeConfig)
	}

	runtimeConfigOptionerFunc func(*runtimeConfig)
)

var (
	_ RuntimeConfigger = (*runtimeConfig)(nil)
	_ json.Marshaler   = (*runtimeConfig)(nil)
	_ object.GetMap    = (*runtimeConfig)(nil)
)

// NewRuntimeConfig is a function.
func NewRuntimeConfig(
	optioners ...runtimeConfigOptioner,
) *runtimeConfig {
	runtimeConfig := &runtimeConfig{
		validateMapRules:            map[string][]map[string]any{},
		node:                        object.URIEmpty,
		kucoinPaginationRequestSize: 0,
	}

	return runtimeConfig.WithOptioners(optioners...)
}

// WithRuntimeConfigValidateMapRules is a function.
func WithRuntimeConfigValidateMapRules(
	validateMapRules map[string][]map[string]any,
) runtimeConfigOptioner {
	return runtimeConfigOptionerFunc(func(
		config *runtimeConfig,
	) {
		config.validateMapRules = validateMapRules
	})
}

// WithRuntimeConfigNode is a function.
func WithRuntimeConfigNode(
	node string,
) runtimeConfigOptioner {
	return runtimeConfigOptionerFunc(func(
		config *runtimeConfig,
	) {
		config.node = node
	})
}

// WithRuntimeConfigKucoinPaginationRequestSize is a function.
func WithRuntimeConfigKucoinPaginationRequestSize(
	kucoinPaginationRequestSize int64,
) runtimeConfigOptioner {
	return runtimeConfigOptionerFunc(func(
		config *runtimeConfig,
	) {
		config.kucoinPaginationRequestSize = kucoinPaginationRequestSize
	})
}

// GetValidateMapRules is a function.
func (config *runtimeConfig) GetValidateMapRules() map[string][]map[string]any {
	return config.validateMapRules
}

// GetNode is a function.
func (config *runtimeConfig) GetNode() string {
	return config.node
}

// GetKucoinPaginationRequestSize is a function.
func (config *runtimeConfig) GetKucoinPaginationRequestSize() int64 {
	return config.kucoinPaginationRequestSize
}

// GetMap is a function.
func (config *runtimeConfig) GetMap() map[string]any {
	return map[string]any{
		"validate_map_rules":             config.GetValidateMapRules(),
		"node":                           config.GetNode(),
		"kucoin_pagination_request_size": config.GetKucoinPaginationRequestSize(),
	}
}

// MarshalJSON is a function.
// read more https://pkg.go.dev/encoding/json#Marshaler
func (config *runtimeConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(config.GetMap())
}

// WithOptioners is a function.
func (config *runtimeConfig) WithOptioners(
	optioners ...runtimeConfigOptioner,
) *runtimeConfig {
	newConfig := config.clone()
	for _, optioner := range optioners {
		optioner.apply(newConfig)
	}

	return newConfig
}

func (config *runtimeConfig) clone() *runtimeConfig {
	newConfig := config

	return newConfig
}

func (optionerFunc runtimeConfigOptionerFunc) apply(
	config *runtimeConfig,
) {
	optionerFunc(config)
}
