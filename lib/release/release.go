package release

import "os"

var (
	Version string
	Build   string
)

func Channel() string {
	channel := os.Getenv("CHANNEL")
	if channel == "" {
		return "dev"
	}
	return channel
}

func DisplayVersion() string {
	if len(Version) < 8 {
		return "dev"
	}
	return Version[:8]
}
