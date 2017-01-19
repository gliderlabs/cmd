package meta

import (
	"fmt"

	"github.com/gliderlabs/gosper/pkg/log"
	"github.com/gliderlabs/ssh"
	"github.com/spf13/cobra"

	"github.com/progrium/cmd/com/cmd"
	"github.com/progrium/cmd/com/core"
	"github.com/progrium/cmd/com/store"
)

var rootHelp = &cmd.MetaCommand{
	Use:    ":help",
	Short:  "Print this help",
	Hidden: true,
	Run: func(cmd *cmd.MetaCommand, sess ssh.Session, args []string) {
		for _, c := range store.Selected().List(sess.User()) {
			cmd.Cmd.Parent().AddCommand(&cobra.Command{
				Use:   c.Name,
				Short: c.Description,
				Run:   func(cmd *cobra.Command, args []string) {},
			})
		}
		cmd.Cmd.Parent().Help()
	},
}

var rootList = &cmd.MetaCommand{
	Use:   ":ls",
	Short: "List installed commands",
	Run: func(meta *cmd.MetaCommand, sess ssh.Session, args []string) {
		fmt.Fprintln(sess, "")
		fmt.Fprintln(sess, "Installed Commands:")
		for _, cmd := range store.Selected().List(sess.User()) {
			fmt.Fprintf(sess, "  %-10s  %s\n", cmd.Name, cmd.Description)
		}
		fmt.Fprintln(sess, "")
	},
}

var rootInstall = &cmd.MetaCommand{
	Use:   ":add <name> <source>",
	Short: "Install a command",
	Run: func(meta *cmd.MetaCommand, sess ssh.Session, args []string) {
		limit := core.Plans[core.DefaultPlan].MaxCmds
		cmds := store.Selected().List(sess.User())
		if limit >= 0 && len(cmds) >= limit {
			fmt.Fprintln(sess, "Unable to install command: command limit for plan reached")
			sess.Exit(1)
			return
		}

		if len(args) < 1 {
			fmt.Fprintln(sess, "Must specify a name")
			sess.Exit(1)
			return
		}
		if len(args) < 2 {
			fmt.Fprintln(sess, "Must specify a source")
			sess.Exit(1)
			return
		}
		cmd := &core.Command{
			Name:   args[0],
			User:   sess.User(),
			Source: args[1],
		}
		if err := cmd.Pull(); err != nil {
			log.Info(err)
			fmt.Fprintln(sess.Stderr(), "Command unable to install:", err)
			sess.Exit(1)
			return
		}
		if err := store.Selected().Put(cmd.User, cmd.Name, cmd); err != nil {
			log.Info(sess, cmd, err)
			fmt.Fprintln(sess.Stderr(), err.Error())
			sess.Exit(255)
			return
		}
		fmt.Fprintln(sess, "Command installed")
	},
}

var rootUninstall = &cmd.MetaCommand{
	Use:   ":rm <name>",
	Short: "Uninstall a command",
	Run: func(meta *cmd.MetaCommand, sess ssh.Session, args []string) {
		if len(args) < 1 {
			fmt.Fprintln(sess, "Must specify a command")
			sess.Exit(1)
			return
		}
		cmd := store.Selected().Get(sess.User(), args[0])
		if cmd == nil {
			fmt.Fprintln(sess, "Command not found")
			sess.Exit(1)
			return
		}
		if err := store.Selected().Delete(cmd.User, cmd.Name); err != nil {
			log.Info(sess, cmd, err)
			fmt.Fprintln(sess.Stderr(), err.Error())
			sess.Exit(255)
			return
		}
		fmt.Fprintln(sess, "Command uninstalled")
	},
}
