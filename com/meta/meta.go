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
		cmd := meta.ForCmd
		for k, v := range cmd.Config {
			fmt.Fprintf(sess, "%s=%s\n", k, v)
		}
	},
	Setup: func(meta *cmd.MetaCommand) {
		meta.Add(metaConfigSet, metaConfigUnset)
	},
}

var metaConfigSet = &cmd.MetaCommand{
	Use:   "set <key=value>...",
	Short: "Manage command configuration",
	Run: func(meta *cmd.MetaCommand, sess ssh.Session, args []string) {
		cmd := meta.ForCmd
		for _, kvp := range args {
			parts := strings.SplitN(kvp, "=", 2)
			if len(parts) < 2 {
				continue
			}
			cmd.Config[parts[0]] = parts[1]
		}
		if err := store.Selected().Put(cmd.User, cmd.Name, cmd); err != nil {
			fmt.Fprintln(sess.Stderr(), err.Error())
			sess.Exit(255)
			return
		}
		fmt.Fprintln(sess, "Config updated.")
	},
}

var metaConfigUnset = &cmd.MetaCommand{
	Use:   "unset <key>",
	Short: "Manage command configuration",
	Run: func(meta *cmd.MetaCommand, sess ssh.Session, args []string) {
		cmd := meta.ForCmd
		for _, key := range args {
			delete(cmd.Config, key)
		}
		if err := store.Selected().Put(cmd.User, cmd.Name, cmd); err != nil {
			fmt.Fprintln(sess.Stderr(), err.Error())
			sess.Exit(255)
			return
		}
		fmt.Fprintln(sess, "Config updated.")
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
		if cmd.IsPublic() {
			fmt.Fprintln(sess, "This command is accessible by everyone.")
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
		if meta.ForCmd.IsPublic() {
			meta.Add(metaAccessPrivate)
		} else {
			meta.Add(metaAccessPublic, metaAccessAdd, metaAccessRemove)
		}
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

var metaAccessPublic = &cmd.MetaCommand{
	Use:   "public",
	Short: "Make command public to all",
	Run: func(meta *cmd.MetaCommand, sess ssh.Session, args []string) {
		cmd := meta.ForCmd
		cmd.MakePublic()
		if err := store.Selected().Put(cmd.User, cmd.Name, cmd); err != nil {
			fmt.Fprintln(sess.Stderr(), err.Error())
			sess.Exit(255)
			return
		}
		fmt.Fprintln(sess, "Command is now public.")
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
