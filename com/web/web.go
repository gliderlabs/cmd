package web

import (
	"net/http"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	"github.com/gorilla/sessions"
)

// Sessions is the global session manager
var Sessions = sessions.NewCookieStore([]byte("secret-secret-replace-me-use-env"))

// RenderTemplate ...
func RenderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	_, filename, _, _ := runtime.Caller(1)
	// filename is the full path on the system it was compiled.
	// that breaks running it anywhere else. so find automata root.
	parts := strings.Split(filename, "/")
	var prefix []string
	for _, part := range parts[1:] {
		if part == "com" {
			break
		}
		prefix = append(prefix, part)
	}
	relativeFilename := filename[len(strings.Join(prefix, "/"))+2:]
	tmplName := tmpl + ".html"
	t, err := template.New("").Funcs(template.FuncMap{"scripts": Scripts}).ParseFiles(
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
