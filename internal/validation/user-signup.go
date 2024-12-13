package validation

import (
	"go-fiber-api-starter/internal/model"
)

func ValidateUserSignup(u *model.UserSignup) []string {
	m := make([]string, 0, 4)
	if !isEmailValid(u.Email) {
		m = append(m, "Email is invalid")
	}
	if !isUserNameValid(u.UserName) {
		m = append(m, "UserName is invalid")
	}
	if !isPasswordValid(u.Password) {
		m = append(m, "Password is invalid")
	}
	if !doesPasswordRepeatMatch(u.Password, u.PasswordRepeat) {
		m = append(m, "PasswordRepeat does not match Password")
	}
	if len(m) > 0 {
		return m
	}
	return nil
}
