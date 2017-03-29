package slack

import (
	"github.com/gliderlabs/comlab/pkg/com"
	"github.com/nlopes/slack"
)

func Post(text string) error {
	api := slack.New(com.GetString("token"))
	params := slack.PostMessageParameters{
		Username: com.GetString("username"),
		IconURL:  com.GetString("icon"),
	}
	attachment := slack.Attachment{
		Text: text,
		/*Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "Feeling",
				Value: "Grateful",
			},
		},*/
	}
	params.Attachments = []slack.Attachment{attachment}
	_, _, err := api.PostMessage(com.GetString("channel"), "", params)
	return err
}
