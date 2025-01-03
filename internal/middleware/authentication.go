package middleware

import (
	"go-fiber-api-starter/internal/db"
	"go-fiber-api-starter/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// Authenticate JWT
func Authn(c *fiber.Ctx) error {
	// Get JWT from header
	authHeader := c.Get("Authorization")
	if len(authHeader) < 7 {
		return fiber.NewError(400, "Authorization header is required")
	}
	bearerPrefix := "Bearer "
	tokenString := authHeader[len(bearerPrefix):]

	// Validate JWT
	payload, err := utils.ParseAndVerifyJWT(tokenString)
	if err != nil {
		return fiber.NewError(401, "Access denied: "+err.Error())
	}

	// Get the user that is requesting access
	jwtUser, err := db.GetUserById(payload.UserId)
	if err != nil {
		return fiber.NewError(500, "Error getting user from database: "+err.Error())
	}

	// Add jwtUser to c.Locals()
	c.Locals("jwtUser", &jwtUser)

	return c.Next()
}
