package auth0

import (
	"fmt"
	"net/http"
	"net/url"
	"text/template"

	"github.com/gliderlabs/comlab/pkg/com"
	"github.com/gliderlabs/cmd/pkg/auth0"
)

func (c *Component) WebTemplateFuncMap(r *http.Request) template.FuncMap {
	return template.FuncMap{
		"auth0": func() string {
			return fmt.Sprintf(`
				var auth0;
				(function() {
					var js = document.createElement("script");
					js.type = "text/javascript";
					js.src = "https://cdn.auth0.com/w2/auth0-7.4.min.js";
					js.onload = function() {
						auth0 = new Auth0({
					    domain:       '%s',
					    clientID:     '%s',
					    callbackURL:  '%s',
							responseType: 'code'
					  });
					};
					document.body.appendChild(js);
				})();
			`,
				com.GetString("domain"),
				com.GetString("client_id"),
				com.GetString("callback_url"))
		},
	}
}

func (c *Component) MatchHTTP(r *http.Request) bool {
	if cb, err := url.Parse(com.GetString("callback_url")); err == nil {
		if r.URL.Path == cb.Path {
			return true
		}
	}
	if logout, err := url.Parse(com.GetString("logout_url")); err == nil {
		if r.URL.Path == logout.Path {
			return true
		}
	}
	return false
}

// ServeHTTP of web.Handler extension point
func (c *Component) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	listeners := AuthListeners()

	if logout, err := url.Parse(com.GetString("logout_url")); err == nil {
		if r.URL.Path == logout.Path {
			if returnTo := r.URL.Query().Get("return"); returnTo != "" {
				http.Redirect(w, r, returnTo, http.StatusFound)
				return
			}
			for _, listener := range listeners {
				if err := listener.WebAuthLogout(w, r); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
			q := url.Values{}
			q.Set("return", r.Referer())
			returnURL := &url.URL{
				Scheme:   logout.Scheme,
				Host:     logout.Host,
				Path:     logout.Path,
				RawQuery: q.Encode(),
			}
			http.Redirect(w, r, auth0.DefaultClient().LogoutURL(returnURL.String()), http.StatusFound)
			return
		}
	}

	token, err := auth0.DefaultClient().NewToken(r.URL.Query().Get("code"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, listener := range listeners {
		if err := listener.WebAuthLogin(w, r, token); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	state := r.URL.Query().Get("state")
	if state != "" {
		http.Redirect(w, r, state, http.StatusFound)
		return
	}
	http.Redirect(w, r, r.Referer(), http.StatusFound)
}
