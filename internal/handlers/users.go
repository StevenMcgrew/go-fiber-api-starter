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
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
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
	userByEmail, err := db.GetUserByEmail(userSignUp.Email)
	if err != nil {
		if err != pgx.ErrNoRows { // some error other than no rows
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Server error on email lookup",
				"data": map[string]any{"errorMessage": err.Error()}})
		}
	} else { // user with this email address was found
		return c.Status(409).JSON(fiber.Map{"status": "fail", "message": "Email address is already in use by another user",
			"data": map[string]any{"errorMessage": userByEmail.Email + "is already in use"}})
	}

	// Check if username is already taken
	userByUsername, err := db.GetUserByUserName(userSignUp.Username)
	if err != nil {
		if err != pgx.ErrNoRows { // some error other than no rows
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Server error on username lookup",
				"data": map[string]any{"errorMessage": err.Error()}})
		}
	} else { // user with this username was found
		return c.Status(409).JSON(fiber.Map{"status": "fail", "message": "Username is already in use by another user",
			"data": map[string]any{"errorMessage": userByUsername.Username + "is already in use"}})
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

	// Verify otp
	if user.OTP != body.OTP {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "The email verification code did not match",
			"data": map[string]any{"errorMessage": "The email verification code did not match"}})
	}

	// Set user status and clear out otp
	user.Status = userstatus.ACTIVE
	user.OTP = ""

	// Save user
	updatedUser, err := db.UpdateUser(&user)
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

	// Set new OTP for user
	user.OTP = utils.RandomSixDigitStr()

	// Save user with new code
	updatedUser, err := db.UpdateUser(&user)
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

	// Set updated fields
	user.Email = userUpdate.Email
	user.Username = userUpdate.Username
	user.Password = userUpdate.Password
	user.Role = userUpdate.Role
	user.Status = userUpdate.Status
	user.ImageUrl = userUpdate.ImageUrl

	// Save to db
	updatedUser, err := db.UpdateUser(user)
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

func DeleteUser(c *fiber.Ctx) error {
	// Type assert user (the user should be in c.Locals() from AttachUserId() middleware).
	// The user only has the Id field set on this route, the other fields are empty or nil.
	user, ok := c.Locals("user").(*models.User)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "c.Locals('user') should be of type '*models.User'",
			"data": map[string]any{"errorMessage": "Incorrect type for c.Locals('user')"}})
	}

	// Delete the user
	if err := db.DeleteUserById(user.Id); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Error deleting user from database",
			"data": map[string]any{"errorMessage": err.Error()}})
	}

	// Send response
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Deleted user from database",
		"data": map[string]any{"userId": user.Id}})
}
