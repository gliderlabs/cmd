package main

import (
	_ "github.com/progrium/cmd/app/builtin"
	_ "github.com/progrium/cmd/app/cmd"
	_ "github.com/progrium/cmd/app/console"
	_ "github.com/progrium/cmd/app/runapi"
	_ "github.com/progrium/cmd/app/store"
	_ "github.com/progrium/cmd/app/store/dynamodb"
	_ "github.com/progrium/cmd/lib/access"
	_ "github.com/progrium/cmd/lib/maint"
	_ "github.com/progrium/cmd/lib/ssh"
	_ "github.com/progrium/cmd/lib/stripe"
	_ "github.com/progrium/cmd/lib/web"

	auth0 "github.com/progrium/cmd/pkg/auth0/com"
)

func init() {
	auth0.Register()
}
