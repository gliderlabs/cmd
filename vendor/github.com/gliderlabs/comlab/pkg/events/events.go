package events

import "sync"

var defaultEmitter = &Emitter{}

// Event interface is used to identify events
type Event interface {
	EventName() string
}

// Signal is a builtin event that is just a name with no other payload
type Signal string

func (s Signal) EventName() string {
	return string(s)
}

// Listener wraps an event handler callback, optionally specifying how
// it receives events from an Emitter
type Listener struct {
	// If not empty, only receive events with this name
	EventName string

	// If true, will fire only once and then listener will be removed
	Once bool

	// If not nil, a callback to determine if handler will be called
	Filter func(Event) bool

	// Callback to handle event
	Handler func(Event)
}

// Emitter is a collection of Listeners that you can Emit events to
type Emitter struct {
	listeners []*Listener
	mu        sync.Mutex
}

// Listen adds a listener to receive events on the default emitter.
//
// If Handler is nil, the listener is ignored.
func Listen(listener *Listener) { defaultEmitter.Listen(listener) }
func (e *Emitter) Listen(listener *Listener) {
	if listener.Handler == nil {
		return
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	e.listeners = append(e.listeners, listener)
}

// Unlisten removes a listener from receiving events on the default emitter.
func Unlisten(listener *Listener) { defaultEmitter.Unlisten(listener) }
func (e *Emitter) Unlisten(listener *Listener) {
	e.mu.Lock()
	defer e.mu.Unlock()
	var i int
	var l *Listener
	var found bool
	for i, l = range e.listeners {
		if l == listener {
			found = true
			break
		}
	}
	if found {
		e.listeners = append(e.listeners[:i], e.listeners[i+1:]...)
	}
}

// Emit triggers the handlers of listeners that match the event, passing
// the event to them.
//
// If it triggers a listener with Once set to true, the listener will be removed
// after triggering.
func Emit(event Event) { defaultEmitter.Emit(event) }
func (e *Emitter) Emit(event Event) {
	e.mu.Lock()
	defer e.mu.Unlock()
	eventName := event.EventName()
	var keep []*Listener
	for _, listener := range e.listeners {
		if listener.EventName == eventName || listener.EventName == "" {
			if listener.Filter != nil {
				if ok := listener.Filter(event); !ok {
					keep = append(keep, listener)
					continue
				}
			}
			listener.Handler(event)
			if !listener.Once {
				keep = append(keep, listener)
			}
		} else {
			keep = append(keep, listener)
		}
	}
	e.listeners = keep
}
