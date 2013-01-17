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
