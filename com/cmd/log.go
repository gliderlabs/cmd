package cmd

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/getsentry/raven-go"
	"github.com/gliderlabs/gosper/pkg/log"
	"github.com/gliderlabs/ssh"
	"github.com/honeycombio/libhoney-go"
	"github.com/spf13/cast"

	"github.com/progrium/cmd/com/core"
)

func fieldProcessor(e log.Event, o interface{}) (log.Event, bool) {
	switch obj := o.(type) {
	case time.Duration:
		return e.Append("dur", cast.ToString(int64(obj/time.Millisecond))), true
	case log.ResponseWriter:
		e = e.Append("bytes", cast.ToString(obj.Size()))
		e = e.Append("status", cast.ToString(obj.Status()))
		return e, true
	case *http.Request:
		e = e.Append("ip", obj.RemoteAddr)
		e = e.Append("method", obj.Method)
		e = e.Append("path", obj.RequestURI)
		return e, true
	case *core.Command:
		e = e.Append("cmd.user", obj.User)
		e = e.Append("cmd.name", obj.Name)
		return e, true
	case ssh.Session:
		e = e.Append("sess.user", obj.User())
		e = e.Append("sess.remoteaddr", obj.RemoteAddr().String())
		e = e.Append("sess.command", strings.Join(obj.Command(), " "))
		return e, true
	}
	return e, false
}

func newRavenLogger(dsn string) *ravenLog {
	r, _ := raven.New(dsn)
	return &ravenLog{r}
}

type ravenLog struct {
	*raven.Client
}

func (c *ravenLog) Log(e log.Event) {
	err, ok := e.Fields["err"]
	if !ok {
		return
	}
	packet := raven.NewPacket(err,
		&raven.User{Username: e.Fields["sess.user"]},
		raven.NewException(errors.New(err), raven.NewStacktrace(4, 3, nil)))
	c.Capture(packet, nil)
}

type honeylog struct{}

func (c *honeylog) Log(e log.Event) {
	ev := libhoney.NewEvent()
	ev.Add(e.Fields)
	ev.Send()
}
