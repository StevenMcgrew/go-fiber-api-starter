package middleware

import (
	"go-fiber-api-starter/internal/enums/jwtuserclaims"
	"go-fiber-api-starter/internal/enums/userrole"
	"go-fiber-api-starter/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
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

	// Validate token
	token, err := utils.ParseAndVerifyJWT(tokenString)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"status": "fail", "message": "Invalid token. Access not allowed.",
			"data": map[string]any{"errorMessage": err.Error()}})
	}

	// Type assert claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "Incorrect claims type for JWT",
			"data": map[string]any{"errorMessage": "JWT claims must be of type 'jwt.MapClaims'"}})
	}

	// Attach token payload to c.Locals
	c.Locals("user", map[string]any{
		jwtuserclaims.ID:     claims[jwtuserclaims.ID],
		jwtuserclaims.ROLE:   claims[jwtuserclaims.ROLE],
		jwtuserclaims.STATUS: claims[jwtuserclaims.STATUS],
	})

	return c.Next()
}

// Only allow users with 'admin' role
func AdminOnly(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(map[string]any)
	if !ok {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "Incorrect type for 'user' in fiber's c.Locals() method",
			"data": map[string]any{"errorMessage": "The type '%T' is incorrect"}})
	}

	if user["role"] != userrole.ADMIN {
		return c.Status(403).JSON(fiber.Map{"status": "fail", "message": "The user does not have the proper role to access this path",
			"data": map[string]any{"errorMessage": "Access denied"}})
	}

	return c.Next()
}
