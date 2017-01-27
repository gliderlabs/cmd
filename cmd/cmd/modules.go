package main

import (
	_ "github.com/progrium/cmd/com/cmd"
	_ "github.com/progrium/cmd/com/console"
	_ "github.com/progrium/cmd/com/core"
	_ "github.com/progrium/cmd/com/meta"
	_ "github.com/progrium/cmd/com/store"
	_ "github.com/progrium/cmd/com/store/dynamodb"
	_ "github.com/progrium/cmd/com/store/filesystem"
	_ "github.com/progrium/cmd/com/stripe"
	_ "github.com/progrium/cmd/com/web"

	access "github.com/progrium/cmd/pkg/access/com"
	auth0 "github.com/progrium/cmd/pkg/auth0/com"
)

func init() {
	auth0.Register()
	access.Register()
}
