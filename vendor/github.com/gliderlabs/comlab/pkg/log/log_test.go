package log

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/spf13/cast"
)

func zeroTime(e Event) Event {
	e.Time = time.Time{}
	return e
}

type logBuffer []Event

func (l *logBuffer) Log(e Event) {
	*l = append(*l, zeroTime(e))
}

func (l *logBuffer) Types() []EventType {
	var types []EventType
	for _, e := range *l {
		types = append(types, e.Type)
	}
	return types
}

type args []interface{}

type extendedError struct {
	Message string
	Details string
	LineNum int
}

func (e extendedError) Error() string {
	return e.Message
}

func TestBuiltinFieldProcessing(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		in   args
		want Fields
	}{
		{in: args{"hello"}, want: Fields{"msg": "hello", "pkg": "testing"}},
		{in: args{"hello", "world"}, want: Fields{"msg": "hello world", "pkg": "testing"}},
		{in: args{"hello", true, "world"}, want: Fields{"msg": "hello world", "data": "true", "pkg": "testing"}},
		{in: args{50.1, 0, 10000}, want: Fields{"data": "50.1 0 10000", "pkg": "testing"}},
		{in: args{Fields{"foo": "bar", "baz": "qux"}}, want: Fields{"foo": "bar", "baz": "qux", "pkg": "testing"}},
		{in: args{Fields{"foo": "bar"}, Fields{"baz": "qux"}}, want: Fields{"foo": "bar", "baz": "qux", "pkg": "testing"}},
		{in: args{fmt.Errorf("an error")}, want: Fields{"err": "an error", "pkg": "testing"}},
		{in: args{fmt.Errorf("not very"), fmt.Errorf("realistic")}, want: Fields{"err": "not very realistic", "pkg": "testing"}},
		{in: args{"a message", fmt.Errorf("and error")}, want: Fields{"err": "and error", "msg": "a message", "pkg": "testing"}},
	} {
		log := newLogger()
		buffer := &logBuffer{}
		log.RegisterObserver(buffer)
		log.Info(test.in...)
		if got := (*buffer)[0].Fields; !reflect.DeepEqual(got, test.want) {
			t.Fatalf("Info(%#v) => %#v; want %#v", test.in, got, test.want)
		}
	}
}

func TestExtendedFieldProcessing(t *testing.T) {
	t.Parallel()
	fieldProcessor := func(e Event, field interface{}) (Event, bool) {
		switch obj := field.(type) {
		case int:
			return e.Append("integer", cast.ToString(obj)), true
		case error:
			if err, ok := obj.(extendedError); ok {
				e = e.Append("err", err.Message)
				e = e.Append("details", err.Details)
				e = e.Append("line", cast.ToString(err.LineNum))
				return e, true
			}
		case time.Time:
			return e.Append("time", obj.Format("15:04:05.000")), true
		}
		return e, false
	}
	for _, test := range []struct {
		in   args
		want Fields
	}{
		{in: args{"hello", 25}, want: Fields{"msg": "hello", "integer": "25", "pkg": "testing"}},
		{in: args{fmt.Errorf("regular error")}, want: Fields{"err": "regular error", "pkg": "testing"}},
		{in: args{extendedError{"extended error", "more details", 25}}, want: Fields{"err": "extended error", "details": "more details", "line": "25", "pkg": "testing"}},
		{in: args{time.Time{}}, want: Fields{"time": "00:00:00.000", "pkg": "testing"}},
	} {
		log := newLogger()
		buffer := &logBuffer{}
		log.RegisterObserver(buffer)
		log.SetFieldProcessor(fieldProcessor)
		log.Info(test.in...)
		if got := (*buffer)[0].Fields; !reflect.DeepEqual(got, test.want) {
			t.Fatalf("Info(%#v) => %#v; want %#v", test.in, got, test.want)
		}
	}
}

func TestLogTypes(t *testing.T) {
	t.Parallel()
	regularLogger := newLogger()
	debugLogger := newLogger()
	debugLogger.SetDebug(true)
	localLogger := newLogger()
	localLogger.SetLocal(true)
	localDebugLogger := newLogger()
	localDebugLogger.SetLocal(true)
	localDebugLogger.SetDebug(true)
	for _, test := range []struct {
		given *logger
		in    []EventType
		want  []EventType
	}{
		{given: regularLogger, in: []EventType{TypeInfo, TypeDebug, TypeInfo, TypeLocal, TypeFatal}, want: []EventType{TypeInfo, TypeInfo, TypeFatal}},
		{given: debugLogger, in: []EventType{TypeInfo, TypeDebug, TypeInfo, TypeLocal, TypeFatal}, want: []EventType{TypeInfo, TypeDebug, TypeInfo, TypeFatal}},
		{given: localLogger, in: []EventType{TypeInfo, TypeDebug, TypeInfo, TypeLocal, TypeFatal}, want: []EventType{TypeInfo, TypeInfo, TypeLocal, TypeFatal}},
		{given: localDebugLogger, in: []EventType{TypeInfo, TypeDebug, TypeInfo, TypeLocal, TypeFatal}, want: []EventType{TypeInfo, TypeDebug, TypeInfo, TypeLocal, TypeFatal}},
	} {
		buffer := &logBuffer{}
		test.given.RegisterObserver(buffer)
		for _, typ := range test.in {
			test.given.log(typ, []interface{}{"hello"})
		}
		test.given.UnregisterObserver(buffer)
		if got := buffer.Types(); !reflect.DeepEqual(got, test.want) {
			t.Fatalf("logTypes(%#v) => %#v; want %#v given %#v", test.in, got, test.want, test.given)
		}
	}
}

func BenchmarkLog(b *testing.B) {
	log := newLogger()
	for i := 0; i < b.N; i++ {
		log.Info("hello", errors.New("error"))
	}
}

func BenchmarkPrintln(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fmt.Sprintln("hello", errors.New("error"))
	}
}
