package cmd

import (
	"os"

	"github.com/gliderlabs/comlab/pkg/com"
	"github.com/gliderlabs/comlab/pkg/log"
	"github.com/honeycombio/libhoney-go"
)

func (c *Component) AppPreStart() error {
	log.SetFieldProcessor(fieldProcessor)

	libhoney.Init(libhoney.Config{
		WriteKey:   com.GetString("honeycomb_key"),
		Dataset:    com.GetString("honeycomb_dataset"),
		SampleRate: 1,
	})
	// when all done, call close
	//defer libhoney.Close()
	hostname, _ := os.Hostname()
	libhoney.AddField("servername", hostname)
	libhoney.AddField("release", os.Getenv("RELEASE"))

	log.RegisterObserver(new(honeylog))
	return nil
}
