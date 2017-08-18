package dev

import "github.com/spf13/cobra"

func (c *Component) RegisterCommands(root *cobra.Command) {
	root.AddCommand(&cobra.Command{
		Use:   "dev",
		Short: "Dev runner",
		Run: func(cmd *cobra.Command, args []string) {
			Run()
		},
	})
}
