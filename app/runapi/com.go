package runapi

import (
	"github.com/gliderlabs/comlab/pkg/com"
)

func init() {
	com.Register("runapi", &Component{})
}

type Component struct{}
