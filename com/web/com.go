package web

import (
	"net/http"
	"strings"

	"github.com/facebookgo/httpdown"

	"github.com/gliderlabs/gosper/pkg/com"
)

func init() {
	com.Register("web", &Component{},
		com.Option("listen_addr", "127.0.0.1:8080", "Address and port to listen on"),
		com.Option("static_dir", "ui/static/", "Directory to serve static files from"),
		com.Option("static_path", "/static", "URL path to serve static files at"))
}

type ScriptProvider interface {
	PageScript() string
}

func Scripts() string {
	var scripts []string
	for _, com := range com.Enabled(new(ScriptProvider), nil) {
		scripts = append(scripts, com.(ScriptProvider).PageScript())
	}
	return strings.Join(scripts, "\n")
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

// Web component
type Component struct {
	http httpdown.Server
}
