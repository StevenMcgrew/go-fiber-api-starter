package utils

import (
	"fmt"
	"go-fiber-api-starter/internal/enums/jwtclaimkeys"
	"go-fiber-api-starter/internal/models"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateUserJWT(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		jwtclaimkeys.USER_ID:     user.Id,
		jwtclaimkeys.USER_TYPE:   user.Role,
		jwtclaimkeys.USER_STATUS: user.Status,
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

func DoesPasswordRepeatMatch(password string, passwordRepeat string) bool {
	return password == passwordRepeat
}
