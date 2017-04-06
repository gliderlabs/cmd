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

var tokensCmd = func(sess cli.Session) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tokens",
		Short: "Manage access tokens",
		RunE: func(c *cobra.Command, args []string) error {
			c.Help()
			return nil
		},
	}
	cli.AddCommand(cmd, tokensListCmd, sess)
	cli.AddCommand(cmd, tokensNew, sess)
	cli.AddCommand(cmd, tokensDelete, sess)
	return cmd
}

var tokensListCmd = func(sess cli.Session) *cobra.Command {
	return &cobra.Command{
		Use:   "ls",
		Short: "List tokens",
		RunE: func(c *cobra.Command, args []string) error {
			cli.Header(sess, "Tokens")
			tokens, err := store.Selected().ListTokens(sess.User())
			if err != nil {
				fmt.Fprintf(sess.Stderr(), err.Error())
				sess.Exit(cli.StatusInternalError)
				return nil
			}
			for _, token := range tokens {
				fmt.Fprintf(sess, "  %-10s  %s %s\n", token.Key, token.Description, token.LastUsedOn)
			}
			fmt.Fprintln(sess, "")
			return nil
		},
	}
}

var tokensNew = func(sess cli.Session) *cobra.Command {
	return &cobra.Command{
		Use:   "new <description>",
		Short: "Create a token",
		RunE: func(c *cobra.Command, args []string) error {
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
				return nil
			}
			fmt.Fprintln(sess, token.Key)
			return nil
		},
	}
}

var tokensDelete = func(sess cli.Session) *cobra.Command {
	return &cobra.Command{
		Use:   "rm <key>",
		Short: "Delete a token",
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) < 1 {
				fmt.Fprintf(sess.Stderr(), "Key name is required")
				sess.Exit(cli.StatusUsageError)
				return nil
			}
			token, _ := store.Selected().GetToken(args[0])
			if token == nil || token.User != sess.User() {
				fmt.Fprintf(sess.Stderr(), "Token not found")
				sess.Exit(cli.StatusError)
				return nil
			}

			if err := store.Selected().DeleteToken(token.Key); err != nil {
				log.Info(sess, token, err)
				fmt.Fprintf(sess.Stderr(), err.Error())
				sess.Exit(cli.StatusInternalError)
				return nil
			}
			fmt.Fprintln(sess, "Token removed")
			return nil
		},
	}
}
