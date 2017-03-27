package builtin

import (
	"github.com/progrium/cmd/lib/cli"
	"github.com/spf13/cobra"
)

var helpCmd = cli.Command{
	Use:    "help",
	Short:  "Print this help",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Root().Help()
	},
}.Init(nil)
