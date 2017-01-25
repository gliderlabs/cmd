package stripe

import "github.com/gliderlabs/comlab/pkg/com"

func init() {
	com.Register("stripe", &Component{},
		com.Option("secret_key", "", "Stripe secret key"),
		com.Option("pub_key", "", "Stripe publishable key"))
}

// Component ...
type Component struct{}
