package log_test

import (
	"strconv"

	"github.com/gliderlabs/comlab/pkg/log"
)

type CustomData struct {
	Foo string
	Bar string
}

func fieldProcessor(e log.Event, field interface{}) (log.Event, bool) {
	switch obj := field.(type) {
	case int:
		return e.Append("num", strconv.Itoa(obj)), true
	case CustomData:
		e = e.Append("custom.foo", obj.Foo)
		e = e.Append("custom.bar", obj.Bar)
		return e, true
	}
	return e, false
}

// FieldProcessor example shows using the field processor callback to support new types.
func Example_fieldProcessor() {
	log.RegisterObserver(LogPrinter{})
	log.SetFieldProcessor(fieldProcessor)
	log.Info("Hello world", 12345, CustomData{"FOO", "BAR"})
	// Output: msg=Hello world, num=12345, custom.foo=FOO, custom.bar=BAR
}
