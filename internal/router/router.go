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

	// Web pages
	app.Get("/", hn.HomePage)
	app.Get("/email-verification-success", hn.EmailVerificationSuccessPage)
	app.Get("/email-verification-failure", hn.EmailVerificationFailurePage)

	// Health check
	app.Get("/health", hn.HealthCheck)

	// Initial grouping of path
	v1 := app.Group("/api/v1")

	// AUTH
	v1.Post("/auth/login", hn.Login)

	// USERS
	v1.Post("/users/", hn.CreateUser)
	v1.Patch("/users/verify-email", hn.VerifyEmail)
	v1.Patch("/users/resend-email-verification", hn.ResendEmailVerification)
	v1.Get("/users/:userId", mw.Authn, mw.AttachUser, hn.GetUser)
	v1.Get("/users/", mw.Authn, hn.GetAllUsers)
	v1.Patch("/users/:userId", mw.Authn, mw.AttachUser, mw.OnlyAdminOrOwner, hn.UpdateUser)
	v1.Delete("/users/:userId", mw.Authn, mw.AttachUserId, mw.OnlyAdminOrOwner, hn.DeleteUser)

	// NOTIFICATIONS
	v1.Get("/users/:userId/notifications", mw.Authn, mw.AttachUser, mw.OnlyAdminOrOwner, hn.GetAllNotificationsForUser)
	v1.Get("/notifications/", mw.Authn, mw.OnlyAdmin, hn.GetAllNotifications)
	v1.Get("/notifications/:noteId", mw.Authn, mw.OnlyAdmin, hn.GetNotification)
	v1.Post("/notifications/", mw.Authn, mw.OnlyAdmin, hn.CreateNotification)
	v1.Delete("/users/:userId/notifications/:noteId", mw.Authn, mw.AttachUser, mw.OnlyAdminOrOwner, hn.DeleteNotification)
}
