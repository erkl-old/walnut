package walnut

import (
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
