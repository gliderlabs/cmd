package builtin

import (
	"github.com/gliderlabs/cmd/lib/cli"
	"github.com/spf13/cobra"
	"github.com/gliderlabs/cmd/app/store"
	"fmt"
)

var sourceCmd = func(sess cli.Session) *cobra.Command {
	return &cobra.Command{
		Use:   "source <name>",
		Short: "Display the command source",
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) < 1 {
				fmt.Fprintln(sess.Stderr(), "Name is a required argument")
				sess.Exit(cli.StatusUsageError)
				return nil
			}
			cmd := store.Selected().Get(sess.User(), args[0])
			if cmd == nil {
				fmt.Fprintln(sess.Stderr(), "Command", cli.Bright(args[0]), "does not exist")
				sess.Exit(cli.StatusUnknownCommand)
				return nil
			}
			fmt.Fprintln(sess.Stderr(), cmd.Source)
			return nil
		},
	}
}