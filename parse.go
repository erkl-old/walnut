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
// a `ErrLineTooLong` error will be returned.
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
			nl := byteIndex(buf[cont:end], '\n')
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
