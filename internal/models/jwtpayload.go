package models

import "github.com/golang-jwt/jwt/v5"

type JwtUser struct {
	UserId uint `json:"userId"`
	jwt.RegisteredClaims
}

type JwtVerifyEmail struct {
	UserId uint   `json:"userId"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}
