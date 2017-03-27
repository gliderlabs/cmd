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
	// TODO
	return []string{}
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
