package cmd

import (
	"net"

	"github.com/gliderlabs/pkg/com"
)

func init() {
	com.Register("cmd", &Component{},
		com.Option("docker_bin", "docker", "command path to use for docker"),
		com.Option("listen_addr", "127.0.0.1:2223", "port to bind on"),
		com.Option("hostkey_pem", "com/cmd/data/id_host", "private key for host verification"),
		com.Option("access_denied_msg",
			"Access Denied: Visit https://goo.gl/forms/CEtAiBUoxCWidAbK2 to request access",
			"message shown when user isn't allowed access"),
		com.Option("config_dir", "local", "directory containing command configs"),
		com.Option("table_name", "", "dynamodb table name for command storage"),
		com.Option("aws_access_key", "", "aws access key for dynamodb store"),
		com.Option("aws_secret_key", "", "aws secret key for dynamodb store"),
		com.Option("aws_region", "us-east-1", "aws region for dynamodb store"),
		com.Option("sentry_dsn", "", ""),
		com.Option("honeycomb_key", "", ""),
		com.Option("honeycomb_dataset", "", ""),
	)
}

type Component struct {
	running  bool
	listener net.Listener
}
