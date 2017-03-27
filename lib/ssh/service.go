package ssh

import (
	"net"

	"github.com/gliderlabs/comlab/pkg/com"
	"github.com/gliderlabs/comlab/pkg/log"
	"github.com/gliderlabs/ssh"
)

func (c *Component) Stop() {
	c.running = false
	if c.listener != nil {
		c.listener.Close()
	}
}

func (c *Component) Serve() {
	server := ssh.Server{}
	server.SetOption(ssh.HostKeyFile(com.GetString("hostkey_pem")))
	server.SetOption(ssh.PublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
		for _, com := range com.Enabled(new(AuthHandler), nil) {
			// just use the first one for now
			return com.(AuthHandler).HandleAuth(ctx, key)
		}
		return false
	}))
	server.Handle(func(sess ssh.Session) {
		for _, com := range com.Enabled(new(SessionHandler), nil) {
			com.(SessionHandler).HandleSSH(sess)
		}
	})

	c.running = true
	var err error
	c.listener, err = net.Listen("tcp", com.GetString("listen_addr"))
	if err != nil {
		panic(err)
	}
	log.Info("listening on", com.GetString("listen_addr"))
	server.Serve(c.listener)
}
