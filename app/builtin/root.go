package builtin

import (
	"github.com/progrium/cmd/app/core"
	"github.com/progrium/cmd/app/store"
	"github.com/progrium/cmd/lib/cli"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.SetHelpTemplate(`{{"\u2318"|gray}} {{"Cmd"|bright}} by Glider Labs
  {{"version:"|gray}} {{version|gray}}

{{.UsageString}}`)
	rootCmd.SetUsageTemplate(`{{"\u25ba"|gray}} {{"Usage"|bright}}{{if .Runnable}}
  ssh {{addr}} [ command | builtin ]{{end}}
{{if .HasExample}}
{{"\u2605"|gray}} {{"Examples"|bright}}
{{ .Example }}
{{end}}
{{"\u2630"|gray}} {{"Commands"|bright}}{{range .UserCommands}}
  {{.Name}}{{end}}
{{if .HasAvailableSubCommands}}
{{"\u2630"|gray}} {{"Builtins"|bright}}{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  :{{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}
{{end}}{{if .HasAvailableFlags}}
{{"\u2691"|gray}} {{"Flags"|bright}}
{{.Flags.FlagUsages | trimRightSpace}}
{{end}}{{if .HasHelpSubCommands}}
Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}
{{end}}{{if .HasAvailableSubCommands}}
Use "ssh {{addr}} [builtin] --help" to learn more about a builtin.{{end}}
`)
	rootCmd.SetUsageFunc(rootUsageFunc)
}

func rootUsageFunc(c *cobra.Command) error {
	sess := cli.ContextSession(c)
	if sess != nil {
		c.SetOutput(sess)
	}
	err := cli.Template(c.OutOrStderr(), c.UsageTemplate(), &rootCmdView{c})
	if err != nil {
		c.Println(err)
	}
	return err
}

var rootCmd = cli.Command{
	Use:    "cmd",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}.Init(nil)

type rootCmdView struct {
	*cobra.Command
}

func (cmd *rootCmdView) UserCommands() []*core.Command {
	sess := cli.ContextSession(cmd.Command)
	return store.Selected().List(sess.User())
}
