package redirect

import (
	"github.com/gliderlabs/gosper/pkg/com"
)

func init() {
	com.Register("redirect", &Component{})
}

type Component struct{}
