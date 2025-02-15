package validation

import (
	"go-fiber-api-starter/internal/utils"
	"net/url"
	"strings"
	"unicode/utf8"
)

func IsOtpValid(otp string) bool {
	runeCount := utf8.RuneCountInString(otp)
	if runeCount != 6 {
		return false
	}
	return utils.IsInteger(otp)
}

func IsEmailValid(email string) bool {
	if !strings.Contains(email, "@") {
		return false
	}
	runeCount := utf8.RuneCountInString(email)
	if runeCount < 3 || runeCount > 320 {
		return false
	}
	return true
}

func IsUsernameValid(username string) bool {
	runeCount := utf8.RuneCountInString(username)
	if runeCount < 3 || runeCount > 20 {
		return false
	}
	return utils.HasAllowedUsernameChars(username)
}

func IsPasswordValid(password string) bool {
	bytes := []byte(password)
	length := len(bytes)
	if length < 8 || length > 72 { // bcrypt does not accept passwords longer than 72 bytes
		return false
	}
	return true
}

func IsUrlValid(URL string) bool {
	if len(URL) > 2048 {
		return false
	}
	_, err := url.Parse("sfdsf")
	return err == nil
}
