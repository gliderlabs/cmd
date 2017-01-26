package console

import (
	"net/http"
	"strings"
	"text/template"

	"github.com/gliderlabs/comlab/pkg/log"
	"github.com/progrium/cmd/com/web"
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
		r.URL.Path == "/" ||
		r.URL.Fragment == "NotFound"
}

func (c *Component) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" || r.URL.Fragment == "NotFound" {
		// temporary handler for notfound
		http.Redirect(w, r, "http://gliderlabs.com/devlog/2016/announcing-cmd-io/", http.StatusFound)
		return
	}
	if r.URL.Path == "/console" {
		http.Redirect(w, r, "/console/", http.StatusMovedPermanently)
		return
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/register", registerHandler)
	mux.HandleFunc("/console/-/billing", billingHandler)
	mux.HandleFunc("/console/", consoleHandler)
	mux.ServeHTTP(w, r)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	user := SessionUser(r)
	if user == nil {
		web.RenderTemplate(w, r, "login", map[string]interface{}{})
		return
	}
	if user.Account.CustomerID == "" {
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

func consoleHandler(w http.ResponseWriter, r *http.Request) {
	user := SessionUser(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
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
