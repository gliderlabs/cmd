package meta

import (
	"fmt"

	"github.com/gliderlabs/gosper/pkg/log"
	"github.com/gliderlabs/ssh"
	uuid "github.com/satori/go.uuid"

	"github.com/progrium/cmd/com/cmd"
	"github.com/progrium/cmd/com/core"
	"github.com/progrium/cmd/com/store"
)

var tokens = &cmd.MetaCommand{
	Use: ":cmd-tokens",
	Run: func(meta *cmd.MetaCommand, sess ssh.Session, args []string) {
		meta.Cmd.Help()
	},
	Setup: func(meta *cmd.MetaCommand) {
		meta.Add(tokensList, tokensNew, tokensDelete)
	},
}

var tokensList = &cmd.MetaCommand{
	Use:   "ls",
	Short: "List tokens",
	Run: func(meta *cmd.MetaCommand, sess ssh.Session, args []string) {
		fmt.Fprintln(sess, "")
		fmt.Fprintln(sess, "Tokens:")
		tokens, err := store.Selected().ListTokens(sess.User())
		if err != nil {
			fmt.Fprintln(sess, err)
		}
		for _, token := range tokens {
			fmt.Fprintf(sess, "  %-10s  %s %s\n", token.Key, token.Description, token.LastUsedOn)
		}
		fmt.Fprintln(sess, "")
	},
}

var tokensNew = &cmd.MetaCommand{
	Use:   "new <description>",
	Short: "Create a token",
	Run: func(meta *cmd.MetaCommand, sess ssh.Session, args []string) {
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
			fmt.Fprintln(sess.Stderr(), err.Error())
			sess.Exit(255)
			return
		}
		fmt.Fprintln(sess, "Token created:", token.Key)
	},
}

var tokensDelete = &cmd.MetaCommand{
	Use:   "rm <key>",
	Short: "Delete a token",
	Run: func(meta *cmd.MetaCommand, sess ssh.Session, args []string) {
		if len(args) < 1 {
			fmt.Fprintln(sess, "Must specify a key")
			sess.Exit(1)
			return
		}
		token, _ := store.Selected().GetToken(args[0])
		if token == nil || token.User != sess.User() {
			fmt.Fprintln(sess, "Token not found")
			sess.Exit(1)
			return
		}

		if err := store.Selected().DeleteToken(token.Key); err != nil {
			log.Info(sess, token, err)
			fmt.Fprintln(sess.Stderr(), err.Error())
			sess.Exit(255)
			return
		}
		fmt.Fprintln(sess, "Token removed")
	},
}
