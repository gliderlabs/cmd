package maint

import (
	"github.com/gliderlabs/comlab/pkg/com"
)

func init() {
	com.Register("maint", &Component{},
		com.Option("active", false, ""),
		com.Option("notice", "cmd.io is currently down for maintenance", "displayed when maintenance active"),
		com.Option("allow", "", "comma separated list of users to allow during maintenance"))
}

// Component ...
type Component struct {
}
