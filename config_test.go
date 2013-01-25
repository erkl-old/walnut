package walnut

import (
	"testing"
	"time"
)

var sample = Config{
	"http.host":       "0.0.0.0",
	"http.port":       int64(8080),
	"greeting.string": "hello",
	"greeting.delay":  2 * time.Second,
	"cake-ratio":      float64(1.0),
	"timestamp":       time.Date(2012, 12, 28, 15, 10, 15, 0, time.UTC),
	"debug-mode":      true,
}

func TestConfigKeys(t *testing.T) {
	actual := sample.Keys()
	expected := []string{
		"cake-ratio",
		"debug-mode",
		"greeting.delay",
		"greeting.string",
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
		t.Fatalf("Config.String(%q) -> %#v, %#v (want %#v, %#v)",
			"undefined", v, err, "", ErrUndefined)
	}

	v, err = sample.String("greeting.delay")
	if v != "" || err != ErrWrongType {
		t.Fatalf("Config.String(%q) -> %#v, %#v (want %#v, %#v)",
			"greeting.delay", v, err, "", ErrWrongType)
	}

	v, err = sample.String("greeting.string")
	if v != "hello" || err != nil {
		t.Fatalf("Config.String(%q) -> %#v, %#v (want %#v, %#v)",
			"greeting.string", v, err, "hello", nil)
	}
}

func TestConfigBool(t *testing.T) {
	v, err := sample.Bool("undefined")
	if v != false || err != ErrUndefined {
		t.Fatalf("Config.Bool(%q) -> %v, %#v (want %#v, %#v)",
			"undefined", v, err, false, ErrUndefined)
	}

	v, err = sample.Bool("cake-ratio")
	if v != false || err != ErrWrongType {
		t.Fatalf("Config.Bool(%q) -> %#v, %#v (want %#v, %#v)",
			"cake-ratio", v, err, false, ErrWrongType)
	}

	v, err = sample.Bool("debug-mode")
	if v != true || err != nil {
		t.Fatalf("Config.Bool(%q) -> %#v, %#v (want %#v, %#v)",
			"debug-mode", v, err, true, nil)
	}
}

func TestConfigInt64(t *testing.T) {
	v, err := sample.Int64("undefined")
	if v != 0 || err != ErrUndefined {
		t.Fatalf("Config.Int64(%q) -> %#v, %#v (want %#v, %#v)",
			"undefined", v, err, 0, ErrUndefined)
	}

	v, err = sample.Int64("greeting.delay")
	if v != 0 || err != ErrWrongType {
		t.Fatalf("Config.Int64(%q) -> %#v, %#v (want %#v, %#v)",
			"greeting.delay", v, err, 0, ErrWrongType)
	}

	v, err = sample.Int64("http.port")
	if v != 8080 || err != nil {
		t.Fatalf("Config.Int64(%q) -> %#v, %#v (want %#v, %#v)",
			"http.port", v, err, 8080, nil)
	}
}

func TestConfigFloat64(t *testing.T) {
	v, err := sample.Float64("undefined")
	if v != 0 || err != ErrUndefined {
		t.Fatalf("Config.Float64(%q) -> %#v, %#v (want %#v, %#v)",
			"undefined", v, err, 0, ErrUndefined)
	}

	v, err = sample.Float64("greeting.delay")
	if v != 0 || err != ErrWrongType {
		t.Fatalf("Config.Float64(%q) -> %#v, %#v (want %#v, %#v)",
			"greeting.delay", v, err, 0, ErrWrongType)
	}

	v, err = sample.Float64("cake-ratio")
	if v != 1.0 || err != nil {
		t.Fatalf("Config.Float64(%q) -> %#v, %#v (want %#v, %#v)",
			"cake-ratio", v, err, 8080, nil)
	}
}

func TestConfigTime(t *testing.T) {
	zero := time.Time{}

	v, err := sample.Time("undefined")
	if v != zero || err != ErrUndefined {
		t.Fatalf("Config.Time(%q) -> %s, %#v (want %s, %#v)",
			"undefined", v, err, 0, ErrUndefined)
	}

	v, err = sample.Time("greeting.delay")
	if v != zero || err != ErrWrongType {
		t.Fatalf("Config.Time(%q) -> %s, %#v (want %s, %#v)",
			"greeting.delay", v, err, 0, ErrWrongType)
	}

	want := time.Date(2012, 12, 28, 15, 10, 15, 0, time.UTC)

	v, err = sample.Time("timestamp")
	if v != want || err != nil {
		t.Fatalf("Config.Time(%q) -> %s, %#v (want %s, %#v)",
			"timestamp", v, err, 8080, nil)
	}
}

func TestConfigDuration(t *testing.T) {
	v, err := sample.Duration("undefined")
	if v != 0 || err != ErrUndefined {
		t.Fatalf("Config.Duration(%q) -> %s, %#v (want %s, %#v)",
			"undefined", v, err, 0, ErrUndefined)
	}

	v, err = sample.Duration("timestamp")
	if v != 0 || err != ErrWrongType {
		t.Fatalf("Config.Duration(%q) -> %s, %#v (want %s, %#v)",
			"timestamp", v, err, 8080, nil)
	}

	v, err = sample.Duration("greeting.delay")
	if v != 2*time.Second || err != nil {
		t.Fatalf("Config.Duration(%q) -> %s, %#v (want %s, %#v)",
			"greeting.delay", v, err, 0, ErrWrongType)
	}
}
