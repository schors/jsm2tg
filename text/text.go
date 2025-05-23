package text

import (
	"strings"
	"unicode"
)

// FindEndOfString - finds the first end of a string
// in a text block and returns the index of the end of the string.
func FindEndOfString(input string) int {
	for i, r := range input {
		if r == '\n' {
			return i
		}
	}

	return 0
}

// DetectSchemeFast - detects if a string has a scheme (like "http://")
// and returns true if it does, along with the scheme.
// RFC 3986: scheme = ALPHA *( ALPHA / DIGIT / "+" / "-" / "." )
func DetectSchemeFast(s string) (bool, string) {
	colon := strings.IndexByte(s, ':')
	if colon <= 0 { // no colon or starts with ':' ⇒ no scheme
		return false, ""
	}

	for i, r := range s[:colon] {
		switch {
		case i == 0 && !unicode.IsLetter(r):
			return false, "" // first char must be a letter
		case i > 0 && !(unicode.IsLetter(r) || unicode.IsDigit(r) ||
			r == '+' || r == '-' || r == '.'):
			return false, "" // invalid char in scheme
		}
	}

	return true, s[:colon]
}

// DetectRune - detects if a string has a specific rune
// (first occurrence) taking into account escape sequences
// and returns true if it does, along with the index of the rune.
func DetectRune(s string, b rune) (bool, int) {
	if len(s) == 0 {
		return false, 0
	}

	skip := false

	for i, r := range s {
		if skip {
			skip = false

			continue
		}

		if r == '\\' {
			skip = true

			continue
		}

		if r == b {
			return true, i
		}
	}

	return false, 0
}
