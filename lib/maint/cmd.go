package maint

import (
	"fmt"

	"github.com/gliderlabs/ssh"
)

func (c *Component) PreprocessOrder() uint {
	return 0
}

func (c *Component) PreprocessSession(sess ssh.Session) (cont bool, msg string) {
	// restrict access when maintenance mode is active
	if Active() && !IsAllowed(sess.User()) {
		fmt.Fprintln(sess, Notice())
		return false, "maintenance"
	}
	return true, ""
}
