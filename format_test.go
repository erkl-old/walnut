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
	// @todo: tests go here
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
