package walnut

// Returns the position of the first `needle` in `haystack`.
func indexOf(haystack []byte, needle byte) int {
	for i, b := range haystack {
		if b == needle {
			return i
		}
	}
	return -1
}

// Cuts a slice of `input` without any leading or trailing whitespace.
func trim(input []byte) []byte {
	for l := len(input); l > 0; l-- {
		if isSpace(input[0]) {
			input = input[1:]
		} else if isSpace(input[l-1]) {
			input = input[:l-1]
		} else {
			break
		}
	}

	return input
}

// Returns true if `b` is a whitespace character.
func isSpace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r'
}
