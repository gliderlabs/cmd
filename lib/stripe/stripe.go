package stripe

import (
	"github.com/gliderlabs/comlab/pkg/com"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/client"
)

func init() {
	stripe.LogLevel = 0
}

func Client() *client.API {
	client := &client.API{}
	client.Init(com.GetString("secret_key"), nil)
	return client
}
