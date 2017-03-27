package ssh

import (
	"net"

	"github.com/gliderlabs/ssh"

	"github.com/gliderlabs/comlab/pkg/com"
)

func init() {
	com.Register("ssh", &Component{},
		com.Option("listen_addr", "127.0.0.1:2223", "port to bind on"),
		com.Option("hostkey_pem", "lib/ssh/data/dev_host", "private key for host verification"),
	)
}

type Component struct {
	running  bool
	listener net.Listener
}

type SessionHandler interface {
	HandleSSH(ssh.Session)
}

type AuthHandler interface {
	HandleAuth(ssh.Context, ssh.PublicKey) bool
}
