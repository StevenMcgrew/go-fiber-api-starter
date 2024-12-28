package middleware

import (
	"go-fiber-api-starter/internal/db"
	"go-fiber-api-starter/internal/models"

	"github.com/gofiber/fiber/v2"
)

func AttachUser(c *fiber.Ctx) error {
	// Parse userId from path
	id, err := c.ParamsInt("userId")
	if err != nil || id == 0 {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "Error parsing userId from URL",
			"data": map[string]any{"errorMessage": "Error parsing userId from URL"}})
	}
	userId := uint(id)

	// Get user
	user, err := db.GetUserById(userId)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Error getting user from database",
			"data": map[string]any{"errorMessage": err.Error()}})
	}

	// Add user to c.Locals()
	c.Locals("user", &user)

	return c.Next()
}

func AttachUserId(c *fiber.Ctx) error {
	// Parse userId from path
	id, err := c.ParamsInt("userId")
	if err != nil || id == 0 {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "Error parsing userId from URL",
			"data": map[string]any{"errorMessage": "Error parsing userId from URL"}})
	}
	userId := uint(id)

	// Add user id to c.Locals()
	user := &models.User{}
	user.Id = userId
	c.Locals("user", user)

	return c.Next()
}
