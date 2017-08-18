package com

// Context interface for dynamically enabled components.
//
// Implement this to create an object to pass into Enabled to filter out
// components that are not enabled, as determined by whatever mechanism you want.
type Context interface {
	ComponentEnabled(name string) bool
}
