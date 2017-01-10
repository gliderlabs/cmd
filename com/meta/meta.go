package meta

import (
	"fmt"
	"strings"

	"github.com/gliderlabs/ssh"

	"github.com/progrium/cmd/com/cmd"
	"github.com/progrium/cmd/com/store"
)

var metaHelp = &cmd.MetaCommand{
	Use:    ":help",
	Short:  "Print this help",
	Hidden: true,
	Run: func(meta *cmd.MetaCommand, sess ssh.Session, args []string) {
		meta.Cmd.Parent().Help()
	},
}

var metaConfig = &cmd.MetaCommand{
	Use:   ":config",
	Short: "Manage command configuration",
	Run: func(meta *cmd.MetaCommand, sess ssh.Session, args []string) {
		fmt.Fprintln(sess.Stderr(), "WARN: meta command", meta.Use, "has been deprecated and replaced with", metaEnv.Use)
		metaEnv.Run(meta, sess, args)
	},
	Setup: func(meta *cmd.MetaCommand) {
		metaEnv.Setup(meta)
	},
}

var metaEnv = &cmd.MetaCommand{
	Use:   ":env",
	Short: "Manage command environment",
	Run: func(meta *cmd.MetaCommand, sess ssh.Session, args []string) {
		cmd := meta.ForCmd
		for k, v := range cmd.Environment {
			fmt.Fprintf(sess, "%s=%s\n", k, v)
		}
	},
	Setup: func(meta *cmd.MetaCommand) {
		meta.Add(metaEnvSet, metaEnvUnset)
	},
}

var metaEnvSet = &cmd.MetaCommand{
	Use:   "set <key=value>...",
	Short: "Manage command environment",
	Run: func(meta *cmd.MetaCommand, sess ssh.Session, args []string) {
		cmd := meta.ForCmd
		for _, kvp := range args {
			parts := strings.SplitN(kvp, "=", 2)
			if len(parts) < 2 {
				continue
			}
			cmd.SetEnv(parts[0], parts[1])
		}
		if err := store.Selected().Put(cmd.User, cmd.Name, cmd); err != nil {
			fmt.Fprintln(sess.Stderr(), err.Error())
			sess.Exit(255)
			return
		}
		fmt.Fprintln(sess, "Env updated.")
	},
}

var metaEnvUnset = &cmd.MetaCommand{
	Use:   "unset <key>",
	Short: "Manage command environment",
	Run: func(meta *cmd.MetaCommand, sess ssh.Session, args []string) {
		cmd := meta.ForCmd
		for _, key := range args {
			delete(cmd.Environment, key)
		}
		if err := store.Selected().Put(cmd.User, cmd.Name, cmd); err != nil {
			fmt.Fprintln(sess.Stderr(), err.Error())
			sess.Exit(255)
			return
		}
		fmt.Fprintln(sess, "Env updated.")
	},
}

var metaAccess = &cmd.MetaCommand{
	Use:   ":access",
	Short: "Manage command access",
	Run: func(meta *cmd.MetaCommand, sess ssh.Session, args []string) {
		cmd := meta.ForCmd
		if len(cmd.ACL) == 0 {
			fmt.Fprintln(sess, "Nobody else has access to this command.")
			return
		}
		fmt.Fprintln(sess, "These GitHub users have access:\n")
		for _, user := range cmd.Admins {
			fmt.Fprintf(sess, "  %s\n", user)
		}
		for _, user := range cmd.ACL {
			fmt.Fprintf(sess, "  %s\n", user)
		}
		fmt.Fprintln(sess, "")
	},
	Setup: func(meta *cmd.MetaCommand) {
		meta.Add(metaAccessPrivate, metaAccessAdd, metaAccessRemove)
	},
}

var metaAccessAdd = &cmd.MetaCommand{
	Use:   "add <user>",
	Short: "Give a GitHub user access",
	Run: func(meta *cmd.MetaCommand, sess ssh.Session, args []string) {
		if len(args) < 1 {
			fmt.Fprintln(sess, "Must specify a GitHub user")
			sess.Exit(1)
			return
		}
		cmd := meta.ForCmd
		cmd.AddAccess(args[0])
		if err := store.Selected().Put(cmd.User, cmd.Name, cmd); err != nil {
			fmt.Fprintln(sess.Stderr(), err.Error())
			sess.Exit(255)
			return
		}
		fmt.Fprintln(sess, "Access granted.")
	},
}

var metaAccessRemove = &cmd.MetaCommand{
	Use:   "rm <user>",
	Short: "Take access for GitHub user",
	Run: func(meta *cmd.MetaCommand, sess ssh.Session, args []string) {
		if len(args) < 1 {
			fmt.Fprintln(sess, "Must specify a GitHub user")
			sess.Exit(1)
			return
		}
		cmd := meta.ForCmd
		cmd.RemoveAccess(args[0])
		if err := store.Selected().Put(cmd.User, cmd.Name, cmd); err != nil {
			fmt.Fprintln(sess.Stderr(), err.Error())
			sess.Exit(255)
			return
		}
		fmt.Fprintln(sess, "Access revoked.")
	},
}

var metaAccessPrivate = &cmd.MetaCommand{
	Use:   "private",
	Short: "Make command private again",
	Run: func(meta *cmd.MetaCommand, sess ssh.Session, args []string) {
		cmd := meta.ForCmd
		cmd.MakePrivate()
		if err := store.Selected().Put(cmd.User, cmd.Name, cmd); err != nil {
			fmt.Fprintln(sess.Stderr(), err.Error())
			sess.Exit(255)
			return
		}
		fmt.Fprintln(sess, "Command is now private.")
	},
}

var metaAdmins = &cmd.MetaCommand{
	Use:   ":admins",
	Short: "Manage command admins",
	Run: func(meta *cmd.MetaCommand, sess ssh.Session, args []string) {
		cmd := meta.ForCmd
		if len(cmd.Admins) == 0 {
			fmt.Fprintln(sess, "Nobody else is admin for this command.")
			return
		}
		fmt.Fprintln(sess, "These GitHub users are also admins:\n")
		for _, user := range cmd.Admins {
			fmt.Fprintf(sess, "  %s\n", user)
		}
		fmt.Fprintln(sess, "")
	},
	Setup: func(meta *cmd.MetaCommand) {
		meta.Add(metaAdminsAdd, metaAdminsRemove)
	},
}

var metaAdminsAdd = &cmd.MetaCommand{
	Use:   "add <user>",
	Short: "Make GitHub user an admin",
	Run: func(meta *cmd.MetaCommand, sess ssh.Session, args []string) {
		if len(args) < 1 {
			fmt.Fprintln(sess, "Must specify a GitHub user")
			sess.Exit(1)
			return
		}
		cmd := meta.ForCmd
		cmd.AddAdmin(args[0])
		if err := store.Selected().Put(cmd.User, cmd.Name, cmd); err != nil {
			fmt.Fprintln(sess.Stderr(), err.Error())
			sess.Exit(255)
			return
		}
		fmt.Fprintln(sess, "Admin granted.")
	},
}

var metaAdminsRemove = &cmd.MetaCommand{
	Use:   "rm <user>",
	Short: "Revoke admin for GitHub user",
	Run: func(meta *cmd.MetaCommand, sess ssh.Session, args []string) {
		if len(args) < 1 {
			fmt.Fprintln(sess, "Must specify a GitHub user")
			sess.Exit(1)
			return
		}
		cmd := meta.ForCmd
		cmd.RemoveAdmin(args[0])
		if err := store.Selected().Put(cmd.User, cmd.Name, cmd); err != nil {
			fmt.Fprintln(sess.Stderr(), err.Error())
			sess.Exit(255)
			return
		}
		fmt.Fprintln(sess, "Admin revoked.")
	},
}
