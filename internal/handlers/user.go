package handlers

import (
	"go-fiber-api-starter/internal/db"
	"go-fiber-api-starter/internal/enums/userstatus"
	"go-fiber-api-starter/internal/enums/usertype"
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

func GetUser(c *fiber.Ctx) error {
	return nil
}

func CreateUser(c *fiber.Ctx) error {

	// Parse
	userSignUp := &models.UserSignUp{}
	if err := c.BodyParser(userSignUp); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Error parsing Sign Up data", "data": err.Error()})
	}

	// Validate
	if warnings := validation.ValidateUserSignUp(userSignUp); warnings != nil {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "One or more invalid inputs", "data": warnings})
	}

	// Check if email is already taken
	userByEmail, err := db.GetUserByEmail(userSignUp.Email)
	if err != nil {
		if err != pgx.ErrNoRows { // some error other than no rows
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Server error on email lookup", "data": err.Error()})
		}
	} else { // user with this email address was found
		return c.Status(409).JSON(fiber.Map{"status": "fail", "message": "Email address is already in use by another user", "data": userByEmail.Email})
	}

	// Check if username is already taken
	userByUsername, err := db.GetUserByUserName(userSignUp.Username)
	if err != nil {
		if err != pgx.ErrNoRows { // some error other than no rows
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Server error on username lookup", "data": err.Error()})
		}
	} else { // user with this username was found
		return c.Status(409).JSON(fiber.Map{"status": "fail", "message": "Username is already in use by another user", "data": userByUsername.Username})
	}

	// Hash password
	pwdBytes, err := bcrypt.GenerateFromPassword([]byte(userSignUp.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Server error when hashing password", "data": err.Error()})
	}

	// Create new user
	u := &models.User{
		Email:    userSignUp.Email,
		Username: userSignUp.Username,
		Password: string(pwdBytes),
		OTP:      utils.RandomSixDigitStr(),
		Role:     usertype.REGULAR,
		Status:   userstatus.UNVERIFIED,
	}

	// Save user
	user, err := db.InsertUser(u)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Server error when saving user to database", "data": err.Error()})
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
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Server error when saving a notification to the database.", "data": err.Error()})
	}

	// Email the OTP
	err = mail.SendEmailCode(user.Email, u.OTP)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "New user was saved to the database, but there was an error sending the email verification", "data": err.Error()})
	}

	// Serialize user
	userResponse := serialization.UserResponse(&user)

	// Respond with 201 and user data
	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "Successfully saved new user", "data": userResponse})
}

func VerifyEmail(c *fiber.Ctx) error {
	// Shape of data in request body
	type reqBody struct {
		Email string `json:"email" form:"email"`
		OTP   string `json:"otp" form:"otp"`
	}
	body := &reqBody{}

	// Get email address and otp from body
	if err := c.BodyParser(body); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Error parsing email verification data", "data": err.Error()})
	}

	// Get user by email
	user, err := db.GetUserByEmail(body.Email)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Server error when getting user from database", "data": err.Error()})
	}

	// Make sure user's current status is "unverified" before continuing
	if user.Status != userstatus.UNVERIFIED {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "This user has already been verified.", "data": serialization.UserResponse(&user)})
	}

	// Check if it's been too long since code was emailed
	expiration := user.CreatedAt.Add(15 * time.Minute)
	if time.Now().After(expiration) {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "The email verification code has expired", "data": ""})
	}

	// Verify otp
	if user.OTP != body.OTP {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "The email verification code did not match", "data": ""})
	}

	// Set user status and clear otp
	user.Status = userstatus.ACTIVE
	user.OTP = ""

	// Save user
	updatedUser, err := db.UpdateUser(&user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Server error when updating user in database", "data": err.Error()})
	}

	// Serialize user
	userResponse := serialization.UserResponse(&updatedUser)

	// Send user in response
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Email has been verified", "data": userResponse})
}

// UpdateUser update user
func UpdateUser(c *fiber.Ctx) error {
	// type UpdateUserInput struct {
	// 	Names string `json:"names"`
	// }
	// var uui UpdateUserInput
	// if err := c.BodyParser(&uui); err != nil {
	// 	return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	// }
	// id := c.Params("id")
	// token := c.Locals("user").(*jwt.Token)

	// if !validToken(token, id) {
	// 	return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Invalid token id", "data": nil})
	// }

	// db := database.Conn
	// var user model.User

	// db.First(&user, id)
	// // user.Names = uui.Names
	// db.Save(&user)

	// return c.JSON(fiber.Map{"status": "success", "message": "User successfully updated", "data": user})
	return nil
}

// DeleteUser delete user
func DeleteUser(c *fiber.Ctx) error {
	// type PasswordInput struct {
	// 	Password string `json:"password"`
	// }
	// var pi PasswordInput
	// if err := c.BodyParser(&pi); err != nil {
	// 	return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	// }
	// id := c.Params("id")
	// token := c.Locals("user").(*jwt.Token)

	// if !validToken(token, id) {
	// 	return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Invalid token id", "data": nil})
	// }

	// if !validUser(id, pi.Password) {
	// 	return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Not valid user", "data": nil})
	// }

	// db := database.Conn
	// var user model.User

	// db.First(&user, id)

	// db.Delete(&user)
	// return c.JSON(fiber.Map{"status": "success", "message": "User successfully deleted", "data": nil})
	return nil
}
