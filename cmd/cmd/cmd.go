package main

import (
	"os"

	"github.com/gliderlabs/comlab/lib/daemon"
	"github.com/mitchellh/panicwrap"
	"github.com/progrium/cmd/com/sentry"
)

func main() {
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

	daemon.Run("cmd")
}
