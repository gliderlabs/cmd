package main

import (
	_ "github.com/progrium/cmd/com/builtin"
	_ "github.com/progrium/cmd/com/cmd"
	_ "github.com/progrium/cmd/com/console"
	_ "github.com/progrium/cmd/com/core"
	_ "github.com/progrium/cmd/com/maintenance"
	_ "github.com/progrium/cmd/com/store"
	_ "github.com/progrium/cmd/com/store/dynamodb"
	_ "github.com/progrium/cmd/com/store/filesystem"
	_ "github.com/progrium/cmd/com/stripe"
	_ "github.com/progrium/cmd/com/web"

	"github.com/progrium/cmd/com/cli"
	"github.com/progrium/cmd/com/sentry"

	access "github.com/progrium/cmd/pkg/access/com"
	auth0 "github.com/progrium/cmd/pkg/auth0/com"
)

var Version string

func init() {
	sentry.Release = Version
	cli.Version = Version
	auth0.Register()
	access.Register()
}
