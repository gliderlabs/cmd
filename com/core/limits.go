package core

import (
	"fmt"
	"time"
)

// ErrMaxRuntimeExceded returned when command runtime exceedes plan limit
var ErrMaxRuntimeExceded = fmt.Errorf("maximum runtime exceded")

const DefaultPlan = "basic"

var Plans = map[string]Plan{
	DefaultPlan: {MaxRuntime: 10 * time.Minute},
}

// Plan describes limits for a specific plan.
type Plan struct {
	MaxRuntime time.Duration
}
