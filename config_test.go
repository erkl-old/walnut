package walnut

import (
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
	// @todo: .(map[string]interface{})

	v, ok := sample.Get("undefined")
	if v != nil || ok != false {
		t.Fatalf("Config.Get(%q) -> %v, %v (want %v, %v)",
			"undefined", v, ok, nil, false)
	}

	v, ok = sample.Get("float64")
	if v.(float64) != 123.45 || ok != true {
		t.Fatalf("Config.Get(%q) -> %#v, %#v (want %#v, %#v)",
			"float64", v, ok, float64(1.0), true)
	}
}

var stringTests = []struct {
	key   string
	value string
	err   error
}{
	{"undefined", "", ErrUndefined},
	{"string", "hello", nil},
	{"bool", "", ErrWrongType},
	{"int64", "", ErrWrongType},
	{"float64", "", ErrWrongType},
	{"time", "", ErrWrongType},
	{"duration", "", ErrWrongType},
}

func TestConfigString(t *testing.T) {
	for _, test := range stringTests {
		v, err := sample.String(test.key)
		if v != test.value || err != test.err {
			t.Errorf("Config.String(%q) -> %#v, %#v (want %#v, %#v)",
				test.key, v, err, test.value, test.err)
		}
	}
}

var boolTests = []struct {
	key   string
	value bool
	err   error
}{
	{"undefined", false, ErrUndefined},
	{"string", false, ErrWrongType},
	{"bool", true, nil},
	{"int64", false, ErrWrongType},
	{"float64", false, ErrWrongType},
	{"time", false, ErrWrongType},
	{"duration", false, ErrWrongType},
}

func TestConfigBool(t *testing.T) {
	for _, test := range boolTests {
		v, err := sample.Bool(test.key)
		if v != test.value || err != test.err {
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
	{"undefined", 0, ErrUndefined},
	{"string", 0, ErrWrongType},
	{"bool", 0, ErrWrongType},
	{"int64", 12345, nil},
	{"float64", 0, ErrWrongType},
	{"time", 0, ErrWrongType},
	{"duration", 0, ErrWrongType},
}

func TestConfigInt64(t *testing.T) {
	for _, test := range int64Tests {
		v, err := sample.Int64(test.key)
		if v != test.value || err != test.err {
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
	{"undefined", 0, ErrUndefined},
	{"string", 0, ErrWrongType},
	{"bool", 0, ErrWrongType},
	{"int64", 0, ErrWrongType},
	{"float64", 123.45, nil},
	{"time", 0, ErrWrongType},
	{"duration", 0, ErrWrongType},
}

func TestConfigFloat64(t *testing.T) {
	for _, test := range float64Tests {
		v, err := sample.Float64(test.key)
		if v != test.value || err != test.err {
			t.Errorf("Config.Float64(%q) -> %#v, %#v (want %#v, %#v)",
				test.key, v, err, test.value, test.err)
		}
	}
}

var timeTests = []struct {
	key   string
	value time.Time
	err   error
}{
	{"undefined", time.Time{}, ErrUndefined},
	{"string", time.Time{}, ErrWrongType},
	{"bool", time.Time{}, ErrWrongType},
	{"int64", time.Time{}, ErrWrongType},
	{"float64", time.Time{}, ErrWrongType},
	{"time", time.Date(2012, 12, 28, 15, 10, 15, 0, time.UTC), nil},
	{"duration", time.Time{}, ErrWrongType},
}

func TestConfigTime(t *testing.T) {
	for _, test := range timeTests {
		v, err := sample.Time(test.key)
		if v != test.value || err != test.err {
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
	{"undefined", 0, ErrUndefined},
	{"string", 0, ErrWrongType},
	{"bool", 0, ErrWrongType},
	{"int64", 0, ErrWrongType},
	{"float64", 0, ErrWrongType},
	{"time", 0, ErrWrongType},
	{"duration", 2 * time.Second, nil},
}

func TestConfigDuration(t *testing.T) {
	for _, test := range durationTests {
		v, err := sample.Duration(test.key)
		if v != test.value || err != test.err {
			t.Errorf("Config.Duration(%q) -> %#v, %#v (want %#v, %#v)",
				test.key, v, err, test.value, test.err)
		}
	}
}
