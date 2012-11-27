package nut

import (
	"testing"
	"time"
)

var stringTests = []struct {
	in  string
	out string
	ok  bool
}{
	// @todo: tests go here
}

var intTests = []struct {
	in  string
	out int64
	ok  bool
}{
	// decimal
	{"0", 0, true},
	{"10", 10, true},
	{"123456789", 123456789, true},

	// hexadecimal
	{"0x02", 2, true},
	{"0xff", 255, true},
	{"0xc", 12, true},

	// octal
	{"010", 8, true},
	{"01234567", 342391, true},
	{"012345678", 0, false},

	// signs
	{"+0", 0, true},
	{"-0", 0, true},
	{"+10", 10, true},
	{"-0x00", 0, true},
	{"-0x10", -16, true},
	{"+01", 1, true},
	{"-010", -8, true},

	// limits
	{"9223372036854775807", 1<<63 - 1, true},
	{"9223372036854775808", 0, false},
	{"9223372036854775809", 0, false},
	{"-9223372036854775807", -(1<<63 - 1), true},
	{"-9223372036854775808", -1 << 63, true},
	{"-9223372036854775809", 0, false},

	{"0x7FFFFFFFFFFFFFFF", 1<<63 - 1, true},
	{"0X8000000000000000", 0, false},
	{"0X8000000000000001", 0, false},
	{"-0x7FFFFFFFFFFFFFFF", -(1<<63 - 1), true},
	{"-0X8000000000000000", -1 << 63, true},
	{"-0X8000000000000001", 0, false},

	{"0777777777777777777777", 1<<63 - 1, true},
	{"01000000000000000000000", 0, false},
	{"01000000000000000000001", 0, false},
	{"-0777777777777777777777", -(1<<63 - 1), true},
	{"-01000000000000000000000", -1 << 63, true},
	{"-01000000000000000000001", 0, false},

	// invalid
	{"", 0, false},
	{"abc", 0, false},
	{"100 blue", 0, false},
	{"-0-", 0, false},
	{"++0", 0, false},
}

var boolTests = []struct {
	in  string
	out bool
	ok  bool
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
	{"y", false, false},
	{"foo", false, false},
}

var durationTests = []struct {
	in  string
	out time.Duration
	ok  bool
}{
	// simple formats
	{"0", 0, true},
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

	// allow (ignore) spaces
	{"1h 1m1s", time.Hour + time.Minute + time.Second, true},
	{"4h 30m", 4*time.Hour + 30*time.Minute, true},
	{"1s    500ms", time.Second + 500*time.Millisecond, true},
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
	{"2 m 3 s", 0, false},
	{"4 d5 h", 0, false},
	{"100", 0, false},
	{"1d 200", 0, false},
	{"3 4 5ms", 0, false},
}

var timeTests = []struct {
	in  string
	out time.Time
	ok  bool
}{
	// @todo: tests go here
	{"1970-01-01 00:00:00 +00:00", time.Time{}, true},

	{"1970-01-01  00:00:00 +00:00", time.Time{}, true},
}

func TestParseString(t *testing.T) {
	for _, test := range stringTests {
		fail := setup(t, "parseString", test.in)
		out, ok := parseString(test.in)

		if ok != test.ok {
			if test.ok {
				fail("parsing failed unexpectedly")
			} else {
				fail("parsing should not succeed")
			}
			continue
		}

		if out != test.out {
			fail("%#v != %#v", out, test.out)
		}
	}
}

func TestParseInt(t *testing.T) {
	for _, test := range intTests {
		fail := setup(t, "parseInt", test.in)
		out, ok := parseInt(test.in)

		if ok != test.ok {
			if test.ok {
				fail("parsing failed unexpectedly")
			} else {
				fail("parsing should not succeed")
			}
			continue
		}

		if out != test.out {
			fail("%#v != %#v", out, test.out)
		}
	}
}

func TestParseBool(t *testing.T) {
	for _, test := range boolTests {
		fail := setup(t, "parseBool", test.in)
		out, ok := parseBool(test.in)

		if ok != test.ok {
			if test.ok {
				fail("parsing failed unexpectedly")
			} else {
				fail("parsing should not succeed")
			}
			continue
		}

		if out != test.out {
			fail("%#v != %#v", out, test)
		}
	}
}

func TestParseDuration(t *testing.T) {
	for _, test := range durationTests {
		fail := setup(t, "parseDuration", test.in)
		out, ok := parseDuration(test.in)

		if ok != test.ok {
			if test.ok {
				fail("parsing failed unexpectedly")
			} else {
				fail("parsing should not succeed")
			}
			continue
		}

		if out != test.out {
			fail("%s != %s", out.String(), test.out.String())
		}
	}
}

func TestParseTime(t *testing.T) {
	for _, test := range timeTests {
		fail := setup(t, "parseTime", test.in)
		out, ok := parseTime(test.in)

		if ok != test.ok {
			if test.ok {
				fail("parsing failed unexpectedly")
			} else {
				fail("parsing should not succeed")
			}
			continue
		}

		if !out.Equal(test.out) {
			fail("%s != %s", out.String(), test.out.String())
		}
	}
}

// Simplify failure reporting.
type failFunc func(format string, values ...interface{})

func setup(t *testing.T, signature string, input interface{}) failFunc {
	return func(format string, values ...interface{}) {
		args := append([]interface{}{signature, input}, values...)
		t.Errorf("%s(%#v): "+format+"\n", args...)
	}
}

// Reduced version of `time.Date`.
func date(y, m, d, H, M, S, ms int) time.Time {
	month := time.Month(m)
	nano := ms * int(time.Millisecond)

	return time.Date(y, month, d, H, M, S, nano, time.UTC)
}
