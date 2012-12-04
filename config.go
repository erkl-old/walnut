package walnut

import (
	"time"
)

type Config interface {
	// Returns a list of all keys, sorted lexographically.
	Keys() []string

	// Returns a new Config instance using a subset of the current
	// configuration, selecting all keys with the supplied prefix.
	Prefix(prefix string) Config

	// Loads values from another config object.
	Extend(other Config) Config

	// Assigns a value to a key. Panics if the value is not of a legal type.
	Set(key string, value interface{})

	// Returns the key's value, or an error if the key is undefined.
	Get(key string) (interface{}, error)

	String(key string) (string, error)
	Bool(key string) (bool, error)
	Float(key string) (float64, error)
	Duration(key string) (time.Duration, error)
	Time(key string) (time.Time, error)
}
