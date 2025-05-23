package parser

import (
	"strings"
	"unicode"

	"github.com/schors/jsm2tg/text"
)

const DefaultJiraCodeType = "java"

// DetectJiraCodeType - detects if a string has a code block (like "{code:java}")
// and returns the code type and the index of the end of the block.
func DetectJiraCodeType(s, b string) (string, int) {
	if len(s) == 0 {
		return DefaultJiraCodeType, 0
	}

	if len(s) < len(b)+2 {
		return DefaultJiraCodeType, 0
	}

	if s[0] != '{' {
		return DefaultJiraCodeType, 0
	}

	if s[len(b)+1] != '}' && s[len(b)+1] != ':' {
		return DefaultJiraCodeType, 0
	}

	if !strings.HasPrefix(s[1:], b) {
		return DefaultJiraCodeType, 0
	}

	if s[len(b)+1] == '}' {
		return DefaultJiraCodeType, len(b) + 1 + 1
	}

	ok, i := text.DetectRune(s[len(b)+1:], '}')
	if ok {
		for _, r := range s[len(b)+1 : len(b)+1+i] {
			if unicode.IsSpace(r) || unicode.IsControl(r) || r == '=' {
				return DefaultJiraCodeType, len(b) + 1 + i + 1
			}
		}

		return s[len(b)+2 : len(b)+1+i], len(b) + 1 + i + 1
	}

	return DefaultJiraCodeType, 0
}

// DetectLeftBlock - detects if a string has a left block (like "{color:red}")
// and returns true if it does, along with the index of the end of the block.
func DetectLeftBlock(s, b string) (bool, int) {
	if len(s) == 0 {
		return false, 0
	}

	if len(s) < len(b)+2 {
		return false, 0
	}

	if s[0] != '{' {
		return false, 0
	}

	if s[len(b)+1] != '}' && s[len(b)+1] != ':' {
		return false, 0
	}

	if !strings.HasPrefix(s[1:], b) {
		return false, 0
	}

	if s[len(b)+1] == '}' {
		return true, len(b) + 1 + 1
	}

	ok, i := text.DetectRune(s[len(b)+1:], '}')
	if ok {
		return true, len(b) + 1 + i + 1
	}

	return false, 0
}

// DetectLeftBlock - detects if a string has a left block (like "{color}")
// and returns true if it does.
func IsLeftBlock(s, b string) bool {
	ok, _ := DetectLeftBlock(s, b)

	return ok
}

// DetectRightBlock - detects if a string has a right block (like "{color}")
// and returns true if it does, along with the index of the end of the block.
func DetectRightBlock(s, b string) (bool, int) {
	if len(s) == 0 {
		return false, 0
	}

	if len(s) < len(b)+2 {
		return false, 0
	}

	if s[0] != '{' {
		return false, 0
	}

	if s[len(b)+1] != '}' {
		return false, 0
	}

	if !strings.HasPrefix(s, "{"+b+"}") {
		return false, 0
	}

	return true, len(b) + 2
}

// DetectListLine - detects if a string is a list line (like "* " or "# ")
// and returns true if it does, along with the index of the end of the line.
func DetectListLine(s string) (bool, int) {
	if len(s) == 0 {
		return false, 0
	}

	if s[0] != '*' && s[0] != '#' {
		return false, 0
	}

	for i, r := range s {
		if r == '*' || r == '#' {
			continue
		}

		if unicode.IsSpace(r) {
			return true, i
		}

		return false, 0
	}

	return false, 0
}
