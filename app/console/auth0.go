package console

import (
	"net/http"

	"github.com/progrium/cmd/lib/web"
	"github.com/progrium/cmd/pkg/auth0"
	"golang.org/x/oauth2"
)

func (c *Component) WebAuthLogin(w http.ResponseWriter, r *http.Request, token *oauth2.Token) error {
	userinfo, err := auth0.DefaultClient().UserInfo(token)
	if err != nil {
		return err
	}

	web.SessionSet(r, w, "_access_token", token.AccessToken)
	web.SessionSet(r, w, "_auth_id", userinfo["user_id"].(string))
	web.SessionSet(r, w, "user_name", userinfo["name"].(string))
	web.SessionSet(r, w, "user_nickname", userinfo["nickname"].(string))
	web.SessionSet(r, w, "user_email", userinfo["email"].(string))
	web.SessionSet(r, w, "user_id", userinfo["user_id"].(string))

	//http.Redirect(w, r, r.Referer(), http.StatusFound)
	return nil
}

func (c *Component) WebAuthLogout(w http.ResponseWriter, r *http.Request) error {
	session, _ := web.Sessions.Get(r, "session")
	delete(session.Values, "_auth_id")
	delete(session.Values, "_access_token")
	delete(session.Values, "user_name")
	delete(session.Values, "user_nickname")
	delete(session.Values, "user_email")
	delete(session.Values, "user_id")
	session.Options.MaxAge = -1
	return session.Save(r, w)
}
