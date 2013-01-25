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
	v, ok := (*c)[key]
	if !ok {
		return "", ErrUndefined
	}

	s, ok := v.(string)
	if !ok {
		return "", ErrWrongType
	}

	return s, nil
}

// Retrieves a bool value. Will return a non-nil error if the key
// either hasn't been defined or is of a different type.
func (c *Config) Bool(key string) (bool, error) {
	v, ok := (*c)[key]
	if !ok {
		return false, ErrUndefined
	}

	b, ok := v.(bool)
	if !ok {
		return false, ErrWrongType
	}

	return b, nil
}

// Retrieves an integer value. Will return a non-nil error if the key
// either hasn't been defined or is of a different type.
func (c *Config) Int64(key string) (int64, error) {
	v, ok := (*c)[key]
	if !ok {
		return 0, ErrUndefined
	}

	i, ok := v.(int64)
	if !ok {
		return 0, ErrWrongType
	}

	return i, nil
}

// Retrieves a float value. Will return a non-nil error if the key
// either hasn't been defined or is of a different type.
func (c *Config) Float64(key string) (float64, error) {
	v, ok := (*c)[key]
	if !ok {
		return 0, ErrUndefined
	}

	f, ok := v.(float64)
	if !ok {
		return 0, ErrWrongType
	}

	return f, nil
}

// Retrieves a duration value. Will return a non-nil error if the key
// either hasn't been defined or is of a different type.
func (c *Config) Duration(key string) (time.Duration, error) {
	v, ok := (*c)[key]
	if !ok {
		return time.Duration(0), ErrUndefined
	}

	d, ok := v.(time.Duration)
	if !ok {
		return time.Duration(0), ErrWrongType
	}

	return d, nil
}

// Retrieves a time value. Will return a non-nil error if the key
// either hasn't been defined or is of a different type.
func (c *Config) Time(key string) (time.Time, error) {
	v, ok := (*c)[key]
	if !ok {
		return time.Time{}, ErrUndefined
	}

	t, ok := v.(time.Time)
	if !ok {
		return time.Time{}, ErrWrongType
	}

	return t, nil
}
