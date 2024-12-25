package handlers

import (
	"fmt"
	"go-fiber-api-starter/internal/db"
	"go-fiber-api-starter/internal/enums/userstatus"
	"go-fiber-api-starter/internal/serialization"
	"go-fiber-api-starter/internal/utils"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// Login get user and password
func Login(c *fiber.Ctx) error {
	// Shape of request body
	type reqBody struct {
		Email    string `json:"email" form:"email"`
		Password string `json:"password" form:"password"`
	}
	body := &reqBody{}

	// Parse body
	if err := c.BodyParser(body); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Error parsing login credentials", "data": err.Error()})
	}

	// Validate inputs
	isValid := true
	if !utils.IsEmailValid(body.Email) {
		isValid = false
	}
	if !utils.IsPasswordValid(body.Password) {
		isValid = false
	}
	if !isValid {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "Login credentials are invalid", "data": ""})
	}

	// Get user by email
	user, err := db.GetUserByEmail(body.Email)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Server error when getting user from database", "data": err.Error()})
	}

	// Make sure user status is 'active'
	if user.Status != userstatus.ACTIVE {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": fmt.Sprintf("Current user status is '%s'. No login allowed.", user.Status), "data": ""})
	}

	// Check password
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "Login credentials are incorrect", "data": ""})
	}

	// Create JWT
	jwt, err := utils.CreateUserJWT(&user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Server error when creating a JWT", "data": err.Error()})
	}

	// Serialize user
	userResponse := serialization.UserResponse(&user)

	// Create response
	response := fiber.Map{
		"status":  "success",
		"message": "User has been logged in",
		"data": map[string]any{
			"user":  userResponse,
			"token": jwt,
		},
	}

	// Send response
	return c.Status(200).JSON(response)
}
