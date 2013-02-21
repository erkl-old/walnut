package walnut

import (
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	_ErrInvalidIndent = "invalid indentation on line %d"
	_ErrInvalidKey    = "invalid key on line %d"
	_ErrInvalidValue  = "unrecognized value on line %d: %q"
	_ErrRedefined     = "%q redefined on line %d (original on line %d)"
)

type definition struct {
	key  string
	val  interface{}
	raw  string
	line int
}

// Parses a configuration file. Panics reading fails, or if the file
// contains a syntax error.
func Load(path string) Config {
	in := make([]byte, 2048)
	r := 0

	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	for {
		if len(in)-r < 512 {
			next := make([]byte, len(in)*2)
			copy(next, in[:r])
			in = next
		}

		n, err := f.Read(in[r:])
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		r += n
	}

	conf, err := Parse(in[:r])
	if err != nil {
		panic(err)
	}

	return conf
}

// Generates a Config instance from a raw configuration file. Returns an
// error if the source contains a syntax error.
func Parse(in []byte) (Config, error) {
	defs, err := parse(in)
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
func parse(in []byte) ([]definition, error) {
	lines := strings.Split(string(in), "\n")

	defs := make([]definition, 0)
	where := make(map[string]int)

	stack := make([]string, 0)
	levels := make([]string, 0)
	allowDeeper := true

	for i, line := range lines {
		if isEmpty(line) {
			continue
		}

		s, k, v := components(line)
		d := depth(levels, s)

		if d < 0 || (d == len(levels) && !allowDeeper) {
			return nil, fmt.Errorf(_ErrInvalidIndent, i+1)
		}

		// trim our history
		if d < len(levels) {
			stack = stack[:d]
			levels = levels[:d]
		}

		stack = append(stack, k)
		levels = append(levels, s)

		// make sure the line specifies a valid key
		if k == "" {
			return nil, fmt.Errorf(_ErrInvalidKey, i+1)
		}

		// does the current line contain an assignment?
		if strings.ContainsRune(line, '=') {
			key := strings.Join(stack, ".")

			// make sure this key doesn't already hold a value
			if prev := where[key]; prev != 0 {
				return nil, fmt.Errorf(_ErrRedefined, key, i+1, prev)
			}

			where[key] = i+1

			parsed, ok := literal(v)
			if !ok {
				return nil, fmt.Errorf(_ErrInvalidValue, i+1, v)
			}

			defs = append(defs, definition{
				key:  key,
				val:  parsed,
				raw:  v,
				line: i + 1,
			})

			allowDeeper = false
			continue
		}

		allowDeeper = true
	}

	return defs, nil
}

// Splits a line into three components: 1) leading whitespace, 2) a key,
// and optionally 3) a raw value.
//
//   components("  foo = bar")  // -> "  ", "foo", "bar"
//   components("foo")          // -> "", "foo", ""
func components(line string) (space, key, value string) {
	for i := 0; i < len(line); i++ {
		if line[i] != ' ' && line[i] != '\t' {
			break
		}
		space = string(line[:i+1])
	}

	if eq := strings.IndexRune(line, '='); eq != -1 {
		key = strings.TrimRight(line[len(space):eq], " \t")
		value = strings.TrimLeft(line[eq+1:], " \t")
	} else {
		key = strings.TrimRight(line[len(space):], " \t")
	}

	return
}

// Traverses a slice of previous indentation levels to see where the subject
// indentation fits in. Returns an integer between 0 and len(context) on
// success, or -1 if subject is not a valid indentation level in this context.
func depth(context []string, subject string) int {
	if subject == "" {
		return 0
	}

	// non-empty indentation without any context is illegal
	if len(context) == 0 {
		return -1
	}

	for i, previous := range context {
		if !strings.HasPrefix(subject, previous) {
			return -1
		}
		if len(subject) == len(previous) {
			return i
		}
	}

	// the subject line is further indented than its parent
	return len(context)
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

// Returns true if the line is completely made up of whitespace, or if the
// line contains only whitespace and a comment rune ('#').
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
