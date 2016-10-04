package core

import (
	"github.com/gliderlabs/gosper/pkg/com"
)

func init() {
	com.Register("core", &Component{},
		com.Option("docker_bin", "docker", "command path to use for docker"))
}

type Component struct{}
