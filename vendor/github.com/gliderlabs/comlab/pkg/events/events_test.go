package events

import "testing"

func (emitter *Emitter) emitAndCapture(events []Event, listeners ...*Listener) int {
	captured := 0
	for _, l := range listeners {
		l.Handler = func(e Event) {
			captured += 1
		}
		emitter.Listen(l)
	}
	for _, e := range events {
		emitter.Emit(e)
	}
	return captured
}

func TestEvents(t *testing.T) {
	t.Parallel()
	emitter := &Emitter{}
	for _, test := range []struct {
		events    []Event
		want      int
		listeners []*Listener
	}{
		{events: []Event{Signal("foo")}, want: 2, listeners: []*Listener{
			&Listener{
				EventName: "foo",
			}, // trigger
			&Listener{
				EventName: "",
			}, // trigger
			&Listener{
				EventName: "bar",
			}, // no trigger
		}},
		{events: []Event{Signal("bar")}, want: 2, listeners: []*Listener{
			&Listener{
				EventName: "bar",
			}, // trigger
			&Listener{
				EventName: "bar",
				Filter: func(e Event) bool {
					return e == Signal("bar")
				},
			}, // trigger
			&Listener{
				EventName: "bar",
				Filter: func(e Event) bool {
					return e == Signal("baz")
				},
			}, // no trigger
		}},
		{events: []Event{Signal("foo"), Signal("bar")}, want: 5, listeners: []*Listener{
			&Listener{
				EventName: "bar",
			}, // 1 trigger
			&Listener{}, // 2 triggers
			&Listener{
				Once: true,
			}, // 1 trigger
			&Listener{
				Filter: func(e Event) bool {
					return e.EventName() == "bar"
				},
			}, // 1 trigger
		}},
	} {
		if got := emitter.emitAndCapture(test.events, test.listeners...); got != test.want {
			t.Fatalf("emitAndCapture(%#v, %#v) = %#v; want %#v", test.events, test.listeners, got, test.want)
		}
	}
}
