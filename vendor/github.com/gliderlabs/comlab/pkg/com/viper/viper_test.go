package viper

import (
	"bytes"
	"testing"
)

func testConfig(t *testing.T) *Config {
	cfg := NewConfig()
	cfg.SetConfigType("toml")
	err := cfg.ReadConfig(bytes.NewBufferString(`
		[component1]
		enabled = true
		boolKey = true
		stringKey = "string"
		intKey = 0

		[component2]
		boolKey = false
		stringKey = ""
		intKey = 100

		[component3]
		enabled = false
		intKey = -100
	`))
	if err != nil {
		t.Fatal(err)
	}
	return cfg
}

func TestComponentEnabled(t *testing.T) {
	t.Parallel()
	cfg := testConfig(t)
	for _, test := range []struct {
		in   string
		want bool
	}{
		{in: "component1", want: true},
		{in: "component2", want: true},
		{in: "component3", want: false},
		{in: "noconfig", want: true},
	} {
		if got := cfg.ComponentEnabled(test.in); got != test.want {
			t.Fatalf("ComponentEnabled(%#v) = %#v; want %#v", test.in, got, test.want)
		}
	}
}

func TestGetString(t *testing.T) {
	t.Parallel()
	cfg := testConfig(t)
	for _, test := range []struct {
		in      string
		wantVal string
		wantSet bool
	}{
		{in: "component1.stringKey", wantVal: "string", wantSet: true},
		{in: "component1.nonKey", wantVal: "", wantSet: false},
		{in: "component2.stringKey", wantVal: "", wantSet: true},
		{in: "component3.intKey", wantVal: "-100", wantSet: true},
		{in: "noconfig.nonKey", wantVal: "", wantSet: false},
	} {
		if gotVal, gotSet := cfg.GetString(test.in); gotVal != test.wantVal || gotSet != test.wantSet {
			t.Fatalf("GetString(%#v) = %#v, %#v; want %#v, %#v", test.in, gotVal, gotSet, test.wantVal, test.wantSet)
		}
	}
}

func TestGetInt(t *testing.T) {
	t.Parallel()
	cfg := testConfig(t)
	for _, test := range []struct {
		in      string
		wantVal int
		wantSet bool
	}{
		{in: "component1.intKey", wantVal: 0, wantSet: true},
		{in: "component1.nonKey", wantVal: 0, wantSet: false},
		{in: "component2.intKey", wantVal: 100, wantSet: true},
		{in: "component3.intKey", wantVal: -100, wantSet: true},
		{in: "component1.boolKey", wantVal: 1, wantSet: true},
		{in: "noconfig.nonKey", wantVal: 0, wantSet: false},
	} {
		if gotVal, gotSet := cfg.GetInt(test.in); gotVal != test.wantVal || gotSet != test.wantSet {
			t.Fatalf("GetInt(%#v) = %#v, %#v; want %#v, %#v", test.in, gotVal, gotSet, test.wantVal, test.wantSet)
		}
	}
}

func TestGetBool(t *testing.T) {
	t.Parallel()
	cfg := testConfig(t)
	for _, test := range []struct {
		in      string
		wantVal bool
		wantSet bool
	}{
		{in: "component1.boolKey", wantVal: true, wantSet: true},
		{in: "component1.nonKey", wantVal: false, wantSet: false},
		{in: "component2.boolKey", wantVal: false, wantSet: true},
		{in: "component1.intKey", wantVal: false, wantSet: true},
		{in: "noconfig.nonKey", wantVal: false, wantSet: false},
	} {
		if gotVal, gotSet := cfg.GetBool(test.in); gotVal != test.wantVal || gotSet != test.wantSet {
			t.Fatalf("GetBool(%#v) = %#v, %#v; want %#v, %#v", test.in, gotVal, gotSet, test.wantVal, test.wantSet)
		}
	}
}
