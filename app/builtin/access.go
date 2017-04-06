package builtin

import (
	"fmt"
	"strings"

	"github.com/gliderlabs/cmd/app/store"
	"github.com/gliderlabs/cmd/lib/cli"
	"github.com/spf13/cobra"
)

var accessListFn = func(sess cli.Session, c *cobra.Command, args []string) error {
	if len(args) < 1 {
		c.Usage()
		sess.Exit(64)
		return nil
	}
	cmd, err := LookupCmd(sess.User(), args[0])
	if err != nil {
		fmt.Fprintln(sess.Stderr(), err.Error())
		sess.Exit(1)
		return nil
	}
	if !cmd.IsAdmin(sess.User()) {
		fmt.Fprintln(sess.Stderr(), "Not allowed")
		sess.Exit(77)
		return nil
	}
	if len(cmd.ACL) == 0 && len(cmd.Admins) == 0 {
		fmt.Fprintln(sess, "Nobody else has access to this command.")
		return nil
	}
	cli.Header(sess, "Users")
	for _, user := range cmd.Admins {
		fmt.Fprintln(sess, user+" (admin)")
	}
	for _, user := range cmd.ACL {
		fmt.Fprintln(sess, user)
	}
	return nil
}

var accessCmd = func(sess cli.Session) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "access <cmd>",
		Short: "Manage command access",
		RunE: func(c *cobra.Command, args []string) error {
			return accessListFn(sess, c, args)
		},
	}
	argCmd := cli.ArgumentCommand(cmd, sess)
	cli.AddCommand(argCmd, accessListCmd, sess)
	cli.AddCommand(argCmd, accessGrantCmd, sess)
	cli.AddCommand(argCmd, accessRevokeCmd, sess)
	return cmd
}

var accessListCmd = func(sess cli.Session) *cobra.Command {
	return &cobra.Command{
		Use:   "ls",
		Short: "List command access",
		RunE: func(c *cobra.Command, args []string) error {
			return accessListFn(sess, c, args)
		},
	}
}

var accessGrantCmd = func(sess cli.Session) *cobra.Command {
	return &cobra.Command{
		Use:     "grant <subject>",
		Aliases: []string{"add"},
		Short:   "Grant command access to a subject",
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) < 2 {
				c.Usage()
				sess.Exit(64)
				return nil
			}
			cmd, err := LookupCmd(sess.User(), args[0])
			if err != nil {
				fmt.Fprintln(sess.Stderr(), err.Error())
				sess.Exit(1)
				return nil
			}
			if !cmd.IsAdmin(sess.User()) {
				fmt.Fprintln(sess.Stderr(), "Not allowed")
				sess.Exit(77)
				return nil
			}
			cli.Status(sess, fmt.Sprintf(
				"Granting %s access to %s",
				cli.Bright(strings.Join(args[1:], ", ")), cli.Bright(cmd.Name)))
			if err := store.Selected().GrantAccess(cmd.User, cmd.Name, args[1:]...); err != nil {
				cli.StatusErr(sess.Stderr(), err.Error())
				sess.Exit(70)
				return nil
			}
			cli.StatusDone(sess)
			return nil
		},
	}
}

var accessRevokeCmd = func(sess cli.Session) *cobra.Command {
	return &cobra.Command{
		Use:     "revoke <subject>",
		Aliases: []string{"rm"},
		Short:   "Revoke command access from a subject",
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) < 2 {
				c.Usage()
				sess.Exit(64)
				return nil
			}
			cmd, err := LookupCmd(sess.User(), args[0])
			if err != nil {
				fmt.Fprintln(sess.Stderr(), err.Error())
				sess.Exit(1)
				return nil
			}
			if !cmd.IsAdmin(sess.User()) {
				fmt.Fprintln(sess.Stderr(), "Not allowed")
				sess.Exit(77)
				return nil
			}
			cli.Status(sess, fmt.Sprintf(
				"Revoking %s access to %s",
				cli.Bright(strings.Join(args[1:], ", ")), cli.Bright(cmd.Name)))
			if err := store.Selected().RevokeAccess(cmd.User, cmd.Name, args[1:]...); err != nil {
				cli.StatusErr(sess.Stderr(), err.Error())
				sess.Exit(70)
				return nil
			}
			cli.StatusDone(sess)
			return nil
		},
	}
}
