package utils

import (
	"fmt"
	"go-fiber-api-starter/internal/config"
	"go-fiber-api-starter/internal/models"
	"math/rand"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func CreateJWT(user *models.User) (string, error) {
	claims := &models.JwtUser{
		UserId: user.Id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.LoginDuration)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := []byte(config.API_SECRET)
	return token.SignedString(secret)
}

// https://pkg.go.dev/github.com/golang-jwt/jwt/v5#Parse
func ParseAndVerifyJWT(tokenString string) (*models.JwtUser, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.JwtUser{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.API_SECRET), nil
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

func SendSuccessJSON(c *fiber.Ctx, code int, data any, msg string) error {
	return c.Status(code).JSON(fiber.Map{
		"status":  "success",
		"code":    code,
		"message": msg,
		"error":   "",
		"data":    data,
	})
}

func SendPaginationJSON(c *fiber.Ctx, data any, pagination *models.Pagination, msg string) error {
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"code":    200,
		"message": msg,
		"error":   "",
		"data":    data,
		"pagination": map[string]any{
			"page":       pagination.Page,
			"perPage":    pagination.PerPage,
			"totalPages": pagination.TotalPages,
			"totalCount": pagination.TotalCount,
			"links": map[string]any{
				"self":     pagination.SelfLink,
				"first":    pagination.FirstLink,
				"previous": pagination.PreviousLink,
				"next":     pagination.NextLink,
				"last":     pagination.LastLink,
			},
		},
	})
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

func HasAllowedUsernameChars(str string) bool {
	return !strings.ContainsFunc(str, func(r rune) bool {
		return (r < 'a' || r > 'z') &&
			(r < 'A' || r > 'Z') &&
			(r < '0' || r > '9') &&
			r != '_'
	})
}
