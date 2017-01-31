package core

import (
	"fmt"
	"time"
)

var (
	// ErrMaxRuntimeExceded returned when command runtime exceedes plan limit
	ErrMaxRuntimeExceded = fmt.Errorf("maximum runtime exceded")
)

const DefaultPlan = "basic"

var Plans = map[string]Plan{
	DefaultPlan: {
		MaxCmds:    10,
		MaxRuntime: 10 * time.Minute,

		ImageSize: 500 << 20, // 500mb
		// 20% of 1 CPU
		CPUPeriod: (50 * time.Millisecond).Nanoseconds() / 1000, // 50000 microseconds
		CPUQuota:  (10 * time.Millisecond).Nanoseconds() / 1000, // 10000 microseconds
		Memory:    300 << 20,                                    // 300mb
	},
}

// Plan describes limits for a specific plan.
type Plan struct {
	MaxCmds    int
	MaxRuntime time.Duration

	ImageSize int64 // size in bytes

	CPUPeriod int64 // length of a period (in microseconds)
	CPUQuota  int64 // total available run-time within a period (in microseconds)

	Memory int64
}
