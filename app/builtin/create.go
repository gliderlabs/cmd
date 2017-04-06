package builtin

import (
	"fmt"
	"io/ioutil"

	"github.com/gliderlabs/comlab/pkg/log"
	"github.com/spf13/cobra"

	"github.com/progrium/cmd/app/billing"
	"github.com/progrium/cmd/app/core"
	"github.com/progrium/cmd/app/store"
	"github.com/progrium/cmd/lib/cli"
)

var createCmd = func(sess cli.Session) *cobra.Command {
	return &cobra.Command{
		Use:   "create <name>",
		Short: "Create a command",
		RunE: func(c *cobra.Command, args []string) error {
			limit := billing.ContextPlan(sess.Context()).MaxCmds
			cmds := store.Selected().List(sess.User())
			if len(cmds) >= limit {
				fmt.Fprintln(sess.Stderr(), "Command limit for plan reached")
				sess.Exit(cli.StatusNoPerm)
				return nil
			}
			if len(args) < 1 {
				fmt.Fprintln(sess.Stderr(), "Name is a required argument")
				sess.Exit(cli.StatusUsageError)
				return nil
			}
			for _, c := range cmds {
				if c.Name == args[0] {
					fmt.Fprintln(sess.Stderr(), "Command", cli.Bright(args[0]), "already exists")
					sess.Exit(cli.StatusCreateError)
					return nil
				}
			}
			cli.Status(sess, "Creating command")
			source, err := ioutil.ReadAll(sess)
			if err != nil {
				cli.StatusErr(sess.Stderr(), err.Error())
				sess.Exit(cli.StatusDataError)
				return nil
			}
			cmd := &core.Command{
				Name:   args[0],
				User:   sess.User(),
				Source: string(source),
			}

			if err := cmd.Build(); err != nil {
				log.Info(err)
				cli.StatusErr(sess.Stderr(), err.Error())
				sess.Exit(cli.StatusInternalError)
				return nil
			}
			if err := store.Selected().Put(cmd.User, cmd.Name, cmd); err != nil {
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
