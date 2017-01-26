package stripe

import (
	"github.com/gliderlabs/comlab/pkg/com"
	"github.com/stripe/stripe-go"
)

func init() {
	com.Register("stripe", &Component{},
		com.Option("secret_key", "", "Stripe secret key"),
		com.Option("pub_key", "", "Stripe publishable key"),
		com.Option("event_endpoint", "/_stripe", "Path to handle Stripe webhooks"))
}

// Component ...
type Component struct{}

type EventListener interface {
	StripeEvents() []string
	StripeReceive(stripe.Event)
}

func EventListeners(eventType string) []EventListener {
	var listeners []EventListener
	for _, com := range com.Enabled(new(EventListener), nil) {
		eventTypes := com.(EventListener).StripeEvents()
		if len(eventTypes) == 0 {
			listeners = append(listeners, com.(EventListener))
			continue
		}
		for _, t := range eventTypes {
			if t == eventType {
				listeners = append(listeners, com.(EventListener))
				break
			}
		}
	}
	return listeners
}
