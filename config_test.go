package walnut

import (
	"fmt"
	"testing"
	"time"
)

var sample = Config{
	"string":   "hello",
	"bool":     true,
	"int64":    int64(12345),
	"float64":  float64(123.45),
	"time":     time.Date(2012, 12, 28, 15, 10, 15, 0, time.UTC),
	"duration": 2 * time.Second,
}

func TestConfigKeys(t *testing.T) {
	actual := sample.Keys()
	expected := []string{
		"bool",
		"duration",
		"float64",
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

	raw := map[string]interface{}(sample)

	for key, want := range raw {
		v, ok := sample.Get(key)
		if v != want || ok != true {
			t.Fatalf("Config.Get(%q) -> %#v, %#v (want %#v, %#v)",
				key, v, ok, want, true)
		}
	}
}

var boolTests = []struct {
	key   string
	value bool
	err   error
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
		v, err := sample.Bool(test.key)
		if v != test.value || !isSameError(err, test.err) {
			t.Errorf("Config.Bool(%q) -> %#v, %#v (want %#v, %#v)",
				test.key, v, err, test.value, test.err)
		}
	}
}

var int64Tests = []struct {
	key   string
	value int64
	err   error
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
		v, err := sample.Int64(test.key)
		if v != test.value || !isSameError(err, test.err) {
			t.Errorf("Config.Int64(%q) -> %#v, %#v (want %#v, %#v)",
				test.key, v, err, test.value, test.err)
		}
	}
}

var float64Tests = []struct {
	key   string
	value float64
	err   error
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
		v, err := sample.Float64(test.key)
		if v != test.value || !isSameError(err, test.err) {
			t.Errorf("Config.Float64(%q) -> %#v, %#v (want %#v, %#v)",
				test.key, v, err, test.value, test.err)
		}
	}
}

var stringTests = []struct {
	key   string
	value string
	err   error
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
		v, err := sample.String(test.key)
		if v != test.value || !isSameError(err, test.err) {
			t.Errorf("Config.String(%q) -> %#v, %#v (want %#v, %#v)",
				test.key, v, err, test.value, test.err)
		}
	}
}

var timeTests = []struct {
	key   string
	value time.Time
	err   error
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
		v, err := sample.Time(test.key)
		if !test.value.Equal(v) || !isSameError(err, test.err) {
			t.Errorf("Config.Time(%q) -> %#v, %#v (want %#v, %#v)",
				test.key, v, err, test.value, test.err)
		}
	}
}

var durationTests = []struct {
	key   string
	value time.Duration
	err   error
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
		v, err := sample.Duration(test.key)
		if v != test.value || !isSameError(err, test.err) {
			t.Errorf("Config.Duration(%q) -> %#v, %#v (want %#v, %#v)",
				test.key, v, err, test.value, test.err)
		}
	}
}

// checks two errors for equality
func isSameError(a, b error) bool {
	if a == nil || b == nil {
		return a == b
	}
	return a.Error() == b.Error()
}
