package models

import (
	"time"
)

type User struct {
	Id        uint       `json:"id" form:"id"`
	Email     string     `json:"email" form:"email"`
	Username  string     `json:"username" form:"username"`
	Password  string     `json:"password" form:"password"`
	OTP       string     `json:"otp" form:"otp"`
	Role      string     `json:"role" form:"role"`
	Status    string     `json:"status" form:"status"`
	ImageUrl  string     `json:"imageUrl" form:"imageUrl"`
	CreatedAt time.Time  `json:"createdAt" form:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt" form:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt" form:"deletedAt"`
}

type UserSignUp struct {
	Email          string `json:"email" form:"email"`
	Username       string `json:"username" form:"username"`
	Password       string `json:"password" form:"password"`
	PasswordRepeat string `json:"passwordRepeat" form:"passwordRepeat"`
}

type UserResponse struct {
	Id        uint       `json:"id" form:"id"`
	Email     string     `json:"email" form:"email"`
	Username  string     `json:"username" form:"username"`
	Role      string     `json:"role" form:"role"`
	Status    string     `json:"status" form:"status"`
	ImageUrl  string     `json:"imageUrl" form:"imageUrl"`
	CreatedAt time.Time  `json:"createdAt" form:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt" form:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt" form:"deletedAt"`
}

type UserUpdate struct {
	Email    string `json:"email" form:"email"`
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
	Role     string `json:"role" form:"role"`
	Status   string `json:"status" form:"status"`
	ImageUrl string `json:"imageUrl" form:"imageUrl"`
}
