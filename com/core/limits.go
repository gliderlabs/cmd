package core

import (
	"context"
	"fmt"
	"time"
)

var (
	// ErrMaxRuntimeExceded returned when command runtime exceedes plan limit
	ErrMaxRuntimeExceded = fmt.Errorf("maximum runtime exceded")
)

const DefaultPlan = "basic"

var Plans = map[string]Plan{
	"basic": {
		MaxCmds:    10,
		MaxRuntime: 30 * time.Second,
		ImageSize:  512 << 20, // 512mb
		Memory:     512 << 20, // 512mb
		DinD:       false,
		// 20% of 1 CPU
		CPUPeriod: (50 * time.Millisecond).Nanoseconds() / 1000, // 50000 microseconds
		CPUQuota:  (10 * time.Millisecond).Nanoseconds() / 1000, // 10000 microseconds
	},
	"plus": {
		MaxCmds:    100,
		MaxRuntime: 5 * time.Minute,
		ImageSize:  2 << 30, // 2gb
		Memory:     2 << 30, // 2gb
		DinD:       false,
		// 20% of 1 CPU
		CPUPeriod: (50 * time.Millisecond).Nanoseconds() / 1000, // 50000 microseconds
		CPUQuota:  (10 * time.Millisecond).Nanoseconds() / 1000, // 10000 microseconds
	},
	"contrib": {
		MaxCmds:    100,
		MaxRuntime: 10 * time.Minute,
		ImageSize:  2 << 30, // 2gb
		Memory:     2 << 30, // 2gb
		DinD:       true,
		// 20% of 1 CPU
		CPUPeriod: (50 * time.Millisecond).Nanoseconds() / 1000, // 50000 microseconds
		CPUQuota:  (10 * time.Millisecond).Nanoseconds() / 1000, // 10000 microseconds
	},
}

func ContextPlan(ctx context.Context) Plan {
	if cp := ctx.Value("plan"); cp != nil {
		return Plans[cp.(string)]
	}
	return Plans[DefaultPlan]
}

// Plan describes limits for a specific plan.
type Plan struct {
	MaxCmds    int
	MaxRuntime time.Duration
	ImageSize  int64 // size in bytes
	CPUPeriod  int64 // length of a period (in microseconds)
	CPUQuota   int64 // total available run-time within a period (in microseconds)
	Memory     int64
	DinD       bool // docker in docker, currently uses host docker
}
