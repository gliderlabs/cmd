package meta

import (
	"fmt"
	"io/ioutil"

	"github.com/gliderlabs/comlab/pkg/log"
	"github.com/gliderlabs/ssh"
	"github.com/spf13/cobra"

	"github.com/progrium/cmd/com/cmd"
	"github.com/progrium/cmd/com/core"
	"github.com/progrium/cmd/com/store"
)

var rootHelp = &cmd.MetaCommand{
	Use:     ":help",
	Aliases: []string{"cmd-help"},
	Short:   "Print this help",
	Hidden:  true,
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
	Use:     ":ls",
	Aliases: []string{"cmd-ls"},
	Short:   "List installed commands",
	Run: func(meta *cmd.MetaCommand, sess ssh.Session, args []string) {
		header(sess, "Installed Commands")
		for _, cmd := range store.Selected().List(sess.User()) {
			fmt.Fprintf(sess, "  %-10s  %s\n", cmd.Name, cmd.Description)
		}
		fmt.Fprintln(sess, "")
	},
}

var rootInstall = &cmd.MetaCommand{
	Use:     ":add <name> <source>",
	Aliases: []string{"cmd-add"},
	Short:   "Install a command",
	Run: func(meta *cmd.MetaCommand, sess ssh.Session, args []string) {
		statusMsg(sess, "Installing")
		limit := core.Plans[core.DefaultPlan].MaxCmds
		cmds := store.Selected().List(sess.User())
		if limit >= 0 && len(cmds) >= limit {
			statusErr(sess.Stderr(), "Unable to install command: command limit for plan reached")
			sess.Exit(1)
			return
		}

		if len(args) < 1 {
			statusErr(sess.Stderr(), "Must specify a name")
			sess.Exit(1)
			return
		}
		if len(args) < 2 {
			statusErr(sess.Stderr(), "Must specify a source")
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
			statusErr(sess.Stderr(), "Command unable to install: "+err.Error())
			sess.Exit(1)
			return
		}
		if err := store.Selected().Put(cmd.User, cmd.Name, cmd); err != nil {
			log.Info(sess, cmd, err)
			statusErr(sess.Stderr(), err.Error())
			sess.Exit(255)
			return
		}
		statusDone(sess)
	},
}

var rootUninstall = &cmd.MetaCommand{
	Use:     ":delete <name>",
	Aliases: []string{"cmd-rm"},
	Short:   "Delete a command",
	Run: func(meta *cmd.MetaCommand, sess ssh.Session, args []string) {
		if len(args) < 1 {
			statusErr(sess.Stderr(), "Must specify a command")
			sess.Exit(1)
			return
		}

		statusMsg(sess, "Deleting")
		cmd := store.Selected().Get(sess.User(), args[0])
		if cmd == nil {
			statusErr(sess.Stderr(), "Command not found")
			sess.Exit(1)
			return
		}
		if err := store.Selected().Delete(cmd.User, cmd.Name); err != nil {
			log.Info(sess, cmd, err)
			statusErr(sess.Stderr(), err.Error())
			sess.Exit(255)
			return
		}
		statusDone(sess)
	},
}

var rootCreate = &cmd.MetaCommand{
	Use:     ":create <name>",
	Aliases: []string{"cmd-create"},
	Short:   "Create a command",
	Run: func(meta *cmd.MetaCommand, sess ssh.Session, args []string) {
		statusMsg(sess, "Creating command")
		limit := core.Plans[core.DefaultPlan].MaxCmds
		cmds := store.Selected().List(sess.User())
		if limit >= 0 && len(cmds) >= limit {
			statusErr(sess.Stderr(), "Unable to create command: command limit for plan reached")
			sess.Exit(1)
			return
		}

		if len(args) < 1 {
			statusErr(sess.Stderr(), "Must specify a name")
			sess.Exit(1)
			return
		}

		source, err := ioutil.ReadAll(sess)
		if err != nil {
			statusErr(sess.Stderr(),
				"Command unable to install: failed to read source: "+err.Error())
			sess.Exit(1)
			return
		}
		cmd := &core.Command{
			Name:   args[0],
			User:   sess.User(),
			Source: string(source),
		}

		if err := cmd.Build(); err != nil {
			log.Info(err)
			statusErr(sess.Stderr(), "Command unable to install: "+err.Error())
			sess.Exit(1)
			return
		}
		if err := store.Selected().Put(cmd.User, cmd.Name, cmd); err != nil {
			log.Info(sess, cmd, err)
			statusErr(sess.Stderr(), err.Error())
			sess.Exit(255)
			return
		}
		statusDone(sess)
	},
}
