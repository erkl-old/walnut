package walnut

import (
	"errors"
	"reflect"
	"testing"
)

var definitionTests = []struct {
	in   string
	defs []definition
	err  error
}{
	{
		"",
		[]definition{},
		nil,
	},
	{
		"\n\t\t\n\n ",
		[]definition{},
		nil,
	},
	{
		"#abc",
		[]definition{},
		nil,
	},
	{
		"a=1",
		[]definition{
			{"a", int64(1), "1", 1},
		},
		nil,
	},
	{
		"b=2\nc=3",
		[]definition{
			{"b", int64(2), "2", 1},
			{"c", int64(3), "3", 2},
		},
		nil,
	},
	{
		"d\n e=4",
		[]definition{
			{"d.e", int64(4), "4", 2},
		},
		nil,
	},
	{
		"foo\n\tbar=5\n\tbaz=6",
		[]definition{
			{"foo.bar", int64(5), "5", 2},
			{"foo.baz", int64(6), "6", 3},
		},
		nil,
	},
	{
		"#\nabc\n def=7\n #\n ghi=8",
		[]definition{
			{"abc.def", int64(7), "7", 3},
			{"abc.ghi", int64(8), "8", 5},
		},
		nil,
	},
	{
		" foo=3",
		nil,
		errors.New("invalid indentation on line 1"),
	},
	{
		"a=1\n b=2",
		nil,
		errors.New("invalid indentation on line 2"),
	},
	{
		"c\n d=3\n  e=4",
		nil,
		errors.New("invalid indentation on line 3"),
	},
	{
		"f\n  g=5\n h=6",
		nil,
		errors.New("invalid indentation on line 3"),
	},
	{
		"i\n\t j=7\n  k=8",
		nil,
		errors.New("invalid indentation on line 3"),
	},
}

func TestDefinitions(t *testing.T) {
	for _, test := range definitionTests {
		defs, err := definitions([]byte(test.in))

		if !reflect.DeepEqual(err, test.err) ||
			!reflect.DeepEqual(defs, test.defs) {
			t.Errorf("definitions(%#v) -> %v, %#v (want %v, %#v)",
				test.in, defs, err, test.defs, test.err)
		}
	}
}
