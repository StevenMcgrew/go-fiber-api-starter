package handlers

import (
	"os"

	"go-fiber-api-starter/internal/db"
	"go-fiber-api-starter/internal/enums/jwtclaimkeys"
	"go-fiber-api-starter/internal/enums/userstatus"
	"go-fiber-api-starter/internal/enums/usertype"
	"go-fiber-api-starter/internal/mail"
	"go-fiber-api-starter/internal/models"
	"go-fiber-api-starter/internal/utils"
	"go-fiber-api-starter/internal/validation"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func GetUser(c *fiber.Ctx) error {
	return nil
}

func CreateUser(c *fiber.Ctx) error {

	// Parse
	userSignup := &models.UserSignup{}
	if err := c.BodyParser(userSignup); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Error parsing signup data", "data": err.Error()})
	}

	// Validate
	if warnings := validation.ValidateUserSignup(userSignup); warnings != nil {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "One or more invalid inputs", "data": warnings})
	}

	// Check if email is already taken
	rowsByEmail, err := db.GetUserByEmail(userSignup.Email)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Server error on email lookup", "data": err.Error()})
	}
	if len(rowsByEmail) > 0 {
		return c.Status(409).JSON(fiber.Map{"status": "fail", "message": "Email address is already in use by another user", "data": userSignup.Email})
	}

	// Check if username is already taken
	rowsByUserName, err := db.GetUserByUserName(userSignup.UserName)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Server error on userName lookup", "data": err.Error()})
	}
	if len(rowsByUserName) > 0 {
		return c.Status(409).JSON(fiber.Map{"status": "fail", "message": "UserName is already in use by another user", "data": userSignup.UserName})
	}

	// Hash password
	pwdBytes, err := bcrypt.GenerateFromPassword([]byte(userSignup.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Server error when hashing password", "data": err.Error()})
	}

	// Create new user
	u := &models.User{}
	u.Email = userSignup.Email
	u.UserName = userSignup.UserName
	u.Password = string(pwdBytes)
	u.UserType = usertype.REGULAR
	u.UserStatus = userstatus.UNVERIFIED
	u.ImageUrl = ""

	// Save user to database
	userRows, err := db.InsertUser(u)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Server error when saving user", "data": err.Error()})
	}
	if len(userRows) == 0 {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "User record was not returned after database insert", "data": userRows})
	}
	user := userRows[0]

	// Create JWT verification link for email verification
	tokenString, err := utils.CreateUserJWT(&user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Server error when creating token for email verification link", "data": err.Error()})
	}
	verificationLink := os.Getenv("API_BASE_URL") + "/api/v1/users/verify/" + tokenString

	// Send verification email with link
	err = mail.SendEmailVerification(user.Email, verificationLink)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "New user was saved to the database, but there was an error sending the email verification", "data": err.Error()})
	}

	// Hide password field
	user.Password = "********"

	// Respond with 201 and user data
	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "Successfully saved new user", "data": user})
}

func VerifyEmail(c *fiber.Ctx) error {
	// Get tokenString from path /api/v1/users/verify/:token
	tokenString := c.Params("token")
	if tokenString == "" {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "Token parameter is missing from url path", "data": map[string]interface{}{"token": tokenString}})
	}

	// Parse/Verify the token
	token, err := utils.ParseAndVerifyJWT(tokenString)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "JWT verification error", "data": err.Error()})
	}

	// Get userId out of JWT claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Error extracting claims from JWT", "data": ""})
	}
	id, ok := claims[jwtclaimkeys.USER_ID].(float64)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Error getting user id from JWT claims", "data": ""})
	}
	userId := uint(id)

	// Get user
	userRows, err := db.GetUserById(userId)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Server error when getting user from database", "data": err.Error()})
	}
	if len(userRows) == 0 {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "User record was not found in database", "data": userRows})
	}
	user := userRows[0]

	// Make sure user's current status is "unverified" before continuing
	if user.UserStatus != userstatus.UNVERIFIED {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "This user has already been verified.", "data": ""})
	}

	// Update userStatus to "active"
	user.UserStatus = userstatus.ACTIVE
	userRows, err = db.UpdateUser(&user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Server error when updating user", "data": err.Error()})
	}
	if len(userRows) == 0 {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "User record was not returned after database update. However, it is likely that the user was still updated.", "data": userRows})
	}

	// Redirect to success webpage
	return c.Redirect("/email-verification-success.html")
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
