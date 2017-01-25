package core

import (
	"github.com/gliderlabs/comlab/pkg/com"
)

func init() {
	com.Register("core", &Component{},
		com.Option("host", "", ""))
}

type Component struct{}
