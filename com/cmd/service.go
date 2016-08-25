package cmd

import (
	"net"

	"github.com/gliderlabs/pkg/com"
	"github.com/gliderlabs/pkg/log"
	"github.com/gliderlabs/pkg/ssh"
)

func (c *Component) Stop() {
	c.running = false
	if c.listener != nil {
		c.listener.Close()
	}
}

func (c *Component) Serve() {
	server := ssh.Server{}
	server.SetOption(ssh.PublicKeyAuth(HandleAuth))
	server.SetOption(ssh.HostKeyFile(com.GetString("hostkey_pem")))
	server.Handle(HandleSSH)

	c.running = true
	var err error
	c.listener, err = net.Listen("tcp", com.GetString("listen_addr"))
	if err != nil {
		panic(err)
	}
	log.Info("listening on", com.GetString("listen_addr"))
	server.Serve(c.listener)
}
