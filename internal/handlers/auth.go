package handlers

import (
	"fmt"
	"go-fiber-api-starter/internal/db"
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

func VerifyEmail(c *fiber.Ctx) error {
	// Query params
	type queryParams struct {
		Token string `query:"token"`
	}
	qParams := &queryParams{}

	// Parse query params
	if err := c.QueryParser(qParams); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Error parsing 'token' query param",
			"data": map[string]any{"errorMessage": err.Error()}})
	}

	// Validate JWT
	token, err := jwt.ParseWithClaims(qParams.Token, &models.JwtVerifyEmail{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET")), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return EmailVerificationFailurePage(c, err.Error())
	}
	if !token.Valid {
		return EmailVerificationFailurePage(c, "The token is invalid")
	}
	payload, ok := token.Claims.(*models.JwtVerifyEmail)
	if !ok {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "The claims in the JWT should be of type *models.JwtVerifyEmail",
			"data": map[string]any{"errorMessage": "Wrong type for JWT claims"}})
	}

	// Get user
	user, err := db.GetUserById(payload.UserId)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Server error when getting user from database",
			"data": map[string]any{"errorMessage": err.Error()}})
	}

	// Determine user status
	var status string
	if user.Status == userstatus.SUSPENDED || user.Status == userstatus.DELETED {
		status = user.Status
	} else {
		status = userstatus.VERIFIED
	}

	// Save to db
	_, err = db.UpdateUser(payload.UserId, &models.UserUpdate{
		Email:    payload.Email,
		Username: user.Username,
		Password: user.Password,
		Role:     user.Role,
		Status:   status,
		ImageUrl: user.ImageUrl,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Server error when updating user in database",
			"data": map[string]any{"errorMessage": err.Error()}})
	}

	// Send to EmailVerificationSuccessPage
	return EmailVerificationSuccessPage(c)
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

	// Create JWT for email verification link
	claims := &models.JwtVerifyEmail{
		UserId: user.Id,
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
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Error sending verification email",
			"data": map[string]any{"errorMessage": err.Error()}})
	}

	// Serialize user
	userResponse := serialization.UserResponse(&user)

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

	// Make sure user status is 'verified'
	if user.Status != userstatus.VERIFIED {
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

func ForgotPassword(c *fiber.Ctx) error {
	// Shape of request body
	type reqBody struct {
		Email string `json:"email" form:"email"`
	}
	body := &reqBody{}

	// Parse
	if err := c.BodyParser(body); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Error parsing login credentials",
			"data": map[string]any{"errorMessage": err.Error()}})
	}

	// Validate
	if !validation.IsEmailValid(body.Email) {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "Email input is invalid",
			"data": map[string]any{"errorMessage": "Invalid input"}})
	}

	// Get user from db
	user, err := db.GetUserByEmail(body.Email)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Error getting user from database",
			"data": map[string]any{"errorMessage": err.Error()}})
	}

	// Create JWT for password reset link
	claims := &models.JwtUser{
		UserId: user.Id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Error when creating JWT for password reset link",
			"data": map[string]any{"errorMessage": err.Error()}})
	}

	// Create password reset link
	link := fmt.Sprintf("%s/api/v1/auth/reset-password/request?token=%s", os.Getenv("API_BASE_URL"), jwtString)

	// Send email
	err = mail.SendPasswordReset(body.Email, link)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Error sending verification email",
			"data": map[string]any{"errorMessage": err.Error()}})
	}

	// Serialize user
	userResponse := serialization.UserResponse(&user)

	// Send user in response
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "A password reset link has been emailed",
		"data": map[string]any{"user": userResponse}})
}

func ResetForgottenPassword(c *fiber.Ctx) error {
	// Shape of request body
	type reqBody struct {
		Token             string `json:"token" form:"token"`
		NewPassword       string `json:"newPassword" form:"newPassword"`
		RepeatNewPassword string `json:"repeatNewPassword" form:"repeatNewPassword"`
	}
	body := &reqBody{}

	// Parse
	if err := c.BodyParser(body); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Error parsing request body",
			"data": map[string]any{"errorMessage": err.Error()}})
	}

	// Validate password inputs
	warnings := make([]string, 0, 2)
	if !validation.IsPasswordValid(body.NewPassword) {
		warnings = append(warnings, "NewPassword is invalid")
	}
	if body.NewPassword != body.RepeatNewPassword {
		warnings = append(warnings, "RepeatNewPassword does not match NewPassword")
	}
	if len(warnings) > 0 {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "One or more inputs are invalid",
			"data": map[string]any{"errorMessage": warnings}})
	}

	// Validate JWT
	token, err := jwt.ParseWithClaims(body.Token, &models.JwtUser{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET")), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "Error parsing JWT",
			"data": map[string]any{"errorMessage": err.Error()}})
	}
	if !token.Valid {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "JWT is invalid",
			"data": map[string]any{"errorMessage": "JWT is invalid"}})
	}
	jwtUser, ok := token.Claims.(*models.JwtUser)
	if !ok {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "Wrong type for JWT claims",
			"data": map[string]any{"errorMessage": "Wrong type for JWT claims"}})
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
	updatedUser, err := db.UpdatePassword(jwtUser.UserId, hashedPassword)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "fail", "message": "Error updating password in database",
			"data": map[string]any{"errorMessage": err.Error()}})
	}

	// Serialize user
	userResponse := serialization.UserResponse(&updatedUser)

	// Send response
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Saved password to database",
		"data": map[string]any{"user": userResponse}})
}
