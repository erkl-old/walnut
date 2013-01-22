package walnut

import (
	"testing"
)

var readBoolTests = []struct {
	in string
	v  bool
	n  int
}{
	{"true", true, 4},
	{"\t false", false, 7},
	{"  on ", true, 4},
	{"offfoo", false, 3},
	{" yes ", true, 4},
	{"\t \tno", false, 5},

	{"blurgh", false, 0},
	{"  blop  ", false, 0},
}

func TestReadBool(t *testing.T) {
	for _, test := range readBoolTests {
		v, n := readBool([]byte(test.in))

		if v != test.v || n != test.n {
			t.Errorf("readBool(%#q) -> %v, %d (want %v, %d)",
				test.in, v, n, test.v, test.n)
		}
	}
}

var readInt64Tests = []struct {
	in string
	v  int64
	n  int
}{
	{"", 0, 0},
	{"  0", 0, 3},
	{"00000000", 0, 8},
	{"00000001", 1, 8},
	{"\t1x", 1, 2},
	{"12345", 12345, 5},
	{" 1 2 3", 1, 2},
	{"-10 ", -10, 3},
	{"- 10", 0, 0},
	{"-10 -10", -10, 3},
	{"--2", 0, 0},
	{"+-3128", 0, 0},
	{"\t-012301 ", -12301, 8},
	{"\t-012301", -12301, 8},
	{"  103.0", 103, 5},
	{"0x31", 0, 1},
	{" 00x0", 0, 3},
	{"9223372036854775807", 1<<63 - 1, 19},
	{"9223372036854775808", 0, 0},
	{"9223372036854775809", 0, 0},
	{"-9223372036854775807", -(1<<63 - 1), 20},
	{"-9223372036854775808", -1 << 63, 20},
	{"-9223372036854775809", 0, 0},
	{" \t", 0, 0},
	{"abc", 0, 0},
}

func TestReadInt64(t *testing.T) {
	for _, test := range readInt64Tests {
		v, n := readInt64([]byte(test.in))

		if v != test.v || n != test.n {
			t.Errorf("readInt64(%#q) -> %d, %d (want %d, %d)",
				test.in, v, n, test.v, test.n)
		}
	}
}
