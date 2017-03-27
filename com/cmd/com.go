package cmd

import (
	"net"

	"github.com/gliderlabs/comlab/pkg/com"
)

func init() {
	com.Register("cmd", &Component{},
		com.Option("listen_addr", "127.0.0.1:2223", "port to bind on"),
		com.Option("hostkey_pem", "com/cmd/data/dev_host", "private key for host verification"),
		com.Option("access_denied_msg",
			"Access Denied: Visit https://alpha.cmd.io/request to request access",
			"message shown when user isn't allowed access"),
		com.Option("honeycomb_key", "", ""),
		com.Option("honeycomb_dataset", "", ""),
		com.Option("gh_team_id", "2144066", "GitHub team ID to allow access to"),
		com.Option("gh_token", "", "GitHub access token"),
	)
}

type Component struct {
	running  bool
	listener net.Listener
}
