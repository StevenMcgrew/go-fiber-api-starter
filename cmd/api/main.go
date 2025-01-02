package main

import (
	"fmt"
	"log"
	"os"

	"go-fiber-api-starter/internal/config"
	"go-fiber-api-starter/internal/db"
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

	app := fiber.New(config.FiberServerConfig)

	db.Connect(os.Getenv("DB_URL"))
	db.ExecuteSqlFile("./internal/db/create-db.sql")
	fmt.Println("Database is ready.")
	router.SetupRoutes(app)

	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
}
