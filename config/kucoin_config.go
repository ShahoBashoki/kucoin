package config

import (
	"encoding/json"

	"github.com/ShahoBashoki/kucoin/object"
)

type (
	// KucoinConfigger is an interface.
	KucoinConfigger interface {
		// GetKey is a function.
		GetKey() string
		// GetPassPhrase is a function.
		GetPassPhrase() string
		// GetSecret is a function.
		GetSecret() string
	}

	// GetKucoinConfigger is an interface.
	GetKucoinConfigger interface {
		// GetKucoinConfigger is a function.
		GetKucoinConfigger() KucoinConfigger
	}

	kucoinConfig struct {
		key        string
		passPhrase string
		secret     string
	}

	kucoinConfigOptioner interface {
		apply(*kucoinConfig)
	}

	kucoinConfigOptionerFunc func(*kucoinConfig)
)

var (
	_ KucoinConfigger = (*kucoinConfig)(nil)
	_ json.Marshaler  = (*kucoinConfig)(nil)
	_ object.GetMap   = (*kucoinConfig)(nil)
)

// NewKucoinConfig is a function.
func NewKucoinConfig(
	optioners ...kucoinConfigOptioner,
) *kucoinConfig {
	kucoinConfig := &kucoinConfig{
		key:        object.URIEmpty,
		passPhrase: object.URIEmpty,
		secret:     object.URIEmpty,
	}

	return kucoinConfig.WithOptioners(optioners...)
}

// WithKucoinConfigKey is a function.
func WithKucoinConfigKey(
	key string,
) kucoinConfigOptioner {
	return kucoinConfigOptionerFunc(func(
		config *kucoinConfig,
	) {
		config.key = key
	})
}

// WithKucoinConfigPassPhrase is a function.
func WithKucoinConfigPassPhrase(
	passPhrase string,
) kucoinConfigOptioner {
	return kucoinConfigOptionerFunc(func(
		config *kucoinConfig,
	) {
		config.passPhrase = passPhrase
	})
}

// WithKucoinConfigSecret is a function.
func WithKucoinConfigSecret(
	secret string,
) kucoinConfigOptioner {
	return kucoinConfigOptionerFunc(func(
		config *kucoinConfig,
	) {
		config.secret = secret
	})
}

// GetKey is a function.
func (config *kucoinConfig) GetKey() string {
	return config.key
}

// GetPassPhrase is a function.
func (config *kucoinConfig) GetPassPhrase() string {
	return config.passPhrase
}

// GetSecret is a function.
func (config *kucoinConfig) GetSecret() string {
	return config.secret
}

// GetMap is a function.
func (config *kucoinConfig) GetMap() map[string]any {
	return map[string]any{
		"key":         config.GetKey(),
		"pass_phrase": config.GetPassPhrase(),
		"secret":      config.GetSecret(),
	}
}

// MarshalJSON is a function.
// read more https://pkg.go.dev/encoding/json#Marshaler
func (config *kucoinConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(config.GetMap())
}

// WithOptioners is a function.
func (config *kucoinConfig) WithOptioners(
	optioners ...kucoinConfigOptioner,
) *kucoinConfig {
	newConfig := config.clone()
	for _, optioner := range optioners {
		optioner.apply(newConfig)
	}

	return newConfig
}

func (config *kucoinConfig) clone() *kucoinConfig {
	newConfig := config

	return newConfig
}

func (optionerFunc kucoinConfigOptionerFunc) apply(
	config *kucoinConfig,
) {
	optionerFunc(config)
}
