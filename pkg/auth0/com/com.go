package auth0

import (
	"net/http"

	"github.com/gliderlabs/comlab/pkg/com"
	"github.com/gliderlabs/cmd/pkg/auth0"
	"golang.org/x/oauth2"
)

func Register() {
	com.Register("auth0", &Component{},
		com.Option("client_id", "foo", "Auth0 client ID"),
		com.Option("client_secret", "", "Auth0 client secret"),
		com.Option("domain", "", "Auth0 domain"),
		com.Option("callback_url", "/_auth/callback", "Auth0 callback URL"),
		com.Option("logout_url", "/_auth/logout", "URL to wrap Auth0 logout"),
		com.Option("api_token", "", "Auth0 API bearer token"))
	auth0.DefaultClientFactory = ClientFactory
}

// Component ...
type Component struct{}

type AuthListener interface {
	WebAuthLogin(http.ResponseWriter, *http.Request, *oauth2.Token) error
	WebAuthLogout(http.ResponseWriter, *http.Request) error
}

func AuthListeners() []AuthListener {
	var listeners []AuthListener
	for _, com := range com.Enabled(new(AuthListener), nil) {
		listeners = append(listeners, com.(AuthListener))
	}
	return listeners
}
