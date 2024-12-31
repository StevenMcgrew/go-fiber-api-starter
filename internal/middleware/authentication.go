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
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "An Authorization header is required to be set for this path",
			"data": map[string]any{"errorMessage": "No value found for Authorization header"}})
	}
	bearerPrefix := "Bearer "
	tokenString := authHeader[len(bearerPrefix):]

	// Validate JWT
	payload, err := utils.ParseAndVerifyJWT(tokenString)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"status": "fail", "message": "Invalid token. Access denied.",
			"data": map[string]any{"errorMessage": err.Error()}})
	}

	// Get the user that is requesting access
	jwtUser, err := db.GetUserById(payload.UserId)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Error getting user from database",
			"data": map[string]any{"errorMessage": err.Error()}})
	}

	// Add jwtUser to c.Locals()
	c.Locals("jwtUser", &jwtUser)

	return c.Next()
}
