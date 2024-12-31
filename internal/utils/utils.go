package utils

import (
	"fmt"
	"go-fiber-api-starter/internal/models"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateJWT(user *models.User) (string, error) {
	claims := &models.JwtUser{
		UserId:           user.Id,
		UserRole:         user.Role,
		UserStatus:       user.Status,
		RegisteredClaims: jwt.RegisteredClaims{},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := []byte(os.Getenv("SECRET"))
	return token.SignedString(secret)
}

// https://pkg.go.dev/github.com/golang-jwt/jwt/v5#Parse
func ParseAndVerifyJWT(tokenString string) (*models.JwtUser, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.JwtUser{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET")), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return nil, fmt.Errorf("JWT parse error: %v", err.Error())
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token: %v", token)
	}
	payload, ok := token.Claims.(*models.JwtUser)
	if !ok {
		return nil, fmt.Errorf("payload of JWT is of the incorrect type")
	}
	return payload, nil
}

func RandomSixDigitStr() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%d", r.Intn(900000)+100000)
}

func IsAlphanumeric(str string) bool {
	return !strings.ContainsFunc(str, func(r rune) bool {
		return (r < 'a' || r > 'z') &&
			(r < 'A' || r > 'Z') &&
			(r < '0' || r > '9')
	})
}

func IsInteger(str string) bool {
	return !strings.ContainsFunc(str, func(r rune) bool {
		return (r < '0' || r > '9')
	})
}
