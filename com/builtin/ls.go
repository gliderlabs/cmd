package builtin

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/progrium/cmd/com/cli"
	"github.com/progrium/cmd/com/store"
)

var listCmd = cli.Command{
	Use:   "ls",
	Short: "List available commands",
	Run: func(c *cobra.Command, args []string) {
		sess := cli.ContextSession(c)
		if ok, _ := c.Flags().GetBool("json"); ok {
			var names []string
			for _, cmd := range store.Selected().List(sess.User()) {
				names = append(names, cmd.Name)
			}
			cli.JSON(sess, names)
			return
		}
		cli.Header(sess, "Your Commands")
		for _, cmd := range store.Selected().List(sess.User()) {
			fmt.Fprintf(sess, "  %-10s  %s\n", cmd.Name, cmd.Description)
		}
		fmt.Fprintln(sess, "")
	},
}.Init(nil,
	cli.Flag{"json", false, "output in JSON", "j", "bool"},
)
