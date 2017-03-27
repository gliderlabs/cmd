package console

import (
	"github.com/stripe/stripe-go"
)

func (c *Component) StripeEvents() []string {
	return []string{
		"invoice.payment_failed",
		//"customer.subscription.deleted",
	}
}

func (c *Component) StripeReceive(event stripe.Event) {
	switch event.Type {
	case "invoice.payment_failed":
		// TODO: email them!
		// if last attempt, cancel subscription
	}
}
