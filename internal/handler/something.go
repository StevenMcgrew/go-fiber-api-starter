package handler

import (
	"github.com/gofiber/fiber/v2"
)

func GetAllSomethings(c *fiber.Ctx) error {
	// db := database.Pool
	// var somethings []model.Something
	// db.Find(&somethings)
	// return c.JSON(fiber.Map{"status": "success", "message": "All somethings", "data": somethings})
	return nil
}

func GetSomething(c *fiber.Ctx) error {
	// id := c.Params("id")
	// db := database.Conn
	// var something model.Something
	// db.Find(&something, id)
	// if something.Title == "" {
	// 	return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No something found with ID", "data": nil})
	// }
	// return c.JSON(fiber.Map{"status": "success", "message": "Something found", "data": something})
	return nil
}

func CreateSomething(c *fiber.Ctx) error {
	// db := database.Conn
	// something := new(model.Something)
	// if err := c.BodyParser(something); err != nil {
	// 	return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't create something", "data": err})
	// }
	// db.Create(&something)
	// return c.JSON(fiber.Map{"status": "success", "message": "Created something", "data": something})
	return nil
}

func DeleteSomething(c *fiber.Ctx) error {
	// id := c.Params("id")
	// db := database.Conn

	// var something model.Something
	// db.First(&something, id)
	// if something.Title == "" {
	// 	return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No something found with ID", "data": nil})
	// }
	// db.Delete(&something)
	// return c.JSON(fiber.Map{"status": "success", "message": "Something successfully deleted", "data": nil})
	return nil
}
