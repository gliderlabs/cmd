package runapi

import (
	"context"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/gliderlabs/ssh"
)

type httpSession struct {
	req         *http.Request
	wc          io.WriteCloser
	stdin       io.ReadCloser
	token       string
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
func (sess *httpSession) RemoteAddr() net.Addr {
	return &net.IPAddr{IP: net.ParseIP(sess.req.RemoteAddr)}
}
func (sess *httpSession) Environ() []string {
	common := []string{
		"SERVER_SOFTWARE=cmd.io",
		"REMOTE_ADDR=" + string(sess.RemoteAddr())}

	if sess.req == nil {
		return common
	}
	return append(common, []string{
		"SERVER_NAME=" + sess.req.Host,
		"SERVER_PROTOCOL=HTTP/1.1",
		"HTTP_HOST=" + sess.req.Host,
		"GATEWAY_INTERFACE=CGI/1.1",
		"REQUEST_METHOD=" + sess.req.Method,
		"QUERY_STRING=" + sess.req.URL.RawQuery,
		"REQUEST_URI=" + sess.req.URL.RequestURI(),
		"PATH_INFO=" + pathInfo,
		"SCRIPT_NAME=" + sess.Command(),
		"SCRIPT_FILENAME=" + "",
		"SERVER_PORT=" + port,
	})
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
