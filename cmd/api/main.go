package main

import (
	"log"

	"go-fiber-api-starter/internal/database"
	"go-fiber-api-starter/internal/router"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {

	// Load env variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := fiber.New()

	database.ConnectDB()
	router.SetupRoutes(app)

	log.Fatal(app.Listen(":3000"))
}
