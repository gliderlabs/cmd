package analytics

import "github.com/gliderlabs/comlab/pkg/com"

func init() {
	com.Register("analytics", &Component{},
		com.Option("tracking_id", "", "Property tracking ID to enable Analytics"),
	)
}

// Component ...
type Component struct{}
