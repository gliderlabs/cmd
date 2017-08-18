package log_test

import (
	"fmt"
	"strings"

	"github.com/gliderlabs/comlab/pkg/log"
)

type LogPrinter struct{}

// This Observer iterates over the event Index to maintain order of fields,
// and creates a list of comma separated key-value pairs to print out.
func (lp LogPrinter) Log(event log.Event) {
	var fields []string
	for _, key := range event.Index {
		fields = append(fields, fmt.Sprintf("%s=%s", key, event.Fields[key]))
	}
	fmt.Println(strings.Join(fields, ", "))
}

// Basic example shows a simple log observer and using the logging functions.
func Example_basic() {
	log.RegisterObserver(LogPrinter{})
	log.Info("Hello world", log.Fields{"foo": "bar"}, 12345)
	// Output: msg=Hello world, foo=bar, data=12345
}
