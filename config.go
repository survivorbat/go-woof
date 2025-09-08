package gowoof

import "github.com/go-viper/mapstructure/v2"

// Config is used internally
type Config struct {
	DecodeConfig *mapstructure.DecoderConfig
}

// Option is used to implement the Functional Option pattern in ParseTable.
type Option func(*Config) error

// WithDecodeConfig allows you to set a custom DecoderConfig that will be used when decoding to a
// struct. The `Result` property is overriden during parsing. `WeaklyTypedInput` is enabled by default.
func WithDecodeConfig(cfg *mapstructure.DecoderConfig) Option {
	return func(c *Config) error {
		c.DecodeConfig = cfg
		return nil
	}
}
