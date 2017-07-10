package console

import (
	"encoding/json"
	"net/http"
	"strings"
	"text/template"

	"github.com/gliderlabs/cmd/lib/access"
	"github.com/gliderlabs/cmd/lib/slack"
	"github.com/gliderlabs/cmd/lib/web"
	"github.com/gliderlabs/cmd/pkg/auth0"
	"github.com/gliderlabs/comlab/pkg/log"
)

func (c *Component) WebTemplateFuncMap(r *http.Request) template.FuncMap {
	return template.FuncMap{
		"title": strings.Title,
	}
}

func (c *Component) MatchHTTP(r *http.Request) bool {
	return strings.HasPrefix(r.URL.Path, "/console") ||
		strings.HasPrefix(r.URL.Path, "/login") ||
		strings.HasPrefix(r.URL.Path, "/register") ||
		strings.HasPrefix(r.URL.Path, "/request") ||
		r.URL.Path == "/" ||
		r.URL.Fragment == "NotFound"
}

func (c *Component) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" || r.URL.Fragment == "NotFound" {
		// temporary handler for notfound
		http.Redirect(w, r, "https://www.cmd.io/", http.StatusFound)
		return
	}
	if r.URL.Path == "/console" {
		http.Redirect(w, r, "/console/", http.StatusMovedPermanently)
		return
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/register", registerHandler)
	mux.HandleFunc("/request", requestAccessHandler)
	mux.HandleFunc("/console/-/billing", billingHandler)
	mux.HandleFunc("/console/-/codes", codesHandler)
	mux.HandleFunc("/console/", consoleHandler)
	mux.ServeHTTP(w, r)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	user := SessionUser(r)
	if user == nil {
		web.RenderTemplate(w, r, "login", map[string]interface{}{})
		return
	}
	if user.Account.CustomerID == "" && access.Check(user.Nickname) {
		web.RenderTemplate(w, r, "login-register", map[string]interface{}{
			"Username": user.Nickname,
		})
		return
	}
	http.Redirect(w, r, "/console/", http.StatusFound)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	user := SessionUser(r)
	if user == nil {
		web.RenderTemplate(w, r, "register", map[string]interface{}{})
		return
	}
	if user.Account.CustomerID == "" {
		err := RegisterUser(user)
		if err != nil {
			log.Info(r, err, user)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	http.Redirect(w, r, "/console/", http.StatusFound)
}

func requestAccessHandler(w http.ResponseWriter, r *http.Request) {
	user := SessionUser(r)
	if user == nil {
		web.RenderTemplate(w, r, "request", map[string]interface{}{})
		return
	}
	if access.Check(user.Nickname) {
		http.Redirect(w, r, "/console/", http.StatusFound)
		return
	}

	if err := slack.InviteToTeam("gliderlabs", user.Email); err != nil {
		log.Info(err, r, user)
	}

	web.RenderTemplate(w, r, "requested", map[string]interface{}{
		"Username": user.Nickname,
		"Email":    user.Email,
	})

}

func consoleHandler(w http.ResponseWriter, r *http.Request) {
	user := SessionUser(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	if !access.Check(user.Nickname) {
		http.Redirect(w, r, "/request", http.StatusFound)
		return
	}
	if user.Account.CustomerID == "" {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	billingInfo, err := GetBillingInfo(user)
	if err != nil {
		log.Info(r, err, log.Fields{"uid": user.ID})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	successFlash := web.SessionValue(r, "success")
	if successFlash != "" {
		web.SessionUnset(r, w, "success")
	}
	errorFlash := web.SessionValue(r, "error")
	if errorFlash != "" {
		web.SessionUnset(r, w, "error")
	}
	web.RenderTemplate(w, r, "console", map[string]interface{}{
		"Username":    user.Nickname,
		"Picture":     user.Picture,
		"BillingInfo": billingInfo,
		"Success":     successFlash,
		"Error":       errorFlash,
	})
}

func billingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	user := SessionUser(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	switch {
	case r.FormValue("email") != "":
		updateEmailHandler(w, r, user)
	case r.FormValue("update-token") != "":
		updatePaymentHandler(w, r, user)
	case r.FormValue("unsubscribe") != "":
		unsubscribeHandler(w, r, user)
	case r.FormValue("subscribe") != "":
		subscribeHandler(w, r, user)
	}
	http.Redirect(w, r, "/console/", http.StatusFound)
}

func codesHandler(w http.ResponseWriter, r *http.Request) {
	user := SessionUser(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	pendingInvites := user.Account.Invites.Pending
	if r.Method == "POST" {
		if len(user.Account.Invites.Pending) <= MaxPendingInvites {
			code, err := GenerateInviteCode(user)
			if err != nil {
				log.Info(r, err, log.Fields{"uid": user.ID})
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			pendingInvites = append(pendingInvites, code)
			err = auth0.DefaultClient().PatchUser(user.ID, auth0.User{
				"app_metadata": map[string]interface{}{
					"invites": map[string]interface{}{
						"pending":    pendingInvites,
						"invited_by": user.Account.Invites.InvitedBy,
					},
				},
			})
			if err != nil {
				log.Info(r, err, log.Fields{"uid": user.ID})
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
	w.Header().Add("content-type", "application/json")
	enc := json.NewEncoder(w)
	err := enc.Encode(map[string]interface{}{
		"pending": pendingInvites,
		"max":     MaxPendingInvites,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
