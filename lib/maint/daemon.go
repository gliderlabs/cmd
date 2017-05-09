package maint

import "github.com/gliderlabs/comlab/pkg/log"

func (c *Component) AppPreStart() error {
	if Active() {
		log.Info(Notice(), log.Fields{"active": "true"})
	}
	return nil
}
