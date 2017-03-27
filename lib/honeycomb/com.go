package honeycomb

import (
	"github.com/gliderlabs/comlab/pkg/com"
)

func init() {
	com.Register("honeycomb", &Component{},
		com.Option("key", "", ""),
		com.Option("dataset", "", ""))
}

type Component struct{}
