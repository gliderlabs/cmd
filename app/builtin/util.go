package builtin

import (
	"fmt"
	"strings"

	"github.com/gliderlabs/cmd/app/core"
	"github.com/gliderlabs/cmd/app/store"
)

func LookupCmd(owner, name string) (*core.Command, error) {
	var (
		userName = owner
		cmdName  = name
	)
	if strings.Contains(cmdName, "/") {
		parts := strings.SplitN(cmdName, "/", 2)
		userName = parts[0]
		cmdName = parts[1]
	}
	cmd := store.Selected().Get(userName, cmdName)
	if cmd == nil {
		return nil, fmt.Errorf("Command not found: %s", name)
	}
	return cmd, nil
}
