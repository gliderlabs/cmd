package cmd

import (
	"github.com/gliderlabs/ssh"
	"github.com/progrium/cmd/com/core"
	"github.com/spf13/cobra"
)

type MetaCommand struct {
	Use     string
	Aliases []string
	Short   string
	Long    string
	Example string
	Hidden  bool
	Run     func(*MetaCommand, ssh.Session, []string)
	Setup   func(*MetaCommand)

	Cmd     *cobra.Command
	Session ssh.Session
	ForCmd  *core.Command
}

func (c *MetaCommand) Add(cmds ...*MetaCommand) *MetaCommand {
	var passCmd *MetaCommand
	for _, cmd := range cmds {
		cmdCopy := *cmd
		cmdCopy.Session = c.Session
		cmdCopy.ForCmd = c.ForCmd
		c.setup().Cmd.AddCommand((&cmdCopy).setup().Cmd)
		passCmd = &cmdCopy
		if passCmd.Setup != nil {
			passCmd.Setup(passCmd)
		}
	}
	return passCmd
}

func (c *MetaCommand) setup() *MetaCommand {
	if c.Cmd != nil {
		return c
	}
	cmd := &cobra.Command{
		Use:     c.Use,
		Aliases: c.Aliases,
		Short:   c.Short,
		Long:    c.Long,
		Example: c.Example,
		Hidden:  c.Hidden,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.SetOutput(c.Session)
			c.Run(c, c.Session, args)
		},
	}
	c.Cmd = cmd
	return c
}
