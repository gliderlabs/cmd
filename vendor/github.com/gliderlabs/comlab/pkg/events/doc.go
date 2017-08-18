/*
Package events implements a simple callback-based event emitter with event
filtering and one-time event handling capabilities.

Although you can create individual events.Emitter structs, it's more common to
use the default emitter and top-level package API to emit and listen for
application-wide events.

The events.Listen and events.Unlisten methods work with events.Listener structs
that wrap the event handling callback. This lets you define several ways to
filter for certain events, as well as specify it as a one-off listener that only
gets called once.

Events emitted need to implement the events.Event interface, allowing arbitrary
event value types as long as they implement the interface. This helps identify
them as an event and identifies the name of the event, useful for further type
assertions.

Handler callbacks need to be aware they are called synchronously and should perform
any time-intensive operation asynchronously. Otherwise it will block the call to
Emit, which shouldn't have to be called asynchronously.

*/
package events
