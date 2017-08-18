package log

import (
	"net/http"
	"os"
	"path"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cast"
)

type EventType int

// Constants used to identify the type of Event received by Observer.
const (
	TypeInfo EventType = iota
	TypeDebug
	TypeLocal
	TypeFatal
)

var defaultLogger = newLogger()

// Event represents a log event, which is what is given to registered
// Observers and the FieldProcessor. Observers use them to actually log,
// FieldProcessor typically adds fields with Append.
//
// Events are treated as immutable so they are passed by value and all operations
// on them will return a modified copy.
type Event struct {
	// EventType enum value
	Type EventType
	// Time of when the event was logged
	Time time.Time
	// Map of key-value fields
	Fields Fields
	// Order of fields added, by key
	Index []string
}

// Append will add a field to Fields and insert new field keys into
// Index. Events are immutable, so Append returns a new Event with the
// field added, which you should use from then on.
func (e Event) Append(key string, value string) Event {
	_, ok := e.Fields[key]
	if ok {
		e.Fields[key] = e.Fields[key] + " " + value
	} else {
		e.Fields[key] = value
		e.Index = append(e.Index, key)
	}
	return e
}

// Remove will remove a field from Fields and Index. Like Append, it
// returns a new Event, as Events are treated as immutable.
func (e Event) Remove(key string) Event {
	delete(e.Fields, key)
	var index []string
	for _, field := range e.Index {
		if field != key {
			index = append(index, field)
		}
	}
	e.Index = index
	return e
}

// Fields is shorthand for string map of strings used by Event for fields, but
// can also be used when logging to explicitly add key-values.
type Fields map[string]string

// FieldProcessor is the function signature for the callback expected by
// SetFieldProcessor. This callback is called for every unknown value
// passed to a logging function. It's given an Event to append to, the field
// value (an argument to Info, Debug, etc), and returns a new Event and whether
// or not it processed the field. Unprocessed fields will be cast to a string
// and appended to a field called "data", unless they are an error.
//
// Fields that are errors are passed to the callback so you have a chance to
// process them into appropriate fields, but if you indicate it was not
// processed, it will create a field called "err" with the string value of its
// Error() method.
//
// There are two types that are handled by the package and not passed to the
// FieldProcessor callback: strings and Fields. Fields are appended directly to
// the Event, and strings are always appended as a field called "msg". Multiple
// strings will be concatenated into the "msg" field in order.
type FieldProcessor func(e Event, field interface{}) (Event, bool)

// Observer is the interface of registerable log observers. They simply receive
// a fully processed Event via the Log method.
type Observer interface {
	Log(e Event)
}

type logger struct {
	sync.Mutex
	debug     bool
	local     bool
	processor FieldProcessor
	observers map[Observer]struct{}
}

func newLogger() *logger {
	return &logger{
		observers: make(map[Observer]struct{}),
	}
}

func (l *logger) log(typ EventType, fields []interface{}) {
	if len(fields) == 0 {
		return
	}
	if typ == TypeDebug && !l.debug {
		return
	}
	if typ == TypeLocal && !l.local {
		return
	}
	e := Event{
		Type:   typ,
		Time:   time.Now(),
		Fields: Fields{"pkg": callerPkg()},
	}
	for _, field := range fields {
		e = l.processField(e, field)
	}
	for o, _ := range l.observers {
		o.Log(e)
	}
}

func (l *logger) processField(e Event, field interface{}) Event {
	switch obj := field.(type) {
	case string:
		return e.Append("msg", obj)
	case Fields:
		for k, v := range obj {
			e = e.Append(k, cast.ToString(v))
		}
		return e
	default:
		if l.processor != nil {
			if ee, ok := l.processor(e, field); ok {
				return ee
			}
		}
		if err, ok := field.(error); ok {
			return e.Append("err", err.Error())
		}
		return e.Append("data", cast.ToString(field))
	}
}

// RegisterObserver adds an Observer to receive processed log events
func RegisterObserver(o Observer) { defaultLogger.RegisterObserver(o) }
func (l *logger) RegisterObserver(o Observer) {
	l.Lock()
	defer l.Unlock()
	l.observers[o] = struct{}{}
}

// UnregisterObserver removes an Observer from receiving log events
func UnregisterObserver(o Observer) { defaultLogger.UnregisterObserver(o) }
func (l *logger) UnregisterObserver(o Observer) {
	l.Lock()
	defer l.Unlock()
	delete(l.observers, o)
}

// SetDebug lets you enable or disable debug mode.
func SetDebug(debug bool) { defaultLogger.SetDebug(debug) }
func (l *logger) SetDebug(debug bool) {
	l.Lock()
	defer l.Unlock()
	l.debug = debug
}

// SetLocal lets you enable or disable local mode
func SetLocal(local bool) { defaultLogger.SetLocal(local) }
func (l *logger) SetLocal(local bool) {
	l.Lock()
	defer l.Unlock()
	l.local = local
}

// SetFieldProcessor lets you specify the field processing callback for your
// application to turn logged values into key-value fields for an event. See
// FieldProcessor for more information.
func SetFieldProcessor(fn FieldProcessor) { defaultLogger.SetFieldProcessor(fn) }
func (l *logger) SetFieldProcessor(fn FieldProcessor) {
	l.Lock()
	defer l.Unlock()
	l.processor = fn
}

// Info logs data intended to be useful to all users in all modes.
func Info(o ...interface{}) { defaultLogger.Info(o...) }
func (l *logger) Info(o ...interface{}) {
	l.Lock()
	defer l.Unlock()
	l.log(TypeInfo, o)
}

// Debug logs data intended to be more verbose for debugging and development.
// Debug only logs when debug mode is enabled.
func Debug(o ...interface{}) { defaultLogger.Debug(o...) }
func (l *logger) Debug(o ...interface{}) {
	l.Lock()
	defer l.Unlock()
	if !l.debug {
		return
	}
	l.log(TypeDebug, o)
}

// Local logs data intended to be seen during local development.
// Local only logs when local mode is enabled.
func Local(o ...interface{}) { defaultLogger.Local(o...) }
func (l *logger) Local(o ...interface{}) {
	l.Lock()
	defer l.Unlock()
	if !l.local {
		return
	}
	l.log(TypeLocal, o)
}

// Fatal logs data that represents a fatal error.
// Fatal will exit with status code 1 after logging.
func Fatal(o ...interface{}) { defaultLogger.Fatal(o...) }
func (l *logger) Fatal(o ...interface{}) {
	l.Lock()
	defer l.Unlock()
	l.log(TypeFatal, o)
	os.Exit(1)
}

func callerPkg() string {
	pc := make([]uintptr, 10)
	runtime.Callers(5, pc)
	f := runtime.FuncForPC(pc[0]).Name()
	base := path.Base(f)
	dir := path.Dir(f)
	dotparts := strings.Split(base, ".")
	pathparts := strings.Split(path.Join(dir, dotparts[0]), "/")
	return pathparts[len(pathparts)-1]
}

// ResponseWriter is the interface for wrapped http.ResponseWriters. Primarily
// it adds Status() and Size() methods (both returning ints).
type ResponseWriter interface {
	loggingResponseWriter
}

// WrapResponseWriter will take an http.ResponseWriter and return it wrapped
// as this package's ResponseWriter.
func WrapResponseWriter(w http.ResponseWriter) ResponseWriter {
	return makeLogger(w)
}
