package console

import (
	"fmt"

	"github.com/gliderlabs/comlab/pkg/com"
)

func (c *Component) AppPreStart() error {
	if com.GetString("slack_token") == "" {
		return fmt.Errorf("config for 'slack_token' must be set")
	}
	return nil
}
