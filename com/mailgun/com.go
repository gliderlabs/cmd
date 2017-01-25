package mailgun

import "github.com/gliderlabs/comlab/pkg/com"

func init() {
	com.Register("mailgun", &Component{},
		com.Option("domain", "", "Mailgun domain"),
		com.Option("api_key", "", "Mailgun API key"),
		com.Option("public_api_key", "", "Mailgun public API key"))
}

// Component ...
type Component struct{}
