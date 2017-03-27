package honeycomb

import (
	"github.com/gliderlabs/comlab/pkg/log"
	"github.com/honeycombio/libhoney-go"
)

type honeylogger struct{}

func (c *honeylogger) Log(e log.Event) {
	ev := libhoney.NewEvent()
	if e.Fields["dur"] == "" {
		ev.Dataset += "_stream"
	}
	ev.Add(e.Fields)
	ev.Send()
}
