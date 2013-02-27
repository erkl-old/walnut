package walnut

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"
)

const (
	errUndefined = "%q is not defined"
	errWrongType = "%q is not the right type (is %s, not %s)"
)

type Config interface {
	// Returns a list of all defined keys, sorted lexographically.
	Keys() []string

	// Selects a subset of the Config. All read operations performed on
	// the returned Config will be prefixed with the prefix.
	Select(prefix string) Config

	// Retrieves an untyped value. The second return value will be false
	// if the value hasn't been defined.
	Get(key string) (interface{}, bool)

	// Retrieves a typed value. Panics if the key doesn't exist, or if its
	// value is of the wrong type.
	Bool(key string) bool
	Int64(key string) int64
	Float64(key string) float64
	Duration(key string) time.Duration
	Time(key string) time.Time
}

// A simple implementation of the Config interface.
type config struct {
	prefix string
	data   map[string]interface{}
}

func (c *config) Keys() []string {
	keys := make([]string, 0)

	for key, _ := range c.data {
		if strings.HasPrefix(key, c.prefix) {
			keys = append(keys, key[len(c.prefix):])
		}
	}

	sort.Strings(keys)

	return keys
}

func (c *config) Select(prefix string) Config {
	return &config{prefix + ".", c.data}
}

func (c *config) Get(key string) (interface{}, bool) {
	v, ok := c.data[c.prefix+key]
	return v, ok
}

func (c *config) Bool(key string) bool {
	v, ok := c.data[c.prefix+key]
	if !ok {
		panic(fmt.Errorf(errUndefined, key))
	}

	b, ok := v.(bool)
	if !ok {
		typ := reflect.TypeOf(v).String()
		panic(fmt.Errorf(errWrongType, key, typ, "bool"))
	}

	return b
}

func (c *config) Int64(key string) int64 {
	v, ok := c.data[c.prefix+key]
	if !ok {
		panic(fmt.Errorf(errUndefined, key))
	}

	i, ok := v.(int64)
	if !ok {
		typ := reflect.TypeOf(v).String()
		panic(fmt.Errorf(errWrongType, key, typ, "int64"))
	}

	return i
}

func (c *config) Float64(key string) float64 {
	v, ok := c.data[c.prefix+key]
	if !ok {
		panic(fmt.Errorf(errUndefined, key))
	}

	f, ok := v.(float64)
	if !ok {
		typ := reflect.TypeOf(v).String()
		panic(fmt.Errorf(errWrongType, key, typ, "float64"))
	}

	return f
}

func (c *config) String(key string) string {
	v, ok := c.data[c.prefix+key]
	if !ok {
		panic(fmt.Errorf(errUndefined, key))
	}

	s, ok := v.(string)
	if !ok {
		typ := reflect.TypeOf(v).String()
		panic(fmt.Errorf(errWrongType, key, typ, "string"))
	}

	return s
}

func (c *config) Time(key string) time.Time {
	v, ok := c.data[c.prefix+key]
	if !ok {
		panic(fmt.Errorf(errUndefined, key))
	}

	t, ok := v.(time.Time)
	if !ok {
		typ := reflect.TypeOf(v).String()
		panic(fmt.Errorf(errWrongType, key, typ, "time.Time"))
	}

	return t
}

func (c *config) Duration(key string) time.Duration {
	v, ok := c.data[c.prefix+key]
	if !ok {
		panic(fmt.Errorf(errUndefined, key))
	}

	d, ok := v.(time.Duration)
	if !ok {
		typ := reflect.TypeOf(v).String()
		panic(fmt.Errorf(errWrongType, key, typ, "time.Duration"))
	}

	return d
}
