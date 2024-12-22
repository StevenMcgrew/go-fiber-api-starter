package validation

import (
	"go-fiber-api-starter/internal/models"
	"go-fiber-api-starter/internal/utils"
)

func ValidateUserSignup(u *models.UserForSignUp) []string {
	m := make([]string, 0, 4)
	if !utils.IsEmailValid(u.Email) {
		m = append(m, "Email is invalid")
	}
	if !utils.IsUserNameValid(u.Username) {
		m = append(m, "UserName is invalid")
	}
	if !utils.IsPasswordValid(u.Password) {
		m = append(m, "Password is invalid")
	}
	if !utils.DoesPasswordRepeatMatch(u.Password, u.PasswordRepeat) {
		m = append(m, "PasswordRepeat does not match Password")
	}
	if len(m) > 0 {
		return m
	}
	return nil
}
