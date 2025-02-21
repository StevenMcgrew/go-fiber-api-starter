package handlers

import (
	"go-fiber-api-starter/internal/db"
	"go-fiber-api-starter/internal/enums/userstatus"
	"go-fiber-api-starter/internal/mail"
	"go-fiber-api-starter/internal/models"
	"go-fiber-api-starter/internal/serialization"
	"go-fiber-api-starter/internal/utils"
	"go-fiber-api-starter/internal/validation"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

func VerifyEmail(c *fiber.Ctx) error {
	// Shape of request body
	type reqBody struct {
		Email            string `json:"email" form:"email"`
		VerificationCode string `json:"verificationCode" form:"verificationCode"`
	}
	body := &reqBody{}

	// Parse body
	if err := c.BodyParser(body); err != nil {
		return fiber.NewError(400, "Error parsing request body: "+err.Error())
	}

	// Validate input
	if !validation.IsOtpValid(body.VerificationCode) {
		return fiber.NewError(400, "Verification code is invalid")
	}

	// Get user
	user, err := db.GetUserByEmail(body.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fiber.NewError(400, "No user was found with that email address. Please sign up first.")
		}
		return fiber.NewError(500, "Error getting user from database: "+err.Error())
	}

	// Check user status
	if user.Status == userstatus.SUSPENDED || user.Status == userstatus.DELETED {
		return fiber.NewError(400, "Cannot perform verification because user is: "+user.Status)
	}
	if user.Status == userstatus.VERIFIED {
		return fiber.NewError(400, "User's email has already been verified")
	}

	// Verify otp matches
	if body.VerificationCode != user.Otp {
		return fiber.NewError(400, "Cannot perform verification because the code did not match")
	}

	// Save to db
	updatedUser, err := db.UpdateUser(user.Id, &models.UserUpdate{
		Email:    user.Email,
		Username: user.Username,
		Otp:      "",
		Role:     user.Role,
		Status:   userstatus.VERIFIED,
		ImageUrl: user.ImageUrl,
	})
	if err != nil {
		return fiber.NewError(500, "Error updating user in database: "+err.Error())
	}

	// Create JWT
	jwt, err := utils.CreateJWT(&user)
	if err != nil {
		return fiber.NewError(500, "Error creating JWT: "+err.Error())
	}

	// Serialize user
	userResponse := serialization.UserResponse(&updatedUser)

	// Create UserLoginResponse
	userLoginResponse := &models.UserLoginResponse{
		Token:        jwt,
		UserResponse: *userResponse,
	}

	// Send user and jwt
	return utils.SendSuccessJSON(c, 200, userLoginResponse, "Verified")
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
		if err == pgx.ErrNoRows {
			return fiber.NewError(400, "No user was found with that email address. Please sign up first.")
		}
		return fiber.NewError(500, "Error getting user from database: "+err.Error())
	}

	// Send verification email
	err = mail.EmailTheVerificationCode(user.Email, user.Otp)
	if err != nil {
		return fiber.NewError(500, "An error occurred when sending email verification: "+err.Error())
	}

	// Serialize user
	userResponse := serialization.UserResponse(&user)

	// Response
	return utils.SendSuccessJSON(c, 200, userResponse, "Sent verification code again")
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
		if err == pgx.ErrNoRows {
			return fiber.NewError(400, "Login credentials are incorrect")
		}
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

func ResetPasswordRequest(c *fiber.Ctx) error {
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
		if err == pgx.ErrNoRows {
			return fiber.NewError(400, "No user was found with that email address")
		}
		return fiber.NewError(500, "Error getting user from database: "+err.Error())
	}

	// Check user status
	if user.Status == userstatus.SUSPENDED || user.Status == userstatus.DELETED {
		return fiber.NewError(400, "Cannot reset password because user is: "+user.Status)
	}

	// Create updated user with new otp
	u := &models.UserUpdate{
		Email:    user.Email,
		Username: user.Username,
		Otp:      utils.RandomSixDigitStr(),
		Role:     user.Role,
		Status:   user.Status,
		ImageUrl: user.ImageUrl,
	}

	// Save updated user
	updatedUser, err := db.UpdateUser(user.Id, u)
	if err != nil {
		return fiber.NewError(500, "Error updating user in database: "+err.Error())
	}

	// Send password reset email
	err = mail.EmailThePasswordResetCode(updatedUser.Email, updatedUser.Otp)
	if err != nil {
		return fiber.NewError(500, "An error occurred when sending email for password reset: "+err.Error())
	}

	// Send response
	return utils.SendSuccessJSON(c, 200, map[string]any{}, "A password reset code has been emailed")
}

func ResetPasswordUpdate(c *fiber.Ctx) error {
	// Shape of request body
	type reqBody struct {
		Email             string `json:"email" form:"email"`
		ResetCode         string `json:"resetCode" form:"resetCode"`
		NewPassword       string `json:"newPassword" form:"newPassword"`
		RepeatNewPassword string `json:"repeatNewPassword" form:"repeatNewPassword"`
	}
	body := &reqBody{}

	// Parse body
	if err := c.BodyParser(body); err != nil {
		return fiber.NewError(400, "Error parsing request body: "+err.Error())
	}

	// Validate inputs
	warnings := make([]string, 0, 4)
	if !validation.IsEmailValid(body.Email) {
		warnings = append(warnings, "Email is invalid.")
	}
	if !validation.IsOtpValid(body.ResetCode) {
		warnings = append(warnings, "The password reset code is invalid.")
	}
	if !validation.IsPasswordValid(body.NewPassword) {
		warnings = append(warnings, "NewPassword is invalid.")
	}
	if body.NewPassword != body.RepeatNewPassword {
		warnings = append(warnings, "RepeatNewPassword does not match NewPassword.")
	}
	if len(warnings) > 0 {
		return fiber.NewError(400, strings.Join(warnings, " "))
	}

	// Get user by email
	user, err := db.GetUserByEmail(body.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fiber.NewError(400, "No user was found with that email address")
		}
		return fiber.NewError(500, "Error getting user from database: "+err.Error())
	}

	// Check user status
	if user.Status == userstatus.SUSPENDED || user.Status == userstatus.DELETED {
		return fiber.NewError(400, "Cannot reset password because user is: "+user.Status)
	}

	// Verify otp matches
	if body.ResetCode != user.Otp {
		return fiber.NewError(400, "Cannot reset password because the reset code did not match")
	}

	// Hash new password
	pwdBytes, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fiber.NewError(500, "Error hashing password: "+err.Error())
	}

	// Set password
	hashedPassword := string(pwdBytes)

	// Save to db
	updatedUser, err := db.UpdatePassword(user.Id, hashedPassword)
	if err != nil {
		return fiber.NewError(500, "Error updating password in database: "+err.Error())
	}

	// Create JWT
	jwt, err := utils.CreateJWT(&updatedUser)
	if err != nil {
		return fiber.NewError(500, "Error creating JWT: "+err.Error())
	}

	// Serialize user
	userResponse := serialization.UserResponse(&updatedUser)

	// Create UserLoginResponse
	userLoginResponse := &models.UserLoginResponse{
		Token:        jwt,
		UserResponse: *userResponse,
	}

	// Send user and jwt
	return utils.SendSuccessJSON(c, 200, userLoginResponse, "Updated user password")
}
