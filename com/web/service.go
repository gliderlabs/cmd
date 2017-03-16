package web

import (
	"net/http"
	"strings"
	"time"

	"github.com/facebookgo/httpdown"
	"github.com/gorilla/context"
	"github.com/progrium/cmd/com/maintenance"

	"github.com/gliderlabs/comlab/pkg/com"
	"github.com/gliderlabs/comlab/pkg/log"
)

// Serve of com.Service extension point
func (c *Component) Serve() {
	server := &http.Server{
		Addr: com.GetString("listen_addr"),
		Handler: context.ClearHandler(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				t := time.Now()
				lw := log.WrapResponseWriter(w)
				defer log.Info(r, lw, time.Now().Sub(t))

				if maintenance.Active() {
					http.Error(lw, maintenance.Notice(), http.StatusServiceUnavailable)
					return
				}

				// serve static
				staticPrefix := com.GetString("static_path") + "/"
				if strings.HasPrefix(r.URL.Path, staticPrefix) {
					http.StripPrefix(staticPrefix,
						http.FileServer(http.Dir(com.GetString("static_dir")))).ServeHTTP(lw, r)
					return
				}

				// serve semantic ui src. TODO: make dev only
				if strings.HasPrefix(r.URL.Path, "/_semantic/") {
					http.StripPrefix("/_semantic/",
						http.FileServer(http.Dir("ui/semantic/src"))).ServeHTTP(lw, r)
					return
				}

				// serve component
				for _, handler := range Handlers() {
					if handler.MatchHTTP(r) {
						handler.ServeHTTP(lw, r)
						return
					}
				}
				r.URL.Fragment = "NotFound"
				for _, handler := range Handlers() {
					if handler.MatchHTTP(r) {
						handler.ServeHTTP(lw, r)
						return
					}
				}
				http.NotFound(lw, r)
			})),
	}
	hd := &httpdown.HTTP{
		StopTimeout: 10 * time.Second,
		KillTimeout: 1 * time.Second,
	}
	var err error
	log.Info("listening on", server.Addr)
	c.http, err = hd.ListenAndServe(server) // TODO: replace with actual serving goroutine
	if err != nil {
		log.Fatal(err)
	}
	c.http.Wait()
}

// Stop of com.Service extension point
func (c *Component) Stop() {
	c.http.Stop()
}
