package cmd

import (
	"github.com/gliderlabs/comlab/pkg/com"
)

func init() {
	com.Register("cmd", &Component{},
		com.Option("access_denied_msg",
			"Access Denied: Visit https://alpha.cmd.io/request to request access",
			"message shown when user isn't allowed access"),
	)
}

type Component struct{}
