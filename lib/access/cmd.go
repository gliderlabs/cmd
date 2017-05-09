package access

import (
	"fmt"

	"github.com/gliderlabs/comlab/pkg/com"
	"github.com/gliderlabs/ssh"
	uuid "github.com/satori/go.uuid"
)

func (c *Component) PreprocessOrder() uint {
	return 20
}

func (c *Component) PreprocessSession(sess ssh.Session) (cont bool, msg string) {
	// check for channel access when user is not a token
	if token := uuid.FromStringOrNil(sess.User()); token == uuid.Nil && !Check(sess.User()) {
		fmt.Fprintln(sess, com.GetString("deny_msg"))
		return false, "channel access denied"
	}
	return true, ""
}
