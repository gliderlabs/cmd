package main

import (
	_ "github.com/gliderlabs/cmd/app/builtin"
	_ "github.com/gliderlabs/cmd/app/cmd"
	_ "github.com/gliderlabs/cmd/app/console"
	_ "github.com/gliderlabs/cmd/app/runapi"
	_ "github.com/gliderlabs/cmd/app/store"
	_ "github.com/gliderlabs/cmd/app/store/dynamodb"
	_ "github.com/gliderlabs/cmd/lib/access"
	_ "github.com/gliderlabs/cmd/lib/crypto"
	_ "github.com/gliderlabs/cmd/lib/github"
	_ "github.com/gliderlabs/cmd/lib/maint"
	_ "github.com/gliderlabs/cmd/lib/slack"
	_ "github.com/gliderlabs/cmd/lib/ssh"
	_ "github.com/gliderlabs/cmd/lib/stripe"
	_ "github.com/gliderlabs/cmd/lib/web"

	auth0 "github.com/gliderlabs/cmd/pkg/auth0/com"
)

func init() {
	auth0.Register()
}
