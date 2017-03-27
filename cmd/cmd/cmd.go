package main

import (
	"os"

	"github.com/mitchellh/panicwrap"
	"github.com/progrium/cmd/lib/daemon"
	"github.com/progrium/cmd/lib/release"
	"github.com/progrium/cmd/lib/sentry"
)

var (
	Version string
	Build   string
)

func main() {
	// TODO: panic wrap should be added to daemon, with hook to allow
	// sentry to handle panic
	if !daemon.LocalMode() {

		exitStatus, err := panicwrap.BasicWrap(sentry.PanicHandler)
		if err != nil {
			// Something went wrong setting up the panic wrapper. Unlikely,
			// but possible.
			panic(err)
		}

		// If exitStatus >= 0, then we're the parent process and the panicwrap
		// re-executed ourselves and completed. Just exit with the proper status.
		if exitStatus >= 0 {
			os.Exit(exitStatus)
		}
	}

	release.Build = Build
	release.Version = Version
	daemon.Run("cmd")
}
