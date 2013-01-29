package walnut

import (
	"testing"
	"time"
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
		v, n := readBool(test.in)

		if v != test.v || n != test.n {
			t.Errorf("readBool(%q) -> %v, %v (want %v, %v)",
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
		v, n := readInt64(test.in)

		if v != test.v || n != test.n {
			t.Errorf("readInt64(%q) -> %v, %v (want %v, %v)",
				test.in, v, n, test.v, test.n)
		}
	}
}

var readFloat64Tests = []struct {
	in string
	v  float64
	n  int
}{
	{"", 0, 0},
	{"0.0", 0, 3},
	{"0000.0000", 0, 9},
	{"123.456", 123.456, 7},
	{"  +12.3", 12.3, 7},
	{"  -12.3", -12.3, 7},
	{"  + 12.3", 0, 0},
	{"10.76-", 10.76, 5},
	{"1.3 ", 1.3, 3},
	{" 0.1", 0.1, 4},
	{"1 \t", 0, 0},
	{"1.", 0, 0},
	{".", 0, 0},
	{" 1.0", 1, 4},
	{"1.0.0", 1, 3},
	{" \t32. 0", 0, 0},
	{" -0", 0, 0},
	{"22.222222222222222", 22.22222222222222, 18},
	{"99999999999999974834176.0", 9.999999999999997e+22, 25},
	{"100000000000000000000000.0", 1e+23, 26},
	{"100000000000000000000001.0", 1.0000000000000001e+23, 26},
	{"100000000000000008388608.0", 1.0000000000000001e+23, 26},
	{"100000000000000016777215.0", 1.0000000000000001e+23, 26},
	{"100000000000000016777216.0", 1.0000000000000003e+23, 26},
	{"-22.222222222222222", -22.22222222222222, 19},
	{"-99999999999999974834176.0", -9.999999999999997e+22, 26},
	{"-100000000000000000000000.0", -1e+23, 27},
	{"-100000000000000000000001.0", -1.0000000000000001e+23, 27},
	{"-100000000000000008388608.0", -1.0000000000000001e+23, 27},
	{"-100000000000000016777215.0", -1.0000000000000001e+23, 27},
	{"-100000000000000016777216.0", -1.0000000000000003e+23, 27},
}

func TestReadFloat64(t *testing.T) {
	for _, test := range readFloat64Tests {
		v, n := readFloat64(test.in)

		if v != test.v || n != test.n {
			t.Errorf("readFloat64(%q) -> %v, %v (want %v, %v)",
				test.in, v, n, test.v, test.n)
		}
	}
}

var readStringTests = []struct {
	in string
	v  string
	n  int
}{
	{"", "", 0},
	{`""`, "", 2},
	{` "日本人"`, "日本人", 12},
	{`"a\nb"`, "a\nb", 6},
	{`"\u00FF"`, "ÿ", 8},
	{"\t\"\\xFF\" ", "\xFF", 7},
	{`"\U00010111"`, "\U00010111", 12},
	{`"\U0001011111"`, "\U0001011111", 14},
	{`"'"`, "'", 3},
	{` "\""`, "\"", 5},
	{`"a""`, "a", 3},
	{`"lone`, "", 0},
	{`hello`, "", 0},
	{`"mismatch'`, "", 0},
	{`"\"`, "", 0},
	{`"a" "b"`, "a", 3},
	{"`a`", "", 0},
	{"'b'", "", 0},
}

func TestReadString(t *testing.T) {
	for _, test := range readStringTests {
		v, n := readString(test.in)

		if v != test.v || n != test.n {
			t.Errorf("readString(%q) -> %q, %v (want %q, %v)",
				test.in, v, n, test.v, test.n)
		}
	}
}

var readTimeTests = []struct {
	in string
	n  int
}{
	{"", 0},
	{"1970-01-01 00:00:00 +0000", 25},
	{"2001-02-03 04:05:06 +0000", 25},
	{"1997-08-28 15:30:27.123 +0000", 29},
	{"1997-08-28 14:07:27 -0123", 25},
	{"01:02:03", 0},
	{"1970-01-01", 0},
	{"1970-01-01 00:00:00", 0},
	{"1970-02-48 00:00:00 +0000", 0},
	{"70-01-01 00:00:00", 0},
	{"1970-01-01 00:00:00 UTC", 0},
}

func TestReadTime(t *testing.T) {
	for _, test := range readTimeTests {
		v, n := readTime(test.in)
		e, _ := time.Parse("2006-01-02 15:04:05 -0700", test.in[:n])

		if !e.Equal(v) || n != test.n {
			t.Errorf("readTime(%q) -> %s, %v (want %s, %v)",
				test.in, v, n, e, test.n)
		}
	}
}

var readDurationTests = []struct {
	in string
	v  time.Duration
	n  int
}{
	{"", 0, 0},
	{"0s", 0, 2},
	{"5s", 5 * time.Second, 2},
	{"37s", 37 * time.Second, 3},
	{"010s", 10 * time.Second, 4},
	{"1sm", time.Second, 2},
	{" 3d\t ", 3 * 24 * time.Hour, 3},
	{"10ns", 10 * time.Nanosecond, 4},
	{"10µs", 10 * time.Microsecond, 5},
	{"10μs", 10 * time.Microsecond, 5},
	{"10us", 10 * time.Microsecond, 4},
	{"10ms", 10 * time.Millisecond, 4},
	{"10s", 10 * time.Second, 3},
	{"10m", 10 * time.Minute, 3},
	{"10h", 10 * time.Hour, 3},
	{"10d", 10 * 24 * time.Hour, 3},
	{"10w", 10 * 7 * 24 * time.Hour, 3},
	{"1h1m1s", time.Hour + time.Minute + time.Second, 6},
	{"1h 1m1s", time.Hour + time.Minute + time.Second, 7},
	{"4h30m", 4*time.Hour + 30*time.Minute, 5},
	{"4h 30m", 4*time.Hour + 30*time.Minute, 6},
	{"1s500ms", time.Second + 500*time.Millisecond, 7},
	{"1s 500ms", time.Second + 500*time.Millisecond, 8},
	{"1w1d24h1440m", 10 * 24 * time.Hour, 12},
	{"1w 1d 24h 1440m", 10 * 24 * time.Hour, 15},
	{"10w -3d", 10 * 7 * 24 * time.Hour, 3},
	{"1d 200", 24 * time.Hour, 2},
	{"1h 1m 1.3s", 1*time.Hour + 1*time.Minute, 5},
	{"-3h", 0, 0},
	{"+5m", 0, 0},
	{"300.5h", 0, 0},
	{"1.2d20m", 0, 0},
	{"1s2h", 0, 0},
	{"1200ms 3s", 0, 0},
	{"4h 5d 6w 7m", 0, 0},
	{"2 m", 0, 0},
	{"4 d5 h", 0, 0},
	{"100", 0, 0},
	{"3 4 5ms", 0, 0},
}

func TestReadDurationTests(t *testing.T) {
	for _, test := range readDurationTests {
		v, n := readDuration(test.in)

		if v != test.v || n != test.n {
			t.Errorf("readTime(%q) -> %s, %v (want %s, %v)",
				test.in, v, n, test.v, test.n)
		}
	}
}
