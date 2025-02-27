package router

import (
	"go-fiber-api-starter/internal/config"
	hn "go-fiber-api-starter/internal/handlers"
	mw "go-fiber-api-starter/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRoutes(app *fiber.App) {

	// Logger
	app.Use(logger.New(config.FiberLoggerConfig))

	// CORS
	app.Use(cors.New(config.FiberCorsConfig))

	// Set some response headers
	app.Use(func(c *fiber.Ctx) error {
		// c.Set("Access-Control-Max-Age", "7200")
		return c.Next()
	})

	// Static
	app.Static("/", "./public", config.FiberStaticConfig)

	// Health check
	app.Get("/health", hn.HealthCheck)

	// Initial grouping of path
	v1 := app.Group("/api/v1")

	// AUTH
	v1.Post("/auth/login", hn.Login)
	v1.Post("/auth/verify-email", hn.VerifyEmail)
	v1.Post("/auth/resend-email-verification", hn.ResendEmailVerification)
	v1.Post("/auth/reset-password/request", hn.ResetPasswordRequest)
	v1.Patch("/auth/reset-password/update", hn.ResetPasswordUpdate)

	// USERS
	v1.Post("/users", hn.CreateUser)
	v1.Get("/users/:userId", mw.AttachUser, hn.GetUser)
	v1.Get("/users", mw.Authn, hn.GetAllUsers)
	v1.Post("/users/email/availability", hn.IsEmailAvailable)
	v1.Post("/users/username/availability", hn.IsUsernameAvailable)
	v1.Patch("/users/:userId", mw.Authn, mw.AttachUser, mw.OnlyAdmin, hn.UpdateUser)
	v1.Patch("/users/:userId/password", mw.Authn, mw.AttachUser, mw.OnlyAdminOrOwner, hn.UpdatePassword)
	v1.Post("/users/:userId/change-email/request", mw.Authn, mw.AttachUser, mw.OnlyAdminOrOwner, hn.ChangeEmailRequest)
	v1.Patch("/users/:userId/change-email/update", mw.Authn, mw.AttachUser, mw.OnlyAdminOrOwner, hn.ChangeEmailUpdate)
	v1.Patch("/users/:userId/username", mw.Authn, mw.AttachUser, mw.OnlyAdminOrOwner, hn.UpdateUsername)
	v1.Delete("/users/:userId", mw.Authn, mw.AttachUser, mw.OnlyAdminOrOwner, hn.SoftDeleteUser)

	// NOTIFICATIONS
	v1.Get("/users/:userId/notifications", mw.Authn, mw.AttachUser, mw.OnlyAdminOrOwner, hn.GetAllNotificationsForUser)
	v1.Get("/notifications", mw.Authn, mw.OnlyAdmin, hn.GetAllNotifications)
	v1.Get("/notifications/:noteId", mw.Authn, mw.OnlyAdmin, hn.GetNotification)
	v1.Post("/notifications", mw.Authn, mw.OnlyAdmin, hn.CreateNotification)
	v1.Delete("/users/:userId/notifications/:noteId", mw.Authn, mw.AttachUser, mw.OnlyAdminOrOwner, hn.DeleteNotification)

	// 404 Not Found
	app.Use(func(c *fiber.Ctx) error {
		return fiber.NewError(404, "Not Found")
	})
}
