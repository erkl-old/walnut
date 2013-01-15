package walnut

import (
	"testing"
)

// configuration input which should parse successfully, returning
// a specific set of definitions
var valid = []struct {
	in string
	d  []def
}{
	{"", []def{}},
	{"\n\t\t\n\n ", []def{}},
	{"#abc", []def{}},
	{"a=1", []def{
		{"a", "1", 1},
	}},
	{"b=2\nc=3", []def{
		{"b", "2", 1},
		{"c", "3", 2},
	}},
	{"d\n e=4", []def{
		{"d.e", "4", 2},
	}},
	{"foo\n\tbar=5\n\tbaz=6", []def{
		{"foo.bar", "5", 2},
		{"foo.baz", "6", 3},
	}},
	{"#\nabc\n def=7\n #\n ghi=8", []def{
		{"abc.def", "7", 3},
		{"abc.ghi", "8", 5},
	}},
}

func TestValidConfigurations(test *testing.T) {
	h := "parseConfig(%#v) ->"

	for _, t := range valid {
		d, line := parseConfig([]byte(t.in))

		if line != 0 {
			test.Errorf(h+" line = %d, want 0", t.in, line)
			continue
		}

		if len(d) != len(t.d) {
			test.Errorf(h+" %v, want %v", t.in, d, t.d)
			continue
		}

		for i := 0; i < len(t.d); i++ {
			if d[i] != t.d[i] {
				test.Errorf(h+" %v, want %v", t.in, d, t.d)
			}
		}
	}
}

// configuration input which should trigger an error on a specific line
var invalid = []struct {
	in string
	l  int
}{
	{" foo=3", 1},
	{"a=1\n b=2", 2},
	{"c\n d=3\n  e=4", 3},
	{"f\n  g=5\n h=6", 3},
	{"i\n\t j=7\n  k=8", 3},
}

func TestInvalidConfigurations(test *testing.T) {
	h := "parseConfig(%#v) ->"

	for _, t := range invalid {
		_, line := parseConfig([]byte(t.in))
		if line != t.l {
			test.Errorf(h+" %d != %d", t.in, line, t.l)
		}
	}
}
