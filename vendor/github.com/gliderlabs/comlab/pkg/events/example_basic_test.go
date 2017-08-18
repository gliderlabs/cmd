package events_test

import (
	"fmt"
	"time"

	"github.com/gliderlabs/comlab/pkg/events"
)

// exported constant allows other packages to refer to event by name
const EventLogin = "login"

// struct used as an event with fields specific to event
type LoginEvent struct {
	Username string
	Time     time.Time
}

// implementing Event interface makes it an event
func (e LoginEvent) EventName() string {
	return EventLogin
}

// Basic example shows a basic event listener and emitting an event struct
func Example_basic() {
	// any package can set up a listener
	events.Listen(&events.Listener{
		EventName: EventLogin,
		Handler: func(e events.Event) {
			login, ok := e.(*LoginEvent)
			if !ok {
				return
			}
			fmt.Println("Hello,", login.Username)
		},
	})
	// any package can emit the event,
	// though it's usually the package that owns the event
	events.Emit(&LoginEvent{
		Username: "progrium",
		Time:     time.Now(),
	})
	// Output: Hello, progrium
}
