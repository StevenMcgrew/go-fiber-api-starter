package utils

import (
	"fmt"
	"go-fiber-api-starter/internal/enums/jwtclaimkeys"
	"go-fiber-api-starter/internal/models"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/golang-jwt/jwt/v5"
)

func CreateUserJWT(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		jwtclaimkeys.USER_ID:     user.Id,
		jwtclaimkeys.USER_TYPE:   user.UserType,
		jwtclaimkeys.USER_STATUS: user.UserStatus,
	})
	return token.SignedString([]byte(os.Getenv("SECRET")))
}

func ParseAndVerifyJWT(tokenString string) (*jwt.Token, error) {
	// https://pkg.go.dev/github.com/golang-jwt/jwt/v5#Parse
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET")), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	if !token.Valid {
		return nil, fmt.Errorf("JWT error: %v", err.Error())
	}
	return token, nil
}

func IsAlphanumeric(str string) bool {
	return !strings.ContainsFunc(str, func(r rune) bool {
		return (r < 'a' || r > 'z') &&
			(r < 'A' || r > 'Z') &&
			(r < '0' || r > '9')
	})
}

func IsEmailValid(email string) bool {
	if !strings.Contains(email, "@") {
		return false
	}
	runeCount := utf8.RuneCountInString(email)
	if runeCount < 3 || runeCount > 320 {
		return false
	}
	return true
}

func IsUserNameValid(userName string) bool {
	runeCount := utf8.RuneCountInString(userName)
	if runeCount < 3 || runeCount > 20 {
		return false
	}
	return IsAlphanumeric(userName)
}

func IsPasswordValid(password string) bool {
	bytes := []byte(password)
	length := len(bytes)
	if length < 8 || length > 72 { // bcrypt does not accept passwords longer than 72 bytes
		return false
	}
	return true
}

func DoesPasswordRepeatMatch(password string, passwordRepeat string) bool {
	return password == passwordRepeat
}
