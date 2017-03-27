package cmd

import (
	"github.com/gliderlabs/comlab/pkg/log"
)

func (c *Component) AppPreStart() error {
	log.SetFieldProcessor(fieldProcessor)
	return nil
}
