package builtin

import (
	"fmt"
	"io/ioutil"

	"github.com/gliderlabs/comlab/pkg/log"
	"github.com/spf13/cobra"

	"github.com/progrium/cmd/app/store"
	"github.com/progrium/cmd/lib/cli"
)

var editCmd = func(sess cli.Session) *cobra.Command {
	return &cobra.Command{
		Use:   "edit <name> [-]",
		Short: "Edit a command",
		Long: `Edit source for an existing command.

	Source will be read from stdin when single "-" provided as last arg.`,
		Example: `  # Edit command with name "cmd" reading source from stdin
	  echo -e '#!cmd alpine\n echo "hello world"' | ssh cmd.io :edit cmd -`,
		RunE: func(c *cobra.Command, args []string) error {

			if len(args) < 2 {
				fmt.Fprintln(sess.Stderr(), "Unsupported: use - to read from stdin")
				c.Usage()
				sess.Exit(cli.StatusUsageError)
				return nil
			}
			cmd, err := LookupCmd(sess.User(), args[0])
			if err != nil {
				fmt.Fprintln(sess.Stderr(), err.Error())
				sess.Exit(cli.StatusError)
				return nil
			}
			cli.Status(sess, "Editing command")
			source, err := ioutil.ReadAll(sess)
			if err != nil {
				cli.StatusErr(sess.Stderr(), err.Error())
				sess.Exit(cli.StatusDataError)
				return nil
			}
			cmd.Source = string(source)
			if err := cmd.Build(); err != nil {
				log.Info(err)
				cli.StatusErr(sess.Stderr(), err.Error())
				sess.Exit(cli.StatusInternalError)
				return nil
			}
			if err := store.Selected().Put(cmd.User, cmd.Name, cmd); err != nil {
				log.Info(sess, cmd, err)
				cli.StatusErr(sess.Stderr(), err.Error())
				sess.Exit(cli.StatusInternalError)
				return nil
			}
			cli.StatusDone(sess)
			return nil
		},
	}
}
