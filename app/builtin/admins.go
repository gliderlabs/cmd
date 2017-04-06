package builtin

import (
	"fmt"
	"strings"

	"github.com/progrium/cmd/app/store"
	"github.com/progrium/cmd/lib/cli"
	"github.com/spf13/cobra"
)

var adminsListFn = func(sess cli.Session, c *cobra.Command, args []string) error {
	if len(args) < 1 {
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
	if !cmd.IsAdmin(sess.User()) {
		fmt.Fprintln(sess.Stderr(), "Not allowed")
		sess.Exit(cli.StatusNoPerm)
		return nil
	}
	if len(cmd.ACL) == 0 && len(cmd.Admins) == 0 {
		fmt.Fprintln(sess, "Nobody else is admin for this command.")
		return nil
	}
	cli.Header(sess, "Admins")
	for _, user := range cmd.Admins {
		fmt.Fprintln(sess, user)
	}
	return nil
}

var adminsCmd = func(sess cli.Session) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "admins <cmd>",
		Short: "Manage command admins",
		RunE: func(c *cobra.Command, args []string) error {
			return adminsListFn(sess, c, args)
		},
	}
	argCmd := cli.ArgumentCommand(cmd, sess)
	cli.AddCommand(argCmd, adminsListCmd, sess)
	cli.AddCommand(argCmd, adminsGrantCmd, sess)
	cli.AddCommand(argCmd, adminsRevokeCmd, sess)
	return cmd
}

var adminsListCmd = func(sess cli.Session) *cobra.Command {
	return &cobra.Command{
		Use:   "ls",
		Short: "List command admins",
		RunE: func(c *cobra.Command, args []string) error {
			return adminsListFn(sess, c, args)
		},
	}
}

var adminsGrantCmd = func(sess cli.Session) *cobra.Command {
	return &cobra.Command{
		Use:     "grant <subject>",
		Aliases: []string{"add"},
		Short:   "Grant command admin to a user",
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) < 2 {
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
			if !cmd.IsAdmin(sess.User()) {
				fmt.Fprintln(sess.Stderr(), "Not allowed")
				sess.Exit(cli.StatusNoPerm)
				return nil
			}
			cli.Status(sess, fmt.Sprintf(
				"Granting %s admin to %s",
				cli.Bright(strings.Join(args[1:], ", ")), cli.Bright(cmd.Name)))
			if err := store.Selected().GrantAdmin(cmd.User, cmd.Name, args[1:]...); err != nil {
				cli.StatusErr(sess.Stderr(), err.Error())
				sess.Exit(cli.StatusInternalError)
				return nil
			}
			cli.StatusDone(sess)
			return nil
		},
	}
}

var adminsRevokeCmd = func(sess cli.Session) *cobra.Command {
	return &cobra.Command{
		Use:     "revoke <user>",
		Aliases: []string{"add"},
		Short:   "Revoke command admin from a user",
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) < 2 {
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
			if !cmd.IsAdmin(sess.User()) {
				fmt.Fprintln(sess.Stderr(), "Not allowed")
				sess.Exit(cli.StatusNoPerm)
				return nil
			}
			cli.Status(sess, fmt.Sprintf(
				"Revoking %s admin to %s",
				cli.Bright(strings.Join(args[1:], ", ")), cli.Bright(cmd.Name)))
			if err := store.Selected().RevokeAdmin(cmd.User, cmd.Name, args[1:]...); err != nil {
				cli.StatusErr(sess.Stderr(), err.Error())
				sess.Exit(cli.StatusInternalError)
				return nil
			}
			cli.StatusDone(sess)
			return nil
		},
	}
}
