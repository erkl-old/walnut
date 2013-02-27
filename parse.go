package walnut

import (
	"fmt"
	"io/ioutil"
	"strings"
)

// String consisting of all characters to be treated as whitespace.
const Space = " \t\r\n\v\f\u0085\u00a0"

// Parses a configuration file. Panics if reading the file fails, or if
// it contains any syntax errors.
func Load(path string) Config {
	in, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	conf, err := Read(in)
	if err != nil {
		panic(err)
	}

	return conf
}

// Generates a Config instance from a raw configuration file. Returns an
// error if the source contains a syntax error.
func Read(in []byte) (Config, error) {
	// generate a slice of lines from the input, while parsing
	// indentation and discarding empty lines
	lines, err := split(in)
	if err != nil {
		return nil, err
	}

	// reduce the lines to a set of assignments
	assignments, err := interpret(lines)
	if err != nil {
		return nil, err
	}

	// generate the key lookup map, while checking for name conflicts
	table, err := initialize(assignments)
	if err != nil {
		return nil, err
	}

	return &config{"", table}, nil
}

const (
	errIndent   = "illegal indentation on line %d"
	errKey      = "illegal key on line %d"
	errValue    = "illegal value on line %d: %q"
	errConflict = "key %q (line %d) collides with %q (line %d)"
)

type line struct {
	index   int
	depth   int
	content string
}

// Splits a raw input file into lines, discarding all empty lines in the
// process. Returns a non-nil error value if an indentation-based syntax error
// was encountered.
func split(in []byte) ([]line, error) {
	raw := strings.Split(string(in), "\n")

	lines := make([]line, 0)
	indents := make([]string, 0)

	for index, content := range raw {
		// discard empty lines
		if isEmpty(content) {
			continue
		}

		indent, content := selectSpace(content)
		level := len(indents)

		// lines should be 1-indexed
		index++

		// measure the line's indentation depth
		switch {
		case indent == "":
			level = 0
		case len(indents) == 0:
			return nil, fmt.Errorf(errIndent, index)
		default:
			for i, prev := range indents {
				if !strings.HasPrefix(indent, prev) {
					return nil, fmt.Errorf(errIndent, index)
				}
				if len(indent) == len(prev) {
					level = i
					break
				}
			}
		}

		indents = append(indents[:level], indent)
		lines = append(lines, line{index, level, content})
	}

	return lines, nil
}

type assignment struct {
	line    int
	key     string
	literal string
	value   interface{}
}

// Transforms a set of lines to a set of key assignments. Resolves key
// hierarchy and parses values, among other things. Returns a non-nil
// error if any line contains an illegal keys or value.
func interpret(lines []line) ([]assignment, error) {
	output := make([]assignment, 0)
	groups := make([]string, 0)

	for _, line := range lines {
		key, rest := selectKey(line.content)
		groups = append(groups[:line.depth], key)

		if isEmpty(rest) {
			continue
		}

		rest, ok := consumeSeparator(rest)
		if key == "" || !ok && !isEmpty(rest) {
			return nil, fmt.Errorf(errKey, line.index)
		}

		value, ok := parseLiteral(rest)
		if !ok {
			return nil, fmt.Errorf(errValue, line.index, rest)
		}

		output = append(output, assignment{
			line.index, strings.Join(groups, "."), rest, value,
		})
	}

	return output, nil
}

// Generates a map with (key -> value) pairs for each assignment. Also checks
// for key name conflicts.
func initialize(in []assignment) (map[string]interface{}, error) {
	out := make(map[string]interface{})

	for i, a := range in {
		for j, b := range in {
			if i == j {
				continue
			}

			if conflicts(a.key, b.key) {
				// always consider the later of the two lines the culprit
				if a.line < b.line {
					a, b = b, a
				}

				return nil, fmt.Errorf(errConflict,
					a.key, a.line, b.key, b.line)
			}
		}

		out[a.key] = a.value
	}

	return out, nil
}

// Returns true if the key b collides with any part of a.
//     conflicts("foo.bar", "foo")      // -> true
//     conflicts("foo.bar", "foo.bar")  // -> true
//     conflicts("foo.bar", "foo.bars") // -> false
func conflicts(a, b string) bool {
	return strings.HasPrefix(a, b) && (len(a) == len(b) || a[len(b)] == '.')
}

// Returns true if the line is either empty or completely made up of
// whitespace, optionally ending with a comment.
func isEmpty(in string) bool {
	for _, ch := range in {
		switch {
		case ch == '#':
			return true
		case !strings.ContainsRune(Space, ch):
			return false
		}
	}
	return true
}

// Takes all valid key runes from the beginning of input.
func selectSpace(in string) (string, string) {
	for i := 0; i < len(in); i++ {
		if !strings.ContainsRune(Space, rune(in[i])) {
			return in[:i], in[i:]
		}
	}
	return in, ""
}

// Takes all valid key runes from the beginning of input.
func selectKey(in string) (string, string) {
	for i := 0; i < len(in); i++ {
		b := in[i]
		if b == '=' || b == '#' || strings.ContainsRune(Space, rune(b)) {
			return in[:i], in[i:]
		}
	}
	return in, ""
}

// Removes all whitespace and up to one '=' rune from the beginning of the
// input string. The second return value signals whether or not an equal sign
// was encountered in the process.
func consumeSeparator(in string) (string, bool) {
	eq := false
	i := 0

	for ; i < len(in); i++ {
		b := in[i]

		if b == '=' && !eq {
			eq = true
		} else if !strings.ContainsRune(Space, rune(in[i])) {
			break
		}
	}

	return in[i:], eq
}

// Attempts to parse s  as any of the known types.
func parseLiteral(in string) (interface{}, bool) {
	if v, n := readBool(in); n > 0 && isEmpty(in[n:]) {
		return v, true
	}
	if v, n := readInt64(in); n > 0 && isEmpty(in[n:]) {
		return v, true
	}
	if v, n := readFloat64(in); n > 0 && isEmpty(in[n:]) {
		return v, true
	}
	if v, n := readString(in); n > 0 && isEmpty(in[n:]) {
		return v, true
	}
	if v, n := readTime(in); n > 0 && isEmpty(in[n:]) {
		return v, true
	}
	if v, n := readDuration(in); n > 0 && isEmpty(in[n:]) {
		return v, true
	}

	return nil, false
}
