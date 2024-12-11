package serialization

import "go-fiber-api-starter/internal/model"

type SerializedUser struct {
	email      string
	username   string
	userType   string
	userStatus string
	imageURL   string
}

func SerializeUser(user *model.User) SerializedUser {
	return SerializedUser{
		email:      user.Email,
		username:   user.Username,
		userType:   user.UserType,
		userStatus: user.UserStatus,
		imageURL:   *user.ImageURL,
	}
}
