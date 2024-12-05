package router

import (
	"go-fiber-api-starter/internal/handler"
	"go-fiber-api-starter/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// SetupRoutes setup router api
func SetupRoutes(app *fiber.App) {

	// Logger
	// https://docs.gofiber.io/api/middleware/logger
	app.Use(logger.New())

	// CORS
	// https://docs.gofiber.io/api/middleware/cors
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Accept,Authorization,Content-Type",
		AllowCredentials: false, // credentials require explicit origins
		MaxAge:           300,
	}))

	// Cache (I'm only adding one header here)
	// see middleware Cache for an alternative at https://docs.gofiber.io/api/middleware/cache
	app.Use(func(c *fiber.Ctx) error {
		c.Set("Cache-Control", "no-cache")
		return c.Next()
	})

	// Static
	// https://docs.gofiber.io/api/app
	app.Static("/", "./public", fiber.Static{
		Index:         "index.html", // sets the file name to serve from "/"
		CacheDuration: -1,           // negative value disables cache, default is 10 * time.Second
	})

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

	product := v1.Group("/product")
	product.Get("/", handler.GetAllProducts)
	product.Get("/:id", handler.GetProduct)
	product.Post("/", middleware.Protected(), handler.CreateProduct)
	product.Delete("/:id", middleware.Protected(), handler.DeleteProduct)

	// "/api/v2" Routes:
	//   put your v2 routes here
}
