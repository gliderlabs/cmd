package builtin

import (
	"fmt"
	"strings"

	"github.com/gliderlabs/cmd/app/store"
	"github.com/gliderlabs/cmd/lib/cli"
	"github.com/gliderlabs/cmd/lib/crypto"
	"github.com/spf13/cobra"
)

var envListFn = func(sess cli.Session, c *cobra.Command, args []string) error {
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
	if len(cmd.Environment) == 0 {
		fmt.Fprintln(sess, "No environment set for this command.")
		return nil
	}
	fields := map[string]interface{}{}
	for k, v := range cmd.Environment {
		fields[cli.Bright(k)] = crypto.Decrypt(v)
	}
	cli.PrintFields(sess, fields, true)
	return nil
}

var envCmd = func(sess cli.Session) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "env <cmd>",
		Short: "Manage command environment",
		Long:  `Without a subcommand, env will run "ls" by default.`,
		RunE: func(c *cobra.Command, args []string) error {
			return envListFn(sess, c, args)
		},
	}
	argCmd := cli.ArgumentCommand(cmd, sess)
	cli.AddCommand(argCmd, envListCmd, sess)
	cli.AddCommand(argCmd, envSetCmd, sess)
	cli.AddCommand(argCmd, envUnsetCmd, sess)
	return cmd
}

var envListCmd = func(sess cli.Session) *cobra.Command {
	return &cobra.Command{
		Use:   "ls",
		Short: "List command environment variables",
		RunE: func(c *cobra.Command, args []string) error {
			return envListFn(sess, c, args)
		},
	}
}

var envSetCmd = func(sess cli.Session) *cobra.Command {
	return &cobra.Command{
		Use:   "set <key=value>...",
		Short: "Set command environment variables",
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) < 2 {
				fmt.Fprintln(sess.Stderr(), "At least one key value pair is required")
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
			var keys []string
			for _, kvp := range args[1:] {
				parts := strings.SplitN(kvp, "=", 2)
				if len(parts) < 2 {
					continue
				}
				keys = append(keys, cli.Bright(parts[0]))
				box, _ := crypto.Encrypt(parts[1])
				cmd.SetEnv(parts[0], box)
			}
			cli.Status(sess, fmt.Sprintf(
				"Setting %s on %s", strings.Join(keys, ", "), cli.Bright(cmd.Name)))
			if err := store.Selected().Put(cmd.User, cmd.Name, cmd); err != nil {
				cli.StatusErr(sess.Stderr(), err.Error())
				sess.Exit(cli.StatusInternalError)
				return nil
			}
			cli.StatusDone(sess)
			return nil
		},
	}
}

var envUnsetCmd = func(sess cli.Session) *cobra.Command {
	return &cobra.Command{
		Use:   "unset <key>...",
		Short: "Unset command environment variables",
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) < 2 {
				fmt.Fprintln(sess.Stderr(), "At least one key is required")
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
			var keys []string
			for _, key := range args[1:] {
				delete(cmd.Environment, key)
				keys = append(keys, cli.Bright(key))
			}
			cli.Status(sess, fmt.Sprintf(
				"Unsetting %s on %s", strings.Join(keys, ", "), cli.Bright(cmd.Name)))
			if err := store.Selected().Put(cmd.User, cmd.Name, cmd); err != nil {
				cli.StatusErr(sess.Stderr(), err.Error())
				sess.Exit(cli.StatusInternalError)
				return nil
			}
			cli.StatusDone(sess)
			return nil
		},
	}
}
