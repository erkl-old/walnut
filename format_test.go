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
	{ "0", 0, true },
	{ "10", 10, true },
	{ "123456789", 123456789, true },

	// hexadecimal
	{ "0x02", 2, true },
	{ "0xff", 255, true },
	{ "0xc", 12, true },

	// octal
	{ "010", 8, true },
	{ "01234567", 342391, true },
	{ "012345678", 0, false },

	// signs
	{ "+0", 0, true },
	{ "-0", 0, true },
	{ "+10", 10, true },
	{ "-0x00", 0, true },
	{ "-0x10", -16, true },
	{ "+01", 1, true },
	{ "-010", -8, true },

	// limits
	{ "9223372036854775807", 1<<63 - 1, true },
	{ "9223372036854775808", 0, false },
	{ "9223372036854775809", 0, false },
	{ "-9223372036854775807", -(1<<63 - 1), true },
	{ "-9223372036854775808", -1 << 63, true },
	{ "-9223372036854775809", 0, false },

	{ "0x7FFFFFFFFFFFFFFF", 1<<63 - 1, true },
	{ "0X8000000000000000", 0, false },
	{ "0X8000000000000001", 0, false },
	{ "-0x7FFFFFFFFFFFFFFF", -(1<<63 - 1), true },
	{ "-0X8000000000000000", -1 << 63, true },
	{ "-0X8000000000000001", 0, false },

	{ "0777777777777777777777", 1<<63 - 1, true },
	{ "01000000000000000000000", 0, false },
	{ "01000000000000000000001", 0, false },
	{ "-0777777777777777777777", -(1<<63 - 1), true },
	{ "-01000000000000000000000", -1 << 63, true },
	{ "-01000000000000000000001", 0, false },

	// invalid
	{ "", 0, false },
	{ "abc", 0, false },
	{ "100 blue", 0, false },
	{ "-0-", 0, false },
	{ "++0", 0, false },
}

var boolTests = []struct {
	in  string
	out bool
	ok  bool
}{
	// truthy
	{ "true", true, true },
	{ "yes", true, true },
	{ "on", true, true },

	// falsy
	{ "false", false, true },
	{ "no", false, true },
	{ "off", false, true },

	// invalid
	{ "", false, false },
	{ "y", false, false },
	{ "foo", false, false },
}

var durationTests = []struct {
	in  string
	out time.Duration
	ok  bool
}{
	// @todo: tests go here
}

var timeTests = []struct {
	in  string
	out *time.Time
	ok  bool
}{
	// @todo: tests go here
}

func TestParseString(t *testing.T) {
	for _, test := range stringTests {
		t.Logf("> parseString(%#v)\n", test.in)
		out, ok := parseString(test.in)

		if ok != test.ok {
			if test.ok {
				t.Fatalf("parsing failed unexpectedly\n")
			} else {
				t.Fatalf("parsing should not succeed\n")
			}
		}

		if out != test.out {
			t.Fatalf("%#v != %#v\n", out, test.out)
		}
	}
}

func TestParseInt(t *testing.T) {
	for _, test := range intTests {
		t.Logf("> parseInt(%#v)\n", test.in)
		out, ok := parseInt(test.in)

		if ok != test.ok {
			if test.ok {
				t.Fatalf("parsing failed unexpectedly\n")
			} else {
				t.Fatalf("parsing should not succeed\n")
			}
		}

		if out != test.out {
			t.Fatalf("%#v != %#v\n", out, test.out)
		}
	}
}

func TestParseBool(t *testing.T) {
	for _, test := range boolTests {
		t.Logf("> parseBool(%#v)\n", test.in)
		out, ok := parseBool(test.in)

		if ok != test.ok {
			if test.ok {
				t.Fatalf("parsing failed unexpectedly\n")
			} else {
				t.Fatalf("parsing should not succeed\n")
			}
		}

		if out != test.out {
			t.Fatalf("%#v != %#v\n", out, test)
		}
	}
}

func TestParseDuration(t *testing.T) {
	for _, test := range durationTests {
		t.Logf("> parseDuration(%#v)\n", test.in)
		out, ok := parseDuration(test.in)

		if ok != test.ok {
			if test.ok {
				t.Fatalf("parsing failed unexpectedly\n")
			} else {
				t.Fatalf("parsing should not succeed\n")
			}
		}

		if out != test.out {
			t.Fatalf("%#v != %#v\n", out.String(), test.out.String())
		}
	}
}

func TestParseTime(t *testing.T) {
	for _, test := range timeTests {
		t.Logf("> parseTime(%#v)\n", test.in)
		out, ok := parseTime(test.in)

		if ok != test.ok {
			if test.ok {
				t.Fatalf("parsing failed unexpectedly\n")
			} else {
				t.Fatalf("parsing should not succeed\n")
			}
		}

		if out == nil {
			if test.out != nil {
				t.Fatalf("returned nil\n")
			}
			return
		}

		if !out.Equal(*test.out) {
			t.Fatalf("%#v != %#v\n", out.String(), test.out.String())
		}
	}
}

// Convenience function for creating a `time.Time` struct of a certain value,
// then returning a pointer to it.
func t(value string) *time.Time {
	t, err := time.Parse("2006-01-02 15:04:05 -07:00", value)
	if err != nil {
		panic(err)
	}

	return &t
}
