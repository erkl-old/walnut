package walnut

import (
	"fmt"
	"io"
)

const (
	_ErrTooLong = "line %d too long (buffer overflow thwarted)"
)

// Defines a "key = value" assignment.
type def struct {
	key   string
	value string
	line  int
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
			nl := indexOf(buf[cont:end], '\n')
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
		if isSpace(b) {
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
		if !isSpace(b) {
			break
		}
		w = line[:i+1]
	}

	k = line[len(w):]

	if eq := indexOf(line, '='); eq != -1 {
		k = k[:eq]
		v = line[eq+1:]
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
		if !hasPrefix(current, parent) {
			return -1
		}
		if len(current) == len(parent) {
			return i
		}
	}

	// the current line is further indented than its parent
	return len(parents)
}
