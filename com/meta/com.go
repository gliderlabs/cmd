package meta

import (
	"github.com/gliderlabs/gosper/pkg/com"
)

func init() {
	com.Register("meta", &Component{})
}

type Component struct{}
