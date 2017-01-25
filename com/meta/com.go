package meta

import (
	"github.com/gliderlabs/comlab/pkg/com"
)

func init() {
	com.Register("meta", &Component{})
}

type Component struct{}
