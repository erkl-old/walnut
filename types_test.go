package walnut

import (
	"testing"
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

func TestParseBool(test *testing.T) {
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

func TestParseFloat(test *testing.T) {
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
