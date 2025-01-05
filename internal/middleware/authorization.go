package middleware

import (
	"go-fiber-api-starter/internal/enums/userrole"
	"go-fiber-api-starter/internal/models"

	"github.com/gofiber/fiber/v2"
)

func OnlyAdmin(c *fiber.Ctx) error {
	// Type assert inquirer
	inquirer, ok := c.Locals("inquirer").(*models.User)
	if !ok {
		return fiber.NewError(500, `Type assertion failed for c.Locals("inquirer")`)
	}

	// Only allow admin
	if inquirer.Role == userrole.ADMIN {
		return c.Next()
	} else {
		return fiber.NewError(403, "The user does not have permission to access this resource")
	}
}

func OnlyAdminOrOwner(c *fiber.Ctx) error {
	// Get the user that is requesting access
	inquirer, ok := c.Locals("inquirer").(*models.User)
	if !ok {
		return fiber.NewError(500, `Type assertion failed for c.Locals("inquirer")`)
	}

	// Get the user to be accessed
	user, ok := c.Locals("user").(*models.User)
	if !ok {
		return fiber.NewError(500, `Type assertion failed for c.Locals("user")`)
	}

	// Allow admin or owner
	isAdmin := (inquirer.Role == userrole.ADMIN)
	isOwner := (inquirer.Id == user.Id)
	if isAdmin || isOwner {
		return c.Next()
	}

	// Reject others
	return fiber.NewError(403, "The user does not have permission to access this resource")
}
