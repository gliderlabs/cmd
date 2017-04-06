package builtin

import (
	"github.com/gliderlabs/cmd/lib/cli"
	"github.com/spf13/cobra"
)

var helpCmd = func(sess cli.Session) *cobra.Command {
	return &cobra.Command{
		Use:    "help",
		Short:  "Print this help",
		Hidden: true,
		RunE: func(c *cobra.Command, args []string) error {
			c.Root().Help()
			return nil
		},
	}
}
