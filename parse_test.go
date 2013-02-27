package walnut

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

var splitTests = []struct {
	in  string
	out []line
	err error
}{
	{"", []line{}, nil},
	{"# comment", []line{}, nil},
	{"\n\na=1\n\n", []line{{3, 0, "a=1"}}, nil},
	{"a=1\nb=2", []line{{1, 0, "a=1"}, {2, 0, "b=2"}}, nil},
	{"a=1\na=1", []line{{1, 0, "a=1"}, {2, 0, "a=1"}}, nil},
	{"a=1\n b=2", []line{{1, 0, "a=1"}, {2, 1, "b=2"}}, nil},
	{"a=1\n\tb=2", []line{{1, 0, "a=1"}, {2, 1, "b=2"}}, nil},
	{"a=1\n\t \n\tb=2", []line{{1, 0, "a=1"}, {3, 1, "b=2"}}, nil},
	{"\n\t\t\n\n ", []line{}, nil},
	{" a=1", nil, fmt.Errorf(errIndent, 1)},
	{"a=1\n  b=2\n c=3", nil, fmt.Errorf(errIndent, 3)},
	{"a=1\n  b=2\n\tc=3", nil, fmt.Errorf(errIndent, 3)},
}

func TestSplit(t *testing.T) {
	for _, test := range splitTests {
		out, err := split([]byte(test.in))
		if !eq(out, test.out) || !eq(err, test.err) {
			t.Errorf("split(%q):", test.in)
			t.Errorf("   got %+v, %v", out, err)
			t.Errorf("  want %+v, %v", test.out, test.err)
		}
	}
}

var interpretTests = []struct {
	in  []line
	out []assignment
	err error
}{
	{
		[]line{{3, 0, "a=1"}},
		[]assignment{{3, "a", "1", int64(1)}},
		nil,
	},
	{
		[]line{{1, 0, "b=2"}, {2, 0, "c=3"}},
		[]assignment{{1, "b", "2", int64(2)}, {2, "c", "3", int64(3)}},
		nil,
	},
	{
		[]line{{1, 0, "d"}, {2, 1, "e=4"}},
		[]assignment{
			{2, "d.e", "4", int64(4)},
		},
		nil,
	},
	{
		[]line{{1, 0, "foo"}, {2, 1, "bar=5"}, {3, 1, "baz=6"}},
		[]assignment{
			{2, "foo.bar", "5", int64(5)},
			{3, "foo.baz", "6", int64(6)},
		},
		nil,
	},
	{
		[]line{{1, 0, "group#snug"}, {3, 1, "key=\"test\"#snug"}},
		[]assignment{
			{3, "group.key", "\"test\"#snug", "test"},
		},
		nil,
	},
	{
		[]line{{1, 0, "bool = true"}},
		[]assignment{
			{1, "bool", "true", true},
		},
		nil,
	},
	{
		[]line{{1, 0, "int64 = 12345"}},
		[]assignment{
			{1, "int64", "12345", int64(12345)},
		},
		nil,
	},
	{
		[]line{{1, 0, "float64 = 123.45"}},
		[]assignment{
			{1, "float64", "123.45", float64(123.45)},
		},
		nil,
	},
	{
		[]line{{1, 0, "string = \"hello\""}},
		[]assignment{
			{1, "string", "\"hello\"", "hello"},
		},
		nil,
	},
	{
		[]line{{1, 0, "time = 2012-01-02 15:30:28.000000000789 +0000"}},
		func() []assignment {
			raw := "2012-01-02 15:30:28.000000000789 +0000"
			t, _ := time.Parse("2006-01-02 15:04:05 -0700", raw)
			return []assignment{{1, "time", raw, t}}
		}(),
		nil,
	},
	{
		[]line{{1, 0, "duration = 10m 20s"}},
		[]assignment{
			{1, "duration", "10m 20s", 10*time.Minute + 20*time.Second},
		},
		nil,
	},
	{
		[]line{{1, 0, "♫ = 123"}},
		[]assignment{
			{1, "♫", "123", int64(123)},
		},
		nil,
	},
	{[]line{{1, 0, "=1"}}, nil, fmt.Errorf(errKey, 1)},
	{[]line{{1, 0, " = 1"}}, nil, fmt.Errorf(errKey, 1)},
	{[]line{{1, 0, "== 1"}}, nil, fmt.Errorf(errKey, 1)},
	{[]line{{1, 0, "a b = 1"}}, nil, fmt.Errorf(errKey, 1)},
	{[]line{{1, 0, "a\tb"}}, nil, fmt.Errorf(errKey, 1)},
	{[]line{{1, 0, "a = 0 0"}}, nil, fmt.Errorf(errValue, 1, "0 0")},
	{[]line{{1, 0, "a == 0"}}, nil, fmt.Errorf(errValue, 1, "= 0")},
}

func TestInterpret(t *testing.T) {
	for _, test := range interpretTests {
		out, err := interpret(test.in)
		if !eq(out, test.out) || !eq(err, test.err) {
			t.Errorf("interpret(%+v):", test.in)
			t.Errorf("   got %+v, %v", out, err)
			t.Errorf("  want %+v, %v", test.out, test.err)
		}
	}
}

var initializeTests = []struct {
	in  []assignment
	out map[string]interface{}
	err error
}{
	{
		[]assignment{{1, "a", "1", int64(1)}},
		map[string]interface{}{
			"a": int64(1),
		},
		nil,
	},
	{
		[]assignment{{1, "foo.bar", "2", int64(2)}, {1, "foo.baz", "3", int64(3)}},
		map[string]interface{}{
			"foo.bar": int64(2),
			"foo.baz": int64(3),
		},
		nil,
	},
	{
		[]assignment{{1, "a", "1", int64(1)}, {2, "a.b", "2", int64(2)}},
		nil,
		fmt.Errorf(errConflict, "a.b", 2, "a", 1),
	},
	{
		[]assignment{{1, "a", "1", int64(1)}, {2, "a", "1", int64(1)}},
		nil,
		fmt.Errorf(errConflict, "a", 2, "a", 1),
	},
	{
		[]assignment{{1, "a.b.c", "1", int64(1)}, {2, "a.b", "2", int64(2)}},
		nil,
		fmt.Errorf(errConflict, "a.b", 2, "a.b.c", 1),
	},
	{
		[]assignment{{1, "a.b", "1", int64(1)}, {2, "a.b.c", "2", int64(2)}},
		nil,
		fmt.Errorf(errConflict, "a.b.c", 2, "a.b", 1),
	},
}

func TestInitialize(t *testing.T) {
	for _, test := range initializeTests {
		out, err := initialize(test.in)
		if !eq(out, test.out) || !eq(err, test.err) {
			t.Errorf("initialize(%+v):", test.in)
			t.Errorf("   got %+v, %v", out, err)
			t.Errorf("  want %+v, %v", test.out, test.err)
		}
	}
}

// shorthand for reflect.DeepEqual
func eq(a, b interface{}) bool {
	return reflect.DeepEqual(a, b)
}
