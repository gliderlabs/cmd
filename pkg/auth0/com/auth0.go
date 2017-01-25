package auth0

import (
	"github.com/gliderlabs/comlab/pkg/com"
	"github.com/progrium/cmd/pkg/auth0"
)

func ClientFactory() *auth0.Client {
	return &auth0.Client{
		ClientID:     com.GetString("client_id"),
		ClientSecret: com.GetString("client_secret"),
		Domain:       com.GetString("domain"),
		CallbackURL:  com.GetString("callback_url"),
		Token:        com.GetString("api_token"),
	}
}
