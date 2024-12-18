package validation

import (
	"strings"
	"unicode/utf8"
)

func isAlphanumeric(str string) bool {
	return !strings.ContainsFunc(str, func(r rune) bool {
		return (r < 'a' || r > 'z') &&
			(r < 'A' || r > 'Z') &&
			(r < '0' || r > '9')
	})
}

func isEmailValid(email string) bool {
	if !strings.Contains(email, "@") {
		return false
	}
	runeCount := utf8.RuneCountInString(email)
	if runeCount < 3 || runeCount > 320 {
		return false
	}
	return true
}

func isUserNameValid(userName string) bool {
	runeCount := utf8.RuneCountInString(userName)
	if runeCount < 3 || runeCount > 20 {
		return false
	}
	return isAlphanumeric(userName)
}

func isPasswordValid(password string) bool {
	bytes := []byte(password)
	length := len(bytes)
	if length < 8 || length > 72 { // bcrypt does not accept passwords longer than 72 bytes
		return false
	}
	return true
}

func doesPasswordRepeatMatch(password string, passwordRepeat string) bool {
	return password == passwordRepeat
}
