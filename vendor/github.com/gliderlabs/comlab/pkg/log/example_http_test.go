package log_test

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gliderlabs/comlab/pkg/log"
)

// our field processor will need to know how to process HTTP related types
func httpFieldProcessor(e log.Event, field interface{}) (log.Event, bool) {
	switch obj := field.(type) {
	case time.Duration:
		return e.Append("dur", strconv.Itoa(int(obj/time.Millisecond))), true
	case log.ResponseWriter:
		e = e.Append("bytes", strconv.Itoa(obj.Size()))
		e = e.Append("status", strconv.Itoa(obj.Status()))
		return e, true
	case *http.Request:
		e = e.Append("ip", obj.RemoteAddr)
		e = e.Append("method", obj.Method)
		e = e.Append("path", obj.RequestURI)
		return e, true
	}
	return e, false
}

// our logging middleware that injects a wrapped ResponseWriter
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		lw := log.WrapResponseWriter(w)
		defer log.Info(r, lw, time.Now().Sub(t))
		next.ServeHTTP(lw, r)
	})
}

// HTTP example shows using the ResponseWriter wrapper for HTTP logging
func Example_http() {
	log.SetFieldProcessor(httpFieldProcessor)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("This is a catch-all route"))
	})

	// wrap our mux with our middleware
	http.ListenAndServe(":8080", loggingMiddleware(mux))
}
