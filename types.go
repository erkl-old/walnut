package walnut

import (
	"strconv"
	"strings"
	"time"
)

// Attempts to parse the input string as a floating point (or integral)
// number.
func ParseFloat(input string) (float64, bool) {
	digits := 0
	dot := -1

	for i, r := range input {
		switch {
		case (r == '-' || r == '+') && i == 0:
			// ignore signs at the beginning of the string
		case r >= '0' && r <= '9':
			digits++
		// only allow a dot if we haven't encountered one before and
		// it's preceeded by at least one digit
		case r == '.' && dot == -1 && digits > 0:
			dot = digits
		default:
			// any other circumstance should result in a parsing error
			return 0, false
		}
	}

	// make sure we have encountered at least one digit, and that there
	// are digits after the dot (if it is present)
	if digits == 0 || dot >= digits {
		return 0, false
	}

	// floating point values are fickle things best left to be dealt with
	// by cleverer people; let strconv do the heavy lifting for us
	value, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return 0, false
	}

	return value, true
}

// Reads ASCII digits from the beginning of a string until either the string
// ends, a non-digit character is encountered, or when an integer overflow
// occurs. Also returns the number of characters consumed if successful.
func digits(input string) (int64, int) {
	total := int64(0)

	for i := 0; i < len(input); i++ {
		char := input[i]

		// stop on the first non-digit character
		if char < '0' || char > '9' {
			return total, i
		}

		digit := int64(char - '0')
		next := (total * 10) + digit

		// check for integer overflows
		if next < total {
			return 0, 0
		}

		total = next
	}

	return 0, 0
}

// Reads a unit of duration from the beginning of a string. Also returns
// the number of bytes read if successful.
func unit(s string) (time.Duration, int) {
	switch {
	case strings.HasPrefix(s, "ns"):
		return time.Nanosecond, 2
	case strings.HasPrefix(s, "μs"):
		return time.Microsecond, 3
	case strings.HasPrefix(s, "µs"):
		return time.Microsecond, 2
	case strings.HasPrefix(s, "us"):
		return time.Microsecond, 2
	case strings.HasPrefix(s, "ms"):
		return time.Millisecond, 2
	case strings.HasPrefix(s, "s"):
		return time.Second, 1
	case strings.HasPrefix(s, "m"):
		return time.Minute, 1
	case strings.HasPrefix(s, "h"):
		return time.Hour, 1
	case strings.HasPrefix(s, "d"):
		return 24 * time.Hour, 1
	case strings.HasPrefix(s, "w"):
		return 7 * 24 * time.Hour, 1
	}

	return 0, 0
}
