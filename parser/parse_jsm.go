package parser

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/schors/jsm2tg/text"
	"github.com/schors/jsm2tg/tg"
)

type tokenType int

const (
	tokenNone tokenType = iota
	tokenItalic
	tokenCitation
	tokenStrike
	tokenUnderline
	tokenCode
	tokenSup
	tokenSub
	tokenBold
	tokenImage
	tokenColor
	tokenPanel
	tokenQuote
	tokenLinkText
	tokenLinkURL
	tokenLinkTextAndURL
	tokenUnsupportedLink
	tokenNoFormat
	tokenCodeBlock
)

type token struct {
	typ               tokenType
	openTag           string // Telegram MarkdownV2 representation, open
	closeTag          string // Telegram MarkdownV2 representation, close
	emergencyCloseTag string // Telegram MarkdownV2 representation, emergency close
}

func (t token) OpenTag() string {
	return t.openTag
}

func (t token) CloseTag() string {
	return t.closeTag
}

func (t token) EmergencyCloseTag() string {
	return t.emergencyCloseTag
}

func (t token) Type() tokenType {
	return t.typ
}

var linkEmergencyCloseTag = "](" + tg.EscapeTelegramLink("https://example.com") + ")"

var tokenMap = map[tokenType]token{
	tokenBold:      {tokenBold, "*", "*", "*"},         // *
	tokenItalic:    {tokenItalic, "_", "_", "_"},       // _
	tokenStrike:    {tokenStrike, "~", "~", "~"},       // -
	tokenUnderline: {tokenUnderline, "__", "__", "__"}, // +
	tokenSup:       {tokenSup, "", "", ""},             // ^
	tokenSub:       {tokenSub, "", "", ""},             // ~
	tokenCitation:  {tokenCitation, "_", "_", "_"},     // ??

	tokenImage: {tokenImage, "", "", ""},   // !
	tokenCode:  {tokenCode, "`", "`", "`"}, // ``

	tokenColor: {tokenColor, "", "", ""}, // {color:xxx}
	tokenPanel: {tokenPanel, "", "", ""}, // {panel:xxx}
	tokenQuote: {tokenQuote, "", "", ""}, // {quote:xxx}

	tokenLinkText:        {tokenLinkText, "", "", linkEmergencyCloseTag},       // [text
	tokenLinkURL:         {tokenLinkURL, "", "", linkEmergencyCloseTag},        // url]
	tokenLinkTextAndURL:  {tokenLinkTextAndURL, "", "", linkEmergencyCloseTag}, // [scheme://uri]
	tokenUnsupportedLink: {tokenUnsupportedLink, "", "", ""},                   // [~username]

	tokenNoFormat:  {tokenNoFormat, "```", "```", "```"},  // no format
	tokenCodeBlock: {tokenCodeBlock, "```", "```", "```"}, // code block
}

type tokenStack []token

func (s *tokenStack) Push(t token) {
	*s = append(*s, t)
}

func (s *tokenStack) Pop() token {
	if len(*s) == 0 {
		return token{tokenNone, "", "", ""}
	}

	t := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]

	return t
}

func (s *tokenStack) TopType() tokenType {
	if len(*s) == 0 {
		return tokenNone
	}

	return (*s)[len(*s)-1].typ
}

func getToken(char rune) token {
	switch char {
	case '*':
		return tokenMap[tokenBold]
	case '_':
		return tokenMap[tokenItalic]
	case '-':
		return tokenMap[tokenStrike]
	case '+':
		return tokenMap[tokenUnderline]
	case '^':
		return tokenMap[tokenSup]
	case '~':
		return tokenMap[tokenSub]
	default:
		return token{typ: tokenNone}
	}
}

func newToken(typ tokenType) token {
	t, ok := tokenMap[typ]
	if !ok {
		return token{tokenNone, "", "", ""}
	}

	return t
}

func ConvertJiraToTgMarkup(input string) string {
	var (
		result strings.Builder
		stack  tokenStack

		linkURL strings.Builder
	)

	lineStart := true
	pr := rune(0)

	i := 0
	for i < len(input) {
		r, sz := utf8.DecodeRuneInString(input[i:])

		if stack.TopType() == tokenNoFormat || stack.TopType() == tokenCodeBlock {
			pr = r

			if r == '\\' && i+sz < len(input) {
				nr, nsz := utf8.DecodeRuneInString(input[i+sz:])
				if unicode.IsSpace(nr) {
					result.WriteString(tg.EscapeTelegramCode(string(r)))

					i += sz

					continue
				}

				result.WriteString(tg.EscapeTelegramCode(string(nr)))

				i += sz + nsz

				continue
			}

			if r == '{' && i+len("{noformat}") <= len(input) && strings.HasPrefix(input[i:], "{noformat}") {
				result.WriteString(stack.Pop().CloseTag())

				i += len("{noformat}")

				continue
			}

			if r == '{' && i+len("{code}") <= len(input) && strings.HasPrefix(input[i:], "{code}") {
				result.WriteString(stack.Pop().CloseTag())

				i += len("{code}")

				continue
			}

			result.WriteString(tg.EscapeTelegramCode(string(r)))

			i += sz

			continue
		}

		// close "image" tag
		if stack.TopType() == tokenImage {
			if r == '!' && pr != '\\' {
				stack.Pop()

				pr = r
				i += sz

				continue
			}

			pr = r
			i += sz

			continue
		}

		// close "code" tag
		if stack.TopType() == tokenCode {
			if r == '}' && pr != '\\' && i+sz < len(input) && strings.HasPrefix(input[i:], "}}") {
				result.WriteString(stack.Pop().CloseTag())

				pr = r
				i += len("}}")

				continue
			}

			result.WriteString(tg.EscapeTelegramCode(string(r)))

			pr = r
			i += sz

			continue
		}

		if stack.TopType() == tokenUnsupportedLink {
			if r == ']' && pr != '\\' {
				stack.Pop()

				pr = r
				i += sz

				continue
			}

			pr = r
			i += sz

			continue
		}

		if r == '\\' && i+sz < len(input) {
			pr = r

			// Handle escaped characters explicitly
			nr, nsz := utf8.DecodeRuneInString(input[i+sz:])
			if unicode.IsSpace(nr) {
				result.WriteString(tg.EscapeTelegram(string(r)))
				if stack.TopType() == tokenLinkTextAndURL ||
					stack.TopType() == tokenLinkURL {
					linkURL.WriteRune(r)
				}

				i += sz

				continue
			}

			result.WriteString(tg.EscapeTelegram(string(nr)))
			if stack.TopType() == tokenLinkTextAndURL ||
				stack.TopType() == tokenLinkURL {
				linkURL.WriteRune(nr)
			}

			i += sz + nsz

			continue
		}

		if r != ']' && (stack.TopType() == tokenLinkTextAndURL ||
			stack.TopType() == tokenLinkURL) {
			linkURL.WriteRune(r)

			if stack.TopType() == tokenLinkURL {
				i += sz

				continue
			}
		}

		switch {
		case i == 0:
			lineStart = true
		case pr == '\n':
			lineStart = true
		default:
			lineStart = false
		}

		if lineStart {
			switch {
			case r == 'h' && i+sz < len(input) && (strings.HasPrefix(input[i:], "h1. ") ||
				strings.HasPrefix(input[i:], "h2. ") || strings.HasPrefix(input[i:], "h3. ") ||
				strings.HasPrefix(input[i:], "h4. ") || strings.HasPrefix(input[i:], "h5. ") ||
				strings.HasPrefix(input[i:], "h6. ")):
				// Handle heading markers
				j := text.FindEndOfString(input[i:])

				result.WriteString("*" + tg.EscapeTelegram(input[i:i+j]) + "*")

				i += j
				pr = r

				continue
			}
		}

	OUTER:
		switch {
		case r == '!':
			// only opening tag
			if nr, nsz := utf8.DecodeRuneInString(input[i+sz:]); nsz == 0 || unicode.IsSpace(nr) {
				result.WriteString(tg.EscapeTelegram(string(r)))
				i += sz

				break
			}

			stack.Push(newToken(tokenImage))
			i += sz
		case r == '\n' && stack.TopType() == tokenQuote:
			result.WriteString("\n>")
			i += sz
		case r == '|' && stack.TopType() == tokenLinkText:
			// result.WriteString(string(']'))

			stack.Pop()
			stack.Push(newToken(tokenLinkURL))

			linkURL.Reset()

			i += sz
		case r == ']' && (stack.TopType() == tokenLinkText ||
			stack.TopType() == tokenLinkTextAndURL || stack.TopType() == tokenLinkURL):

			if stack.TopType() == tokenLinkText {
				linkURL.WriteString(tg.EscapeTelegramLink("https://example.com"))
			}

			if stack.TopType() == tokenLinkURL {
				result.WriteString(string(']'))
			}

			if stack.TopType() == tokenLinkTextAndURL ||
				stack.TopType() == tokenLinkText {

				result.WriteString(string(r))
			}

			result.WriteString("(" + tg.EscapeTelegramLink(linkURL.String()) + ")")

			stack.Pop()

			linkURL.Reset()

			i += sz
		case r == '[' && i+sz < len(input) && (stack.TopType() != tokenLinkText &&
			stack.TopType() != tokenLinkTextAndURL && stack.TopType() != tokenLinkURL &&
			stack.TopType() != tokenUnsupportedLink):
			linkURL.Reset()

			if input[i+sz] == '^' ||
				input[i+sz] == '#' ||
				input[i+sz] == '~' ||
				strings.HasPrefix(input[i+sz:], "file://") {
				stack.Push(newToken(tokenUnsupportedLink))

				i += sz

				break
			}

			okScheme, _ := text.DetectSchemeFast(input[i+sz:])
			okDelimiter, _ := text.DetectRune(input[i+sz:], '|')
			if okScheme && !okDelimiter {
				stack.Push(newToken(tokenLinkTextAndURL))
				result.WriteString(string(r))

				i += sz

				break
			}

			if !okDelimiter {
				stack.Push(newToken(tokenUnsupportedLink))

				i += sz

				break
			}

			stack.Push(newToken(tokenLinkText))
			result.WriteString(string(r))
			i += sz
		case r == '{' && i+len("{noformat}") <= len(input) && stack.TopType() == tokenNone && strings.HasPrefix(input[i:], "{noformat}"):
			// only opening tag
			currentToken := newToken(tokenNoFormat)

			stack.Push(currentToken)
			result.WriteString(currentToken.OpenTag())

			i += len("{noformat}")
		case r == '{' && i+len("{code}") <= len(input) && stack.TopType() == tokenNone && IsLeftBlock(input[i:], "code"):
			// only opening tag
			currentToken := newToken(tokenCodeBlock)

			stack.Push(currentToken)

			if !lineStart {
				result.WriteString("\n")
			}

			result.WriteString(currentToken.OpenTag())

			lang, j := DetectJiraCodeType(input[i:], "code")
			if lang != "" {
				result.WriteString(tg.EscapeTelegramCode(lang) + "\n")
			}

			i += j
		case r == '{' && i+sz < len(input) && stack.TopType() == tokenNone && strings.HasPrefix(input[i:], "{{"):
			// only opening tag
			currentToken := newToken(tokenCode)

			stack.Push(currentToken)
			result.WriteString(currentToken.OpenTag())

			i += len("{{")

		case r == '{' && i+sz < len(input) && (stack.TopType() == tokenPanel ||
			stack.TopType() == tokenColor || stack.TopType() == tokenQuote):
			switch {
			case stack.TopType() == tokenPanel:
				ok, j := DetectRightBlock(input[i:], "panel")
				if ok {
					stack.Pop()

					i += j

					break OUTER
				}
			case stack.TopType() == tokenColor:
				ok, j := DetectRightBlock(input[i:], "color")
				if ok {
					stack.Pop()

					i += j

					break OUTER
				}
			case stack.TopType() == tokenQuote:
				if strings.HasPrefix(input[i:], "{quote}") {
					stack.Pop()

					result.WriteString("\n")

					i += len("{quote}")

					break OUTER
				}
			}

			result.WriteString(tg.EscapeTelegram(string(r)))
			i += sz
		case r == '{' && i+sz < len(input):
			ok, j := DetectLeftBlock(input[i:], "color")
			if ok {
				stack.Push(newToken(tokenColor))

				i += j

				break
			}

			ok, j = DetectLeftBlock(input[i:], "panel")
			if ok {
				stack.Push(newToken(tokenPanel))

				i += j

				break
			}

			ok, j = DetectLeftBlock(input[i:], "anchor")
			if ok {
				i += j

				break
			}

			if strings.HasPrefix(input[i:], "{quote}") {
				stack.Push(newToken(tokenQuote))

				result.WriteString("\n>")

				i += len("{quote}")

				break
			}

			result.WriteString(tg.EscapeTelegram(string(r)))
			i += sz
		case r == '?' && i+sz < len(input) && strings.HasPrefix(input[i:], "??"):
			pattern := "??"
			currentToken := newToken(tokenCitation)

			if stack.TopType() == currentToken.Type() {
				result.WriteString(currentToken.CloseTag())
				stack.Pop()

				i += len(pattern)

				break
			}

			stack.Push(currentToken)
			result.WriteString(currentToken.OpenTag())

			i += len(pattern)
		case r == '#' && i+sz < len(input) && lineStart && stack.TopType() == tokenNone:
			ok, j := DetectListLine(input[i:])
			if ok {
				result.WriteString(input[i : i+j])
				i += j

				break
			}

			result.WriteString(tg.EscapeTelegram(string(r)))
			i += sz
		case r == '*' || r == '_' || r == '-' || r == '+' || r == '^' || r == '~':
			if r == '*' && lineStart && stack.TopType() == tokenNone {
				ok, j := DetectListLine(input[i:])
				if ok {
					result.WriteString(input[i : i+j])
					i += j

					break
				}
			}

			currentToken := getToken(r)

			if stack.TopType() == currentToken.Type() {
				result.WriteString(currentToken.CloseTag())
				stack.Pop()

				i += sz

				break
			}

			stack.Push(currentToken)
			result.WriteString(currentToken.OpenTag())

			i += sz
		default:
			result.WriteString(tg.EscapeTelegram(string(r)))

			i += sz
		}

		pr = r
	}

	for len(stack) > 0 {
		result.WriteString(stack.Pop().EmergencyCloseTag())
	}

	return result.String()
}
