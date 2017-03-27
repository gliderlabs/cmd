package honeycomb

import (
	"os"

	"github.com/gliderlabs/comlab/pkg/com"
	"github.com/gliderlabs/comlab/pkg/log"
	"github.com/honeycombio/libhoney-go"
)

func (c *Component) AppPreStart() error {
	libhoney.Init(libhoney.Config{
		WriteKey:   com.GetString("key"),
		Dataset:    com.GetString("dataset"),
		SampleRate: 1,
	})
	hostname, _ := os.Hostname()
	libhoney.AddField("servername", hostname)
	libhoney.AddField("release", os.Getenv("RELEASE"))

	log.RegisterObserver(new(honeylogger))
	return nil
}
