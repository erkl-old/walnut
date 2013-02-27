package walnut

import (
	"fmt"
	"testing"
	"time"
)

var sample = &config{
	"",
	map[string]interface{}{
		"string":   "hello",
		"bool":     true,
		"int64":    int64(12345),
		"float64":  float64(123.45),
		"time":     time.Date(2012, 12, 28, 15, 10, 15, 0, time.UTC),
		"duration": 2 * time.Second,
		"foo.def":  "hello",
		"foo.abc":  "bye",
	},
}

func TestConfigSelect(t *testing.T) {
	v, _ := sample.Select("foo").Get("def")
	if v != sample.data["foo.def"] {
		t.Errorf(`Config.Select("foo").Get("def") != Config.Get("foo.def")`)
	}

	got := sample.Select("foo").Keys()
	want := []string{"abc", "def"}

	if !eq(got, want) {
		t.Errorf("sample.Select(\"foo\").Keys():")
		t.Errorf("   got %v", got)
		t.Errorf("  want %v", want)
	}
}

func TestConfigKeys(t *testing.T) {
	got := sample.Keys()
	want := []string{
		"bool",
		"duration",
		"float64",
		"foo.abc",
		"foo.def",
		"int64",
		"string",
		"time",
	}

	if !eq(got, want) {
		t.Errorf("sample.Keys():")
		t.Errorf("   got %v", got)
		t.Errorf("  want %v", want)
	}
}

func TestConfigGet(t *testing.T) {
	got, ok := sample.Get("undefined")
	if got != nil || ok != false {
		t.Errorf("sample.Get(\"undefined\"):")
		t.Errorf("   got %#v, %#v", got, ok)
		t.Errorf("   got %#v, %#v", nil, false)
	}

	for key, want := range sample.data {
		got, ok := sample.Get(key)
		if got != want || ok != true {
			t.Errorf("sample.Get(%q):", key)
			t.Errorf("   got %#v, %#v", got, ok)
			t.Errorf("  want %#v, %#v", want, true)
		}
	}
}

var boolTests = []struct {
	key  string
	want bool
	err  error
}{
	{"undefined", false, fmt.Errorf(errUndefined, "undefined")},
	{"string", false, fmt.Errorf(errWrongType, "string", "string", "bool")},
	{"bool", true, nil},
	{"int64", false, fmt.Errorf(errWrongType, "int64", "int64", "bool")},
	{"float64", false, fmt.Errorf(errWrongType, "float64", "float64", "bool")},
	{"time", false, fmt.Errorf(errWrongType, "time", "time.Time", "bool")},
	{"duration", false, fmt.Errorf(errWrongType, "duration", "time.Duration", "bool")},
}

func TestConfigBool(t *testing.T) {
	for _, test := range boolTests {
		func() {
			defer shouldPanic(t, "Config.Bool", test.key, test.err)
			if got := sample.Bool(test.key); got != test.want {
				t.Errorf("Config.Bool(%q):", test.key)
				t.Errorf("   got %#v", got)
				t.Errorf("  want %#v", test.want)
			}
		}()
	}
}

var int64Tests = []struct {
	key  string
	want int64
	err  error
}{
	{"undefined", 0, fmt.Errorf(errUndefined, "undefined")},
	{"string", 0, fmt.Errorf(errWrongType, "string", "string", "int64")},
	{"bool", 0, fmt.Errorf(errWrongType, "bool", "bool", "int64")},
	{"int64", 12345, nil},
	{"float64", 0, fmt.Errorf(errWrongType, "float64", "float64", "int64")},
	{"time", 0, fmt.Errorf(errWrongType, "time", "time.Time", "int64")},
	{"duration", 0, fmt.Errorf(errWrongType, "duration", "time.Duration", "int64")},
}

func TestConfigInt64(t *testing.T) {
	for _, test := range int64Tests {
		func() {
			defer shouldPanic(t, "Config.Int64", test.key, test.err)
			if got := sample.Int64(test.key); got != test.want {
				t.Errorf("Config.Int64(%q):", test.key)
				t.Errorf("   got %#v", got)
				t.Errorf("  want %#v", test.want)
			}
		}()
	}
}

var float64Tests = []struct {
	key  string
	want float64
	err  error
}{
	{"undefined", 0, fmt.Errorf(errUndefined, "undefined")},
	{"string", 0, fmt.Errorf(errWrongType, "string", "string", "float64")},
	{"bool", 0, fmt.Errorf(errWrongType, "bool", "bool", "float64")},
	{"int64", 0, fmt.Errorf(errWrongType, "int64", "int64", "float64")},
	{"float64", 123.45, nil},
	{"time", 0, fmt.Errorf(errWrongType, "time", "time.Time", "float64")},
	{"duration", 0, fmt.Errorf(errWrongType, "duration", "time.Duration", "float64")},
}

func TestConfigFloat64(t *testing.T) {
	for _, test := range float64Tests {
		func() {
			defer shouldPanic(t, "Config.Float64", test.key, test.err)
			if got := sample.Float64(test.key); got != test.want {
				t.Errorf("Config.Float64(%q):", test.key)
				t.Errorf("   got %#v", got)
				t.Errorf("  want %#v", test.want)
			}
		}()
	}
}

var stringTests = []struct {
	key  string
	want string
	err  error
}{
	{"undefined", "", fmt.Errorf(errUndefined, "undefined")},
	{"string", "hello", nil},
	{"bool", "", fmt.Errorf(errWrongType, "bool", "bool", "string")},
	{"int64", "", fmt.Errorf(errWrongType, "int64", "int64", "string")},
	{"float64", "", fmt.Errorf(errWrongType, "float64", "float64", "string")},
	{"time", "", fmt.Errorf(errWrongType, "time", "time.Time", "string")},
	{"duration", "", fmt.Errorf(errWrongType, "duration", "time.Duration", "string")},
}

func TestConfigString(t *testing.T) {
	for _, test := range stringTests {
		func() {
			defer shouldPanic(t, "Config.String", test.key, test.err)
			if got := sample.String(test.key); got != test.want {
				t.Errorf("Config.String(%q):", test.key)
				t.Errorf("   got %#v", got)
				t.Errorf("  want %#v", test.want)
			}
		}()
	}
}

var timeTests = []struct {
	key  string
	want time.Time
	err  error
}{
	{"undefined", time.Time{}, fmt.Errorf(errUndefined, "undefined")},
	{"string", time.Time{}, fmt.Errorf(errWrongType, "string", "string", "time.Time")},
	{"bool", time.Time{}, fmt.Errorf(errWrongType, "bool", "bool", "time.Time")},
	{"int64", time.Time{}, fmt.Errorf(errWrongType, "int64", "int64", "time.Time")},
	{"float64", time.Time{}, fmt.Errorf(errWrongType, "float64", "float64", "time.Time")},
	{"time", time.Date(2012, 12, 28, 15, 10, 15, 0, time.UTC), nil},
	{"duration", time.Time{}, fmt.Errorf(errWrongType, "duration", "time.Duration", "time.Time")},
}

func TestConfigTime(t *testing.T) {
	for _, test := range timeTests {
		func() {
			defer shouldPanic(t, "Config.Time", test.key, test.err)
			if got := sample.Time(test.key); got != test.want {
				t.Errorf("Config.Time(%q):", test.key)
				t.Errorf("   got %#v", got)
				t.Errorf("  want %#v", test.want)
			}
		}()
	}
}

var durationTests = []struct {
	key  string
	want time.Duration
	err  error
}{
	{"undefined", 0, fmt.Errorf(errUndefined, "undefined")},
	{"string", 0, fmt.Errorf(errWrongType, "string", "string", "time.Duration")},
	{"bool", 0, fmt.Errorf(errWrongType, "bool", "bool", "time.Duration")},
	{"int64", 0, fmt.Errorf(errWrongType, "int64", "int64", "time.Duration")},
	{"float64", 0, fmt.Errorf(errWrongType, "float64", "float64", "time.Duration")},
	{"time", 0, fmt.Errorf(errWrongType, "time", "time.Time", "time.Duration")},
	{"duration", 2 * time.Second, nil},
}

func TestConfigDuration(t *testing.T) {
	for _, test := range durationTests {
		func() {
			defer shouldPanic(t, "Config.Duration", test.key, test.err)
			if got := sample.Duration(test.key); got != test.want {
				t.Errorf("Config.Duration(%q):", test.key)
				t.Errorf("   got %#v", got)
				t.Errorf("  want %#v", test.want)
			}
		}()
	}
}

func shouldPanic(t *testing.T, method, key string, want error) {
	r := recover()
	switch {
	case r == nil && want != nil:
		fallthrough
	case r != nil && !eq(r.(error), want):
		t.Errorf(method+"(%q):", key)
		t.Errorf("  recovered %v", r)
		t.Errorf("       want %v", want)
	}
}
