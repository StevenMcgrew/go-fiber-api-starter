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
		return fiber.NewError(500, `Type assertion failed for c.Locals("jwtUser")`)
	}

	// Only allow admin
	if jwtUser.Role == userrole.ADMIN {
		return c.Next()
	} else {
		return fiber.NewError(403, "The user does not have permission to access this resource")
	}
}

func OnlyAdminOrOwner(c *fiber.Ctx) error {
	// Get the user that is requesting access
	jwtUser, ok := c.Locals("jwtUser").(*models.User)
	if !ok {
		return fiber.NewError(500, `Type assertion failed for c.Locals("jwtUser")`)
	}

	// Get the user to be accessed
	user, ok := c.Locals("user").(*models.User)
	if !ok {
		return fiber.NewError(500, `Type assertion failed for c.Locals("user")`)
	}

	// Only allow admin or owner
	isAdmin := (jwtUser.Role == userrole.ADMIN)
	isOwner := (jwtUser.Id == user.Id)
	if isAdmin || isOwner {
		return c.Next()
	} else {
		return fiber.NewError(403, "The user does not have permission to access this resource")
	}
}
