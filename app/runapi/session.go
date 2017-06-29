package runapi

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/gliderlabs/cmd/lib/release"
	"github.com/gliderlabs/ssh"
)

const (
	ServerProtocol   = "HTTP/1.1"
	GatewayInterface = "CGI/1.1"
)

type httpSession struct {
	req         *http.Request
	wc          io.WriteCloser
	token       string
	isWebSocket bool
	ctx         context.Context
	cmd         []string
}

func (sess *httpSession) Write(p []byte) (n int, err error) {
	return sess.wc.Write(p)
}
func (sess *httpSession) Read(data []byte) (int, error) {
	return sess.req.Body.Read(data)
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
	if sess.req == nil {
		return nil
	}
	sh := strings.Split(sess.req.Host, ":")
	port := ""
	if len(sh) > 1 {
		port = sh[1]
	}
	return []string{
		fmt.Sprintf("SERVER_NAME=%s", release.Hostname()),
		fmt.Sprintf("SERVER_PROTOCOL=%s", ServerProtocol),
		fmt.Sprintf("HTTP_HOST=%s", release.Hostname()),
		fmt.Sprintf("GATEWAY_INTERFACE=%s", GatewayInterface),
		fmt.Sprintf("REQUEST_METHOD=%s", sess.req.Method),
		fmt.Sprintf("QUERY_STRING=%s", sess.req.URL.RawQuery),
		fmt.Sprintf("REQUEST_URI=%s", sess.req.URL.RequestURI()),
		fmt.Sprintf("PATH_INFO=%s", sess.req.URL.Path),
		fmt.Sprintf("SCRIPT_NAME=%s", sess.Command()[0]),
		fmt.Sprintf("SERVER_PORT=%s", port),
		fmt.Sprintf("CONTENT_TYPE=%s", sess.req.Header.Get("Content-Type")),
		fmt.Sprintf("CONTENT_LENGTH=%s", strconv.Itoa(int(sess.req.ContentLength))),
	}
}

func (sess *httpSession) Command() []string {
	return append([]string(nil), sess.cmd...)
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
	sess.req.Body.Close()
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
