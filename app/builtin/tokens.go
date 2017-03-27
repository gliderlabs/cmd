package builtin

import (
	"fmt"

	"github.com/gliderlabs/comlab/pkg/log"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/cobra"

	"github.com/progrium/cmd/app/core"
	"github.com/progrium/cmd/app/store"
	"github.com/progrium/cmd/lib/cli"
)

var tokensCmd = cli.Command{
	Use:   "tokens",
	Short: "Manage access tokens",
	Run: func(c *cobra.Command, args []string) {
		c.Help()
	},
}.Init(nil)

var tokensListCmd = cli.Command{
	Use:   "ls",
	Short: "List tokens",
	Run: func(c *cobra.Command, args []string) {
		sess := cli.ContextSession(c)
		cli.Header(sess, "Tokens")
		tokens, err := store.Selected().ListTokens(sess.User())
		if err != nil {
			fmt.Fprintf(sess.Stderr(), err.Error())
			sess.Exit(cli.StatusInternalError)
			return
		}
		for _, token := range tokens {
			fmt.Fprintf(sess, "  %-10s  %s %s\n", token.Key, token.Description, token.LastUsedOn)
		}
		fmt.Fprintln(sess, "")
	},
}.Init(tokensCmd)

var tokensNew = cli.Command{
	Use:   "new <description>",
	Short: "Create a token",
	Run: func(c *cobra.Command, args []string) {
		sess := cli.ContextSession(c)
		var desc string
		if len(args) > 1 {
			desc = args[0]
		}

		token := &core.Token{
			Key:         uuid.NewV4().String(),
			User:        sess.User(),
			Description: desc,
		}

		if err := store.Selected().PutToken(token); err != nil {
			log.Info(sess, token, err)
			fmt.Fprintf(sess.Stderr(), err.Error())
			sess.Exit(cli.StatusInternalError)
			return
		}
		fmt.Fprintln(sess, token.Key)
	},
}.Init(tokensCmd)

var tokensDelete = cli.Command{
	Use:   "rm <key>",
	Short: "Delete a token",
	Run: func(c *cobra.Command, args []string) {
		sess := cli.ContextSession(c)
		if len(args) < 1 {
			fmt.Fprintf(sess.Stderr(), "Key name is required")
			sess.Exit(cli.StatusUsageError)
			return
		}
		token, _ := store.Selected().GetToken(args[0])
		if token == nil || token.User != sess.User() {
			fmt.Fprintf(sess.Stderr(), "Token not found")
			sess.Exit(cli.StatusError)
			return
		}

		if err := store.Selected().DeleteToken(token.Key); err != nil {
			log.Info(sess, token, err)
			fmt.Fprintf(sess.Stderr(), err.Error())
			sess.Exit(cli.StatusInternalError)
			return
		}
		fmt.Fprintln(sess, "Token removed")
	},
}.Init(tokensCmd)
