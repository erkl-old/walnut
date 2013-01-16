package walnut

import (
	"strings"
	"testing"
	"testing/iotest"
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

func TestParseValid(test *testing.T) {
	h := "parse(%#v) ->"

	for _, t := range valid {
		r := iotest.OneByteReader(strings.NewReader(t.in))
		d, err := parse(r, make([]byte, 1024))

		if err != nil {
			test.Errorf(h+" err = %v, want nil", t.in, err)
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

// configuration input which should trigger an indentation error
var indentErrors = []struct {
	in  string
	err string
}{
	{" foo=3", "invalid indentation on line 1"},
	{"a=1\n b=2", "invalid indentation on line 2"},
	{"c\n d=3\n  e=4", "invalid indentation on line 3"},
	{"f\n  g=5\n h=6", "invalid indentation on line 3"},
	{"i\n\t j=7\n  k=8", "invalid indentation on line 3"},
}

func TestParseIndentErrors(test *testing.T) {
	h := "parse(%#v) ->"

	for _, t := range indentErrors {
		r := iotest.OneByteReader(strings.NewReader(t.in))
		_, err := parse(r, make([]byte, 1024))

		if err == nil || err.Error() != t.err {
			test.Errorf(h+" %q != %q", t.in, err, t.err)
		}
	}
}
