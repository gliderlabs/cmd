package maintenance

import (
	"strings"

	"github.com/gliderlabs/comlab/pkg/com"
	"github.com/gliderlabs/comlab/pkg/log"
)

func init() {
	com.Register("maintenance", &Component{},
		com.Option("active", false, ""),
		com.Option("notice", "cmd.io is currently down for maintenance", "displayed when maintenance active"),
		com.Option("allow", "", "comma separated list of users to allow during maintenance"))
}

// Component ...
type Component struct {
}

func (c *Component) AppPreStart() error {
	if Active() {
		log.Info(Notice(), log.Fields{"active": "true"})
	}
	return nil
}

// Active returns current maintenance state
func Active() bool {
	return com.GetBool("active")
}

// Allowed returns a slice of users which are allowed access during maintenance
func Allowed() []string {
	return strings.Split(com.GetString("allow"), ",")
}

// IsAllowed returns true if name is allowed access during maintenance
func IsAllowed(name string) bool {
	for _, u := range Allowed() {
		if u == name {
			return true
		}
	}
	return false
}

// Notice printed during maintenance
func Notice() string {
	return com.GetString("notice")
}
