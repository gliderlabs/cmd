package github

import "github.com/gliderlabs/comlab/pkg/com"

func init() {
	com.Register("github", &Component{},
		com.Option("endpoint", "/_github", "webhook endpoint"),
	)
}

// Component ...
type Component struct{}
