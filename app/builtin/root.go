package builtin

import (
	"github.com/progrium/cmd/app/core"
	"github.com/progrium/cmd/app/store"
	"github.com/progrium/cmd/lib/cli"
	"github.com/spf13/cobra"
)

const (
	helpTemplate = `{{"\u2318"|gray}} {{"Cmd"|bright}} by Glider Labs
  {{"version:"|gray}} {{version|gray}}

{{.UsageString}}`
	usageTemplate = `{{"\u25ba"|gray}} {{"Usage"|bright}}{{if .Runnable}}
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
`
)

func rootUsageFunc(c *cobra.Command) error {
	err := cli.Template(c.OutOrStderr(), c.UsageTemplate(), &rootCmdView{c})
	if err != nil {
		c.Println(err)
	}
	return err
}

var rootCmd = func(sess cli.Session) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "cmd",
		Hidden: true,
		RunE: func(c *cobra.Command, args []string) error {
			c.Help()
			return nil
		},
	}
	cmd.Annotations = map[string]string{
		"user": sess.User(),
	}
	return cmd
}

type rootCmdView struct {
	*cobra.Command
}

func (cmd *rootCmdView) UserCommands() []*core.Command {
	return store.Selected().List(cmd.Annotations["user"])
}
