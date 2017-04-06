package sentry

import (
	"bytes"
	"errors"
	"os"
	"strings"

	raven "github.com/getsentry/raven-go"
	"github.com/gliderlabs/comlab/pkg/com"
	"github.com/gliderlabs/comlab/pkg/log"
	"github.com/maruel/panicparse/stack"
	"github.com/gliderlabs/cmd/lib/release"
)

func init() {
	com.Register("sentry", &Component{},
		com.Option("dsn", "", "dsn for sentry project"))
	log.RegisterObserver(&Component{})
}

// Component ...
type Component struct{}

// Client returns a configured ravent client
func Client() *raven.Client {
	client, _ := raven.New(com.GetString("dsn"))
	if client != nil {
		client.SetEnvironment(release.Channel())
		client.SetRelease(release.Version)
	}
	return client
}

// Log captures events describing an error and sends them to sentry
func (c *Component) Log(e log.Event) {
	err, ok := e.Fields["err"]
	if !ok {
		return
	}

	packet := raven.NewPacket(err,
		&raven.User{Username: e.Fields["sess.user"]},
		raven.NewException(errors.New(err), raven.NewStacktrace(4, 3, nil)))
	Client().Capture(packet, nil)
}

func PanicHandler(output string) {
	title := output[strings.Index(output, "\n"):]
	in := bytes.NewBufferString(output)
	goroutines, _ := stack.ParseDump(in, os.Stdout)

	var cause stack.Goroutine
	for _, gr := range goroutines {
		if gr.First {
			cause = gr
			break
		}
	}

	var trace raven.Stacktrace
	for _, call := range cause.Stack.Calls {
		frame := &raven.StacktraceFrame{
			Function:     call.Func.Name(),
			AbsolutePath: call.SourcePath,
			Module:       call.Func.String(),
			Lineno:       call.Line,
			Filename:     call.SourceName(),
		}
		trace.Frames = append(trace.Frames, frame)
	}

	packet := raven.NewPacket(title, &trace)
	packet.Level = raven.FATAL
	_, err := Client().Capture(packet, nil)
	<-err
	os.Exit(1)
}
