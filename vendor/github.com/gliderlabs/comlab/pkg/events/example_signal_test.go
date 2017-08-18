package events_test

import (
	"fmt"

	"github.com/gliderlabs/comlab/pkg/events"
)

const SignalHello = "Hello world"

// Signal example shows using Signal type as event, which has no payload but its name
func Example_signal() {
	events.Listen(&events.Listener{
		EventName: SignalHello,
		Handler: func(e events.Event) {
			fmt.Println(e)
		},
	})
	events.Emit(events.Signal(SignalHello))
	// Output: Hello world
}
