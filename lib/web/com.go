package web

import (
	"net/http"
	"text/template"

	"github.com/gliderlabs/comlab/pkg/com"
)

func init() {
	com.Register("web", &Component{},
		com.Option("listen_addr", "0.0.0.0:8080", "Address and port to listen on"),
		com.Option("static_dir", "ui/static/", "Directory to serve static files from"),
		com.Option("static_path", "/static", "URL path to serve static files at"),
		com.Option("cookie_secret", "", "Random string to use for session cookies"),
		com.Option("tls_addr", "0.0.0.0:4443", "Address and port to listen for TLS on"),
		com.Option("tls_cert", "", "Path to TLS cert file"),
		com.Option("tls_key", "", "Path to TLS key file"),
	)
}

// Handler extension point for matching and handling HTTP requests
type Handler interface {
	MatchHTTP(r *http.Request) bool
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

// Handlers accessor for web.Handler extension point
func Handlers() []Handler {
	var handlers []Handler
	for _, com := range com.Enabled(new(Handler), nil) {
		handlers = append(handlers, com.(Handler))
	}
	return handlers
}

type TemplateFuncProvider interface {
	WebTemplateFuncMap(r *http.Request) template.FuncMap
}

func TemplateFuncMap(r *http.Request) template.FuncMap {
	funcMap := template.FuncMap{}
	for _, com := range com.Enabled(new(TemplateFuncProvider), nil) {
		for k, v := range com.(TemplateFuncProvider).WebTemplateFuncMap(r) {
			funcMap[k] = v
		}
	}
	return funcMap
}

// Web component
type Component struct {
	http  *http.Server
	https *http.Server
	cert  *CertReloader
}
