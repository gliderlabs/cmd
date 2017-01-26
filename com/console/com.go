package console

import "github.com/gliderlabs/comlab/pkg/com"

func init() {
	com.Register("console", &Component{})
}

type Component struct {
}
