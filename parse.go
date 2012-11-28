package walnut

import (
	"strconv"
	"time"
)

const (
	TypeString = iota // = 0
	TypeInt
	TypeBool
	TypeDuration
	TypeTime

	TypeNone = -1
)

func DetectType(input string) int {
	if _, ok := ParseString(input); ok {
		return TypeString
	}

	if _, ok := ParseInt(input); ok {
		return TypeInt
	}

	if _, ok := ParseBool(input); ok {
		return TypeBool
	}

	if _, ok := ParseDuration(input); ok {
		return TypeDuration
	}

	if _, ok := ParseTime(input); ok {
		return TypeTime
	}

	return TypeNone
}

func ParseString(input string) (string, bool) {
	if input == "" || input[0] != '"' {
		return "", false
	}

	value, err := strconv.Unquote(input)
	if err != nil {
		return "", false
	}

	return value, true
}

func ParseInt(input string) (int64, bool) {
	value, err := strconv.ParseInt(input, 0, 64)
	if err != nil {
		return 0, false
	}

	return value, true
}

func ParseBool(input string) (bool, bool) {
	switch input {
	case "true", "yes", "on":
		return true, true
	case "false", "no", "off":
		return false, true
	}

	return false, false
}

func ParseDuration(input string) (time.Duration, bool) {
	// don't bother parsing an empty string
	if input == "" {
		return 0, false
	}

	total := time.Duration(0)
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
		char := s[i]

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

func readUnit(s string) (time.Duration, int) {
	switch {
	case hasPrefix(s, "ns"):
		return time.Nanosecond, 2
	case hasPrefix(s, "μs"), hasPrefix(s, "µs"):
		return time.Microsecond, 3
	case hasPrefix(s, "us"):
		return time.Microsecond, 2
	case hasPrefix(s, "ms"):
		return time.Millisecond, 2
	case hasPrefix(s, "s"):
		return time.Second, 1
	case hasPrefix(s, "m"):
		return time.Minute, 1
	case hasPrefix(s, "h"):
		return time.Hour, 1
	case hasPrefix(s, "d"):
		return 24 * time.Hour, 1
	case hasPrefix(s, "w"):
		return 7 * 24 * time.Hour, 1
	}

	return 0, 0
}

func hasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[0:len(prefix)] == prefix
}

func ParseTime(input string) (time.Time, bool) {
	t, err := time.Parse("2006-01-02 15:04:05 -0700", input)
	if err != nil {
		return time.Time{}, false
	}

	return t, true
}
