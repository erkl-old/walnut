package walnut

import (
	"strings"
	"time"
)

// Reads ASCII digits from the beginning of a string until the string
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
