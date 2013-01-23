package walnut

import (
	"bytes"
	"fmt"
	"io"
)

const (
	_ErrTooLong       = "line %d too long (buffer overflow thwarted)"
	_ErrInvalidIndent = "invalid indentation on line %d"
	_ErrInvalidValue  = "unrecognized value on line %d: %q"
)

// Defines a "key = value" assignment.
type def struct {
	key   string
	value string
	line  int
}

// Generates a Config from an io.Reader instance.
func Parse(r io.Reader, lineLimit int) (Config, error) {
	defs, err := parse(r, make([]byte, lineLimit))
	if err != nil {
		return nil, err
	}

	conf := make(Config)

	for _, d := range defs {
		v, ok := parseLiteral(d.value)
		if !ok {
			return nil, fmt.Errorf(_ErrInvalidValue, d.line, d.value)
		}

		conf[d.key] = v
	}

	return conf, nil
}

// Processes a string, extracting a literal value if one can be found.
func parseLiteral(s string) (interface{}, bool) {
	b := []byte(s)

	if v, n := readBool(b); n != 0 && isEmpty(b[n:]) {
		return v, true
	}
	if v, n := readInt64(b); n != 0 && isEmpty(b[n:]) {
		return v, true
	}
	if v, n := readFloat64(b); n != 0 && isEmpty(b[n:]) {
		return v, true
	}
	if v, n := readString(b); n != 0 && isEmpty(b[n:]) {
		return v, true
	}
	if v, n := readTime(b); n != 0 && isEmpty(b[n:]) {
		return v, true
	}
	if v, n := readDuration(b); n != 0 && isEmpty(b[n:]) {
		return v, true
	}

	return nil, false
}

// Reads the output of `r`, and puts together a list of key/value
// definitions.
func parse(r io.Reader, buf []byte) ([]def, error) {
	lines, err := readLines(r, buf)
	if err != nil {
		return nil, err
	}

	defs := make([]def, 0)

	indents := make([][]byte, 0)
	stack := make([][]byte, 0)
	allowChild := true

	for i, line := range lines {
		if isEmpty(line) {
			continue
		}

		w, k, v := split(line)
		d := depth(indents, w)

		if d == -1 || (d >= len(indents) && !allowChild) {
			return nil, fmt.Errorf(_ErrInvalidIndent, i+1)
		}

		// trim redundant indentation info
		if d < len(indents) {
			indents = indents[:d]
			stack = stack[:d]
		}

		indents = append(indents, w)
		stack = append(stack, k)

		// contains an assignment
		if v != nil {
			defs = append(defs, def{
				key:   resolve(stack...),
				value: string(v),
				line:  i + 1,
			})

			allowChild = false
			continue
		}

		allowChild = true
	}

	return defs, nil
}

// Generates a slice of lines from a reader. Each line must fit in `buf`, or
// an error will be returned.
func readLines(r io.Reader, buf []byte) ([][]byte, error) {
	lines := make([][]byte, 0)
	start, cont := 0, 0

	for {
		if cont == len(buf) {
			return nil, fmt.Errorf(_ErrTooLong, len(lines)+1)
		}

		n, err := r.Read(buf[cont:])
		end := cont + n

		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		for {
			nl := bytes.IndexByte(buf[cont:end], '\n')
			if nl == -1 {
				break
			}

			lines = push(lines, buf[start:cont+nl])

			start = cont + 1
			cont += nl + 1
		}

		if start == 0 {
			cont = end
		} else {
			cont = copy(buf, buf[start:end])
			start = 0
		}
	}

	return push(lines, buf[start:cont]), nil
}

// Safely appends a line to a slice of lines.
func push(lines [][]byte, line []byte) [][]byte {
	dup := make([]byte, len(line))
	copy(dup, line)

	return append(lines, dup)
}

// Returns `true` if the line is blank, or is commented out.
func isEmpty(line []byte) bool {
	for _, b := range line {
		if b == ' ' || b == '\t' {
			continue
		}
		if b == '#' {
			break
		}
		return false
	}
	return true
}

// Returns the prefix whitespace and "key" and "value" components
// of a "key = value" line.
func split(line []byte) (w, k, v []byte) {
	for i, b := range line {
		if b != ' ' && b != '\t' {
			break
		}
		w = line[:i+1]
	}

	if eq := bytes.IndexByte(line, '='); eq != -1 {
		k = line[len(w):eq]
		v = line[eq+1:]
	} else {
		k = line[len(w):]
	}

	return w, k, v
}

// Determines a line's indentation depth by looking at its leading whitespace
// and the leading whitespace of previous lines. Returns an int < 0 when the
// indentation is invalid.
func depth(parents [][]byte, current []byte) int {
	if len(current) == 0 {
		return 0
	} else if len(parents) == 0 {
		// the base indentation level must be empty
		return -1
	}

	for i, parent := range parents {
		if !bytes.HasPrefix(current, parent) {
			return -1
		}
		if len(current) == len(parent) {
			return i
		}
	}

	// the current line is further indented than its parent
	return len(parents)
}

// Generates a string containing each key in `stack`, separated by a dot.
func resolve(stack ...[]byte) string {
	size := len(stack) - 1
	for _, key := range stack {
		size += len(key)
	}

	joined := make([]byte, size)
	offset := 0

	for i, key := range stack {
		if i > 0 {
			joined[offset] = '.'
			offset++
		}

		offset += copy(joined[offset:], key)
	}

	return string(joined)
}
