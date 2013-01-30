package walnut

import (
	"fmt"
	"reflect"
	"sort"
	"time"
)

var (
	_ErrUndefined = "%q is not defined"
	_ErrWrongType = "%q is not the right type (is %s, not %s)"
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

// Retrieves a bool value. Will return a non-nil error if the key
// either hasn't been defined or is of a different type.
func (c *Config) Bool(key string) (bool, error) {
	v, ok := (*c)[key]
	if !ok {
		return false, fmt.Errorf(_ErrUndefined, key)
	}

	b, ok := v.(bool)
	if !ok {
		typ := reflect.TypeOf(v).String()
		return false, fmt.Errorf(_ErrWrongType, key, typ, "bool")
	}

	return b, nil
}

// Retrieves an integer value. Will return a non-nil error if the key
// either hasn't been defined or is of a different type.
func (c *Config) Int64(key string) (int64, error) {
	v, ok := (*c)[key]
	if !ok {
		return 0, fmt.Errorf(_ErrUndefined, key)
	}

	i, ok := v.(int64)
	if !ok {
		typ := reflect.TypeOf(v).String()
		return i, fmt.Errorf(_ErrWrongType, key, typ, "int64")
	}

	return i, nil
}

// Retrieves a float value. Will return a non-nil error if the key
// either hasn't been defined or is of a different type.
func (c *Config) Float64(key string) (float64, error) {
	v, ok := (*c)[key]
	if !ok {
		return 0, fmt.Errorf(_ErrUndefined, key)
	}

	f, ok := v.(float64)
	if !ok {
		typ := reflect.TypeOf(v).String()
		return f, fmt.Errorf(_ErrWrongType, key, typ, "float64")
	}

	return f, nil
}

// Retrieves a string value. Will return a non-nil error if the key
// either hasn't been defined or is of a different type.
func (c *Config) String(key string) (string, error) {
	v, ok := (*c)[key]
	if !ok {
		return "", fmt.Errorf(_ErrUndefined, key)
	}

	s, ok := v.(string)
	if !ok {
		typ := reflect.TypeOf(v).String()
		return s, fmt.Errorf(_ErrWrongType, key, typ, "string")
	}

	return s, nil
}

// Retrieves a time value. Will return a non-nil error if the key
// either hasn't been defined or is of a different type.
func (c *Config) Time(key string) (time.Time, error) {
	v, ok := (*c)[key]
	if !ok {
		return time.Time{}, fmt.Errorf(_ErrUndefined, key)
	}

	t, ok := v.(time.Time)
	if !ok {
		typ := reflect.TypeOf(v).String()
		return t, fmt.Errorf(_ErrWrongType, key, typ, "time.Time")
	}

	return t, nil
}

// Retrieves a duration value. Will return a non-nil error if the key
// either hasn't been defined or is of a different type.
func (c *Config) Duration(key string) (time.Duration, error) {
	v, ok := (*c)[key]
	if !ok {
		return time.Duration(0), fmt.Errorf(_ErrUndefined, key)
	}

	d, ok := v.(time.Duration)
	if !ok {
		typ := reflect.TypeOf(v).String()
		return d, fmt.Errorf(_ErrWrongType, key, typ, "time.Duration")
	}

	return d, nil
}
