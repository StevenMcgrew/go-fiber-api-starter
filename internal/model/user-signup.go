package model

type UserSignup struct {
	Email          string `json:"email" form:"email"`
	UserName       string `json:"userName" form:"userName"`
	Password       string `json:"password" form:"password"`
	PasswordRepeat string `json:"passwordRepeat" form:"passwordRepeat"`
}
