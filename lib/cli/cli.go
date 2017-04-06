package cli

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

const (
	noArgsFlag     = "--noargs"
	argCmdTmplName = "__argcmd"
)

type Command cobra.Command

type CommandFactory func(Session) *cobra.Command

type Flag struct {
	Name      string
	Value     interface{}
	Usage     string
	Shorthand string
	Kind      string
}

func AddFlag(cmd *cobra.Command, flag Flag) {
	flags := cmd.Flags()
	switch flag.Kind {
	// TODO: more non-string Kinds
	case "bool":
		flags.BoolP(flag.Name, flag.Shorthand, flag.Value.(bool), flag.Usage)
	default:
		// empty/other Kind assume string
		flags.StringP(flag.Name, flag.Shorthand, flag.Value.(string), flag.Usage)
	}
}

type Error struct {
	Err    error
	Status int
}

func (e Error) Error() string {
	return e.Err.Error()
}

func AddCommand(parent *cobra.Command, factory CommandFactory, sess Session) *cobra.Command {
	cmd := factory(sess)
	cmd.SetOutput(sess)
	cmd.SetHelpFunc(helpFunc)
	cmd.SetUsageFunc(usageFunc)
	cmd.SetUsageTemplate(usageTemplate)
	cmd.SetHelpTemplate(helpTemplate)
	cmd.DisableFlagParsing = true
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true
	cmd.PreRunE = func(c *cobra.Command, args []string) error {
		err := c.Flags().Parse(popArgFix(args))
		if err != nil {
			copyAnyArgCmdChildren(c)
			return Error{err, StatusUsageError}
		}
		// only show help if NArgs < 2 (argcmd will have 1, parent will have 2)
		if help, _ := c.Flags().GetBool("help"); help && c.Flags().NArg() < 2 {
			copyAnyArgCmdChildren(c)
			c.Help()
			return Error{errors.New(""), -1} // cancel run
		}
		return nil
	}
	if cmd.Annotations[argCmdTmplName] != "" {
		cmd.RunE = argCmdWrapper(cmd.RunE, sess)
	}
	parent.AddCommand(cmd)
	return cmd
}

func ArgumentCommand(cmd *cobra.Command, sess Session) *cobra.Command {
	if cmd.Annotations == nil {
		cmd.Annotations = map[string]string{}
	}
	cmd.Annotations[argCmdTmplName] = argCmdTmplName
	return addArgCmd(cmd, sess, argCmdTmplName)
}

func fixArgCmdArgs(cmd *cobra.Command, args []string) []string {
	newargs := []string{args[0]}
	if len(args) > 1 {
		newargs = append(newargs, []string{args[1], args[0]}...)
	}
	if len(args) > 2 {
		newargs = append(newargs, args[2:]...)
	}
	return append(strings.Split(cmd.CommandPath(), " ")[1:], newargs...)
}

func argCmdWrapper(handle func(*cobra.Command, []string) error, sess Session) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			copyAnyArgCmdChildren(cmd)
			return Error{fmt.Errorf("Usage error"), StatusUsageError}
		}
		if len(args) < 2 {
			return handle(cmd, args)
		}
		argCmdTmpl := cmd.Commands()[0]
		argCmd := addArgCmd(cmd, sess, args[0])
		copyChildren(argCmd, argCmdTmpl)
		cmd.Root().SetArgs(fixArgCmdArgs(cmd, args))
		return argCmd.Execute()
	}
}

func addArgCmd(cmd *cobra.Command, sess Session, name string) *cobra.Command {
	argCmdFactory := func(sess Session) *cobra.Command {
		argCmd := &cobra.Command{
			Use:    name,
			Hidden: true,
		}
		return argCmd
	}
	argCmd := AddCommand(cmd, argCmdFactory, sess)
	return argCmd
}

func copyChildren(to, from *cobra.Command) {
	for _, sub := range from.Commands() {
		subCopy := *sub
		to.AddCommand(&subCopy)
	}
}

func copyAnyArgCmdChildren(cmd *cobra.Command) {
	if cmd.Annotations[argCmdTmplName] != "" {
		copyChildren(cmd, cmd.Commands()[0]) // from argCmdTmpl
	}
}

// if cobra gets empty args, it loads from os.Args
// and we never want this.
func pushArgFix(args []string) []string {
	if len(args) == 0 {
		return []string{noArgsFlag}
	}
	return args
}

func popArgFix(args []string) []string {
	if len(args) == 1 && args[0] == noArgsFlag {
		return []string{}
	}
	return args
}

func Execute(cmd *cobra.Command, sess Session, args []string) error {
	root := *cmd // copy
	root.SilenceUsage = true
	root.SilenceErrors = true
	root.SetArgs(pushArgFix(args))
	cmd, err := root.ExecuteC()
	if cliErr, ok := err.(Error); ok {
		if cliErr.Status == -1 {
			return nil
		}
		if cliErr.Status == StatusUsageError {
			cmd.Help()
			return nil
		}
	}
	return err
}
