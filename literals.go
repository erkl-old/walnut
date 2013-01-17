package walnut

import (
	"regexp"
	"strconv"
)

var (
	_TruthyRegexp = regexp.MustCompile(`^[ \t]*(true|yes|on)`)
	_FalsyRegexp  = regexp.MustCompile(`^[ \t]*(false|no|off)`)
	_IntRegexp    = regexp.MustCompile(`^[ \t]*([\+\-]?\d+)`)
	_FloatRegexp  = regexp.MustCompile(`^[ \t]*([\+\-]?\d+(?:\.\d+)?)`)
)

// Attempts to extract a boolean value from the beginning of `in`.
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
