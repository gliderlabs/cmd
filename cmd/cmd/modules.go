package main

import (
	_ "github.com/progrium/cmd/com/cmd"
	_ "github.com/progrium/cmd/com/core"
	_ "github.com/progrium/cmd/com/meta"
	_ "github.com/progrium/cmd/com/redirect"
	_ "github.com/progrium/cmd/com/store"
	_ "github.com/progrium/cmd/com/store/dynamodb"
	_ "github.com/progrium/cmd/com/store/filesystem"
	_ "github.com/progrium/cmd/com/web"
)
