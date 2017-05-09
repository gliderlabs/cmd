package maint

import (
	"strings"

	"github.com/gliderlabs/comlab/pkg/com"
)

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
