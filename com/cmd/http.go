package cmd

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/gliderlabs/comlab/pkg/log"
	"github.com/gorilla/websocket"
	"github.com/progrium/cmd/com/core"
	"github.com/progrium/cmd/com/store"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func (c *Component) MatchHTTP(r *http.Request) bool {
	return strings.HasPrefix(r.URL.Path, "/run/")
}

func getToken(r *http.Request) string {
	if token, ok := r.URL.Query()["access_token"]; ok {
		return token[0]
	}
	auth := r.Header.Get("Authorization")
	fields := strings.SplitN(auth, " ", 2)
	if len(fields) != 2 {
		return ""
	}
	if fields[0] == "Basic" {
		b, _ := base64.StdEncoding.DecodeString(fields[1])
		// remove colon which may be present even if a password was not provided
		return strings.TrimSuffix(string(b), ":")
	}
	return fields[1]
}

func (c *Component) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	token, _ := store.Selected().GetToken(getToken(r))
	if token == nil {
		http.Error(w, "unauthorized token", http.StatusUnauthorized)
		return
	}
	parts := strings.SplitN(r.URL.Path, "/", 5)
	if len(parts) < 3 {
		http.Error(w, "path missing user and/or cmd", http.StatusBadRequest)
		return
	}
	cmd := store.Selected().Get(parts[2], parts[3])
	if cmd == nil {
		http.Error(w, "cmd not found", http.StatusNotFound)
		return
	}

	if !cmd.HasAccess(token.Key) {
		http.Error(w, "unauthorized token", http.StatusUnauthorized)
		return
	}

	var ow io.WriteCloser
	if ow = outputWriter(w, r); ow == nil {
		// TODO: errors in outputWriter need to be handled
		return
	}
	defer ow.Close()
	stream := &core.Stream{
		Stdin:  ioutil.NopCloser(strings.NewReader("")),
		Stdout: ow,
		Stderr: ow,
	}

	var args []string
	if len(parts) > 4 {
		args = strings.Split(parts[4], "+")
	}

	// TODO: put exit status in resp headers / stream trailers
	if status := cmd.Run(stream, token.Key, args); status != 0 {
		fmt.Fprintf(ow, "exit status: %d", status)
		return
	}
}

func outputWriter(w http.ResponseWriter, r *http.Request) io.WriteCloser {
	if websocket.IsWebSocketUpgrade(r) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Debug("upgrade failed", err)
			return nil
		}
		return &gorillaWSAdapter{sync.Mutex{}, conn}
	}
	if _, stream := r.URL.Query()["stream"]; stream {
		if f, ok := w.(http.Flusher); ok {
			return &flushWriter{f, w}
		}
	}
	return &flushWriter{nil, w}
}

type flushWriter struct {
	f http.Flusher
	w io.Writer
}

func (fw *flushWriter) Close() error {
	return nil
}

func (fw *flushWriter) Write(p []byte) (n int, err error) {
	n, err = fw.w.Write(p)
	if fw.f != nil {
		fw.f.Flush()
	}
	return
}

type gorillaWSAdapter struct {
	sync.Mutex
	*websocket.Conn
}

func (ws *gorillaWSAdapter) Write(p []byte) (int, error) {
	ws.Lock()
	defer ws.Unlock()
	return len(p), ws.WriteMessage(websocket.TextMessage, p)
}

func (ws *gorillaWSAdapter) Close() error {
	ws.Lock()
	defer ws.Unlock()
	ws.WriteMessage(websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	return ws.Conn.Close()
}
