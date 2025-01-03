package main

import (
	"fmt"
	"log"

	"go-fiber-api-starter/internal/config"
	"go-fiber-api-starter/internal/db"
	"go-fiber-api-starter/internal/router"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New(config.FiberServerConfig)

	db.Connect(config.DB_URL)
	db.ExecuteSqlFile("./internal/db/create-db.sql")
	fmt.Println("Database is ready.")
	router.SetupRoutes(app)

	log.Fatal(app.Listen(":" + config.API_PORT))
}
