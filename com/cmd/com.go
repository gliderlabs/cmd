package cmd

import (
	"net"

	"github.com/gliderlabs/pkg/com"
)

func init() {
	com.Register("cmd", &Component{},
		com.Option("docker_bin", "docker", "command path to use for docker"),
		com.Option("listen_addr", "127.0.0.1:2223", "port to bind on"),
		com.Option("config_dir", "local", "directory containing command configs"),
		com.Option("hostkey_pem", "com/cmd/data/id_host", "private key for host verification"))
}

type Component struct {
	running  bool
	listener net.Listener
}
