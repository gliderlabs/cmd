package meta

import (
	"github.com/progrium/cmd/com/cmd"
)

func (c *Component) MetaCommands() []*cmd.MetaCommand {
	return []*cmd.MetaCommand{
		metaHelp,
		metaConfig,
		metaEnv,
		metaAccess,
		metaAdmins,
	}
}

func (c *Component) RootCommands() []*cmd.MetaCommand {
	return []*cmd.MetaCommand{
		rootHelp,
		rootInstall,
		rootUninstall,
		rootList,
		rootCreate,
		tokens,
	}
}
