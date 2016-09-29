package cmd

import (
	"errors"
	"fmt"
	"strings"
	"time"

	raven "github.com/getsentry/raven-go"
	"github.com/gliderlabs/pkg/log"
	"github.com/gliderlabs/pkg/ssh"
	libhoney "github.com/honeycombio/libhoney-go"
	"github.com/mgutz/ansi"
)

func init() {
	if !LocalMode() {
		ansi.DisableColors(true)
	}
}

var (
	gray   = ansi.ColorFunc("black+h")
	cyan   = ansi.ColorFunc("cyan")
	red    = ansi.ColorFunc("red")
	yellow = ansi.ColorFunc("yellow")
	reset  = ansi.ColorFunc("reset")
	bright = ansi.ColorFunc("white+h")
)

type logging struct{}

func (c *logging) Log(e log.Event) {
	color := reset
	switch e.Type {
	case log.TypeLocal:
		color = yellow
	case log.TypeFatal:
		color = red
	case log.TypeInfo:
		if DebugMode() {
			color = bright
		}
	}
	if _, ok := e.Fields["err"]; ok {
		color = red
	}
	pkg := e.Fields["pkg"]
	e = e.Remove("pkg")
	var parts []string
	for _, key := range e.Index {
		if key == "msg" {
			parts = append([]string{e.Fields[key]}, parts...)
		} else {
			parts = append(parts, fmt.Sprintf("%s=%v", key, e.Fields[key]))
		}
	}
	fmt.Println(gray(e.Time.Format("15:04:05.000")), cyan("["+pkg+"]"), color(strings.Join(parts, " ")))
}

func fieldProcessor(e log.Event, o interface{}) (log.Event, bool) {
	switch obj := o.(type) {
	case time.Duration:
		return e.Append("dur", obj.String()), true
	case *Command:
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
	fmt.Println(e.Fields)
	packet := raven.NewPacket(err,
		&raven.User{Username: e.Fields["sess.user"]},
		raven.NewException(errors.New(err), raven.NewStacktrace(4, 3, nil)))
	c.Capture(packet, nil)
}

type honeylog struct{}

func (c *honeylog) Log(e log.Event) {

	ev := libhoney.NewEvent()
	for k, v := range e.Fields {
		ev.AddField(k, v)
	}
	// send the event
	ev.Send()
}
