package web

import (
	"crypto/tls"
	"net/http"
	"strings"
	"sync"
	"time"

	"context"

	"github.com/gliderlabs/cmd/lib/maint" // TODO: remove dep via hook
	gcontext "github.com/gorilla/context"

	"github.com/gliderlabs/comlab/pkg/com"
	"github.com/gliderlabs/comlab/pkg/log"
)

// Serve of com.Service extension point
func (c *Component) Serve() {
	handler := gcontext.ClearHandler(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t := time.Now()
			lw := log.WrapResponseWriter(w)
			defer func() {
				log.Info(r, lw, time.Now().Sub(t))
			}()

			if maint.Active() {
				http.Error(lw, maint.Notice(), http.StatusServiceUnavailable)
				return
			}

			// redirect to https if available
			// if r.URL.Scheme != "https" && c.https != nil {
			// 	u := r.URL
			// 	u.Host = r.Host
			// 	u.Scheme = "https"
			// 	http.Redirect(w, r, u.String(), 302)
			// 	return
			// }

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
		}),
	)

	c.http = &http.Server{
		Addr:    com.GetString("listen_addr"),
		Handler: handler,
	}
	if c.cert != nil {
		c.https = &http.Server{
			Addr:    com.GetString("tls_addr"),
			Handler: handler,
			TLSConfig: &tls.Config{
				GetCertificate: c.cert.GetCertificate,
			},
		}
		go func() {
			for {
				time.Sleep(1 * time.Hour)
				err := c.cert.Reload()
				if err != nil {
					log.Info("unable to reload cert", err)
				}
			}
		}()
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {

		log.Info("http listening on", c.http.Addr)
		if err := c.http.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
		wg.Done()
	}()
	if c.https != nil {
		wg.Add(1)
		go func() {
			log.Info("https listening on", c.https.Addr)
			if err := c.https.ListenAndServeTLS("", ""); err != nil {
				log.Fatal(err)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

// Stop of com.Service extension point
func (c *Component) Stop() {
	if c.https != nil {
		c.https.Shutdown(context.Background())
	}
	c.http.Shutdown(context.Background())
}
