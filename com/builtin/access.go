package builtin

import (
	"fmt"

	"github.com/progrium/cmd/com/cli"
	"github.com/progrium/cmd/com/store"
	"github.com/spf13/cobra"
)

var accessListFn = func(c *cobra.Command, args []string) {
	sess := cli.ContextSession(c)
	if len(args) < 1 {
		c.Usage()
		sess.Exit(64)
		return
	}
	cmd, err := LookupCmd(sess.User(), args[0])
	if err != nil {
		fmt.Fprintln(sess.Stderr(), err.Error())
		sess.Exit(1)
		return
	}
	if !cmd.IsAdmin(sess.User()) {
		fmt.Fprintln(sess.Stderr(), "Not allowed")
		sess.Exit(77)
		return
	}
	if len(cmd.ACL) == 0 && len(cmd.Admins) == 0 {
		fmt.Fprintln(sess, "Nobody else has access to this command.")
		return
	}
	cli.Header(sess, "Users")
	for _, user := range cmd.Admins {
		fmt.Fprintln(sess, user+" (admin)")
	}
	for _, user := range cmd.ACL {
		fmt.Fprintln(sess, user)
	}
}

var accessCmd = cli.Command{
	Use:   "access <cmd>",
	Short: "Manage command access",
	Run:   cli.ArgCmd(accessListFn),
}.Init(nil)

var accessListCmd = cli.Command{
	Use:   "ls",
	Short: "List command access",
	Run:   accessListFn,
}.Init(accessCmd)

var accessGrantCmd = cli.Command{
	Use:     "grant <subject>",
	Aliases: []string{"add"},
	Short:   "Grant command access to a subject",
	Run: func(c *cobra.Command, args []string) {
		sess := cli.ContextSession(c)
		if len(args) < 2 {
			c.Usage()
			sess.Exit(64)
			return
		}
		cmd, err := LookupCmd(sess.User(), args[0])
		if err != nil {
			fmt.Fprintln(sess.Stderr(), err.Error())
			sess.Exit(1)
			return
		}
		if !cmd.IsAdmin(sess.User()) {
			fmt.Fprintln(sess.Stderr(), "Not allowed")
			sess.Exit(77)
			return
		}
		cli.Status(sess, fmt.Sprintf(
			"Granting %s access to %s", cli.Bright(args[1]), cli.Bright(cmd.Name)))
		cmd.AddAccess(args[1])
		if err := store.Selected().Put(cmd.User, cmd.Name, cmd); err != nil {
			cli.StatusErr(sess.Stderr(), err.Error())
			sess.Exit(70)
			return
		}
		cli.StatusDone(sess)
	},
}.Init(accessCmd)

var accessRevokeCmd = cli.Command{
	Use:     "revoke <subject>",
	Aliases: []string{"rm"},
	Short:   "Revoke command access from a subject",
	Run: func(c *cobra.Command, args []string) {
		sess := cli.ContextSession(c)
		if len(args) < 2 {
			c.Usage()
			sess.Exit(64)
			return
		}
		cmd, err := LookupCmd(sess.User(), args[0])
		if err != nil {
			fmt.Fprintln(sess.Stderr(), err.Error())
			sess.Exit(1)
			return
		}
		if !cmd.IsAdmin(sess.User()) {
			fmt.Fprintln(sess.Stderr(), "Not allowed")
			sess.Exit(77)
			return
		}
		cli.Status(sess, fmt.Sprintf(
			"Revoking %s access to %s", cli.Bright(args[1]), cli.Bright(cmd.Name)))
		cmd.RemoveAccess(args[1])
		if err := store.Selected().Put(cmd.User, cmd.Name, cmd); err != nil {
			cli.StatusErr(sess.Stderr(), err.Error())
			sess.Exit(70)
			return
		}
		cli.StatusDone(sess)
	},
}.Init(accessCmd)
