package walnut

import (
	"fmt"
	"strings"
)

const (
	_ErrInvalidIndent = "invalid indentation on line %d"
	_ErrInvalidValue  = "unrecognized value on line %d: %q"
)

type definition struct {
	key  string
	val  interface{}
	raw  string
	line int
}

// Generates a Config instance from a raw configuration file. Returns an
// error if the source contains a syntax error.
func Parse(in []byte) (Config, error) {
	defs, err := definitions(in)
	if err != nil {
		return nil, err
	}

	m := make(map[string]interface{})
	for _, def := range defs {
		m[def.key] = def.val
	}

	return Config(m), nil
}

// Generates a set of definitions from a raw configuration file. Returns an
// error if the source contains a syntax error.
func definitions(in []byte) ([]definition, error) {
	defs := make([]definition, 0)

	stack := make([]string, 0)
	levels := make([]string, 0)
	isLeaf := false

	lines := strings.Split(string(in), "\n")

	for i, line := range lines {
		if isEmpty(line) {
			continue
		}

		s, k, v := components(line)
		d := depth(levels, s)

		if d < 0 || (d == len(levels) && isLeaf) {
			return nil, fmt.Errorf(_ErrInvalidIndent, i+1)
		}

		// trim redundant indentation levels
		if d < len(levels) {
			stack = stack[:d]
			levels = levels[:d]
		}

		stack = append(stack, k)
		levels = append(levels, s)

		// contains an assignment
		if v != "" {
			value, ok := literal(v)
			if !ok {
				return nil, fmt.Errorf(_ErrInvalidValue, i+1, v)
			}

			defs = append(defs, definition{
				key:  strings.Join(stack, "."),
				val:  value,
				raw:  v,
				line: i + 1,
			})

			isLeaf = true
			continue
		}

		isLeaf = false
	}

	return defs, nil
}

// Splits a line into three components; 1) leading whitespace, 2) a key
// string, and 3) a raw value string.
func components(line string) (space, key, value string) {
	for i := 0; i < len(line); i++ {
		if line[i] != ' ' && line[i] != '\t' {
			break
		}
		space = string(line[:i+1])
	}

	if eq := strings.IndexRune(line, '='); eq != -1 {
		key = string(line[len(space):eq])
		value = string(line[eq+1:])
	} else {
		key = string(line[len(space):])
	}

	return
}

// ...
func depth(parents []string, current string) int {
	// ...
	if len(current) == 0 {
		return 0
	}

	// ...
	if len(parents) == 0 {
		return -1
	}

	for i, parent := range parents {
		if !strings.HasPrefix(current, parent) {
			return -1
		}
		if len(current) == len(parent) {
			return i
		}
	}

	// the current line is further indented than its parent
	return len(parents)
}

// Tries to extract a literal value from a string.
func literal(s string) (interface{}, bool) {
	if v, n := readBool(s); n != 0 && isEmpty(s[n:]) {
		return v, true
	}
	if v, n := readInt64(s); n != 0 && isEmpty(s[n:]) {
		return v, true
	}
	if v, n := readFloat64(s); n != 0 && isEmpty(s[n:]) {
		return v, true
	}
	if v, n := readString(s); n != 0 && isEmpty(s[n:]) {
		return v, true
	}
	if v, n := readTime(s); n != 0 && isEmpty(s[n:]) {
		return v, true
	}
	if v, n := readDuration(s); n != 0 && isEmpty(s[n:]) {
		return v, true
	}

	return nil, false
}

// Returns true if the line is completely made up of whitespace, or if all
// non-whitespace characters appear after the comment rune ('#').
func isEmpty(s string) bool {
	for i := 0; i < len(s); i++ {
		b := s[i]

		if b == ' ' || b == '\t' || b == '\r' {
			continue
		}
		if b == '#' {
			break
		}

		return false
	}

	return true
}
