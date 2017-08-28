package web

import (
	"net/http"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	"github.com/gliderlabs/comlab/pkg/com"
	"github.com/gorilla/sessions"
)

// Sessions is the global session manager
// BUG: this is not set properly i don't think. needs to happen in app lifecycle
var Sessions = sessions.NewCookieStore([]byte(com.GetString("cookie_secret")))

func SessionValue(r *http.Request, key string) string {
	session, err := Sessions.Get(r, "session")
	if err != nil {
		return ""
	}
	val, exists := session.Values[key]
	if !exists {
		return ""
	}
	return val.(string)
}

func SessionSet(r *http.Request, w http.ResponseWriter, key, value string) error {
	session, err := Sessions.Get(r, "session")
	if err != nil {
		return err
	}
	session.Values[key] = value
	if err := session.Save(r, w); err != nil {
		return err
	}
	return nil
}

func SessionUnset(r *http.Request, w http.ResponseWriter, key string) error {
	session, err := Sessions.Get(r, "session")
	if err != nil {
		return err
	}
	delete(session.Values, key)
	if err := session.Save(r, w); err != nil {
		return err
	}
	return nil
}

func SessionDel(r *http.Request, w http.ResponseWriter) error {
	session, err := Sessions.Get(r, "session")
	if err != nil {
		return err
	}
	session.Options.MaxAge = -1
	if err := session.Save(r, w); err != nil {
		return err
	}
	return nil
}

// RenderTemplate ...
func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, data interface{}) {
	_, filename, _, _ := runtime.Caller(1)
	// filename is the full path on the system it was compiled.
	// that breaks running it anywhere else. so find automata root.
	parts := strings.Split(filename, "/")
	var prefix []string
	for _, part := range parts[1:] {
		// TODO: fix this, it should not be dependent on file structure
		if part == "app" || part == "lib" {
			break
		}
		prefix = append(prefix, part)
	}
	relativeFilename := filename[len(strings.Join(prefix, "/"))+2:]
	tmplName := tmpl + ".html"
	t, err := template.New("").Funcs(TemplateFuncMap(r)).ParseFiles(
		filepath.Join(filepath.Dir(relativeFilename), "ui", "html", tmplName))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.ExecuteTemplate(w, tmplName, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
