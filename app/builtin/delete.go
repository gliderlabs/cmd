package builtin

import (
	"fmt"

	"github.com/gliderlabs/comlab/pkg/log"
	"github.com/spf13/cobra"

	"github.com/progrium/cmd/app/store"
	"github.com/progrium/cmd/lib/cli"
)

var deleteCmd = func(sess cli.Session) *cobra.Command {
	return &cobra.Command{
		Use:     "delete <name>",
		Aliases: []string{"rm"},
		Short:   "Delete a command",
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) < 1 {
				fmt.Fprintln(sess.Stderr(), "Name is a required argument")
				sess.Exit(cli.StatusUsageError)
				return nil
			}
			cmd, err := LookupCmd(sess.User(), args[0])
			if err != nil {
				fmt.Fprintln(sess.Stderr(), err.Error())
				sess.Exit(cli.StatusError)
				return nil
			}
			cli.Status(sess, "Deleting command")
			if err := store.Selected().Delete(cmd.User, cmd.Name); err != nil {
				log.Info(sess, cmd, err)
				cli.StatusErr(sess.Stderr(), err.Error())
				sess.Exit(cli.StatusInternalError)
				return nil
			}
			cli.StatusDone(sess)
			return nil
		},
	}
}
