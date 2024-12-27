package middleware

import (
	"go-fiber-api-starter/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// Authenticate JWT
func Authn(c *fiber.Ctx) error {
	// Get token from header
	authHeader := c.Get("Authorization")
	if len(authHeader) < 7 {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "An Authorization header is required to be set for this path",
			"data": map[string]any{"errorMessage": "No value found for Authorization header"}})
	}
	bearerPrefix := "Bearer "
	tokenString := authHeader[len(bearerPrefix):]

	// Validate payload
	payload, err := utils.ParseAndVerifyJWT(tokenString)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"status": "fail", "message": "Invalid token. Access denied.",
			"data": map[string]any{"errorMessage": err.Error()}})
	}

	// Add payload to c.Locals()
	c.Locals("jwtPayload", payload)

	return c.Next()
}
