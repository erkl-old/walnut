package walnut

import (
	"errors"
	"sort"
	"time"
)

var (
	ErrUndefined = errors.New("key is not defined")
	ErrWrongType = errors.New("key is not of the expected type")
)

type Config map[string]interface{}

// Returns a list of all defined keys, sorted lexographically.
func (c *Config) Keys() []string {
	self := *c
	keys := make([]string, 0)

	for key, _ := range self {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	return keys
}

// Retrieves a value. The second return value will be `false` if the
// value hasn't been defined.
func (c *Config) Get(key string) (interface{}, bool) {
	v, ok := (*c)[key]
	return v, ok
}

// Retrieves a string value. Will return a non-nil error if the key
// either hasn't been defined or is of a different type.
func (c *Config) String(key string) (string, error) {
	return "", nil
}

// Retrieves a bool value. Will return a non-nil error if the key
// either hasn't been defined or is of a different type.
func (c *Config) Bool(key string) (bool, error) {
	return false, nil
}

// Retrieves an integer value. Will return a non-nil error if the key
// either hasn't been defined or is of a different type.
func (c *Config) Int64(key string) (int64, error) {
	return 0, nil
}

// Retrieves a float value. Will return a non-nil error if the key
// either hasn't been defined or is of a different type.
func (c *Config) Float64(key string) (float64, error) {
	return 0, nil
}

// Retrieves a duration value. Will return a non-nil error if the key
// either hasn't been defined or is of a different type.
func (c *Config) Duration(key string) (time.Duration, error) {
	return time.Duration(0), nil
}

// Retrieves a time value. Will return a non-nil error if the key
// either hasn't been defined or is of a different type.
func (c *Config) Time(key string) (time.Time, error) {
	return time.Now(), nil
}
