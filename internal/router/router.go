package router

import (
	"go-fiber-api-starter/internal/config"
	"go-fiber-api-starter/internal/handlers"
	"go-fiber-api-starter/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// SetupRoutes setup router api
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

	// Health check
	app.Get("/health", handlers.HealthCheck)

	// Initial grouping of API and version route paths
	api := app.Group("/api", logger.New())
	v1 := api.Group("/v1")
	// v2 := api.Group("/v2")

	// "/api/v1" Routes:
	auth := v1.Group("/auth")
	auth.Post("/login", handlers.Login)

	user := v1.Group("/users")
	user.Post("/", handlers.CreateUser)
	user.Get("verify/:token", handlers.VerifyEmail)
	user.Get("/:id", handlers.GetUser)
	user.Patch("/:id", middleware.Authn(), handlers.UpdateUser)
	user.Delete("/:id", middleware.Authn(), handlers.DeleteUser)

	something := v1.Group("/somethings")
	something.Get("/", handlers.GetAllSomethings)
	something.Get("/:id", handlers.GetSomething)
	something.Post("/", middleware.Authn(), handlers.CreateSomething)
	something.Delete("/:id", middleware.Authn(), handlers.DeleteSomething)

	// // "/api/v2" Routes:
	// //   put your v2 routes here
}
