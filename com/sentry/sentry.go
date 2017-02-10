package sentry

import (
	"errors"

	raven "github.com/getsentry/raven-go"
	"github.com/gliderlabs/comlab/pkg/com"
	"github.com/gliderlabs/comlab/pkg/log"
)

func init() {
	com.Register("sentry", &Component{},
		com.Option("dsn", "", "dsn for sentry project"),
		com.Option("environment", "dev", "environment tag"))
	log.RegisterObserver(&Component{})
}

var Release string

// Component ...
type Component struct{}

// Client returns a configured ravent client
func Client() *raven.Client {
	client, _ := raven.New(com.GetString("dsn"))
	if client != nil {
		client.SetEnvironment(com.GetString("environment"))
		client.SetRelease(Release)
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
