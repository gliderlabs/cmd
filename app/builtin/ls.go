package builtin

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gliderlabs/cmd/app/store"
	"github.com/gliderlabs/cmd/lib/cli"
)

var listCmd = func(sess cli.Session) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls",
		Short: "List available commands",
		RunE: func(c *cobra.Command, args []string) error {
			if ok, _ := c.Flags().GetBool("json"); ok {
				var names []string
				for _, cmd := range store.Selected().List(sess.User()) {
					names = append(names, cmd.Name)
				}
				cli.JSON(sess, names)
				return nil
			}
			cli.Header(sess, "Your Commands")
			for _, cmd := range store.Selected().List(sess.User()) {
				fmt.Fprintf(sess, "  %-10s  %s\n", cmd.Name, cmd.Description)
			}
			fmt.Fprintln(sess, "")
			return nil
		},
	}
	cli.AddFlag(cmd, cli.Flag{"json", false, "output in JSON", "j", "bool"})
	return cmd
}
