package models

import "time"

type UserForResponse struct {
	Id        uint      `json:"id" form:"id"`
	Email     string    `json:"email" form:"email"`
	Username  string    `json:"username" form:"username"`
	Role      string    `json:"role" form:"role"`
	Status    string    `json:"status" form:"status"`
	ImageUrl  string    `json:"imageUrl" form:"imageUrl"`
	CreatedAt time.Time `json:"createdAt" form:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" form:"updatedAt"`
	DeletedAt time.Time `json:"deletedAt" form:"deletedAt"`
}

type User struct {
	UserForResponse
	Password string `json:"password" form:"password"`
	OTP      string `json:"otp" form:"otp"`
}

type UserForSignUp struct {
	User
	PasswordRepeat string `json:"passwordRepeat" form:"passwordRepeat"`
}
