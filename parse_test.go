package walnut

import (
	"testing"
)

// configuration input which should parse successfully, returning
// a specific set of definitions
var parseTests = []struct {
	in string
	d  []definition
}{
	{"", []definition{}},
	{"\n\t\t\n\n ", []definition{}},
	{"#abc", []definition{}},
	{"a=1", []definition{
		{"a", "1", 1},
	}},
	{"b=2\nc=3", []definition{
		{"b", "2", 1},
		{"c", "3", 2},
	}},
	{"d\n e=4", []definition{
		{"d.e", "4", 2},
	}},
	{"foo\n\tbar=5\n\tbaz=6", []definition{
		{"foo.bar", "5", 2},
		{"foo.baz", "6", 3},
	}},
	{"#\nabc\n def=7\n #\n ghi=8", []definition{
		{"abc.def", "7", 3},
		{"abc.ghi", "8", 5},
	}},
}

func TestValidConfigParsing(test *testing.T) {
	h := "parse(%#v) ->"

	for _, t := range parseTests {
		d, err := parse([]byte(t.in))

		if err != nil {
			test.Errorf(h+" error: %#v", t.in, err.Error())
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
