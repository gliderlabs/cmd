package console

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gliderlabs/comlab/pkg/log"
	"github.com/goware/emailx"
	"github.com/leekchan/accounting"
	"github.com/progrium/cmd/lib/mailgun"
	"github.com/progrium/cmd/lib/stripe"
	"github.com/progrium/cmd/lib/web"
	"github.com/progrium/cmd/pkg/auth0"
	stripelib "github.com/stripe/stripe-go"
)

type BillingInfo struct {
	Plan         string
	Email        string
	Period       string
	NextPayment  time.Time
	NextAmount   string
	CardType     string
	CardLastFour string
	CardExpMonth uint8
	CardExpYear  uint16
	History      []BillingInvoice
}

type BillingInvoice struct {
	Paid         bool
	Date         time.Time
	Amount       string
	Description  string
	CardType     string
	CardLastFour string
}

func GetBillingInfo(user *User) (BillingInfo, error) {
	var err error
	df := accounting.Accounting{Symbol: "$", Precision: 2}
	cparams := &stripelib.CustomerParams{}
	cparams.Expand("default_source")
	customer, err := stripe.Client().Customers.Get(user.Account.CustomerID, cparams)
	if err != nil {
		return BillingInfo{}, err
	}
	if user.Account.SubscriptionID == "" {
		return BillingInfo{
			Plan:  "Basic",
			Email: customer.Email,
		}, nil
	}
	sub, err := stripe.Client().Subs.Get(user.Account.SubscriptionID, nil)
	if err != nil {
		return BillingInfo{}, err
	}
	if sub.Status == "canceled" {
		return BillingInfo{
			Plan:  "Basic",
			Email: customer.Email,
		}, nil
	}
	iparams := &stripelib.InvoiceListParams{
		Customer: customer.ID,
	}
	iparams.Expand("data.charge")
	invoices := stripe.Client().Invoices.List(iparams)
	var history []BillingInvoice
	var lastCard stripelib.Card
	for invoices.Next() {
		invoice := invoices.Invoice()
		if !invoice.Attempted {
			continue
		}
		if invoice.Charge.Source.Type != stripelib.PaymentSourceCard {
			continue
		}
		history = append(history, BillingInvoice{
			Paid:         invoice.Paid,
			Date:         time.Unix(invoice.Date, 0),
			Amount:       df.FormatMoney(invoice.Amount / 100),
			Description:  invoice.Lines.Values[0].Plan.Meta["name"],
			CardType:     string(invoice.Charge.Source.Card.Brand),
			CardLastFour: invoice.Charge.Source.Card.LastFour,
		})
		if lastCard.LastFour == "" {
			lastCard = *invoice.Charge.Source.Card
		}
	}
	if invoices.Err() != nil {
		return BillingInfo{}, invoices.Err()
	}
	var card stripelib.Card
	if customer.DefaultSource != nil && customer.DefaultSource.Type == stripelib.PaymentSourceCard {
		card = *customer.DefaultSource.Card
	} else {
		card = lastCard
	}
	return BillingInfo{
		Email:        customer.Email,
		Plan:         sub.Plan.Meta["name"],
		Period:       string(sub.Plan.Interval),
		NextPayment:  time.Unix(sub.PeriodEnd, 0),
		NextAmount:   df.FormatMoney(sub.Plan.Amount / 100),
		CardType:     string(card.Brand),
		CardLastFour: card.LastFour,
		CardExpMonth: card.Month,
		CardExpYear:  card.Year,
		History:      history,
	}, nil
}

func updateEmailHandler(w http.ResponseWriter, r *http.Request, user *User) {
	err := emailx.Validate(emailx.Normalize(r.FormValue("email")))
	if err == emailx.ErrInvalidFormat {
		web.SessionSet(r, w, "error", "Invalid email format.")
		return
	}
	if err == emailx.ErrUnresolvableHost {
		web.SessionSet(r, w, "error", "Unresolvable email host.")
		return
	}
	params := &stripelib.CustomerParams{
		Email: emailx.Normalize(r.FormValue("email")),
	}
	_, err = stripe.Client().Customers.Update(user.Account.CustomerID, params)
	if err != nil {
		log.Info(r, err, log.Fields{"uid": user.ID, "svc": "stripe"})
		web.SessionSet(r, w, "error", err.Error())
		return
	}
	web.SessionSet(r, w, "success", "Your billing email has been updated.")
}

func updatePaymentHandler(w http.ResponseWriter, r *http.Request, user *User) {
	params := &stripelib.CustomerParams{}
	if err := params.SetSource(r.FormValue("update-token")); err != nil {
		log.Info(r, err, log.Fields{"uid": user.ID, "svc": "stripe"})
		web.SessionSet(r, w, "error", err.Error())
		return
	}
	_, err := stripe.Client().Customers.Update(user.Account.CustomerID, params)
	if err != nil {
		log.Info(r, err, log.Fields{"uid": user.ID, "svc": "stripe"})
		web.SessionSet(r, w, "error", err.Error())
		return
	}
	web.SessionSet(r, w, "success", "Your payment method has been updated.")
}

func unsubscribeHandler(w http.ResponseWriter, r *http.Request, user *User) {
	_, err := stripe.Client().Subs.Cancel(user.Account.SubscriptionID, nil)
	if err != nil {
		log.Info(r, err, log.Fields{"uid": user.ID, "svc": "stripe"})
		web.SessionSet(r, w, "error", err.Error())
		return
	}
	err = auth0.DefaultClient().PatchUser(user.ID, auth0.User{
		"app_metadata": map[string]interface{}{
			"subscription_id": "",
		},
	})
	if err != nil {
		log.Info(r, err, log.Fields{"uid": user.ID, "svc": "auth0"})
		web.SessionSet(r, w, "error", err.Error())
		return
	}
	web.SessionSet(r, w, "success", "Your subscription is canceled. You now have a Basic account.")
}

func subscribeHandler(w http.ResponseWriter, r *http.Request, user *User) {
	sub, err := stripe.Client().Subs.New(&stripelib.SubParams{
		Customer: user.Account.CustomerID,
		Plan:     "cmd-plus-" + r.FormValue("period"),
		Token:    r.FormValue("stripe-token"),
	})
	if err != nil {
		log.Info(r, err, log.Fields{"uid": user.ID, "svc": "stripe"})
		web.SessionSet(r, w, "error", err.Error())
		return
	}
	err = auth0.DefaultClient().PatchUser(user.ID, auth0.User{
		"app_metadata": map[string]interface{}{
			"subscription_id": sub.ID,
			"plan":            "plus",
		},
	})
	if err != nil {
		log.Info(r, err, log.Fields{"uid": user.ID, "svc": "auth0"})
		web.SessionSet(r, w, "error", err.Error())
		return
	}
	err = mailgun.SendText(
		user.Email,
		"Jeff Lindsay <jeff@gliderlabs.com>",
		"Cmd.io Plus Upgrade",
		fmt.Sprintf(`Hello, %s!

Thanks for upgrading and supporting our work. Be sure to join our
Slack community if you haven't already:

http://slack.gliderlabs.com

-jeff`, strings.Split(user.Name, " ")[0]))
	if err != nil {
		log.Info(r, err, log.Fields{"uid": user.ID, "svc": "mailgun"})
		web.SessionSet(r, w, "error", err.Error())
		return
	}
	web.SessionSet(r, w, "success", "Congrats! You now have a Plus account.")
}
