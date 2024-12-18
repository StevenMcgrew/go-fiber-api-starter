package handlers

import (
	"os"
	"strconv"

	"go-fiber-api-starter/internal/db"
	"go-fiber-api-starter/internal/enums/userstatus"
	"go-fiber-api-starter/internal/enums/usertype"
	"go-fiber-api-starter/internal/models"
	"go-fiber-api-starter/internal/validation"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func GetUser(c *fiber.Ctx) error {
	return nil
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
	user := &models.User{}
	user.Email = userSignup.Email
	user.UserName = userSignup.UserName
	user.Password = string(pwdBytes)
	user.UserType = usertype.REGULAR
	user.UserStatus = userstatus.UNVERIFIED
	user.ImageUrl = ""

	// Save user to database
	userRows, err := db.InsertUser(user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Server error when saving user", "data": err.Error()})
	}
	if len(userRows) < 0 {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "User record was not returned after database insert", "data": userRows})
	}

	// Create JWT for verification link in email
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":     userRows[0].Id,
		"userType":   userRows[0].UserType,
		"userStatus": userRows[0].UserStatus,
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

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
