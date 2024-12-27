package models

import "github.com/golang-jwt/jwt/v5"

type JwtPayload struct {
	UserId     uint   `json:"userId"`
	UserRole   string `json:"userRole"`
	UserStatus string `json:"userStatus"`
	jwt.RegisteredClaims
}
