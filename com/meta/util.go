package meta

import (
	"fmt"
	"io"

	"github.com/mgutz/ansi"
)

var (
	gray   = ansi.ColorFunc("black+h")
	cyan   = ansi.ColorFunc("cyan")
	red    = ansi.ColorFunc("red")
	yellow = ansi.ColorFunc("yellow")
	reset  = ansi.ColorFunc("reset")
	bright = ansi.ColorFunc("white+h")
)

func header(w io.Writer, text string) (int, error) {
	return fmt.Fprintln(w, gray("==="), bright(text))
}

func statusMsg(w io.Writer, message string) (int, error) {
	return fmt.Fprintf(w, "%s... ", message)
}

func statusErr(w io.Writer, message string) (int, error) {
	return fmt.Fprintln(w, red("error:"), message)
}

func statusDone(w io.Writer) (int, error) {
	return fmt.Fprintln(w, "done")
}

func printFields(w io.Writer, fields map[string]interface{}, colon bool) {
	longest := 0
	for k := range fields {
		if len(k) > longest {
			longest = len(k) + 2
		}
	}
	for k, v := range fields {
		if colon {
			k = k + ":"
		}
		fmt.Fprintf(w, fmt.Sprintf("%%-%ds  %%s\n", longest), k, v)
	}
}
