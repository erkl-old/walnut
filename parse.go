package walnut

import (
	"strings"
)

const whitespace = " \t\n\v\f\r\u0085\u00A0"

// Outlines a "key = value" assignment.
type configDefinition struct {
	key, value string
	line       int
}

// Generates a map of resolved keys and raw string values from a byte slice.
// If the second argument is not 0, an indentation error was detected on
// that line (1 being the first line).
func parseConfig(buf []byte) ([]configDefinition, int) {
	lines := strings.Split(string(buf), "\n")
	raw := make([]configDefinition, 0)

	// collapse lines without any content
	for i, line := range lines {
		lines[i] = collapse(line)
	}

	parents := make([]string, 0)
	indents := make([]string, 0)
	first := true

	for n, line := range lines {
		if line == "" {
			continue
		}

		// line numbers should be 1-indexed
		n++

		k := selectKey(line)
		i := selectIndentation(line)
		d := calculateDepth(indents, i)

		// check for invalid indentation
		if d == -1 || (d == len(indents) && !first) {
			return nil, n
		}

		// trim now redundant levels
		if d < len(indents) {
			parents = parents[:d]
			indents = indents[:d]
		}

		// push the key and indentation onto their respective stacks
		parents = append(parents, k)
		indents = append(indents, i)

		// if the line contains an assignment, record the value
		if strings.ContainsRune(line, '=') {
			raw = append(raw, configDefinition{
				key:   strings.Join(parents, "."),
				value: selectValue(line),
				line:  n,
			})

			first = false
			continue
		}

		first = true
	}

	return raw, 0
}

// Trims trailing whitespace or, in the case of comment lines, returns
// an empty string.
func collapse(input string) string {
	s := strings.TrimRight(input, whitespace)

	for _, r := range s {
		if strings.ContainsRune(whitespace, r) {
			continue
		}

		// comment detected, blank this line
		if r == '#' {
			break
		}

		// if the first non-whitespace character @todo
		return input
	}

	return ""
}

// Returns the "key" component from a "key = value" string.
func selectKey(input string) string {
	if eq := strings.IndexRune(input, '='); eq != -1 {
		input = input[:eq]
	}

	return strings.Trim(input, whitespace)
}

// Returns the "value" component from a "key = value" string.
func selectValue(input string) string {
	if eq := strings.IndexRune(input, '='); eq != -1 {
		input = input[eq+1:]
	}

	return strings.Trim(input, whitespace)
}

// Returns the string's whitespace prefix.
func selectIndentation(input string) string {
	end := strings.IndexFunc(input, func(r rune) bool {
		return strings.IndexRune(whitespace, r) == -1
	})

	if end == -1 {
		return ""
	}

	return input[:end]
}

// Given a list of previous indentation levels, finds the provided indentation
// level's depth value. A depth of 0 represents the lowest possible level of
// indentation. Returns -1 on errors caused by illegal indentation.
func calculateDepth(parents []string, current string) int {
	if current == "" {
		return 0
	}

	// the base indentation level must be an empty string
	if len(parents) == 0 {
		return -1
	}

	for i, prefix := range parents {
		switch {
		case current == prefix:
			return i
		case !strings.HasPrefix(current, prefix):
			return -1
		}
	}

	// if we get this far, the current line is further indented
	// than its parent
	return len(parents)
}
