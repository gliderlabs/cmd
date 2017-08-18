package comlab

import (
	"log"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Short: "Comlab CLI",
}

func Run() {
	for _, provider := range CommandProviders() {
		provider.RegisterCommands(RootCmd)
	}
	if err := RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
