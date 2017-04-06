package builtin

import (
	"github.com/gliderlabs/comlab/pkg/com"
	"github.com/gliderlabs/cmd/lib/cli"
)

func init() {
	com.Register("builtin", &Component{})
}

type Component struct{}

func (c *Component) BuiltinCommands() []cli.CommandFactory {
	return []cli.CommandFactory{
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

func Commands() []cli.CommandFactory {
	var cmds []cli.CommandFactory
	for _, com := range com.Enabled(new(Provider), nil) {
		cmds = append(cmds, com.(Provider).BuiltinCommands()...)
	}
	return cmds
}

type Provider interface {
	BuiltinCommands() []cli.CommandFactory
}
