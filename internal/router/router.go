package router

import (
	"go-fiber-api-starter/internal/config"
	"go-fiber-api-starter/internal/handler"
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
	app.Get("/health", handler.HealthCheck)

	// Initial grouping of API and version route paths
	api := app.Group("/api", logger.New())
	v1 := api.Group("/v1")
	// v2 := api.Group("/v2")

	// "/api/v1" Routes:
	auth := v1.Group("/auth")
	auth.Post("/login", handler.Login)

	user := v1.Group("/user")
	user.Get("/:id", handler.GetUser)
	user.Post("/", handler.CreateUser)
	user.Patch("/:id", middleware.Protected(), handler.UpdateUser)
	user.Delete("/:id", middleware.Protected(), handler.DeleteUser)

	something := v1.Group("/something")
	something.Get("/", handler.GetAllSomethings)
	something.Get("/:id", handler.GetSomething)
	something.Post("/", middleware.Protected(), handler.CreateSomething)
	something.Delete("/:id", middleware.Protected(), handler.DeleteSomething)

	// // "/api/v2" Routes:
	// //   put your v2 routes here
}
