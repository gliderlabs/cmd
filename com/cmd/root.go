package cmd

import (
	"fmt"

	"github.com/gliderlabs/pkg/ssh"
	"github.com/spf13/cobra"
)

const rootUsageTmpl = `Usage:{{if .Runnable}}
  ssh <user>@{{.UseLine}}{{ if .HasAvailableSubCommands}} [command]{{end}}{{end}}{{if gt .Aliases 0}}

Aliases:
  {{.NameAndAliases}}
{{end}}{{if .HasExample}}

Examples:
{{ .Example }}{{end}}{{ if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if .IsAvailableCommand}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{ if .HasAvailableLocalFlags}}

Additional help topics:{{range .Commands}}{{if .IsHelpCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{ if .HasAvailableSubCommands }}

Use "[command] --help" for help about a meta command.{{end}}

`

var rootHelp = &MetaCommand{
	Use:   ":help",
	Short: "Print this help",
  Hidden: true,
	Run: func(cmd *MetaCommand, sess ssh.Session, args []string) {
    for _, c := range Store.List(sess.User()) {
      cmd.Cmd.Parent().AddCommand(&cobra.Command{
        Use: c.Name,
        Short: c.Description,
        Run: func(cmd *cobra.Command, args []string) {},
      })
    }
		cmd.Cmd.Parent().Help()
	},
}

var rootList = &MetaCommand{
	Use:   ":ls",
	Short: "List installed commands",
	Run: func(meta *MetaCommand, sess ssh.Session, args []string) {
		fmt.Fprintln(sess, "")
    fmt.Fprintln(sess, "Installed Commands:")
    for _, cmd := range Store.List(sess.User()) {
      fmt.Fprintf(sess, "  %-10s  %s\n", cmd.Name, cmd.Description)
    }
    fmt.Fprintln(sess, "")
	},
}

var rootInstall = &MetaCommand{
	Use:   ":add <name> <source>",
	Short: "Install a command",
	Run: func(meta *MetaCommand, sess ssh.Session, args []string) {
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
		cmd := &Command{
      Name: args[0],
      User: sess.User(),
      Source: args[1],
      Config: make(map[string]string),
    }
    if err := cmd.Pull(); err != nil {
      fmt.Fprintln(sess.Stderr(), "Command unable to install")
			sess.Exit(1)
			return
  	}
    if err := Store.Put(cmd.User, cmd.Name, cmd); err != nil {
			fmt.Fprintln(sess.Stderr(), err.Error())
			sess.Exit(255)
			return
		}
    fmt.Fprintln(sess, "Command installed")
	},
}

var rootUninstall = &MetaCommand{
	Use:   ":rm <name>",
	Short: "Uninstall a command",
	Run: func(meta *MetaCommand, sess ssh.Session, args []string) {
    if len(args) < 1 {
      fmt.Fprintln(sess, "Must specify a command")
      sess.Exit(1)
      return
    }
		cmd := Store.Get(sess.User(), args[0])
    if cmd == nil {
      fmt.Fprintln(sess, "Command not found")
      sess.Exit(1)
      return
    }
    if err := Store.Delete(cmd.User, cmd.Name); err != nil {
      fmt.Fprintln(sess.Stderr(), err.Error())
      sess.Exit(255)
      return
    }
    fmt.Fprintln(sess, "Command uninstalled")
	},
}
