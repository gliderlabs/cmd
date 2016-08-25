package cmd

import (
	"fmt"
  "strings"

	"github.com/gliderlabs/pkg/ssh"
)

const metaUsageTmpl = `Usage:{{if .Runnable}}{{if not .HasParent }}
  ssh <user>@cmd.io {{.UseLine}}{{ if .HasAvailableSubCommands}}:[command]{{end}}{{else}}
  {{.UseLine}}{{ if .HasAvailableSubCommands}} [command]{{end}}{{end}}{{end}}{{ if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if .IsAvailableCommand}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{ if .HasAvailableLocalFlags}}

Additional help topics:{{range .Commands}}{{if .IsHelpCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{ if .HasAvailableSubCommands }}

Use "[command] --help" for help about a meta command.{{end}}

`

var metaHelp = &MetaCommand{
	Use:    ":help",
	Short:  "Print this help",
	Hidden: true,
	Run: func(meta *MetaCommand, sess ssh.Session, args []string) {
		meta.Cmd.Parent().Help()
	},
}

var metaConfig = &MetaCommand{
	Use:   ":config",
	Short: "Manage command configuration",
	Run: func(meta *MetaCommand, sess ssh.Session, args []string) {
    cmd := meta.ForCmd
    for k, v := range cmd.Config {
      fmt.Fprintf(sess, "%s=%s\n", k, v)
    }
	},
}

var metaConfigSet = &MetaCommand{
	Use:   "set <key=value>...",
	Short: "Manage command configuration",
	Run: func(meta *MetaCommand, sess ssh.Session, args []string) {
    cmd := meta.ForCmd
    for _, kvp := range args {
      parts := strings.SplitN(kvp, "=", 2)
      if len(parts) < 2 {
        continue
      }
      cmd.Config[parts[0]] = parts[1]
    }
		if err := Store.Put(cmd.User, cmd.Name, cmd); err != nil {
      fmt.Fprintln(sess.Stderr(), err.Error())
  		sess.Exit(255)
      return
    }
    fmt.Fprintln(sess, "Config updated.")
	},
}

var metaConfigUnset = &MetaCommand{
	Use:   "unset <key>",
	Short: "Manage command configuration",
	Run: func(meta *MetaCommand, sess ssh.Session, args []string) {
    cmd := meta.ForCmd
    for _, key := range args {
      delete(cmd.Config, key)
    }
		if err := Store.Put(cmd.User, cmd.Name, cmd); err != nil {
      fmt.Fprintln(sess.Stderr(), err.Error())
  		sess.Exit(255)
      return
    }
    fmt.Fprintln(sess, "Config updated.")
	},
}

var metaAccess = &MetaCommand{
	Use:   ":access",
	Short: "Manage command access",
	Run: func(meta *MetaCommand, sess ssh.Session, args []string) {
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
}

var metaAccessAdd = &MetaCommand{
  Use:  "add <user>",
  Short: "Give a GitHub user access",
  Run: func(meta *MetaCommand, sess ssh.Session, args []string) {
    if len(args) < 1 {
      fmt.Fprintln(sess, "Must specify a GitHub user")
      sess.Exit(1)
      return
    }
    cmd := meta.ForCmd
    cmd.AddAccess(args[0])
    if err := Store.Put(cmd.User, cmd.Name, cmd); err != nil {
      fmt.Fprintln(sess.Stderr(), err.Error())
  		sess.Exit(255)
      return
    }
    fmt.Fprintln(sess, "Access granted.")
  },
}

var metaAccessRemove = &MetaCommand{
  Use:  "rm <user>",
  Short: "Take access for GitHub user",
  Run: func(meta *MetaCommand, sess ssh.Session, args []string) {
    if len(args) < 1 {
      fmt.Fprintln(sess, "Must specify a GitHub user")
      sess.Exit(1)
      return
    }
    cmd := meta.ForCmd
    cmd.RemoveAccess(args[0])
    if err := Store.Put(cmd.User, cmd.Name, cmd); err != nil {
      fmt.Fprintln(sess.Stderr(), err.Error())
  		sess.Exit(255)
      return
    }
    fmt.Fprintln(sess, "Access revoked.")
  },
}

var metaAccessPublic = &MetaCommand{
  Use:  "public",
  Short: "Make command public to all",
  Run: func(meta *MetaCommand, sess ssh.Session, args []string) {
    cmd := meta.ForCmd
    cmd.MakePublic()
    if err := Store.Put(cmd.User, cmd.Name, cmd); err != nil {
      fmt.Fprintln(sess.Stderr(), err.Error())
  		sess.Exit(255)
      return
    }
    fmt.Fprintln(sess, "Command is now public.")
  },
}

var metaAccessPrivate = &MetaCommand{
  Use:  "private",
  Short: "Make command private again",
  Run: func(meta *MetaCommand, sess ssh.Session, args []string) {
    cmd := meta.ForCmd
    cmd.MakePrivate()
    if err := Store.Put(cmd.User, cmd.Name, cmd); err != nil {
      fmt.Fprintln(sess.Stderr(), err.Error())
  		sess.Exit(255)
      return
    }
    fmt.Fprintln(sess, "Command is now private.")
  },
}

var metaAdmins = &MetaCommand{
	Use:   ":admins",
	Short: "Manage command admins",
	Run: func(meta *MetaCommand, sess ssh.Session, args []string) {
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
}

var metaAdminsAdd = &MetaCommand{
  Use:  "add <user>",
  Short: "Make GitHub user an admin",
  Run: func(meta *MetaCommand, sess ssh.Session, args []string) {
    if len(args) < 1 {
      fmt.Fprintln(sess, "Must specify a GitHub user")
      sess.Exit(1)
      return
    }
    cmd := meta.ForCmd
    cmd.AddAdmin(args[0])
    if err := Store.Put(cmd.User, cmd.Name, cmd); err != nil {
      fmt.Fprintln(sess.Stderr(), err.Error())
  		sess.Exit(255)
      return
    }
    fmt.Fprintln(sess, "Admin granted.")
  },
}

var metaAdminsRemove = &MetaCommand{
  Use:  "rm <user>",
  Short: "Revoke admin for GitHub user",
  Run: func(meta *MetaCommand, sess ssh.Session, args []string) {
    if len(args) < 1 {
      fmt.Fprintln(sess, "Must specify a GitHub user")
      sess.Exit(1)
      return
    }
    cmd := meta.ForCmd
    cmd.RemoveAdmin(args[0])
    if err := Store.Put(cmd.User, cmd.Name, cmd); err != nil {
      fmt.Fprintln(sess.Stderr(), err.Error())
  		sess.Exit(255)
      return
    }
    fmt.Fprintln(sess, "Admin revoked.")
  },
}
