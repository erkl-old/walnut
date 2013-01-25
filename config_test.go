package walnut

import (
	"testing"
	"time"
)

var sample = Config{
	"http.host":    "0.0.0.0",
	"http.port":    int64(8080),
	"greet.string": "hello",
	"greet.delay":  2 * time.Second,
	"cake-ratio":   float64(1.0),
	"timestamp":    time.Date(2012, 12, 28, 15, 10, 15, 0, time.UTC),
	"debug-mode":   true,
}

func TestConfigKeys(t *testing.T) {
	actual := sample.Keys()
	expected := []string{
		"cake-ratio",
		"debug-mode",
		"greet.delay",
		"greet.string",
		"http.host",
		"http.port",
		"timestamp",
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
		t.Fatalf("Config.Get(%q) -> %v, %v (want %v, %v)",
			"undefined", v, ok, nil, false)
	}

	v, ok = sample.Get("cake-ratio")
	if v.(float64) != 1.0 || ok != true {
		t.Fatalf("Config.Get(%q) -> %v, %v (want %v, %v)",
			"cake-ratio", v, ok, float64(1.0), true)
	}
}

func TestConfigString(t *testing.T) {
	v, err := sample.String("undefined")
	if v != "" || err != ErrUndefined {
		t.Fatalf("Config.String(%q) -> %q, %#v (want %q, %#v)",
			"undefined", v, err, "", ErrUndefined)
	}

	v, err = sample.String("greet.delay")
	if v != "" || err != ErrWrongType {
		t.Fatalf("Config.String(%q) -> %q, %#v (want %q, %#v)",
			"greet.delay", v, err, "", ErrWrongType)
	}

	v, err = sample.String("greet.string")
	if v != "hello" || err != nil {
		t.Fatalf("Config.String(%q) -> %q, %#v (want %q, %#v)",
			"greet.string", v, err, "hello", nil)
	}
}

func TestConfigBool(t *testing.T) {
	v, err := sample.Bool("undefined")
	if v != false || err != ErrUndefined {
		t.Fatalf("Config.Bool(%q) -> %q, %#v (want %q, %#v)",
			"undefined", v, err, false, ErrUndefined)
	}

	v, err = sample.Bool("cake-ratio")
	if v != false || err != ErrWrongType {
		t.Fatalf("Config.Bool(%q) -> %q, %#v (want %q, %#v)",
			"cake-ratio", v, err, false, ErrWrongType)
	}

	v, err = sample.Bool("debug-mode")
	if v != true || err != nil {
		t.Fatalf("Config.Bool(%q) -> %q, %#v (want %q, %#v)",
			"debug-mode", v, err, true, nil)
	}
}
