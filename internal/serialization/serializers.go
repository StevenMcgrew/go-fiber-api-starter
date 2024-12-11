package serialization

import (
	"go-fiber-api-starter/internal/model"
)

type SerializedUser struct {
	Email      string `json:"email"`
	Username   string `json:"username"`
	UserType   string `json:"userType"`
	UserStatus string `json:"userStatus"`
	ImageURL   string `json:"imageURL"`
}

func SerializeUser(user *model.User) SerializedUser {
	return SerializedUser{
		Email:      user.Email,
		Username:   user.Username,
		UserType:   user.UserType,
		UserStatus: user.UserStatus,
		ImageURL:   *user.ImageURL,
	}
}
