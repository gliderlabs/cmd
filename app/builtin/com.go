package builtin

import (
	"github.com/gliderlabs/comlab/pkg/com"
	"github.com/spf13/cobra"
)

func init() {
	com.Register("builtin", &Component{})
}

type Component struct{}

func (c *Component) BuiltinCommands() []*cobra.Command {
	return []*cobra.Command{
		envCmd,
		listCmd,
		createCmd,
		accessCmd,
		adminsCmd,
		importCmd,
		deleteCmd,
		editCmd,
		tokensCmd,
	}
}

func Commands() []*cobra.Command {
	var cmds []*cobra.Command
	for _, com := range com.Enabled(new(Provider), nil) {
		cmds = append(cmds, com.(Provider).BuiltinCommands()...)
	}
	return cmds
}

type Provider interface {
	BuiltinCommands() []*cobra.Command
}
