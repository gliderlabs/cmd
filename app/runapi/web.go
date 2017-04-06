package runapi

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/gliderlabs/cmd/app/console"
	"github.com/gliderlabs/cmd/app/store"
)

const runPrefix = "/run/"

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func (c *Component) MatchHTTP(r *http.Request) bool {
	return strings.HasPrefix(r.URL.Path, runPrefix)
}

func parseToken(r *http.Request) string {
	if token, ok := r.URL.Query()["access_token"]; ok {
		return token[0]
	}
	user, _, _ := r.BasicAuth()
	return user
}

func parsePath(r *http.Request) (string, string, []string) {
	path := strings.TrimPrefix(r.URL.Path, runPrefix)
	parts := strings.SplitN(path, "/", 3)
	if len(parts) < 2 {
		// not enough parts
		return "", "", []string{}
	}
	if len(parts) > 2 {
		// args
		if strings.Contains(parts[2], "/") {
			// args via path parts
			args := strings.Split(strings.Replace(parts[2], "+", " ", -1), "/")
			return parts[0], parts[1], args
		}
		return parts[0], parts[1], strings.Split(parts[2], "+")
	}
	// no args
	return parts[0], parts[1], []string{}
}

func (c *Component) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	token, _ := store.Selected().GetToken(parseToken(r))
	if token == nil {
		http.Error(w, "unauthorized token", http.StatusUnauthorized)
		return
	}
	owner, cmdName, args := parsePath(r)
	if owner == "" || cmdName == "" {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	cmd := store.Selected().Get(owner, cmdName)
	if cmd == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if !cmd.HasAccess(token.Key) {
		http.Error(w, "unauthorized token", http.StatusUnauthorized)
		return
	}

	var wc io.WriteCloser
	var isWebSocket bool
	if websocket.IsWebSocketUpgrade(r) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		wc = &gorillaWSAdapter{sync.Mutex{}, conn}
		isWebSocket = true
	} else {
		_, shouldStream := r.URL.Query()["stream"]
		if f, ok := w.(http.Flusher); ok && shouldStream {
			wc = &flushWriter{f, w}
		} else {
			wc = &flushWriter{nil, w}
		}
	}
	ctx := context.Background()
	u, err := console.LookupNickname(token.User)
	if err == nil {
		context.WithValue(ctx, "plan", u.Account.Plan)
	}
	session := &httpSession{
		req:         r,
		wc:          wc,
		stdin:       ioutil.NopCloser(strings.NewReader("")),
		token:       token.Key,
		isWebSocket: isWebSocket,
		ctx:         ctx,
	}
	defer session.Close()

	// TODO: put exit status in resp headers / stream trailers
	if status := cmd.Run(session, args); status != 0 {
		fmt.Fprintf(session, "exit status: %d", status)
		return
	}
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
