package nut

import (
	"strconv"
	"time"
)

func parseString(input string) (string, bool) {
	return "", false
}

func parseInt(input string) (int64, bool) {
	i, err := strconv.ParseInt(input, 0, 64)
	if err != nil {
		return 0, false
	}

	return i, true
}

func parseBool(input string) (bool, bool) {
	switch input {
	case "true", "yes", "on":
		return true, true
	case "false", "no", "off":
		return false, true
	}

	return false, false
}

func parseDuration(input string) (time.Duration, bool) {
	// don't bother parsing an empty string
	if input == "" {
		return 0, false
	}

	// special case
	if input == "0" {
		return 0, true
	}

	var total time.Duration
	var lastUnit time.Duration

	for {
		// consume leading whitespace
		for input != "" && input[0] == ' ' {
			input = input[1:]
		}

		// stop with nothing left to parse
		if input == "" {
			break
		}

		// parse the number and the unit; fail if either of them
		// are missing
		num, digits := readDigits(input)
		unit, chars := readUnit(input[digits:])

		if digits == 0 || chars == 0 {
			return 0, false
		}

		// make sure the components are ordered by unit, in
		// descending order
		if unit >= lastUnit && lastUnit != 0 {
			return 0, false
		}

		// detect integer overflows
		next := total + (time.Duration(num) * unit)
		if next < total {
			return 0, false
		}

		lastUnit = unit
		input = input[digits+chars:]

		total = next
	}

	return total, true
}

func readDigits(s string) (int64, int) {
	total := int64(0)

	for i := 0; i < len(s); i++ {
		ch := s[i]

		// make sure we're looking at a digit
		if ch < '0' || ch > '9' {
			return total, i
		}

		digit := int64(ch - '0')
		next := (total * 10) + digit

		// check for integer overflows
		if next < total {
			return 0, 0
		}

		total = next
	}

	return 0, 0
}

func readUnit(s string) (time.Duration, int) {
	switch {
	case begins(s, "ns"):
		return time.Nanosecond, 2
	case begins(s, "μs"), begins(s, "µs"):
		return time.Microsecond, 3
	case begins(s, "us"):
		return time.Microsecond, 2
	case begins(s, "ms"):
		return time.Millisecond, 2
	case begins(s, "s"):
		return time.Second, 1
	case begins(s, "m"):
		return time.Minute, 1
	case begins(s, "h"):
		return time.Hour, 1
	case begins(s, "d"):
		return 24*time.Hour, 1
	case begins(s, "w"):
		return 7*24*time.Hour, 1
	}

	return 0, 0
}

func begins(s, prefix string) bool {
	return len(s) >= len(prefix) && s[0:len(prefix)] == prefix
}

func parseTime(input string) (*time.Time, bool) {
	return nil, false
}
