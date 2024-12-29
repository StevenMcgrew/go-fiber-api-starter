package validation

import (
	"go-fiber-api-starter/internal/models"
)

func ValidateUserSignUp(u *models.UserSignUp) []string {
	m := make([]string, 0, 4)
	if !IsEmailValid(u.Email) {
		m = append(m, "Email is invalid")
	}
	if !IsUsernameValid(u.Username) {
		m = append(m, "UserName is invalid")
	}
	if !IsPasswordValid(u.Password) {
		m = append(m, "Password is invalid")
	}
	if u.Password != u.PasswordRepeat {
		m = append(m, "PasswordRepeat does not match Password")
	}
	if len(m) > 0 {
		return m
	}
	return nil
}
