package walnut

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	reInt   = regexp.MustCompile(`^[\+\-]?\d+`)
	reFloat = regexp.MustCompile(`^[\+\-]?\d+\.\d+`)
	reTime  = regexp.MustCompile(
		`^\d{4}\-\d{2}\-\d{2} \d{2}:\d{2}:\d{2}(?:\.\d+)? [\-\+]\d{4}`)
)

const maxDuration = 1<<63 - 1

// Attempts to extract a boolean from the beginning of
// the input string.
func readBool(in string) (bool, int) {
	switch {
	case strings.HasPrefix(in, "true"):
		return true, 4
	case strings.HasPrefix(in, "false"):
		return false, 5
	}

	return false, 0
}

// Attempts to extract an integer from the beginning of
// the input string.
func readInt64(in string) (int64, int) {
	m := reInt.FindStringSubmatchIndex(in)
	if m == nil {
		return 0, 0
	}

	v, err := strconv.ParseInt(in[m[0]:m[1]], 10, 64)
	if err != nil {
		return 0, 0
	}

	return v, m[1]
}

// Attempts to extract a floating point value from the
// beginning of the input string.
func readFloat64(in string) (float64, int) {
	m := reFloat.FindStringSubmatchIndex(in)
	if m == nil {
		return 0, 0
	}

	v, err := strconv.ParseFloat(in[m[0]:m[1]], 64)
	if err != nil {
		return 0, 0
	}

	return v, m[1]
}

// Attempts to extract a string value from the beginning of
// the input string.
func readString(s string) (string, int) {
	if len(s) < 2 || s[0] != '"' {
		return "", 0
	}

	i, end := 1, -1
	escaped := false

	for end < 0 {
		if i == len(s) {
			// found EOL reached before the closing quote
			return "", 0
		}

		switch b := s[i]; {
		case b < 0x20:
			return "", 0
		case escaped:
			escaped = false
		case b == '\\':
			escaped = true
		case b == '"':
			end = i
		}

		i++
	}

	v, err := strconv.Unquote(s[:end+1])
	if err != nil {
		return "", 0
	}

	return v, end + 1
}

// Attempts to extract a timestamp from the beginning of `in`.
func readDuration(in string) (time.Duration, int) {
	var total, prev time.Duration
	var offset int

	for {
		num, unit, n := readDurationPartial(in[offset:])
		if n == 0 {
			break
		}

		switch {
		// units must appear in descending order (greatest first),
		// and only once each
		case prev != 0 && unit >= prev:
			return 0, 0
		// make sure this component doesn't single-handedly overflow
		// time.Duration values
		case num > 0 && unit > maxDuration/num:
			return 0, 0
		// would adding this component to the current total cause it
		// to overflow?
		case num*unit > maxDuration-total:
			return 0, 0
		}

		offset += n
		total += num * unit
		prev = unit
	}

	return total, offset
}

var durations = []struct {
	name  string
	value time.Duration
}{
	{"ns", time.Nanosecond},
	{"μs", time.Microsecond}, // \u03bc
	{"µs", time.Microsecond}, // \u00b5
	{"us", time.Microsecond},
	{"ms", time.Millisecond},
	{"s", time.Second},
	{"m", time.Minute},
	{"h", time.Hour},
	{"d", 24 * time.Hour},
	{"w", 7 * 24 * time.Hour},
}

func readDurationPartial(in string) (num, unit time.Duration, n int) {
	// skip whitespace
	for n < len(in) && strings.ContainsRune(Space, rune(in[n])) {
		n++
	}

	s := n
	for n < len(in) && ('0' <= in[n] && in[n] <= '9') {
		n++
	}

	unsigned, err := strconv.ParseUint(in[s:n], 10, 63)
	if err != nil {
		return 0, 0, 0
	}

	for _, unit := range durations {
		if strings.HasPrefix(in[n:], unit.name) {
			return time.Duration(unsigned), unit.value, n + len(unit.name)
		}
	}

	return 0, 0, 0
}

// Attempts to extract a timestamp from the beginning of `in`.
func readTime(in string) (time.Time, int) {
	m := reTime.FindStringSubmatchIndex(in)
	if m == nil {
		return time.Time{}, 0
	}

	v, err := time.Parse("2006-01-02 15:04:05 -0700", in[m[0]:m[1]])
	if err != nil {
		return time.Time{}, 0
	}

	return v, m[1]
}
