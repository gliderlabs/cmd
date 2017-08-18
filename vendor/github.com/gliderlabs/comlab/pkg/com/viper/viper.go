package viper

import (
	"fmt"

	v "github.com/spf13/viper"
)

// Config wraps Viper and implements ConfigProvider
type Config struct {
	*v.Viper
}

// NewConfig creates a Config based on viper.New()
func NewConfig() *Config {
	return &Config{v.New()}
}

// ComponentEnabled determines if a component is enabled by looking up a
// boolean value called "enabled" under the component namespace. If it's not set,
// it defaults to true.
func (c *Config) ComponentEnabled(name string) bool {
	fqn := fmt.Sprintf("%s.enabled", name)
	if !c.Viper.IsSet(fqn) {
		return true
	}
	return c.Viper.GetBool(fqn)
}

func (c *Config) GetString(key string) (string, bool) {
	if c.Viper.IsSet(key) {
		return c.Viper.GetString(key), true
	}
	return "", false
}

func (c *Config) GetInt(key string) (int, bool) {
	if c.Viper.IsSet(key) {
		return c.Viper.GetInt(key), true
	}
	return 0, false
}

func (c *Config) GetBool(key string) (bool, bool) {
	if c.Viper.IsSet(key) {
		return c.Viper.GetBool(key), true
	}
	return false, false
}
