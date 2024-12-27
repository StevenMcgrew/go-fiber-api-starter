package validation

import (
	"go-fiber-api-starter/internal/enums/userrole"
	"go-fiber-api-starter/internal/enums/userstatus"
	"go-fiber-api-starter/internal/models"
)

func ValidateUserUpdate(u *models.UserUpdate) []string {
	m := make([]string, 0, 4)
	if !IsEmailValid(u.Email) {
		m = append(m, "Email address input is invalid")
	}
	if !IsUserNameValid(u.Username) {
		m = append(m, "UserName input is invalid")
	}
	if !IsPasswordValid(u.Password) {
		m = append(m, "Password input is invalid")
	}
	if u.Role != userrole.ADMIN &&
		u.Role != userrole.REGULAR {
		m = append(m, "User role input is invalid")
	}
	if u.Status != userstatus.ACTIVE &&
		u.Status != userstatus.DELETED &&
		u.Status != userstatus.SUSPENDED &&
		u.Status != userstatus.UNVERIFIED {
		m = append(m, "User status input is invalid")
	}
	if !IsUrlValid(u.ImageUrl) {
		m = append(m, "Image URL input is invalid")
	}
	if len(m) > 0 {
		return m
	}
	return nil
}
