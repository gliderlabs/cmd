package release

import (
	"fmt"
	"os"
)

const (
	ChannelDev    = "dev"
	ChannelAlpha  = "alpha"
	ChannelBeta   = "beta"
	ChannelStable = "stable"
)

var (
	Version string
	Build   string
)

func Channel() string {
	channel := os.Getenv("CHANNEL")
	if channel == "" {
		return ChannelDev
	}
	return channel
}

func DisplayVersion() string {
	if len(Version) < 8 {
		return ChannelDev
	}
	return Version[:8]
}

func Hostname() string {
	if Channel() == ChannelDev {
		return "localhost"
	}
	if Channel() == ChannelStable {
		return "cmd.io"
	}
	return fmt.Sprintf("%s.cmd.io", Channel())
}
