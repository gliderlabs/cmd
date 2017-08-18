package com

import (
	"github.com/spf13/cast"
)

// ConfigProvider is the interface expected by SetConfig and is used to look
// up configuration. It includes Context, so a ConfigProvider must define how
// the underlying configuration specifies if components are enabled.
//
// The getter methods return a value and a boolean of whether it was found. This
// informs com whether to use the value returned here or use the default value.
type ConfigProvider interface {
	Context
	GetString(key string) (string, bool)
	GetInt(key string) (int, bool)
	GetBool(key string) (bool, bool)
}

type mapConfig map[string]interface{}

func (c mapConfig) ComponentEnabled(name string) bool {
	v, ok := c[name+".enabled"]
	if !ok {
		return true
	}
	return cast.ToBool(v)
}

func (c mapConfig) GetString(key string) (string, bool) {
	v, ok := c[key]
	return cast.ToString(v), ok
}

func (c mapConfig) GetInt(key string) (int, bool) {
	v, ok := c[key]
	return cast.ToInt(v), ok
}

func (c mapConfig) GetBool(key string) (bool, bool) {
	v, ok := c[key]
	return cast.ToBool(v), ok
}
