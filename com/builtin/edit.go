package builtin

import (
	"fmt"
	"io/ioutil"

	"github.com/gliderlabs/comlab/pkg/log"
	"github.com/spf13/cobra"

	"github.com/progrium/cmd/com/cli"
	"github.com/progrium/cmd/com/store"
)

var editCmd = cli.Command{
	Use:   "edit <name> [-]",
	Short: "Edit a command",
	Long: `Edit source for an existing command.

Source will be read from stdin when single "-" provided as last arg.`,
	Example: `  # Edit command with name "cmd" reading source from stdin
  echo -e '#!cmd alpine\n echo "hello world"' | ssh cmd.io :edit cmd -`,
	Run: func(c *cobra.Command, args []string) {
		sess := cli.ContextSession(c)
		if len(args) < 2 {
			fmt.Fprintln(sess.Stderr(), "Unsupported: use - to read from stdin")
			c.Usage()
			sess.Exit(cli.StatusUsageError)
			return
		}
		cmd, err := LookupCmd(sess.User(), args[0])
		if err != nil {
			fmt.Fprintln(sess.Stderr(), err.Error())
			sess.Exit(cli.StatusError)
			return
		}
		cli.Status(sess, "Editing command")
		source, err := ioutil.ReadAll(sess)
		if err != nil {
			cli.StatusErr(sess.Stderr(), err.Error())
			sess.Exit(cli.StatusDataError)
			return
		}
		cmd.Source = string(source)
		if err := cmd.Build(); err != nil {
			log.Info(err)
			cli.StatusErr(sess.Stderr(), err.Error())
			sess.Exit(cli.StatusInternalError)
			return
		}
		if err := store.Selected().Put(cmd.User, cmd.Name, cmd); err != nil {
			log.Info(sess, cmd, err)
			cli.StatusErr(sess.Stderr(), err.Error())
			sess.Exit(cli.StatusInternalError)
			return
		}
		cli.StatusDone(sess)
	},
}.Init(nil)
