package models

import "github.com/golang-jwt/jwt/v5"

type JwtUser struct {
	UserId     uint   `json:"userId"`
	UserRole   string `json:"userRole"`
	UserStatus string `json:"userStatus"`
	jwt.RegisteredClaims
}

type JwtVerifyEmail struct {
	UserId uint   `json:"userId"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}
