package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var Version string

var noArgs = []string{"--noargs"}

type Command cobra.Command

type Flag struct {
	Name      string
	Value     interface{}
	Usage     string
	Shorthand string
	Kind      string
}

func addFlag(flags *pflag.FlagSet, flag Flag) {
	switch flag.Kind {
	// TODO: more non-string Kinds
	case "bool":
		flags.BoolP(flag.Name, flag.Shorthand, flag.Value.(bool), flag.Usage)
	default:
		// empty/other Kind assume string
		flags.StringP(flag.Name, flag.Shorthand, flag.Value.(string), flag.Usage)
	}
}

func (cmd Command) Init(parent *cobra.Command, flags ...Flag) *cobra.Command {
	cobraCmd := cobra.Command(cmd)
	cobraCmd.SetHelpFunc(helpFunc)
	cobraCmd.SetUsageFunc(usageFunc)
	cobraCmd.SetUsageTemplate(usageTemplate)
	cobraCmd.SetHelpTemplate(helpTemplate)
	cobraCmd.DisableFlagParsing = true
	if parent != nil {
		parent.AddCommand(&cobraCmd)
	}
	wrappedRun := cobraCmd.Run
	cobraCmd.Run = func(c *cobra.Command, args []string) {
		sess := ContextSession(c)
		if sess != nil {
			c.SetOutput(sess)
		}
		for _, flag := range flags {
			addFlag(c.Flags(), flag)
		}
		if len(args) == 1 && args[0] == noArgs[0] {
			args = []string{}
		}
		err := c.Flags().Parse(args)
		if err != nil {
			fmt.Fprintln(sess.Stderr(), err)
			sess.Exit(StatusUsageError)
			return
		}
		// only show help if NArgs < 2 (argcmd will have 1, parent will have 2)
		if help, _ := c.Flags().GetBool("help"); help && c.Flags().NArg() < 2 {
			c.Help()
			return
		}
		wrappedRun(c, args)
	}
	return &cobraCmd
}

type commandHandler func(*cobra.Command, []string)

var ArgCmd = func(handle commandHandler) commandHandler {
	return func(cmd *cobra.Command, args []string) {
		sess := ContextSession(cmd)
		if len(args) < 1 {
			cmd.Usage()
			sess.Exit(StatusUsageError)
			return
		}
		argCmd := Command{
			Use:    args[0],
			Hidden: true,
			Run: func(c *cobra.Command, a []string) {
				a = append([]string{args[0]}, a...)
				handle(c, a)
			},
		}.Init(nil)
		for _, sub := range cmd.Commands() {
			subCopy := *sub
			argCmd.AddCommand(&subCopy)
		}
		cmd.AddCommand(argCmd)
		newargs := []string{args[0]}
		if len(args) > 1 {
			newargs = append(newargs, []string{args[1], args[0]}...)
		}
		if len(args) > 2 {
			newargs = append(newargs, args[2:]...)
		}
		cmd.Root().SetArgs(append(strings.Split(cmd.CommandPath(), " ")[1:], newargs...))
		cmd.Root().Execute()
	}
}

func Execute(cmd cobra.Command, cmds []*cobra.Command, ctx context.Context, args []string) error {
	root := cmd
	key := defaultCtxRegistry.Add(ctx)
	root.Annotations = map[string]string{
		"_ctx": key,
	}
	defer defaultCtxRegistry.Clear(key)
	for _, sub := range cmds {
		subCopy := *sub
		root.AddCommand(&subCopy)
	}
	if len(args) == 0 {
		args = noArgs
	}
	root.SetArgs(args)
	return root.Execute()
}
