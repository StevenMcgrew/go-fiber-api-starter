package handler

import (
	"go-fiber-api-starter/internal/database"
	"go-fiber-api-starter/internal/model"

	"github.com/gofiber/fiber/v2"
)

func GetAllSomethings(c *fiber.Ctx) error {
	db := database.DB
	var somethings []model.Something
	db.Find(&somethings)
	return c.JSON(fiber.Map{"status": "success", "message": "All somethings", "data": somethings})
}

func GetSomething(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB
	var something model.Something
	db.Find(&something, id)
	if something.Title == "" {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No something found with ID", "data": nil})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Something found", "data": something})
}

func CreateSomething(c *fiber.Ctx) error {
	db := database.DB
	something := new(model.Something)
	if err := c.BodyParser(something); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't create something", "data": err})
	}
	db.Create(&something)
	return c.JSON(fiber.Map{"status": "success", "message": "Created something", "data": something})
}

func DeleteSomething(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB

	var something model.Something
	db.First(&something, id)
	if something.Title == "" {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No something found with ID", "data": nil})
	}
	db.Delete(&something)
	return c.JSON(fiber.Map{"status": "success", "message": "Something successfully deleted", "data": nil})
}
