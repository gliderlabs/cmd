package runapi

import (
	"context"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/gliderlabs/cmd/lib/release"
	"github.com/gliderlabs/ssh"
)

const (
	SERVER_PROTOCOL   = "HTTP/1.1"
	GATEWAY_INTERFACE = "CGI/1.1"
)

type httpSession struct {
	req         *http.Request
	wc          io.WriteCloser
	stdin       io.ReadCloser
	token       string
	cmdName     string
	isWebSocket bool
	ctx         context.Context
}

func (sess *httpSession) Write(p []byte) (n int, err error) {
	return sess.wc.Write(p)
}
func (sess *httpSession) Read(data []byte) (int, error) {
	return sess.stdin.Read(data)
}
func (sess *httpSession) PublicKey() ssh.PublicKey {
	return nil
}
func (sess *httpSession) Exit(code int) error {
	return nil
}
func (sess *httpSession) Permissions() ssh.Permissions {
	return ssh.Permissions{}
}
func (sess *httpSession) Context() context.Context {
	return sess.ctx
}
func (sess *httpSession) User() string {
	return sess.token
}
func (sess *httpSession) CmdName() string {
	return sess.cmdName
}
func (sess *httpSession) RemoteAddr() net.Addr {
	return &net.IPAddr{IP: net.ParseIP(sess.req.RemoteAddr)}
}
func (sess *httpSession) Environ() []string {
	if sess.req == nil {
		return nil
	}
	sh := strings.Split(sess.req.Host, ":")
	port := ""
	if len(sh) > 1 {
		port = sh[1]
	}
	return []string{
		"SERVER_NAME=" + release.Hostname(),
		"SERVER_PROTOCOL=" + SERVER_PROTOCOL,
		"HTTP_HOST=" + release.Hostname(),
		"GATEWAY_INTERFACE=" + GATEWAY_INTERFACE,
		"REQUEST_METHOD=" + sess.req.Method,
		"QUERY_STRING=" + sess.req.URL.RawQuery,
		"REQUEST_URI=" + sess.req.URL.RequestURI(),
		"PATH_INFO=" + sess.req.URL.Path,
		"SCRIPT_NAME=" + sess.CmdName(),
		"SERVER_PORT=" + port,
		"CONTENT_TYPE=" + sess.req.Header.Get("Content-Type"),
		"CONTENT_LENGTH=" + strconv.Itoa(int(sess.req.ContentLength)),
	}
}

func (sess *httpSession) Command() []string {
	// TODO
	return []string{}
}
func (sess *httpSession) Pty() (ssh.Pty, <-chan ssh.Window, bool) {
	var winch chan ssh.Window
	//if sess.isWebSocket {
	if false {
		win := ssh.Window{Width: 80, Height: 40}
		winch = make(chan ssh.Window, 1)
		winch <- win
		return ssh.Pty{
			Term:   "xterm",
			Window: win,
		}, winch, true
	}
	return ssh.Pty{}, winch, false
}

func (sess *httpSession) Close() error {
	sess.stdin.Close()
	return sess.wc.Close()
}
func (sess *httpSession) CloseWrite() error {
	return sess.wc.Close()
}
func (sess *httpSession) SendRequest(name string, wantReply bool, payload []byte) (bool, error) {
	return false, nil
}
func (sess *httpSession) Stderr() io.ReadWriter {
	return &struct {
		io.Reader
		io.Writer
	}{strings.NewReader(""), sess.wc}
}
