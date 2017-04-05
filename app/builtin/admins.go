package builtin

import (
	"fmt"
	"strings"

	"github.com/progrium/cmd/app/store"
	"github.com/progrium/cmd/lib/cli"
	"github.com/spf13/cobra"
)

var adminsListFn = func(c *cobra.Command, args []string) {
	sess := cli.ContextSession(c)
	if len(args) < 1 {
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
	if !cmd.IsAdmin(sess.User()) {
		fmt.Fprintln(sess.Stderr(), "Not allowed")
		sess.Exit(cli.StatusNoPerm)
		return
	}
	if len(cmd.ACL) == 0 && len(cmd.Admins) == 0 {
		fmt.Fprintln(sess, "Nobody else is admin for this command.")
		return
	}
	cli.Header(sess, "Admins")
	for _, user := range cmd.Admins {
		fmt.Fprintln(sess, user)
	}
}

var adminsCmd = cli.Command{
	Use:   "admins <cmd>",
	Short: "Manage command admins",
	Run:   cli.ArgCmd(adminsListFn),
}.Init(nil)

var adminsListCmd = cli.Command{
	Use:   "ls",
	Short: "List command admins",
	Run:   adminsListFn,
}.Init(adminsCmd)

var adminsGrantCmd = cli.Command{
	Use:     "grant <user>",
	Aliases: []string{"add"},
	Short:   "Grant command admin to a user",
	Run: func(c *cobra.Command, args []string) {
		sess := cli.ContextSession(c)
		if len(args) < 2 {
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
		if !cmd.IsAdmin(sess.User()) {
			fmt.Fprintln(sess.Stderr(), "Not allowed")
			sess.Exit(cli.StatusNoPerm)
			return
		}
		cli.Status(sess, fmt.Sprintf(
			"Granting %s admin to %s",
			cli.Bright(strings.Join(args[1:], ", ")), cli.Bright(cmd.Name)))
		if err := store.Selected().GrantAdmin(cmd.User, cmd.Name, args[1:]...); err != nil {
			cli.StatusErr(sess.Stderr(), err.Error())
			sess.Exit(cli.StatusInternalError)
			return
		}
		cli.StatusDone(sess)
	},
}.Init(adminsCmd)

var adminsRevokeCmd = cli.Command{
	Use:     "revoke <user>",
	Aliases: []string{"rm"},
	Short:   "Revoke command admin from a user",
	Run: func(c *cobra.Command, args []string) {
		sess := cli.ContextSession(c)
		if len(args) < 2 {
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
		if !cmd.IsAdmin(sess.User()) {
			fmt.Fprintln(sess.Stderr(), "Not allowed")
			sess.Exit(cli.StatusNoPerm)
			return
		}
		cli.Status(sess, fmt.Sprintf(
			"Revoking %s admin to %s",
			cli.Bright(strings.Join(args[1:], ", ")), cli.Bright(cmd.Name)))
		if err := store.Selected().RevokeAdmin(cmd.User, cmd.Name, args[1:]...); err != nil {
			cli.StatusErr(sess.Stderr(), err.Error())
			sess.Exit(cli.StatusInternalError)
			return
		}
		cli.StatusDone(sess)
	},
}.Init(adminsCmd)
