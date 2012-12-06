package walnut

import (
	"testing"
)

// configuration input which should parse successfully, returning
// a specific set of definitions
var valid = []struct {
	in string
	d  []configDefinition
}{
	{"", []configDefinition{}},
	{"\n\t\t\n\n ", []configDefinition{}},
	{"#abc", []configDefinition{}},
	{"a=1", []configDefinition{
		{"a", "1", 1},
	}},
	{"b=2\nc=3", []configDefinition{
		{"b", "2", 1},
		{"c", "3", 2},
	}},
	{"d\n e=4", []configDefinition{
		{"d.e", "4", 2},
	}},
	{"foo\n\tbar=5\n\tbaz=6", []configDefinition{
		{"foo.bar", "5", 2},
		{"foo.baz", "6", 3},
	}},
	{"#\nabc\n def=7\n #\n ghi=8", []configDefinition{
		{"abc.def", "7", 3},
		{"abc.ghi", "8", 5},
	}},
}

func Test_ValidConfigurations(test *testing.T) {
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
			ok := d[i].key == t.d[i].key &&
				d[i].value == t.d[i].value &&
				d[i].line == t.d[i].line

			if !ok {
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

func Test_InvalidConfigurations(test *testing.T) {
	h := "parseConfig(%#v) ->"

	for _, t := range invalid {
		_, line := parseConfig([]byte(t.in))
		if line != t.l {
			test.Errorf(h+" %d != %d", t.in, line, t.l)
		}
	}
}
