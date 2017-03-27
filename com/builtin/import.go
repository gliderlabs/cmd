package builtin

import (
	"fmt"

	"github.com/gliderlabs/comlab/pkg/log"
	"github.com/spf13/cobra"

	"github.com/progrium/cmd/com/cli"
	"github.com/progrium/cmd/com/core"
	"github.com/progrium/cmd/com/store"
)

var importCmd = cli.Command{
	Use:     "import <name> <source>",
	Aliases: []string{"add"},
	Hidden:  true,
	Short:   "Import a command from Docker image",
	Run: func(c *cobra.Command, args []string) {
		sess := cli.ContextSession(c)
		limit := core.ContextPlan(cli.Context(c)).MaxCmds
		cmds := store.Selected().List(sess.User())
		if len(cmds) >= limit {
			fmt.Fprintln(sess.Stderr(), "Command limit for plan reached")
			sess.Exit(cli.StatusNoPerm)
			return
		}
		if len(args) < 1 {
			fmt.Fprintln(sess.Stderr(), "Name is a required argument")
			sess.Exit(cli.StatusUsageError)
			return
		}
		if len(args) < 2 {
			fmt.Fprintln(sess.Stderr(), "Source is a required argument")
			sess.Exit(cli.StatusUsageError)
			return
		}
		for _, c := range cmds {
			if c.Name == args[0] {
				fmt.Fprintln(sess.Stderr(), "Command", cli.Bright(args[0]), "already exists")
				sess.Exit(cli.StatusCreateError)
				return
			}
		}
		cli.Status(sess, "Importing command")
		cmd := &core.Command{
			Name:   args[0],
			User:   sess.User(),
			Source: args[1],
		}
		if err := cmd.Pull(cli.Context(c)); err != nil {
			log.Info(err)
			cli.StatusErr(sess.Stderr(), "Command unable to install: "+err.Error())
			sess.Exit(cli.StatusError)
			return
		}
		if err := store.Selected().Put(cmd.User, cmd.Name, cmd); err != nil {
			log.Info(sess, cmd, err)
			cli.StatusErr(sess.Stderr(), err.Error())
			sess.Exit(cli.StatusInternalError)
			return
		}
		cli.StatusDone(sess)
	},
}.Init(nil)
