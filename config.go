package gowoof

import "github.com/go-viper/mapstructure/v2"

// Config is used internally
type Config struct {
	DecodeConfig *mapstructure.DecoderConfig
	Vertical     bool
}

// Option is used to implement the Functional Option pattern in ParseTable.
type Option func(*Config) error

// WithDecodeConfig allows you to set a custom DecoderConfig that will be used when decoding to a
// struct. The `Result` property is overridden during parsing. `WeaklyTypedInput` is enabled by default.
func WithDecodeConfig(cfg *mapstructure.DecoderConfig) Option {
	return func(c *Config) error {
		c.DecodeConfig = cfg
		return nil
	}
}

// Vertical chnanges the parsing from a row-based approach to a column-based one. This means
// that instead of this:
//
// | ID | Name | Job |
// | 1  | Foo | Dev |
// | 2  | Bar | Ops |
//
// Input is expected to be this:
//
// | ID    | 1   | 2   |
// | Name  | Foo | Bar |
// | Job   | Dev | Ops |
//
// In case you want to do that ¯\_(ツ)_/¯
func Vertical() Option {
	return func(c *Config) error {
		c.Vertical = true
		return nil
	}
}
