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
	}
	params.Attachments = []slack.Attachment{attachment}
	_, _, err := api.PostMessage(com.GetString("channel"), "", params)
	return err
}

func PostAttachment(attachment slack.Attachment, params *slack.PostMessageParameters) error {
	api := slack.New(com.GetString("token"))
	if params == nil {
		params = &slack.PostMessageParameters{}
	}
	if params.Username == "" {
		params.Username = com.GetString("username")
	}
	if params.IconURL == "" {
		params.IconURL = com.GetString("icon")
	}
	params.Attachments = []slack.Attachment{attachment}
	_, _, err := api.PostMessage(com.GetString("channel"), "", *params)
	return err
}
