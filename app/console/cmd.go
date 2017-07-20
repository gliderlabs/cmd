package console

import (
	"fmt"

	"github.com/gliderlabs/cmd/lib/cli"
	"github.com/gliderlabs/cmd/lib/release"
	"github.com/gliderlabs/ssh"
)

func (c *Component) PreprocessOrder() uint {
	return 10
}

func (c *Component) PreprocessSession(sess ssh.Session) (cont bool, msg string) {
	// check for first time user
	if user := ContextUser(sess.Context()); user != nil {
		if user.Account.CustomerID == "" {
			fmt.Fprintf(sess, cli.Bright("\nWelcome, %s!\n\n"), sess.User())
			fmt.Fprintln(sess, "We noticed this is your first login. So far so good!")
			fmt.Fprintln(sess, "Would you mind logging in via the web interface?")
			fmt.Fprintln(sess, "This way we can properly set up your account:\n")
			fmt.Fprintf(sess, cli.Bright("https://%s/login\n\n"), release.Hostname())
			fmt.Fprintln(sess, "Then you can come back and use SSH as usual. Thanks!\n")
			return false, "first time login"
		}
	}
	return true, ""
}
