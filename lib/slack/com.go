package slack

import (
	"github.com/gliderlabs/comlab/pkg/com"
)

func init() {
	com.Register("slack", &Component{},
		com.Option("token", "", "Slack token"),
		com.Option("username", "gliderbot", "Username to post as"),
		com.Option("icon", "http://i.imgur.com/9P6bSVv.png", "Default icon to use"),
		com.Option("channel", "cmd", "Channel to post to"),
	)
}

// Component ...
type Component struct{}
