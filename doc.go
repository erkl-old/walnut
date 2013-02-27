// walnut is a sensible configuration format designed with readability and
// compatibility with Go in mind. In addition to the expected primitives, it
// also provides first class support for durations and timestamps.
//
//
// Comments
//
// Comments are indicated by the '#' character (0x23).
//
//     # comments can go on lines of their own
//     foo # or after key groups
//       bar = 123 # or after value keys
//
//
// Indentation
//
// Indentation in walnut is significant -- indented lines extend the parent
// key group, each line joined by a dot. The snippet below illustrates two
// ways of defining the same key.
//
//     abc.def = 123
//     abc
//       def = 123
//
// Indentation for any line must line up with, and optionally extend, one of
// its ancestors.
//
//     parent
//         first = "ok"
//        second = "not ok"
//
// Finally, the first key in any configuration file must not be indented.
//
//
// Keys
//
// Keys fall under two categories: key groups and value keys. The former don't
// have any values assigned to them, but instead serve to namespace further
// indented keys below them.
//
// Keys may contain any rune, except whitespace characters and control
// characters.
//
//     simple = true
//     ♫ = "we're cool with unicode"
//
// A key may only be defined once. Additionally, a value key may not have any
// child keys.
//
//     foo.bar = 1
//     foo.bar.baz = 2 # this key will collide with "key.bar"
//     foo.bars = 3    # this line will not
//
//
// Types
//
// Boolean values are simply expressed as either "true" or "false".
//
//     enabled = true
//     active = false
//
// All numbers must be written in base-10. Integers and floats are distinguished
// by the decimal point, which is required for floats, as is at least one digit
// on either side of it.
//
//     neg = -47
//     zero = 0.0
//     pi = 3.1415
//
// Any valid double-quoted Go string literals is also a valid walnut string.
//
//     plain = "hello"
//     unicode = "\u2603"
//     EOL = "\r\n"
//
// Durations consist of 1..N (value, unit) pairs, where value is a positive
// base-10 integer and unit is one of "ns", "us" (or "µs"), "ms", "s", "m",
// "h", "d" or "w". These pairs must be ordered by the magnitude of their
// units, in descending order, and each unit may only appear once. Pairs may
// be separated by whitespace.
//
//     delay = 2m 30s
//     interval = 750ms
//
// Times are represented in a "2006-01-02 15:04:05 -0700" format. The fraction
// after the second is optional.
//
//     timestamp = 2013-02-25 17:07:46.409 +0100
package walnut
