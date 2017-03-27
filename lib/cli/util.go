package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"

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

func Gray(s interface{}) string {
	if str, ok := s.(string); ok {
		return gray(str)
	}
	return ansi.LightBlack
}

func Reset(s interface{}) string {
	if str, ok := s.(string); ok {
		return reset(str)
	}
	return ansi.Reset
}

func Bright(s interface{}) string {
	if str, ok := s.(string); ok {
		return bright(str)
	}
	return ansi.LightWhite
}

func JSON(w io.Writer, obj interface{}) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(obj)
}

func Header(w io.Writer, text string) (int, error) {
	return fmt.Fprintln(w, gray("==="), bright(text))
}

func Status(w io.Writer, message string) (int, error) {
	return fmt.Fprintf(w, "%s... ", message)
}

func StatusDone(w io.Writer) (int, error) {
	return fmt.Fprintln(w, "done")
}

func StatusErr(w io.Writer, message string) (int, error) {
	return fmt.Fprintln(w, red("error:"), message)
}

func PrintFields(w io.Writer, fields map[string]interface{}, colon bool) {
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

func NewTable(w io.Writer) *tabwriter.Writer {
	return tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
}
