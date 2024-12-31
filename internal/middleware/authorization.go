package middleware

import (
	"go-fiber-api-starter/internal/enums/userrole"
	"go-fiber-api-starter/internal/models"

	"github.com/gofiber/fiber/v2"
)

func OnlyAdmin(c *fiber.Ctx) error {
	// Type assert jwtUser
	jwtUser, ok := c.Locals("jwtUser").(*models.User)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Type assert failed for c.Locals",
			"data": map[string]any{"errorMessage": "Incorrect type for c.Locals"}})
	}

	// Only allow admin
	if jwtUser.Role == userrole.ADMIN {
		return c.Next()
	} else {
		return c.Status(403).JSON(fiber.Map{"status": "fail", "message": "The user does not have permission to access this path",
			"data": map[string]any{"errorMessage": "Access denied"}})
	}
}

func OnlyAdminOrOwner(c *fiber.Ctx) error {
	// Get the user that is requesting access
	jwtUser, ok := c.Locals("jwtUser").(*models.User)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Type assert failed for c.Locals",
			"data": map[string]any{"errorMessage": "Incorrect type for c.Locals"}})
	}

	// Get the user to be accessed
	user, ok := c.Locals("user").(*models.User)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "c.Locals('user') should be of type '*models.User'",
			"data": map[string]any{"errorMessage": "Incorrect type for c.Locals('user')"}})
	}

	// Only allow admin or owner
	isAdmin := (jwtUser.Role == userrole.ADMIN)
	isOwner := (jwtUser.Id == user.Id)
	if isAdmin || isOwner {
		return c.Next()
	} else {
		return c.Status(403).JSON(fiber.Map{"status": "fail", "message": "The user does not have permission to access this path",
			"data": map[string]any{"errorMessage": "Access denied"}})
	}

}
