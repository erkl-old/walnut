package walnut

import (
	"bytes"
	"regexp"
	"strconv"
	"time"
)

var (
	_TruthyRegexp = regexp.MustCompile(`^[ \t]*(true|yes|on)`)
	_FalsyRegexp  = regexp.MustCompile(`^[ \t]*(false|no|off)`)
	_IntRegexp    = regexp.MustCompile(`^[ \t]*([\+\-]?\d+)`)
	_FloatRegexp  = regexp.MustCompile(`^[ \t]*([\+\-]?\d+\.\d+)`)
	_TimeRegexp   = regexp.MustCompile(
		`^[ \t]*(\d{4}\-\d{2}\-\d{2} \d{2}:\d{2}:\d{2}(?:\.\d+)? [\-\+]\d{4})`)

	_MaxInt64    int64         = 1<<63 - 1
	_MaxDuration time.Duration = 1<<63 - 1
)

// Attempts to extract a string literal from the beginning of `in`.
func readBool(in []byte) (bool, int) {
	if m := _TruthyRegexp.FindIndex(in); m != nil {
		return true, m[1]
	}
	if m := _FalsyRegexp.FindIndex(in); m != nil {
		return false, m[1]
	}

	return false, 0
}

// Attempts to extract a signed integer from the beginning of `in`.
func readInt64(in []byte) (int64, int) {
	m := _IntRegexp.FindSubmatchIndex(in)
	if m == nil {
		return 0, 0
	}

	num := string(in[m[2]:m[3]])
	v, err := strconv.ParseInt(num, 10, 64)
	if err != nil {
		return 0, 0
	}

	return v, m[3]
}

// Attempts to extract a floating point value from the beginning of `in`.
func readFloat64(in []byte) (float64, int) {
	m := _FloatRegexp.FindSubmatchIndex(in)
	if m == nil {
		return 0, 0
	}

	slice := string(in[m[2]:m[3]])
	v, err := strconv.ParseFloat(slice, 64)
	if err != nil {
		return 0, 0
	}

	return v, m[3]
}

// Attempts to extract a timestamp from the beginning of `in`.
func readString(in []byte) (string, int) {
	start := 0
	for start < len(in) && (in[start] == ' ' || in[start] == '\t') {
		start++
	}

	if len(in)-start < 2 || in[start] != '"' {
		return "", 0
	}

	i := start + 1 // jump the first double quote
	end := -1
	escaped := false

	for end == -1 {
		if i == len(in) {
			// end of input reached before finding a closing quote
			return "", 0
		}

		b := in[i]

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

	v, err := strconv.Unquote(string(in[start : end+1]))
	if err != nil {
		return "", 0
	}

	return v, end + 1
}

// Attempts to extract a timestamp from the beginning of `in`.
func readTime(in []byte) (time.Time, int) {
	m := _TimeRegexp.FindSubmatchIndex(in)
	if m == nil {
		return time.Time{}, 0
	}

	slice := string(in[m[2]:m[3]])
	v, err := time.Parse("2006-01-02 15:04:05 -0700", slice)
	if err != nil {
		return time.Time{}, 0
	}

	return v, m[3]
}

// Attempts to extract a timestamp from the beginning of `in`.
func readDuration(in []byte) (time.Duration, int) {
	offset := 0
	total := time.Duration(0)

	for {
		v, n := readDurationPartial(in[offset:])
		if n == 0 {
			break
		}

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
	name []byte
	dur  time.Duration
}{
	{[]byte("ns"), time.Nanosecond},
	{[]byte("μs"), time.Microsecond}, // \u03bc
	{[]byte("µs"), time.Microsecond}, // \u00b5
	{[]byte("us"), time.Microsecond},
	{[]byte("ms"), time.Millisecond},
	{[]byte("s"), time.Second},
	{[]byte("m"), time.Minute},
	{[]byte("h"), time.Hour},
	{[]byte("d"), 24 * time.Hour},
	{[]byte("w"), 7 * 24 * time.Hour},
}

func readDurationPartial(in []byte) (time.Duration, int) {
	i, end := 0, len(in)

	// skip whitespace
	for i < end && (in[i] == ' ' || in[i] == '\t') {
		i++
	}

	value := int64(0)
	start := i

	for ; i < end && ('0' <= in[i] && in[i] <= '9'); i++ {
		digit := int64(in[i] - '0')

		// guard against overflow
		if digit > _MaxInt64-(value*10) {
			return 0, 0
		}

		value = (value * 10) + digit
	}

	// did we find any digits?
	if i == start {
		return 0, 0
	}

	for _, unit := range timeUnits {
		if bytes.HasPrefix(in[i:], unit.name) {
			return time.Duration(value) * unit.dur, i + len(unit.name)
		}
	}

	return 0, 0
}
