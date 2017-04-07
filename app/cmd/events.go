package cmd

import (
	"github.com/gliderlabs/cmd/app/console"
	"github.com/gliderlabs/cmd/lib/github"
	"github.com/gliderlabs/cmd/lib/sentry"
	"github.com/gliderlabs/cmd/lib/slack"
	"github.com/gliderlabs/comlab/pkg/events"
	"github.com/gliderlabs/comlab/pkg/log"

	slacklib "github.com/nlopes/slack"
)

func init() {
	events.Listen(&events.Listener{
		EventName: console.EventFirstLogin,
		Handler: func(event events.Event) {
			err := slack.Post("New first time login just now! 🎉")
			if err != nil {
				log.Info(err)
			}
		},
	})
	events.Listen(&events.Listener{
		EventName: console.EventNewSubscriber,
		Handler: func(event events.Event) {
			err := slack.Post("New Plus subscriber just now! 🎉")
			if err != nil {
				log.Info(err)
			}
		},
	})
	events.Listen(&events.Listener{
		EventName: sentry.EventNewIssue,
		Handler: func(event events.Event) {
			issue := event.(sentry.IssueEvent)
			var fields []slacklib.AttachmentField
			for _, f := range issue.Event.Tags {
				fields = append(fields, slacklib.AttachmentField{
					Title: f[0],
					Value: f[1],
				})
			}
			err := slack.PostAttachment(slacklib.Attachment{
				Color:     "danger",
				Title:     issue.Message,
				TitleLink: issue.URL,
				Fields:    fields,
				Footer:    issue.Level,
			}, nil)
			if err != nil {
				log.Info(err)
			}
		},
	})
	events.Listen(&events.Listener{
		EventName: github.EventStatus,
		Handler: func(event events.Event) {
			status := event.(github.StatusEvent)
			branches := []string{}
			foundBranch := ""
			for _, branch := range status.Branches {
				if len(branches) == 0 {
					foundBranch = String(branch.Name)
					break
				}
				for _, b := range branches {
					if String(branch.Name) == b {
						foundBranch = b
					}
				}
			}
			if foundBranch == "" {
				return
			}
			/*if String(status.State) == "pending" {
				return
			}*/
			colors := map[string]string{
				"success": "good",
				"failure": "danger",
				"error":   "danger",
				"pending": "",
			}
			err := slack.PostAttachment(slacklib.Attachment{
				Color:      colors[String(status.State)],
				AuthorName: foundBranch,
				Title:      String(status.State),
				TitleLink:  String(status.TargetURL),
				Footer:     "CircleCI",
			}, nil)
			if err != nil {
				log.Info(err)
			}
		},
	})
}

func String(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
