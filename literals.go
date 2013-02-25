package walnut

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	_TruthyRegexp = regexp.MustCompile(`^[ \t]*on`)
	_FalsyRegexp  = regexp.MustCompile(`^[ \t]*off`)
	_IntRegexp    = regexp.MustCompile(`^[ \t]*([\+\-]?\d+)`)
	_FloatRegexp  = regexp.MustCompile(`^[ \t]*([\+\-]?\d+\.\d+)`)
	_TimeRegexp   = regexp.MustCompile(
		`^[ \t]*(\d{4}\-\d{2}\-\d{2} \d{2}:\d{2}:\d{2}(?:\.\d+)? [\-\+]\d{4})`)

	_MaxInt64    int64         = 1<<63 - 1
	_MaxDuration time.Duration = 1<<63 - 1
)

// Attempts to extract a string literal from the beginning of `in`.
func readBool(s string) (bool, int) {
	if m := _TruthyRegexp.FindStringIndex(s); m != nil {
		return true, m[1]
	}
	if m := _FalsyRegexp.FindStringIndex(s); m != nil {
		return false, m[1]
	}

	return false, 0
}

// Attempts to extract a signed integer from the beginning of `in`.
func readInt64(s string) (int64, int) {
	m := _IntRegexp.FindStringSubmatchIndex(s)
	if m == nil {
		return 0, 0
	}

	num := s[m[2]:m[3]]
	v, err := strconv.ParseInt(num, 10, 64)
	if err != nil {
		return 0, 0
	}

	return v, m[3]
}

// Attempts to extract a floating point value from the beginning of `in`.
func readFloat64(s string) (float64, int) {
	m := _FloatRegexp.FindStringSubmatchIndex(s)
	if m == nil {
		return 0, 0
	}

	slice := s[m[2]:m[3]]
	v, err := strconv.ParseFloat(slice, 64)
	if err != nil {
		return 0, 0
	}

	return v, m[3]
}

// Attempts to extract a timestamp from the beginning of `in`.
func readString(s string) (string, int) {
	start := 0
	for start < len(s) && (s[start] == ' ' || s[start] == '\t') {
		start++
	}

	if len(s)-start < 2 || s[start] != '"' {
		return "", 0
	}

	i := start + 1 // jump the first double quote
	end := -1
	escaped := false

	for end == -1 {
		if i == len(s) {
			// end of input reached before finding a closing quote
			return "", 0
		}

		b := s[i]

		switch {
		case b <= 0x20:
			// control characters aren't inside a string literal
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

	v, err := strconv.Unquote(s[start : end+1])
	if err != nil {
		return "", 0
	}

	return v, end + 1
}

// Attempts to extract a timestamp from the beginning of `in`.
func readTime(s string) (time.Time, int) {
	m := _TimeRegexp.FindStringSubmatchIndex(s)
	if m == nil {
		return time.Time{}, 0
	}

	slice := s[m[2]:m[3]]
	v, err := time.Parse("2006-01-02 15:04:05 -0700", slice)
	if err != nil {
		return time.Time{}, 0
	}

	return v, m[3]
}

// Attempts to extract a timestamp from the beginning of `in`.
func readDuration(s string) (time.Duration, int) {
	var total, prev time.Duration
	var offset int

	for {
		num, unit, n := readDurationPartial(s[offset:])
		if n == 0 {
			break
		}

		// time units must appear in descending order (greatest first),
		// and only once each
		if prev != 0 && unit >= prev {
			return 0, 0
		}

		prev = unit
		v := time.Duration(num) * unit

		// guard against integer overflow
		if v > _MaxDuration-total {
			return 0, 0
		}

		total += v
		offset += n
	}

	return total, offset
}

var timeUnits = []struct {
	name string
	dur  time.Duration
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

func readDurationPartial(s string) (num int64, unit time.Duration, n int) {
	i, end := 0, len(s)

	// skip whitespace
	for i < end && (s[i] == ' ' || s[i] == '\t') {
		i++
	}

	start := i

	for ; i < end && ('0' <= s[i] && s[i] <= '9'); i++ {
		digit := int64(s[i] - '0')

		// guard against overflow
		if digit > _MaxInt64-(num*10) {
			return 0, 0, 0
		}

		num = (num * 10) + digit
	}

	// did we find any digits?
	if i == start {
		return 0, 0, 0
	}

	for _, unit := range timeUnits {
		if strings.HasPrefix(s[i:], unit.name) {
			return num, unit.dur, i + len(unit.name)
		}
	}

	return 0, 0, 0
}
