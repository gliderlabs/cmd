package builtin

import (
	"fmt"

	"github.com/gliderlabs/comlab/pkg/log"
	"github.com/spf13/cobra"

	"github.com/progrium/cmd/lib/cli"
	"github.com/progrium/cmd/app/store"
)

var deleteCmd = cli.Command{
	Use:     "delete <name>",
	Aliases: []string{"rm"},
	Short:   "Delete a command",
	Run: func(c *cobra.Command, args []string) {
		sess := cli.ContextSession(c)
		if len(args) < 1 {
			fmt.Fprintln(sess.Stderr(), "Name is a required argument")
			sess.Exit(cli.StatusUsageError)
			return
		}
		cmd, err := LookupCmd(sess.User(), args[0])
		if err != nil {
			fmt.Fprintln(sess.Stderr(), err.Error())
			sess.Exit(cli.StatusError)
			return
		}
		cli.Status(sess, "Deleting command")
		if err := store.Selected().Delete(cmd.User, cmd.Name); err != nil {
			log.Info(sess, cmd, err)
			cli.StatusErr(sess.Stderr(), err.Error())
			sess.Exit(cli.StatusInternalError)
			return
		}
		cli.StatusDone(sess)
	},
}.Init(nil)
