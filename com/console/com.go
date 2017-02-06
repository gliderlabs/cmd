package console

import "github.com/gliderlabs/comlab/pkg/com"

func init() {
	com.Register("console", &Component{},
		com.Option("slack_token", "", "Slack API token"))
}

type Component struct {
}
