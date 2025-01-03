package middleware

import (
	"go-fiber-api-starter/internal/db"

	"github.com/gofiber/fiber/v2"
)

func AttachUser(c *fiber.Ctx) error {
	// Parse userId from path
	id, err := c.ParamsInt("userId")
	if err != nil || id == 0 {
		return fiber.NewError(400, "Error parsing path parameter: "+err.Error())
	}
	userId := uint(id)

	// Get user
	user, err := db.GetUserById(userId)
	if err != nil {
		return fiber.NewError(500, "Error getting user from database: "+err.Error())
	}

	// Add user to c.Locals()
	c.Locals("user", &user)

	return c.Next()
}
