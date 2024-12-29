package handlers

import (
	"go-fiber-api-starter/internal/db"
	"go-fiber-api-starter/internal/enums/userrole"
	"go-fiber-api-starter/internal/enums/userstatus"
	"go-fiber-api-starter/internal/mail"
	"go-fiber-api-starter/internal/models"
	"go-fiber-api-starter/internal/serialization"
	"go-fiber-api-starter/internal/utils"
	"go-fiber-api-starter/internal/validation"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(c *fiber.Ctx) error {
	// Parse
	userSignUp := &models.UserSignUp{}
	if err := c.BodyParser(userSignUp); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Error parsing Sign Up data",
			"data": map[string]any{"errorMessage": err.Error()}})
	}

	// Validate
	if warnings := validation.ValidateUserSignUp(userSignUp); warnings != nil {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "One or more invalid inputs",
			"data": map[string]any{"errorMessage": warnings}})
	}

	// Check if email is already taken
	msg, err := db.CheckEmailAvailability(userSignUp.Email)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": msg,
			"data": map[string]any{"errorMessage": err.Error()}})
	}

	// Check if username is already taken
	msg, err = db.CheckUsernameAvailability(userSignUp.Username)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": msg,
			"data": map[string]any{"errorMessage": err.Error()}})
	}

	// Hash password
	pwdBytes, err := bcrypt.GenerateFromPassword([]byte(userSignUp.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Server error when hashing password",
			"data": map[string]any{"errorMessage": err.Error()}})
	}

	// Create new user
	u := &models.User{
		Email:    userSignUp.Email,
		Username: userSignUp.Username,
		Password: string(pwdBytes),
		OTP:      utils.RandomSixDigitStr(),
		Role:     userrole.REGULAR,
		Status:   userstatus.UNVERIFIED,
	}

	// Save user
	user, err := db.InsertUser(u)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Server error when saving user to database",
			"data": map[string]any{"errorMessage": err.Error()}})
	}

	// Create notification
	n := &models.Notification{
		TextContent: "You have not verified your email address yet",
		HasViewed:   false,
		UserId:      user.Id,
	}

	// Save notification
	_, err = db.InsertNotification(n)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Server error when saving a notification to the database.",
			"data": map[string]any{"errorMessage": err.Error()}})
	}

	// Email the OTP
	err = mail.SendEmailCode(user.Email, u.OTP)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "New user was saved to the database, but there was an error sending the email verification",
			"data": map[string]any{"errorMessage": err.Error()}})
	}

	// Serialize user
	userResponse := serialization.UserResponse(&user)

	// Respond with 201 and user data
	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "Successfully saved new user",
		"data": map[string]any{"user": userResponse}})
}

func GetAllUsers(c *fiber.Ctx) error {
	return nil
}

func GetUser(c *fiber.Ctx) error {
	// Type assert user (the user should be in c.Locals() from AttachUser() middleware)
	user, ok := c.Locals("user").(*models.User)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "c.Locals('user') should be of type '*models.User'",
			"data": map[string]any{"errorMessage": "Incorrect type for c.Locals('user')"}})
	}

	// Get the user that is requesting access
	userRequestingAccess, ok := c.Locals("jwtPayload").(*models.JwtPayload)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Error type asserting jwtPayload",
			"data": map[string]any{"errorMessage": "Incorrect type for c.Locals('jwtPayload')"}})
	}

	// Hide email address if not admin or owner
	isAdmin := (userRequestingAccess.UserRole == userrole.ADMIN)
	isOwner := (userRequestingAccess.UserId == user.Id)
	if !isAdmin && !isOwner {
		user.Email = "********"
	}

	// Serialize user
	userResponse := serialization.UserResponse(user)

	// Send user
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Retrieved user from database",
		"data": map[string]any{"user": userResponse}})
}

func UpdateUser(c *fiber.Ctx) error {
	// Parse
	userUpdate := &models.UserUpdate{}
	if err := c.BodyParser(userUpdate); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Error parsing request body",
			"data": map[string]any{"errorMessage": err.Error()}})
	}

	// Validate
	if warnings := validation.ValidateUserUpdate(userUpdate); warnings != nil {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "One or more invalid inputs",
			"data": map[string]any{"errorMessage": warnings}})
	}

	// Type assert user (the user should be in c.Locals() from AttachUser() middleware)
	user, ok := c.Locals("user").(*models.User)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "c.Locals('user') should be of type '*models.User'",
			"data": map[string]any{"errorMessage": "Incorrect type for c.Locals('user')"}})
	}

	// Check email availability, if new email
	if userUpdate.Email != user.Email {
		msg, err := db.CheckEmailAvailability(userUpdate.Email)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": msg,
				"data": map[string]any{"errorMessage": err.Error()}})
		}
	}

	// Check username availability, if new username
	if userUpdate.Username != user.Username {
		msg, err := db.CheckUsernameAvailability(userUpdate.Username)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": msg,
				"data": map[string]any{"errorMessage": err.Error()}})
		}
	}

	// Hash password
	pwdBytes, err := bcrypt.GenerateFromPassword([]byte(userUpdate.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Server error when hashing password",
			"data": map[string]any{"errorMessage": err.Error()}})
	}

	// Set password
	userUpdate.Password = string(pwdBytes)

	// Save to db
	updatedUser, err := db.UpdateUser(user.Id, userUpdate)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Error updating user in database",
			"data": map[string]any{"errorMessage": err.Error()}})
	}

	// Serialize user
	userResponse := serialization.UserResponse(&updatedUser)

	// Send response
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Saved updated user to database",
		"data": map[string]any{"user": userResponse}})
}

func UpdatePassword(c *fiber.Ctx) error {
	// Shape of data in request body
	type reqBody struct {
		CurrentPassword   string `json:"currentPassword" form:"currentPassword"`
		NewPassword       string `json:"newPassword" form:"newPassword"`
		RepeatNewPassword string `json:"repeatNewPassword" form:"repeatNewPassword"`
	}
	body := &reqBody{}

	// Parse
	if err := c.BodyParser(body); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "Error parsing password data",
			"data": map[string]any{"errorMessage": err.Error()}})
	}

	// Validate inputs
	warnings := make([]string, 0, 3)
	if !validation.IsPasswordValid(body.CurrentPassword) {
		warnings = append(warnings, "Current password is invalid")
	}
	if !validation.IsPasswordValid(body.NewPassword) {
		warnings = append(warnings, "New password is invalid")
	}
	if body.NewPassword != body.RepeatNewPassword {
		warnings = append(warnings, "NewPassword and RepeatNewPassword do not match")
	}
	if len(warnings) > 0 {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "One or more invalid inputs",
			"data": map[string]any{"errorMessage": warnings}})
	}

	// Type assert user (the user should be in c.Locals() from AttachUser() middleware)
	user, ok := c.Locals("user").(*models.User)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "c.Locals('user') should be of type '*models.User'",
			"data": map[string]any{"errorMessage": "Incorrect type for c.Locals('user')"}})
	}

	// Type assert payload (the jwt payload should be in c.Locals() from Authn() middleware)
	payload, ok := c.Locals("jwtPayload").(*models.JwtPayload)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "c.Locals('jwtPayload') should be of type '*models.JwtPayload'",
			"data": map[string]any{"errorMessage": "Incorrect type for c.Locals('jwtPayload')"}})
	}

	// Check password, if not admin (this is an extra security check in addition to the middleware checks)
	if payload.UserRole != userrole.ADMIN {
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.CurrentPassword)); err != nil {
			return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "Only admin or a user with a correct password are allowed to change passwords",
				"data": map[string]any{"errorMessage": "Password input is incorrect"}})
		}
	}

	// Hash new password
	pwdBytes, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Server error when hashing password",
			"data": map[string]any{"errorMessage": err.Error()}})
	}

	// Set password
	hashedPassword := string(pwdBytes)

	// Save to db
	updatedUser, err := db.UpdatePassword(user.Id, hashedPassword)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Error updating password in database",
			"data": map[string]any{"errorMessage": err.Error()}})
	}

	// Serialize user
	userResponse := serialization.UserResponse(&updatedUser)

	// Send response
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Saved password to database",
		"data": map[string]any{"user": userResponse}})
}

func SoftDeleteUser(c *fiber.Ctx) error {
	// Type assert user (the user should be in c.Locals() from AttachUserId() middleware).
	// The user only has the Id field set from AttachUserId().
	user, ok := c.Locals("user").(*models.User)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "c.Locals('user') should be of type '*models.User'",
			"data": map[string]any{"errorMessage": "Incorrect type for c.Locals('user')"}})
	}

	// Soft delete the user
	if err := db.SoftDeleteUser(user.Id); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Error deleting user from database",
			"data": map[string]any{"errorMessage": err.Error()}})
	}

	// Send response
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Deleted user from database",
		"data": map[string]any{"userId": user.Id}})
}
