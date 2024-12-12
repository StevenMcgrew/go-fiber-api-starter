package model

type User struct {
	Email          string `json:"email"`
	UserName       string `json:"userName"`
	Password       string `json:"password"`
	PasswordRepeat string `json:"passwordRepeat"`
	UserType       string `json:"userType"`
	UserStatus     string `json:"userStatus"`
	ImageUrl       string `json:"imageUrl"`
}
