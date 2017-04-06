package cmd

import (
	"github.com/gliderlabs/comlab/pkg/events"
	"github.com/gliderlabs/comlab/pkg/log"
	"github.com/gliderlabs/cmd/app/console"
	"github.com/gliderlabs/cmd/lib/slack"
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
}
