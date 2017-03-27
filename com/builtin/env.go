package builtin

import (
	"fmt"
	"strings"

	"github.com/progrium/cmd/com/cli"
	"github.com/progrium/cmd/com/store"
	"github.com/spf13/cobra"
)

var envListFn = func(c *cobra.Command, args []string) {
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
	if len(cmd.Environment) == 0 {
		fmt.Fprintln(sess, "No environment set for this command.")
		return
	}
	fields := map[string]interface{}{}
	for k, v := range cmd.Environment {
		fields[cli.Bright(k)] = v
	}
	cli.PrintFields(sess, fields, true)
}

var envCmd = cli.Command{
	Use:   "env <cmd>",
	Short: "Manage command environment",
	Long:  `Without a subcommand, env will run "ls" by default.`,
	Run:   cli.ArgCmd(envListFn),
}.Init(nil)

var envListCmd = cli.Command{
	Use:   "ls",
	Short: "List command environment variables",
	Run:   envListFn,
}.Init(envCmd)

var envSetCmd = cli.Command{
	Use:   "set <key=value>...",
	Short: "Set command environment variables",
	Run: func(c *cobra.Command, args []string) {
		sess := cli.ContextSession(c)
		if len(args) < 2 {
			fmt.Fprintln(sess.Stderr(), "At least one key value pair is required")
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
		var keys []string
		for _, kvp := range args[1:] {
			parts := strings.SplitN(kvp, "=", 2)
			if len(parts) < 2 {
				continue
			}
			keys = append(keys, cli.Bright(parts[0]))
			cmd.SetEnv(parts[0], parts[1])
		}
		cli.Status(sess, fmt.Sprintf(
			"Setting %s on %s", strings.Join(keys, ", "), cli.Bright(cmd.Name)))
		if err := store.Selected().Put(cmd.User, cmd.Name, cmd); err != nil {
			cli.StatusErr(sess.Stderr(), err.Error())
			sess.Exit(cli.StatusInternalError)
			return
		}
		cli.StatusDone(sess)
	},
}.Init(envCmd)

var envUnsetCmd = cli.Command{
	Use:   "unset <key>...",
	Short: "Unset command environment variables",
	Run: func(c *cobra.Command, args []string) {
		sess := cli.ContextSession(c)
		if len(args) < 2 {
			fmt.Fprintln(sess.Stderr(), "At least one key is required")
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
			return
		}
		cli.StatusDone(sess)
	},
}.Init(envCmd)
