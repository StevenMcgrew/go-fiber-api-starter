package handlers

import (
	"go-fiber-api-starter/internal/db"
	"go-fiber-api-starter/internal/enums/userstatus"
	"go-fiber-api-starter/internal/serialization"
	"go-fiber-api-starter/internal/utils"
	"go-fiber-api-starter/internal/validation"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *fiber.Ctx) error {
	// Shape of request body
	type reqBody struct {
		Email    string `json:"email" form:"email"`
		Password string `json:"password" form:"password"`
	}
	body := &reqBody{}

	// Parse body
	if err := c.BodyParser(body); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Error parsing login credentials",
			"data": map[string]any{"errorMessage": err.Error()}})
	}

	// Validate inputs
	isValid := true
	if !validation.IsEmailValid(body.Email) {
		isValid = false
	}
	if !validation.IsPasswordValid(body.Password) {
		isValid = false
	}
	if !isValid {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "Login credentials are invalid",
			"data": map[string]any{"errorMessage": "Login credentials are invalid"}})
	}

	// Get user by email
	user, err := db.GetUserByEmail(body.Email)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Server error when getting user from database",
			"data": map[string]any{"errorMessage": err.Error()}})
	}

	// Make sure user status is 'active'
	if user.Status != userstatus.ACTIVE {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "Cannot login users who have an account status of " + user.Status,
			"data": map[string]any{"errorMessage": "Cannot login users who have an account status of " + user.Status}})
	}

	// Check password
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "Login credentials are incorrect",
			"data": map[string]any{"errorMessage": "Login credentials are incorrect"}})
	}

	// Create JWT
	jwt, err := utils.CreateJWT(&user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Server error when creating a JWT",
			"data": map[string]any{"errorMessage": err.Error()}})
	}

	// Serialize user
	userResponse := serialization.UserResponse(&user)

	// Send user and jwt
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "User has been logged in",
		"data": map[string]any{"user": userResponse, "token": jwt}})
}
