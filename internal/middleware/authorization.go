package middleware

import (
	"go-fiber-api-starter/internal/enums/userrole"
	"go-fiber-api-starter/internal/models"

	"github.com/gofiber/fiber/v2"
)

func OnlyAdmin(c *fiber.Ctx) error {
	// Type assert payload
	payload, ok := c.Locals("jwtPayload").(*models.JwtUser)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "c.Locals('jwtPayload') should be of type '*models.JwtPayload'",
			"data": map[string]any{"errorMessage": "Incorrect type for c.Locals('jwtPayload')"}})
	}

	// Only allow admin
	if payload.UserRole == userrole.ADMIN {
		return c.Next()
	} else {
		return c.Status(403).JSON(fiber.Map{"status": "fail", "message": "The user does not have permission to access this path",
			"data": map[string]any{"errorMessage": "Access denied"}})
	}
}

func OnlyAdminOrOwner(c *fiber.Ctx) error {
	// Get the user that is requesting access
	userRequestingAccess, ok := c.Locals("jwtPayload").(*models.JwtUser)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "c.Locals('jwtPayload') should be of type '*models.JwtPayload'",
			"data": map[string]any{"errorMessage": "Incorrect type for c.Locals('jwtPayload')"}})
	}

	// Get the user to be accessed
	userToBeAccessed, ok := c.Locals("user").(*models.User)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "c.Locals('user') should be of type '*models.User'",
			"data": map[string]any{"errorMessage": "Incorrect type for c.Locals('user')"}})
	}

	isAdmin := (userRequestingAccess.UserRole == userrole.ADMIN)
	isOwner := (userRequestingAccess.UserId == userToBeAccessed.Id)

	// Only allow admin or owner
	if isAdmin || isOwner {
		return c.Next()
	} else {
		return c.Status(403).JSON(fiber.Map{"status": "fail", "message": "The user does not have permission to access this path",
			"data": map[string]any{"errorMessage": "Access denied"}})
	}

}
