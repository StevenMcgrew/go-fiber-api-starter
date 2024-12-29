package handlers

import (
	"go-fiber-api-starter/internal/db"
	"go-fiber-api-starter/internal/enums/userstatus"
	"go-fiber-api-starter/internal/mail"
	"go-fiber-api-starter/internal/models"
	"go-fiber-api-starter/internal/serialization"
	"go-fiber-api-starter/internal/utils"
	"go-fiber-api-starter/internal/validation"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func VerifyEmail(c *fiber.Ctx) error {
	// Shape of data in request body
	type reqBody struct {
		Email string `json:"email" form:"email"`
		OTP   string `json:"otp" form:"otp"`
	}
	body := &reqBody{}

	// Parse body
	if err := c.BodyParser(body); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Error parsing email verification data",
			"data": map[string]any{"errorMessage": err.Error()}})
	}

	// Validate inputs
	warnings := make([]string, 0, 2)
	if !validation.IsEmailValid(body.Email) {
		warnings = append(warnings, "Email is invalid")
	}
	if !validation.IsOtpValid(body.OTP) {
		warnings = append(warnings, "Verification code is invalid")
	}
	if len(warnings) > 0 {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "One or more invalid inputs",
			"data": map[string]any{"errorMessage": warnings}})
	}

	// Get user by email
	user, err := db.GetUserByEmail(body.Email)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Server error when getting user from database",
			"data": map[string]any{"errorMessage": err.Error()}})
	}

	// Make sure user's current status is "unverified" before continuing
	if user.Status != userstatus.UNVERIFIED {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "The user's current status is '" + user.Status + "'",
			"data": map[string]any{"errorMessage": "This user has already been verified"}})
	}

	// Check if it's been too long since code was emailed
	expiration := user.CreatedAt.Add(15 * time.Minute)
	if time.Now().After(expiration) {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "The email verification code has expired",
			"data": map[string]any{"errorMessage": "The email verification code has expired"}})
	}

	// Check if otp matches
	if user.OTP != body.OTP {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "The email verification code did not match",
			"data": map[string]any{"errorMessage": "The email verification code did not match"}})
	}

	// Update user
	updatedUser, err := db.UpdateUser(user.Id, &models.UserUpdate{
		Email:    user.Email,
		Username: user.Username,
		Password: user.Password,
		OTP:      "",
		Role:     user.Role,
		Status:   userstatus.ACTIVE,
		ImageUrl: user.ImageUrl,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Server error when updating user in database",
			"data": map[string]any{"errorMessage": err.Error()}})
	}

	// Serialize user
	userResponse := serialization.UserResponse(&updatedUser)

	// Send user in response
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Email has been verified",
		"data": map[string]any{"user": userResponse}})
}

func ResendEmailVerification(c *fiber.Ctx) error {
	// Shape of request body
	type reqBody struct {
		Email string `json:"email" form:"email"`
	}
	body := &reqBody{}

	// Parse body
	if err := c.BodyParser(body); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Error parsing email address from request body",
			"data": map[string]any{"errorMessage": err.Error()}})
	}

	// Validate input
	if !validation.IsEmailValid(body.Email) {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "Email address is invalid",
			"data": map[string]any{"errorMessage": "Email address is invalid"}})
	}

	// Get user by email
	user, err := db.GetUserByEmail(body.Email)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Server error when getting user from database",
			"data": map[string]any{"errorMessage": err.Error()}})
	}

	// Make sure user status is 'unverified' before continuing
	if user.Status != userstatus.UNVERIFIED {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "The user's current status is '" + user.Status + "'",
			"data": map[string]any{"errorMessage": "This user has already been verified"}})
	}

	// Update user with new OTP
	updatedUser, err := db.UpdateUser(user.Id, &models.UserUpdate{
		Email:    user.Email,
		Username: user.Username,
		Password: user.Password,
		OTP:      utils.RandomSixDigitStr(),
		Role:     user.Role,
		Status:   user.Status,
		ImageUrl: user.ImageUrl,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Server error when updating user in database",
			"data": map[string]any{"errorMessage": err.Error()}})
	}

	// Resend email verification
	err = mail.SendEmailCode(body.Email, updatedUser.OTP)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Server error when resending the email verification",
			"data": map[string]any{"errorMessage": err.Error()}})
	}

	// Serialize user
	userResponse := serialization.UserResponse(&updatedUser)

	// Send user in response
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Another email verification has been sent",
		"data": map[string]any{"user": userResponse}})
}

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
