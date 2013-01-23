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
}

func TestConfigKeys(t *testing.T) {
	actual := sample.Keys()
	expected := []string{
		"cake-ratio",
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
		t.Fatalf("Config.Get(%q) -> %v, %v (want %v, %v)", v, ok, nil, false)
	}

	v, ok = sample.Get("cake-ratio")
	if v.(float64) != 1.0 || ok == false {
		t.Fatalf("Config.Get(%q) -> %v, %v (want %v, %v)",
			v, ok, float64(1.0), true)
	}
}
