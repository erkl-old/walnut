package walnut

import (
	"regexp"
	"strconv"
	"time"
)

var (
	_TruthyRegexp = regexp.MustCompile(`^[ \t]*(true|yes|on)`)
	_FalsyRegexp  = regexp.MustCompile(`^[ \t]*(false|no|off)`)
	_IntRegexp    = regexp.MustCompile(`^[ \t]*([\+\-]?\d+)`)
	_FloatRegexp  = regexp.MustCompile(`^[ \t]*([\+\-]?\d+(?:\.\d+)?)`)
	_TimeRegexp   = regexp.MustCompile(
		`^[ \t]*(\d{4}\-\d{2}\-\d{2} \d{2}:\d{2}:\d{2}(?:\.\d+)? [\-\+]\d{4})`)
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
