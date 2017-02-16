package meta

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gliderlabs/ssh"

	"github.com/progrium/cmd/com/cmd"
	"github.com/progrium/cmd/com/core"
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
		fmt.Fprintln(sess.Stderr(),
			"WARN: meta command", meta.Use, "has been deprecated and replaced with", metaEnv.Use)
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
		fields := map[string]interface{}{}
		for k, v := range cmd.Environment {
			fields[bright(k)] = v
		}
		printFields(sess, fields, true)
	},
	Setup: func(meta *cmd.MetaCommand) {
		meta.Add(metaEnvSet, metaEnvUnset)
	},
}

var metaEnvSet = &cmd.MetaCommand{
	Use:   "set <key=value>...",
	Short: "Manage command environment",
	Run: func(meta *cmd.MetaCommand, sess ssh.Session, args []string) {
		if len(args) == 0 {
			statusErr(sess.Stderr(), "No keys provided")
			meta.Cmd.Usage()
			sess.Exit(1)
			return
		}
		cmd := meta.ForCmd
		var keys []string
		for _, kvp := range args {
			parts := strings.SplitN(kvp, "=", 2)
			if len(parts) < 2 {
				continue
			}
			keys = append(keys, bright(parts[0]))
			cmd.SetEnv(parts[0], parts[1])
		}
		statusMsg(sess, fmt.Sprintf(
			"Setting %s on %s", strings.Join(keys, ", "), bright(cmd.Name)))
		if err := store.Selected().Put(cmd.User, cmd.Name, cmd); err != nil {
			statusErr(sess.Stderr(), err.Error())
			sess.Exit(255)
			return
		}
		statusDone(sess)
	},
}

var metaEnvUnset = &cmd.MetaCommand{
	Use:   "unset <key>...",
	Short: "Manage command environment",
	Run: func(meta *cmd.MetaCommand, sess ssh.Session, args []string) {
		if len(args) == 0 {
			statusErr(sess.Stderr(), "No keys provided")
			meta.Cmd.Usage()
			sess.Exit(1)
			return
		}
		cmd := meta.ForCmd
		var keys []string
		for _, key := range args {
			delete(cmd.Environment, key)
			keys = append(keys, bright(key))
		}
		statusMsg(sess, fmt.Sprintf(
			"Unsetting %s on %s", strings.Join(keys, ", "), bright(cmd.Name)))
		if err := store.Selected().Put(cmd.User, cmd.Name, cmd); err != nil {
			statusErr(sess.Stderr(), err.Error())
			sess.Exit(255)
			return
		}
		statusDone(sess)
	},
}

var metaAccess = &cmd.MetaCommand{
	Use:   ":access",
	Short: "Manage command access",
	Run: func(meta *cmd.MetaCommand, sess ssh.Session, args []string) {
		cmd := meta.ForCmd
		if len(cmd.ACL) == 0 && len(cmd.Admins) == 0 {
			fmt.Fprintln(sess, "Nobody else has access to this command.")
			return
		}
		header(sess, "Users")
		for _, user := range cmd.Admins {
			fmt.Fprintln(sess, user+" (admin)")
		}
		for _, user := range cmd.ACL {
			fmt.Fprintln(sess, user)
		}
	},
	Setup: func(meta *cmd.MetaCommand) {
		meta.Add(metaAccessPrivate, metaAccessAdd, metaAccessRemove)
	},
}

var metaAccessAdd = &cmd.MetaCommand{
	Use:     "grant <subject>",
	Aliases: []string{"add"},
	Short:   "Grant access to a command",
	Run: func(meta *cmd.MetaCommand, sess ssh.Session, args []string) {
		if len(args) < 1 {
			statusErr(sess.Stderr(), "Must specify a subject")
			sess.Exit(1)
			return
		}
		cmd := meta.ForCmd
		statusMsg(sess, fmt.Sprintf(
			"Granting %s access to %s", bright(args[0]), bright(cmd.Name)))
		cmd.AddAccess(args[0])
		if err := store.Selected().Put(cmd.User, cmd.Name, cmd); err != nil {
			statusErr(sess.Stderr(), err.Error())
			sess.Exit(255)
			return
		}
		statusDone(sess)
	},
}

var metaAccessRemove = &cmd.MetaCommand{
	Use:     "revoke <subject>",
	Aliases: []string{"rm"},
	Short:   "Revoke access to a command",
	Run: func(meta *cmd.MetaCommand, sess ssh.Session, args []string) {
		if len(args) < 1 {
			statusErr(sess.Stderr(), "Must specify a subject")
			sess.Exit(1)
			return
		}
		cmd := meta.ForCmd
		statusMsg(sess, fmt.Sprintf(
			"Revoking %s access to %s", bright(args[0]), bright(cmd.Name)))
		cmd.RemoveAccess(args[0])
		if err := store.Selected().Put(cmd.User, cmd.Name, cmd); err != nil {
			statusErr(sess.Stderr(), err.Error())
			sess.Exit(255)
			return
		}
		statusDone(sess)
	},
}

var metaAccessPrivate = &cmd.MetaCommand{
	Use:   "private",
	Short: "Make command private again",
	Run: func(meta *cmd.MetaCommand, sess ssh.Session, args []string) {
		cmd := meta.ForCmd
		statusMsg(sess, "Revoking external access to "+bright(cmd.Name))
		cmd.MakePrivate()
		if err := store.Selected().Put(cmd.User, cmd.Name, cmd); err != nil {
			statusErr(sess.Stderr(), err.Error())
			sess.Exit(255)
			return
		}
		statusDone(sess)
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
		header(sess, "Admins")
		for _, user := range cmd.Admins {
			fmt.Fprintln(sess, user)
		}
	},
	Setup: func(meta *cmd.MetaCommand) {
		meta.Add(metaAdminsAdd, metaAdminsRemove)
	},
}

var metaAdminsAdd = &cmd.MetaCommand{
	Use:     "grant <user>",
	Aliases: []string{"add"},
	Short:   "Make GitHub user an admin",
	Run: func(meta *cmd.MetaCommand, sess ssh.Session, args []string) {
		if len(args) < 1 {
			statusErr(sess.Stderr(), "Must specify a user")
			sess.Exit(1)
			return
		}
		cmd := meta.ForCmd
		statusMsg(sess, fmt.Sprintf(
			"Granting %s admin to %s", bright(args[0]), bright(cmd.Name)))
		cmd.AddAdmin(args[0])
		if err := store.Selected().Put(cmd.User, cmd.Name, cmd); err != nil {
			statusErr(sess.Stderr(), err.Error())
			sess.Exit(255)
			return
		}
		statusDone(sess)
	},
}

var metaAdminsRemove = &cmd.MetaCommand{
	Use:     "revoke <user>",
	Aliases: []string{"rm"},
	Short:   "Revoke admin for GitHub user",
	Run: func(meta *cmd.MetaCommand, sess ssh.Session, args []string) {
		if len(args) < 1 {
			statusErr(sess.Stderr(), "Must specify a user")
			sess.Exit(1)
			return
		}
		cmd := meta.ForCmd
		statusMsg(sess, fmt.Sprintf(
			"Revoking %s admin to %s", bright(args[0]), bright(cmd.Name)))
		cmd.RemoveAdmin(args[0])
		if err := store.Selected().Put(cmd.User, cmd.Name, cmd); err != nil {
			statusErr(sess.Stderr(), err.Error())
			sess.Exit(255)
			return
		}
		statusDone(sess)
	},
}

const metaUsageTmpl = `Usage:{{if .Runnable}}
  ssh <user>@cmd.io {{.UseLine}}{{ if .HasAvailableSubCommands}} [command]{{end}}
{{end}}{{ if .HasAvailableSubCommands}}

Available Sub Commands:{{range .Commands}}{{if .IsAvailableCommand}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{ if .HasAvailableLocalFlags}}

Additional help topics:{{range .Commands}}{{if .IsHelpCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{ if .HasAvailableSubCommands }}

Use "[command] --help" for help about a meta command.{{end}}

`

func wrapMeta(c *cmd.MetaCommand) *cmd.MetaCommand {
	use := c.Use + " <name>"
	return &cmd.MetaCommand{
		Use:   use,
		Short: c.Short,
		Run: func(meta *cmd.MetaCommand, sess ssh.Session, args []string) {
			cmd, err := getCmd(sess.User(), args)
			if err != nil {
				statusErr(sess.Stderr(), err.Error())
				newCmd := *c
				newCmd.Session = sess
				newCmd.Use = use
				newCmd.Setup(&newCmd)
				newCmd.Cmd.SetOutput(sess)
				newCmd.Cmd.SetUsageTemplate(metaUsageTmpl)
				newCmd.Cmd.Usage()
				sess.Exit(1)
				return
			}
			newCmd := *c
			newCmd.Session = sess
			newCmd.Use = use
			newCmd.ForCmd = cmd
			newCmd.Setup(&newCmd)
			newCmd.Cmd.SetUsageTemplate(metaUsageTmpl)
			newCmd.Cmd.SetArgs(args[1:])
			newCmd.Cmd.SetOutput(sess)
			newCmd.Cmd.Execute()
		},
	}
}

func getCmd(user string, args []string) (*core.Command, error) {
	if len(args) < 1 {
		return nil, errors.New("Command name missing")
	}

	var cmd *core.Command
	name := args[0]
	if strings.Contains(name, "/") {
		parts := strings.SplitN(name, "/", 2)
		cmd = store.Selected().Get(parts[0], parts[1])
	} else {
		cmd = store.Selected().Get(user, name)
	}
	if cmd == nil {
		return nil, errors.New("Command not found: " + name)
	}
	if !cmd.IsAdmin(user) {
		return nil, errors.New("Not allowed")
	}

	return cmd, nil
}
