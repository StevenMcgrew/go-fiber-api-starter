package router

import (
	"go-fiber-api-starter/internal/config"
	"go-fiber-api-starter/internal/handlers"
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
	app.Get("/", handlers.HomePage)
	app.Get("/email-verification-success", handlers.EmailVerificationSuccessPage)
	app.Get("/email-verification-failure", handlers.EmailVerificationFailurePage)

	// Health check
	app.Get("/health", handlers.HealthCheck)

	// Initial grouping of route paths
	api := app.Group("/api", logger.New())
	v1 := api.Group("/v1")

	auth := v1.Group("/auth") // /api/v1/auth
	auth.Post("/login", handlers.Login)

	user := v1.Group("/users") // /api/v1/users
	user.Post("/", handlers.CreateUser)
	user.Put("/verify-email", handlers.VerifyEmail)
	user.Put("/resend-email-verification", handlers.ResendEmailVerification)
	// user.Get("/", mw.Authn, handlers.GetAllUsers)
	user.Get("/:id", mw.Authn, handlers.GetUser)
	user.Patch("/:id", mw.Authn, handlers.UpdateUser)
	user.Delete("/:id", mw.Authn, handlers.DeleteUser)

	something := v1.Group("/somethings")
	something.Get("/", mw.Authn, mw.AdminOnly, handlers.GetAllSomethings)
	something.Get("/:id", mw.Authn, handlers.GetSomething)
	something.Post("/", mw.Authn, handlers.CreateSomething)
	something.Delete("/:id", mw.Authn, handlers.DeleteSomething)
}
