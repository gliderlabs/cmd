package cli

import (
	"github.com/spf13/cobra"
)

const helpTemplate = `{{.Short|trim|bright}}

{{if .Long}}{{.Long|trim}}

{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`

const usageTemplate = `{{"\u25ba"|gray}} {{"Usage"|bright}}{{if .Runnable}}
  ssh {{addr}} {{.UseLine|replace "cmd " ":"}}{{end}}{{if .HasAvailableSubCommands}} [subcommand]{{end}}
{{if gt .Aliases 0}}
Aliases:
  {{.NameAndAliases}}
{{end}}{{if .HasExample}}
{{"\u2605"|gray}} {{"Examples"|bright}}
{{ .Example }}
{{end}}{{if .HasAvailableSubCommands}}
{{"\u2630"|gray}} {{"Subcommands"|bright}}{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}
{{end}}{{if .HasAvailableFlags}}
{{"\u2691"|gray}} {{"Flags"|bright}}
{{.Flags.FlagUsages | trimRightSpace}}
{{end}}{{if .HasHelpSubCommands}}
Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}
{{end}}{{if .HasAvailableSubCommands}}
Use "{{.UseLine|replace "cmd " ":"}} [subcommand] --help" to learn more about a subcommand.{{end}}
`

func usageFunc(c *cobra.Command) error {
	//c.mergePersistentFlags()
	//c.ParseFlags()
	sess := ContextSession(c)
	if sess != nil {
		c.SetOutput(sess)
	}
	err := Template(c.OutOrStderr(), c.UsageTemplate(), c)
	if err != nil {
		c.Println(err)
	}
	return err
}

func helpFunc(cmd *cobra.Command, args []string) {
	sess := ContextSession(cmd)
	if sess != nil {
		cmd.SetOutput(sess)
	}
	if len(args) < 1 {
		err := Template(cmd.OutOrStdout(), cmd.HelpTemplate(), cmd)
		if err != nil {
			cmd.Println(err)
		}
		return
	}
	argCmd := &cobra.Command{
		Use:    args[1],
		Hidden: true,
		Run: func(cmd *cobra.Command, _ []string) {
			cmd.Help()
		},
	}
	for _, c := range cmd.Commands() {
		argCmd.AddCommand(c)
	}
	cmd.AddCommand(argCmd)
	leafCmd, _, err := cmd.Root().Find(args)
	if err != nil {
		panic(err)
	}
	leafCmd.Help()
}
