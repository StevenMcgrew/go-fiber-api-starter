package handlers

import (
	"strconv"

	"go-fiber-api-starter/internal/models"
	"go-fiber-api-starter/internal/validation"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func validToken(t *jwt.Token, id string) bool {
	n, err := strconv.Atoi(id)
	if err != nil {
		return false
	}

	claims := t.Claims.(jwt.MapClaims)
	uid := int(claims["user_id"].(float64))

	return uid == n
}

func validUser(id string, p string) bool {
	// db := database.Conn
	// var user model.User
	// db.First(&user, id)
	// if user.Username == "" {
	// 	return false
	// }
	// if !CheckPasswordHash(p, user.Password) {
	// 	return false
	// }
	return true
}

func GetUser(c *fiber.Ctx) error {
	// id := c.Params("id")
	// db := database.Conn
	// var user model.User
	// db.Find(&user, id)
	// if user.Username == "" {
	// 	return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No user found with ID", "data": nil})
	// }
	// return c.JSON(fiber.Map{"status": "success", "message": "User found", "data": user})
	return nil
}

func CreateUser(c *fiber.Ctx) error {

	// Parse
	userSignup := &models.UserSignup{}
	if err := c.BodyParser(userSignup); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Error parsing signup data", "data": err})
	}

	// Validate
	if warnings := validation.ValidateUserSignup(userSignup); warnings != nil {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "One or more invalid inputs", "data": warnings})
	}

	// Check if email is already taken

	// Check if username is already taken

	// Hash password

	// Set missing data
	// user.UserType = usertype.REGULAR
	// user.UserStatus = userstatus.UNVERIFIED
	// user.ImageUrl = ""

	// Save user to database

	// Create JWT for email verification

	// Create URL link for email verification

	// Send verification email with link

	// Hide or remove some fields
	// user.Password = "encrypted"

	// Respond with 201 and user data

	// hash, err := hashPassword(user.Password)
	// if err != nil {
	// 	return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't hash password", "data": err})
	// }

	// user.Password = hash
	// if err := database.DB.Create(&user).Error; err != nil {
	// 	return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't create user", "data": err})
	// }

	// newUser := serialization.SerializeUser(user)
	return c.JSON(fiber.Map{"status": "success", "message": "Created user", "data": userSignup})
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
