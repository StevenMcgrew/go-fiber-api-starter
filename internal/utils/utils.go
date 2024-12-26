package utils

import (
	"fmt"
	"go-fiber-api-starter/internal/enums/jwtuserclaims"
	"go-fiber-api-starter/internal/models"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateUserJWT(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		jwtuserclaims.ID:     user.Id,
		jwtuserclaims.ROLE:   user.Role,
		jwtuserclaims.STATUS: user.Status,
	})
	secret := []byte(os.Getenv("SECRET"))
	return token.SignedString(secret)
}

// https://pkg.go.dev/github.com/golang-jwt/jwt/v5#Parse
func ParseAndVerifyJWT(tokenString string) (*jwt.Token, error) {
	if tokenString == "" { // Need to check for empty string because jwt.Parse will return a nil token instead of setting an error
		return nil, fmt.Errorf("JWT parse error: the tokenString is empty")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET")), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	if err != nil {
		return nil, fmt.Errorf("JWT parse error: %v", err.Error())
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token: %v", token)
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
