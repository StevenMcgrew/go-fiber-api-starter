package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// Authenticate
func Authn() fiber.Handler {
	a := func(c *fiber.Ctx) error {
		return nil
	}
	return a
}

// Authorize
func Authz(c *fiber.Ctx) {
	return
}
