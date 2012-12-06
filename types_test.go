package walnut

import (
	"testing"
	"time"
)

// test suite for ParseBool
var boolTests = []struct {
	in string
	b  bool
	ok bool
}{
	// truthy
	{"true", true, true},
	{"yes", true, true},
	{"on", true, true},

	// falsy
	{"false", false, true},
	{"no", false, true},
	{"off", false, true},

	// invalid
	{"", false, false},
	{"x", false, false},
	{"1", false, false},
}

func Test_ParseBool(test *testing.T) {
	h := "ParseBool(%#v) ->"

	for _, t := range boolTests {
		b, ok := ParseBool(t.in)

		switch {
		case ok != t.ok:
			test.Errorf(h+" ok != %#v", t.in, t.ok)
		case b != t.b:
			test.Errorf(h+" %#v, want %#v", t.in, b, t.b)
		}
	}
}

// test suite for ParseFloat
var floatTests = []struct {
	in string
	f  float64
	ok bool
}{
	// integers
	{"1", 1, true},
	{"987654321", 987654321, true},
	{"123456700", 1.234567e+08, true},

	// decimals
	{"1.3", 1.3, true},
	{"100.0", 100, true},
	{"38.002", 38.002, true},
	{"0.1", 0.1, true},

	// signs
	{"-0", 0, true},
	{"+1", 1, true},
	{"-2", -2, true},
	{"-0.1", -0.1, true},
	{"-3.4", -3.4, true},

	// long
	{"100000000000000000000000", 1e+23, true},
	{"99999999999999974834176", 9.999999999999997e+22, true},
	{"100000000000000000000001", 1.0000000000000001e+23, true},
	{"100000000000000008388608", 1.0000000000000001e+23, true},
	{"100000000000000016777215", 1.0000000000000001e+23, true},
	{"100000000000000016777216", 1.0000000000000003e+23, true},
	{"22.222222222222222", 22.22222222222222, true},

	// invalid
	{"", 0, false},
	{"1a", 0, false},
	{"0x30", 0, false},
	{"1.1.", 0, false},
	{"+-0", 0, false},
	{"-0-", 0, false},
}

func Test_ParseFloat(test *testing.T) {
	h := "ParseFloat(%#v) ->"

	for _, t := range floatTests {
		f, ok := ParseFloat(t.in)

		switch {
		case ok != t.ok:
			test.Errorf(h+" ok != %#v", t.in, t.ok)
		case f != t.f:
			test.Errorf(h+" %#v, want %#v", t.in, f, t.f)
		}
	}
}

// test suite for ParseString
var stringTests = []struct {
	in string
	s  string
	ok bool
}{
	// valid
	{`""`, "", true},
	{`"hello"`, "hello", true},
	{`"日本人"`, "日本人", true},
	{`"a\nb"`, "a\nb", true},
	{`"\u00FF"`, "ÿ", true},
	{`"\xFF"`, "\xFF", true},
	{`"\U00010111"`, "\U00010111", true},
	{`"\U0001011111"`, "\U0001011111", true},
	{`"'"`, "'", true},
	{`"\""`, "\"", true},

	// invalid
	{``, "", false},
	{`"lone`, "", false},
	{`hello`, "", false},
	{`"mismatch'`, "", false},
	{`"\"`, "", false},
	{`"\1"`, "", false},
	{`"\19"`, "", false},
	{`"\129"`, "", false},
	{"`a`", "", false},
	{"'b'", "", false},
}

func Test_ParseString(test *testing.T) {
	h := "ParseString(%#v) ->"

	for _, t := range stringTests {
		s, ok := ParseString(t.in)

		switch {
		case ok != t.ok:
			test.Errorf(h+" ok != %#v", t.in, t.ok)
		case s != t.s:
			test.Errorf(h+" %#v, want %#v", t.in, s, t.s)
		}
	}
}

// test suite for ParseTime
var timeTests = []struct {
	in string
	ok bool
}{
	// basic
	{"1970-01-01 00:00:00 +0000", true},
	{"2001-02-03 04:05:06 +0000", true},
	{"1997-08-28 15:30:27.123 +0000", true},
	{"1997-08-28 14:07:27 -0123", true},

	// invalid
	{"", false},
	{"01:02:03", false},
	{"1970-01-01", false},
	{"1970-01-01 00:00:00", false},
	{"1970-02-48 00:00:00 +0000", false},
	{"70-01-01 00:00:00", false},
	{"1970-01-01 00:00:00 UTC", false},
}

func Test_ParseTime(test *testing.T) {
	h := "ParseTime(%#v) ->"

	for _, t := range timeTests {
		e, _ := time.Parse("2006-01-02 15:04:05 -0700", t.in)
		a, ok := ParseTime(t.in)

		switch {
		case ok != t.ok:
			test.Errorf(h+" ok != %#v", t.in, t.ok)
		case !a.Equal(e):
			test.Errorf(h+" %s, want %s", t.in, a.String(), e.String())
		}
	}
}

// test suite for ParseDuration
var durationTests = []struct {
	in string
	d  time.Duration
	ok bool
}{
	// simple formats
	{"0s", 0, true},
	{"5s", 5 * time.Second, true},
	{"37s", 37 * time.Second, true},
	{"010s", 10 * time.Second, true},
	{"3d", 3 * 24 * time.Hour, true},

	// all units
	{"10ns", 10 * time.Nanosecond, true},
	{"10µs", 10 * time.Microsecond, true},
	{"10μs", 10 * time.Microsecond, true},
	{"10us", 10 * time.Microsecond, true},
	{"10ms", 10 * time.Millisecond, true},
	{"10s", 10 * time.Second, true},
	{"10m", 10 * time.Minute, true},
	{"10h", 10 * time.Hour, true},
	{"10d", 10 * 24 * time.Hour, true},
	{"10w", 10 * 7 * 24 * time.Hour, true},

	// mixed units
	{"1h1m1s", time.Hour + time.Minute + time.Second, true},
	{"4h30m", 4*time.Hour + 30*time.Minute, true},
	{"1s500ms", time.Second + 500*time.Millisecond, true},
	{"1w1d24h1440m", 10 * 24 * time.Hour, true},

	// allow (ignore) spaces between components
	{"1h 1m1s", time.Hour + time.Minute + time.Second, true},
	{"4h 30m", 4*time.Hour + 30*time.Minute, true},
	{"1s 500ms", time.Second + 500*time.Millisecond, true},
	{"1w 1d 24h 1440m", 10 * 24 * time.Hour, true},

	// disallow signs and decimal values
	{"-3h", 0, false},
	{"+5m", 0, false},
	{"300.5h", 0, false},
	{"1h 1m 1.3s", 0, false},
	{"10w -3d", 0, false},
	{"1.2d20m", 0, false},

	// units out of order
	{"1s2m", 0, false},
	{"1200ms 3s", 0, false},
	{"4h 5d 6w 7m", 0, false},

	// other invalid formats
	{"", 0, false},
	{"1sm", 0, false},
	{" 1h 1m1s", 0, false},
	{"4h  30m", 0, false},
	{"2 m 3 s", 0, false},
	{"4 d5 h", 0, false},
	{"100", 0, false},
	{"1d 200", 0, false},
	{"3 4 5ms", 0, false},
}

func Test_ParseDuration(test *testing.T) {
	h := "ParseDuration(%#v) ->"

	for _, t := range durationTests {
		d, ok := ParseDuration(t.in)

		switch {
		case ok != t.ok:
			test.Errorf(h+" ok != %#v", t.in, t.ok)
		case d != t.d:
			test.Errorf(h+" %#v, want %#v", t.in, d, t.d)
		}
	}
}
