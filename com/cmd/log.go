package cmd

import (
	"fmt"
	"strings"

	"github.com/gliderlabs/pkg/log"
	"github.com/mgutz/ansi"
)

func init() {
	if !LocalMode() {
		ansi.DisableColors(true)
	}
}

var (
	gray   = ansi.ColorFunc("black+h")
	cyan   = ansi.ColorFunc("cyan")
	red    = ansi.ColorFunc("red")
	yellow = ansi.ColorFunc("yellow")
	reset  = ansi.ColorFunc("reset")
	bright = ansi.ColorFunc("white+h")
)

type logging struct{}

func (c *logging) Log(e log.Event) {
	color := reset
	switch e.Type {
	case log.TypeLocal:
		color = yellow
	case log.TypeFatal:
		color = red
	case log.TypeInfo:
		if DebugMode() {
			color = bright
		}
	}
	if _, ok := e.Fields["err"]; ok {
		color = red
	}
	pkg := e.Fields["pkg"]
	e = e.Remove("pkg")
	var parts []string
	for _, key := range e.Index {
		if key == "msg" {
			parts = append([]string{e.Fields[key]}, parts...)
		} else {
			parts = append(parts, fmt.Sprintf("%s=%v", key, e.Fields[key]))
		}
	}
	fmt.Println(gray(e.Time.Format("15:04:05.000")), cyan("["+pkg+"]"), color(strings.Join(parts, " ")))
}
