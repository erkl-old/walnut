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

	keys := sample.Select("foo").Keys()
	if len(keys) != 2 || keys[0] != "abc" || keys[1] != "def" {
		t.Errorf(`Config.Select("foo").Keys() != []{"abc","def"}`)
	}
}

func TestConfigKeys(t *testing.T) {
	actual := sample.Keys()
	expected := []string{
		"bool",
		"duration",
		"float64",
		"foo.abc",
		"foo.def",
		"int64",
		"string",
		"time",
	}

	if len(actual) != len(expected) {
		t.Fatalf("Config.Keys() -> %v (want %v)", actual, expected)
	}

	for i := 0; i < len(expected); i++ {
		if actual[i] != expected[i] {
			t.Fatalf("Config.Keys() -> %v (want %v)", actual, expected)
		}
	}
}

func TestConfigGet(t *testing.T) {
	v, ok := sample.Get("undefined")
	if v != nil || ok != false {
		t.Fatalf("Config.Get(%q) -> %#v, %#v (want %#v, %#v)",
			"undefined", v, ok, nil, false)
	}

	raw := map[string]interface{}(sample.data)

	for key, want := range raw {
		v, ok := sample.Get(key)
		if v != want || ok != true {
			t.Fatalf("Config.Get(%q) -> %#v, %#v (want %#v, %#v)",
				key, v, ok, want, true)
		}
	}
}

var boolTests = []struct {
	k   string
	v   bool
	err error
}{
	{"undefined", false, fmt.Errorf(_ErrUndefined, "undefined")},
	{"string", false, fmt.Errorf(_ErrWrongType, "string", "string", "bool")},
	{"bool", true, nil},
	{"int64", false, fmt.Errorf(_ErrWrongType, "int64", "int64", "bool")},
	{"float64", false, fmt.Errorf(_ErrWrongType, "float64", "float64", "bool")},
	{"time", false, fmt.Errorf(_ErrWrongType, "time", "time.Time", "bool")},
	{"duration", false, fmt.Errorf(_ErrWrongType, "duration", "time.Duration", "bool")},
}

func TestConfigBool(t *testing.T) {
	for _, test := range boolTests {
		func() {
			defer shouldPanic(t, "Config.Bool", test.k, test.err)
			if v := sample.Bool(test.k); v != test.v {
				t.Errorf("Config.Bool(%q) -> %#v (want %#v)",
					test.k, v, test.v)
			}
		}()
	}
}

var int64Tests = []struct {
	k   string
	v   int64
	err error
}{
	{"undefined", 0, fmt.Errorf(_ErrUndefined, "undefined")},
	{"string", 0, fmt.Errorf(_ErrWrongType, "string", "string", "int64")},
	{"bool", 0, fmt.Errorf(_ErrWrongType, "bool", "bool", "int64")},
	{"int64", 12345, nil},
	{"float64", 0, fmt.Errorf(_ErrWrongType, "float64", "float64", "int64")},
	{"time", 0, fmt.Errorf(_ErrWrongType, "time", "time.Time", "int64")},
	{"duration", 0, fmt.Errorf(_ErrWrongType, "duration", "time.Duration", "int64")},
}

func TestConfigInt64(t *testing.T) {
	for _, test := range int64Tests {
		func() {
			defer shouldPanic(t, "Config.Int64", test.k, test.err)
			if v := sample.Int64(test.k); v != test.v {
				t.Errorf("Config.Int64(%q) -> %#v (want %#v)",
					test.k, v, test.v)
			}
		}()
	}
}

var float64Tests = []struct {
	k   string
	v   float64
	err error
}{
	{"undefined", 0, fmt.Errorf(_ErrUndefined, "undefined")},
	{"string", 0, fmt.Errorf(_ErrWrongType, "string", "string", "float64")},
	{"bool", 0, fmt.Errorf(_ErrWrongType, "bool", "bool", "float64")},
	{"int64", 0, fmt.Errorf(_ErrWrongType, "int64", "int64", "float64")},
	{"float64", 123.45, nil},
	{"time", 0, fmt.Errorf(_ErrWrongType, "time", "time.Time", "float64")},
	{"duration", 0, fmt.Errorf(_ErrWrongType, "duration", "time.Duration", "float64")},
}

func TestConfigFloat64(t *testing.T) {
	for _, test := range float64Tests {
		func() {
			defer shouldPanic(t, "Config.Float64", test.k, test.err)
			if v := sample.Float64(test.k); v != test.v {
				t.Errorf("Config.Float64(%q) -> %#v (want %#v)",
					test.k, v, test.v)
			}
		}()
	}
}

var stringTests = []struct {
	k   string
	v   string
	err error
}{
	{"undefined", "", fmt.Errorf(_ErrUndefined, "undefined")},
	{"string", "hello", nil},
	{"bool", "", fmt.Errorf(_ErrWrongType, "bool", "bool", "string")},
	{"int64", "", fmt.Errorf(_ErrWrongType, "int64", "int64", "string")},
	{"float64", "", fmt.Errorf(_ErrWrongType, "float64", "float64", "string")},
	{"time", "", fmt.Errorf(_ErrWrongType, "time", "time.Time", "string")},
	{"duration", "", fmt.Errorf(_ErrWrongType, "duration", "time.Duration", "string")},
}

func TestConfigString(t *testing.T) {
	for _, test := range stringTests {
		func() {
			defer shouldPanic(t, "Config.String", test.k, test.err)
			if v := sample.String(test.k); v != test.v {
				t.Errorf("Config.String(%q) -> %#v (want %#v)",
					test.k, v, test.v)
			}
		}()
	}
}

var timeTests = []struct {
	k   string
	v   time.Time
	err error
}{
	{"undefined", time.Time{}, fmt.Errorf(_ErrUndefined, "undefined")},
	{"string", time.Time{}, fmt.Errorf(_ErrWrongType, "string", "string", "time.Time")},
	{"bool", time.Time{}, fmt.Errorf(_ErrWrongType, "bool", "bool", "time.Time")},
	{"int64", time.Time{}, fmt.Errorf(_ErrWrongType, "int64", "int64", "time.Time")},
	{"float64", time.Time{}, fmt.Errorf(_ErrWrongType, "float64", "float64", "time.Time")},
	{"time", time.Date(2012, 12, 28, 15, 10, 15, 0, time.UTC), nil},
	{"duration", time.Time{}, fmt.Errorf(_ErrWrongType, "duration", "time.Duration", "time.Time")},
}

func TestConfigTime(t *testing.T) {
	for _, test := range timeTests {
		func() {
			defer shouldPanic(t, "Config.Time", test.k, test.err)
			if v := sample.Time(test.k); v != test.v {
				t.Errorf("Config.Time(%q) -> %#v (want %#v)",
					test.k, v, test.v)
			}
		}()
	}
}

var durationTests = []struct {
	k   string
	v   time.Duration
	err error
}{
	{"undefined", 0, fmt.Errorf(_ErrUndefined, "undefined")},
	{"string", 0, fmt.Errorf(_ErrWrongType, "string", "string", "time.Duration")},
	{"bool", 0, fmt.Errorf(_ErrWrongType, "bool", "bool", "time.Duration")},
	{"int64", 0, fmt.Errorf(_ErrWrongType, "int64", "int64", "time.Duration")},
	{"float64", 0, fmt.Errorf(_ErrWrongType, "float64", "float64", "time.Duration")},
	{"time", 0, fmt.Errorf(_ErrWrongType, "time", "time.Time", "time.Duration")},
	{"duration", 2 * time.Second, nil},
}

func TestConfigDuration(t *testing.T) {
	for _, test := range durationTests {
		func() {
			defer shouldPanic(t, "Config.Duration", test.k, test.err)
			if v := sample.Duration(test.k); v != test.v {
				t.Errorf("Config.Duration(%q) -> %#v (want %#v)",
					test.k, v, test.v)
			}
		}()
	}
}

func shouldPanic(t *testing.T, method, key string, want error) {
	r := recover()
	switch {
	case r == nil && want != nil:
		fallthrough
	case r != nil && !isSameError(r.(error), want):
		t.Errorf(method+"(%q) recover -> %#v (want %#v)", key, r, want)
	}
}

func isSameError(a, b error) bool {
	if a == nil || b == nil {
		return a == b
	}
	return a.Error() == b.Error()
}
