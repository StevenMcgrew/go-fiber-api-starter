package handlers

import (
	"fmt"
	"go-fiber-api-starter/internal/config"
	"go-fiber-api-starter/internal/db"
	"go-fiber-api-starter/internal/enums/userstatus"
	"go-fiber-api-starter/internal/mail"
	"go-fiber-api-starter/internal/models"
	"go-fiber-api-starter/internal/serialization"
	"go-fiber-api-starter/internal/utils"
	"go-fiber-api-starter/internal/validation"
	"strings"
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
		return fiber.NewError(400, "Error parsing query parameter: "+err.Error())
	}

	// Validate JWT
	token, err := jwt.ParseWithClaims(qParams.Token, &models.JwtVerifyEmail{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.API_SECRET), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return FailedToVerifyEmailPage(c, err.Error())
	}
	if !token.Valid {
		return FailedToVerifyEmailPage(c, "The token is invalid")
	}
	payload, ok := token.Claims.(*models.JwtVerifyEmail)
	if !ok {
		return fiber.NewError(400, "Type assertion failed for token claims")
	}

	// Get user
	user, err := db.GetUserById(payload.UserId)
	if err != nil {
		return fiber.NewError(500, "Error getting user from database: "+err.Error())
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
		Role:     user.Role,
		Status:   status,
		ImageUrl: user.ImageUrl,
	})
	if err != nil {
		return fiber.NewError(500, "Error updating user in database: "+err.Error())
	}

	// Send to EmailVerificationSuccessPage
	return SuccessfullyVerifiedEmailPage(c)
}

func ResendEmailVerification(c *fiber.Ctx) error {
	// Shape of request body
	type reqBody struct {
		Email string `json:"email" form:"email"`
	}
	body := &reqBody{}

	// Parse body
	if err := c.BodyParser(body); err != nil {
		return fiber.NewError(400, "Error parsing request body: "+err.Error())
	}

	// Validate input
	if !validation.IsEmailValid(body.Email) {
		return fiber.NewError(400, "Email address is invalid")
	}

	// Get user by email
	user, err := db.GetUserByEmail(body.Email)
	if err != nil {
		return fiber.NewError(500, "Error getting user from database: "+err.Error())
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
	jwtString, err := token.SignedString([]byte(config.API_SECRET))
	if err != nil {
		return fiber.NewError(400, "Error creating JWT: "+err.Error())
	}

	// Create email verification link
	link := fmt.Sprintf("%s/api/v1/auth/verify-email/?token=%s", config.API_BASE_URL, jwtString)

	// Send verification email
	err = mail.SendEmailVerification(body.Email, link)
	if err != nil {
		return fiber.NewError(500, "Error sending email: "+err.Error())
	}

	// Serialize user
	userResponse := serialization.UserResponse(&user)

	// Send user in response
	return utils.SendSuccessJSON(c, 200, userResponse, "A verification email has been sent")
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
		return fiber.NewError(400, "Error parsing request body: "+err.Error())
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
		return fiber.NewError(400, "Login credentials are invalid")
	}

	// Get user by email
	user, err := db.GetUserByEmail(body.Email)
	if err != nil {
		return fiber.NewError(500, "Error getting user from database: "+err.Error())
	}

	// Make sure user status is 'verified'
	if user.Status != userstatus.VERIFIED {
		return fiber.NewError(400, "Cannot login users who have an account status of "+user.Status)
	}

	// Check password
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		return fiber.NewError(400, "Login credentials are incorrect")
	}

	// Create JWT
	jwt, err := utils.CreateJWT(&user)
	if err != nil {
		return fiber.NewError(500, "Error creating JWT: "+err.Error())
	}

	// Serialize user
	userResponse := serialization.UserResponse(&user)

	// Create UserLoginResponse
	userLoginResponse := &models.UserLoginResponse{
		Token:        jwt,
		UserResponse: *userResponse,
	}

	// Send user and jwt
	return utils.SendSuccessJSON(c, 200, userLoginResponse, "Auth token created for user")
}

func ForgotPassword(c *fiber.Ctx) error {
	// Shape of request body
	type reqBody struct {
		Email string `json:"email" form:"email"`
	}
	body := &reqBody{}

	// Parse
	if err := c.BodyParser(body); err != nil {
		return fiber.NewError(400, "Error parsing request body: "+err.Error())
	}

	// Validate
	if !validation.IsEmailValid(body.Email) {
		return fiber.NewError(400, "Email input is invalid")
	}

	// Get user from db
	user, err := db.GetUserByEmail(body.Email)
	if err != nil {
		return fiber.NewError(500, "Error getting user from database: "+err.Error())
	}

	// Create JWT for password reset link
	claims := &models.JwtUser{
		UserId: user.Id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtString, err := token.SignedString([]byte(config.API_SECRET))
	if err != nil {
		return fiber.NewError(500, "Error creating JWT: "+err.Error())
	}

	// Create password reset link
	link := fmt.Sprintf("%s/api/v1/auth/reset-password/request?token=%s", config.API_BASE_URL, jwtString)

	// Send email
	err = mail.SendPasswordReset(body.Email, link)
	if err != nil {
		return fiber.NewError(500, "Error emailing password reset: "+err.Error())
	}

	// Serialize user
	userResponse := serialization.UserResponse(&user)

	// Send user in response
	return utils.SendSuccessJSON(c, 200, userResponse, "A password reset link has been emailed")
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
		return fiber.NewError(400, "Error parsing request body: "+err.Error())
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
		return fiber.NewError(400, strings.Join(warnings, " "))
	}

	// Validate JWT
	token, err := jwt.ParseWithClaims(body.Token, &models.JwtUser{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.API_SECRET), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return fiber.NewError(500, "Error parsing JWT: "+err.Error())
	}
	if !token.Valid {
		return fiber.NewError(400, "JWT is invalid")
	}
	jwtUser, ok := token.Claims.(*models.JwtUser)
	if !ok {
		return fiber.NewError(400, "Type assertion failed for token claims")
	}

	// Hash new password
	pwdBytes, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fiber.NewError(500, "Error hashing password: "+err.Error())
	}

	// Set password
	hashedPassword := string(pwdBytes)

	// Save to db
	updatedUser, err := db.UpdatePassword(jwtUser.UserId, hashedPassword)
	if err != nil {
		return fiber.NewError(500, "Error updating password in database: "+err.Error())
	}

	// Serialize user
	userResponse := serialization.UserResponse(&updatedUser)

	// Send response
	return utils.SendSuccessJSON(c, 200, userResponse, "Saved password to database")
}
