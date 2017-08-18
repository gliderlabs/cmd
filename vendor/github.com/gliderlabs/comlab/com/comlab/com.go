package comlab

import (
	"github.com/gliderlabs/comlab/pkg/com"
	"github.com/spf13/cobra"
)

func init() {
	com.Register("comlab", &struct{}{})
}

type CommandProvider interface {
	RegisterCommands(root *cobra.Command)
}

func CommandProviders() []CommandProvider {
	var providers []CommandProvider
	for _, com := range com.Enabled(new(CommandProvider), nil) {
		providers = append(providers, com.(CommandProvider))
	}
	return providers
}
