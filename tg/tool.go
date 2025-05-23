package tg

import "strings"

var TelegramSpecialChars = "_*[]()~`>#+-=|{}.!\\"

func EscapeTelegram(text string) string {
	var escaped strings.Builder

	for _, ch := range text {
		if strings.ContainsRune(TelegramSpecialChars, ch) {
			escaped.WriteRune('\\')
		}

		escaped.WriteRune(ch)
	}

	return escaped.String()
}

func EscapeTelegramCode(text string) string {
	var escaped strings.Builder

	for _, ch := range text {
		if ch == '`' || ch == '\\' {
			escaped.WriteRune('\\')
		}

		escaped.WriteRune(ch)
	}

	return escaped.String()
}

func EscapeTelegramLink(text string) string {
	var escaped strings.Builder

	for _, ch := range text {
		if ch == '(' || ch == '\\' {
			escaped.WriteRune('\\')
		}

		escaped.WriteRune(ch)
	}

	return escaped.String()
}
