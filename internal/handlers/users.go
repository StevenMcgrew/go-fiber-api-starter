package handlers

import (
	"fmt"
	"go-fiber-api-starter/internal/db"
	"go-fiber-api-starter/internal/enums/userrole"
	"go-fiber-api-starter/internal/enums/userstatus"
	"go-fiber-api-starter/internal/mail"
	"go-fiber-api-starter/internal/models"
	"go-fiber-api-starter/internal/serialization"
	"go-fiber-api-starter/internal/utils"
	"go-fiber-api-starter/internal/validation"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(c *fiber.Ctx) error {
	// Parse
	userSignUp := &models.UserSignUp{}
	if err := c.BodyParser(userSignUp); err != nil {
		return fiber.NewError(400, "Error parsing request body: "+err.Error())
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
		TextContent: "Welcome! Thanks for signing up.",
		HasViewed:   false,
		UserId:      user.Id,
	}

	// Save notification
	_, err = db.InsertNotification(n)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Server error when saving a notification to the database.",
			"data": map[string]any{"errorMessage": err.Error()}})
	}

	// Create JWT for email verification link
	claims := &models.JwtVerifyEmail{
		UserId: user.Id,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Error when creating JWT for email verification link",
			"data": map[string]any{"errorMessage": err.Error()}})
	}

	// Create email verification link
	link := fmt.Sprintf("%s/api/v1/auth/verify-email/?token=%s", os.Getenv("API_BASE_URL"), jwtString)

	// Send verification email
	err = mail.SendEmailVerification(user.Email, link)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "New user was saved to the database, but there was an error sending the email verification",
			"data": map[string]any{"errorMessage": err.Error()}})
	}

	// Serialize user
	userResponse := serialization.UserResponse(&user)

	// Create JSON response and send
	return utils.SendSuccessJSON(c, 201, userResponse, "Saved new user")
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

	// Type assert jwtUser
	jwtUser, ok := c.Locals("jwtUser").(*models.User)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Error type asserting jwtUser",
			"data": map[string]any{"errorMessage": "Incorrect type for c.Locals('jwtUser')"}})
	}

	// Hide email address if not admin or owner
	isAdmin := (jwtUser.Role == userrole.ADMIN)
	isOwner := (jwtUser.Id == user.Id)
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

	// Type assert jwtUser (the jwt jwtUser should be in c.Locals() from Authn() middleware)
	jwtUser, ok := c.Locals("jwtUser").(*models.User)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Wrong type in c.Locals",
			"data": map[string]any{"errorMessage": "Wrong type in c.Locals"}})
	}

	// Check password, if not admin (this is an extra security check)
	if jwtUser.Role != userrole.ADMIN {
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

func UpdateEmail(c *fiber.Ctx) error {
	// Shape of request body
	type reqBody struct {
		Email string `json:"email" form:"email"`
	}
	body := &reqBody{}

	// Parse
	if err := c.BodyParser(body); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "Error parsing email data",
			"data": map[string]any{"errorMessage": err.Error()}})
	}

	// Validate
	if !validation.IsEmailValid(body.Email) {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "Username input is invalid",
			"data": map[string]any{"errorMessage": "Invalid input"}})
	}

	// Check if email is already taken
	msg, err := db.CheckEmailAvailability(body.Email)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": msg,
			"data": map[string]any{"errorMessage": err.Error()}})
	}

	// Parse userId from path
	id, err := c.ParamsInt("userId")
	if err != nil || id == 0 {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "Error parsing userId from URL",
			"data": map[string]any{"errorMessage": "Error parsing userId from URL"}})
	}
	userId := uint(id)

	// Create JWT for email verification link
	claims := &models.JwtVerifyEmail{
		UserId: userId,
		Email:  body.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Error when creating JWT for email verification link",
			"data": map[string]any{"errorMessage": err.Error()}})
	}

	// Create email verification link
	link := fmt.Sprintf("%s/api/v1/auth/verify-email/?token=%s", os.Getenv("API_BASE_URL"), jwtString)

	// Send verification email
	err = mail.SendEmailVerification(body.Email, link)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "There was an error sending the email verification",
			"data": map[string]any{"errorMessage": err.Error()}})
	}

	// Send response
	return c.Status(200).JSON(fiber.Map{"status": "success",
		"message": fmt.Sprintf("A verification link has been emailed to %s", body.Email), "data": ""})
}

func UpdateUsername(c *fiber.Ctx) error {
	// Shape of data in request body
	type reqBody struct {
		Username string `json:"username" form:"username"`
	}
	body := &reqBody{}

	// Parse
	if err := c.BodyParser(body); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "Error parsing username data",
			"data": map[string]any{"errorMessage": err.Error()}})
	}

	// Validate
	if !validation.IsUsernameValid(body.Username) {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "Username input is invalid",
			"data": map[string]any{"errorMessage": "Invalid input"}})
	}

	// Check if username is already taken
	msg, err := db.CheckUsernameAvailability(body.Username)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": msg,
			"data": map[string]any{"errorMessage": err.Error()}})
	}

	// Parse userId from path
	id, err := c.ParamsInt("userId")
	if err != nil || id == 0 {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "Error parsing userId from URL",
			"data": map[string]any{"errorMessage": "Error parsing userId from URL"}})
	}
	userId := uint(id)

	// Save to db
	updatedUser, err := db.UpdateUsername(userId, body.Username)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Error updating username in database",
			"data": map[string]any{"errorMessage": err.Error()}})
	}

	// Serialize user
	userResponse := serialization.UserResponse(&updatedUser)

	// Send response
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Saved username to database",
		"data": map[string]any{"user": userResponse}})
}

func SoftDeleteUser(c *fiber.Ctx) error {
	// Parse userId from path
	id, err := c.ParamsInt("userId")
	if err != nil || id == 0 {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "Error parsing userId from URL",
			"data": map[string]any{"errorMessage": "Error parsing userId from URL"}})
	}
	userId := uint(id)

	// Soft delete the user
	if err := db.SoftDeleteUser(userId); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Error deleting user from database",
			"data": map[string]any{"errorMessage": err.Error()}})
	}

	// Send response
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Deleted user from database",
		"data": map[string]any{"userId": userId}})
}
