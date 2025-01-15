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

	// Cache (I'm only adding one header here)
	// see middleware Cache for an alternative at https://docs.gofiber.io/api/middleware/cache
	app.Use(func(c *fiber.Ctx) error {
		c.Set("Cache-Control", "no-cache")
		return c.Next()
	})

	// Static
	app.Static("/", "./public", config.FiberStaticConfig)

	// PAGES
	app.Get("/", hn.HomePage)
	app.Get("/signup", hn.SignUpPage)
	app.Get("/login", hn.LogInPage)

	// Health check
	app.Get("/health", hn.HealthCheck)

	// Initial grouping of path
	v1 := app.Group("/api/v1")

	// AUTH
	v1.Post("/auth/login", hn.Login)
	v1.Get("/auth/verify-email", hn.VerifyEmail)
	v1.Patch("/auth/resend-email-verification", hn.ResendEmailVerification)
	v1.Post("/auth/forgot-password", hn.ForgotPassword)
	v1.Get("/auth/reset-password/request", hn.ResetPasswordPage)
	v1.Post("/auth/reset-password/update", hn.ResetForgottenPassword)

	// USERS
	v1.Post("/users", hn.CreateUser)
	v1.Get("/users/:userId", mw.AttachUser, hn.GetUser)
	v1.Get("/users", mw.Authn, hn.GetAllUsers)
	v1.Post("/users/email/availability", hn.IsEmailAvailable)
	v1.Post("/users/username/availability", hn.IsUsernameAvailable)
	v1.Patch("/users/:userId", mw.Authn, mw.AttachUser, mw.OnlyAdmin, hn.UpdateUser)
	v1.Patch("/users/:userId/password", mw.Authn, mw.AttachUser, mw.OnlyAdminOrOwner, hn.UpdatePassword)
	v1.Patch("/users/:userId/email", mw.Authn, mw.AttachUser, mw.OnlyAdminOrOwner, hn.UpdateEmail)
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
