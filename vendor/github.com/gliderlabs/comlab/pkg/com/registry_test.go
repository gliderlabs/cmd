package com

import (
	"reflect"
	"testing"
)

type datastore interface {
	Get(key string) interface{}
	Set(key string, value interface{})
}

type searchable interface {
	Search(query string) []interface{}
}

type observer interface {
	ValueChanged(key string, from, to interface{})
}

type enabledContext []string

func (c enabledContext) ComponentEnabled(name string) bool {
	for _, enabled := range c {
		if name == enabled {
			return true
		}
	}
	return false
}

type comAll struct{}

func (c comAll) Get(key string) interface{}                    { return nil }
func (c comAll) Set(key string, value interface{})             {}
func (c comAll) Search(query string) []interface{}             { return nil }
func (c comAll) ValueChanged(key string, from, to interface{}) {}

type comSearchableStore struct{}

func (c comSearchableStore) Get(key string) interface{}        { return nil }
func (c comSearchableStore) Set(key string, value interface{}) {}
func (c comSearchableStore) Search(query string) []interface{} { return nil }

type comObserver struct{}

func (c comObserver) ValueChanged(key string, from, to interface{}) {}

type scenario struct {
	coms     map[string]interface{}
	cfg      map[string]bool
	register []string
}

func (s scenario) newRegistry() registry {
	r := newRegistry()
	if len(s.register) > 0 {
		for _, name := range s.register {
			r.Register(name, s.coms[name])
		}
	} else {
		for k, v := range s.coms {
			r.Register(k, v)
		}
	}
	cfg := mapConfig{}
	for k, v := range s.cfg {
		cfg[k] = v
	}
	r.SetConfig(cfg)
	return r
}

func (s scenario) named(names []string) []interface{} {
	var coms []interface{}
	for _, name := range names {
		com, ok := s.coms[name]
		if !ok {
			continue
		}
		coms = append(coms, com)
	}
	return coms
}

func TestSelect(t *testing.T) {
	t.Parallel()
	scenario1 := scenario{
		coms: map[string]interface{}{
			"com1": new(comAll),
			"com2": new(comObserver),
			"com3": new(comObserver),
		},
		cfg: map[string]bool{},
	}
	scenario2 := scenario{
		coms: map[string]interface{}{
			"com1": new(comObserver),
			"com2": new(comObserver),
		},
		cfg: map[string]bool{
			"com1.enabled": false,
		},
	}
	for _, test := range []struct {
		given   scenario
		inName  string
		inIface interface{}
		want    interface{}
	}{
		{given: scenario1, inName: "", inIface: nil, want: nil},
		{given: scenario1, inName: "com1", inIface: nil, want: scenario1.coms["com1"]},
		{given: scenario1, inName: "com1", inIface: new(searchable), want: scenario1.coms["com1"]},
		{given: scenario1, inName: "com2", inIface: new(datastore), want: nil},
		{given: scenario1, inName: "com2", inIface: nil, want: scenario1.coms["com2"]},

		{given: scenario2, inName: "com1", inIface: nil, want: nil},
		{given: scenario2, inName: "com1", inIface: new(observer), want: nil},
		{given: scenario2, inName: "com0", inIface: nil, want: nil},
	} {
		registry := test.given.newRegistry()
		if got := registry.Select(test.inName, test.inIface); got != test.want {
			t.Fatalf("Select(%#v, %#v) = %#v; want %#v given %#v",
				test.inName, test.inIface, got, test.want, test.given)
		}
	}
}

func TestEnabled(t *testing.T) {
	t.Parallel()
	scenario1 := scenario{
		coms: map[string]interface{}{
			"com1": new(comAll),
			"com2": new(comObserver),
		},
		register: []string{"com1", "com2"},
		cfg:      map[string]bool{},
	}
	scenario2 := scenario{
		coms: map[string]interface{}{
			"com1": new(comObserver),
			"com2": new(comObserver),
			"com3": new(comObserver),
		},
		register: []string{"com1", "com2", "com3"},
		cfg: map[string]bool{
			"com1.enabled": false,
		},
	}
	for _, test := range []struct {
		given   scenario
		inIface interface{}
		inCtx   Context
		want    []string
	}{
		{given: scenario1, inIface: nil, inCtx: nil, want: []string{}},
		{given: scenario1, inIface: new(observer), inCtx: nil, want: []string{"com1", "com2"}},
		{given: scenario1, inIface: new(observer), inCtx: enabledContext{"com2"}, want: []string{"com2"}},
		{given: scenario1, inIface: new(searchable), inCtx: nil, want: []string{"com1"}},

		{given: scenario2, inIface: new(observer), inCtx: nil, want: []string{"com2", "com3"}},
		{given: scenario2, inIface: new(observer), inCtx: enabledContext{"com1", "com2"}, want: []string{"com2"}},
	} {
		registry := test.given.newRegistry()
		want := test.given.named(test.want)
		if got := registry.Enabled(test.inIface, test.inCtx); !reflect.DeepEqual(got, want) {
			t.Fatalf("Enabled(%#v, %#v) = %#v; want %#v given %#v",
				test.inIface, test.inCtx, got, want, test.given)
		}
	}
}
