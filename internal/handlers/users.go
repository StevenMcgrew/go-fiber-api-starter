package handlers

import (
	"fmt"
	"go-fiber-api-starter/internal/config"
	"go-fiber-api-starter/internal/db"
	"go-fiber-api-starter/internal/enums/userrole"
	"go-fiber-api-starter/internal/enums/userstatus"
	"go-fiber-api-starter/internal/mail"
	"go-fiber-api-starter/internal/models"
	"go-fiber-api-starter/internal/serialization"
	"go-fiber-api-starter/internal/utils"
	"go-fiber-api-starter/internal/validation"
	"math"
	"strings"
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
		return fiber.NewError(400, strings.Join(warnings, " "))
	}

	// Check if email is already taken
	err := db.CheckEmailAvailability(userSignUp.Email)
	if err != nil {
		return fiber.NewError(400, err.Error())
	}

	// Check if username is already taken
	err = db.CheckUsernameAvailability(userSignUp.Username)
	if err != nil {
		return fiber.NewError(400, err.Error())
	}

	// Hash password
	pwdBytes, err := bcrypt.GenerateFromPassword([]byte(userSignUp.Password), bcrypt.DefaultCost)
	if err != nil {
		return fiber.NewError(500, "Error hashing password: "+err.Error())
	}

	// Create new user
	u := &models.User{
		Email:    userSignUp.Email,
		Username: userSignUp.Username,
		Password: string(pwdBytes),
		Otp:      utils.RandomSixDigitStr(),
		Role:     userrole.REGULAR,
		Status:   userstatus.UNVERIFIED,
		ImageUrl: "/default-profile-pic.png",
	}

	// Save user
	user, err := db.InsertUser(u)
	if err != nil {
		return fiber.NewError(500, "Error saving user to database: "+err.Error())
	}

	// Send verification email
	err = mail.SendEmailVerification(user.Email, user.Otp)
	if err != nil {
		return fiber.NewError(500, "Saved new user, but an error occurred when sending email verification: "+err.Error())
	}

	// Serialize user
	userResponse := serialization.UserResponse(&user)

	// Response
	return utils.SendSuccessJSON(c, 201, userResponse, "Saved new user")
}

// ?page=5&per_page=20&query=where.status.eq.verified.and.role.eq.regular.orderby.id.asc,created_at.desc
func GetAllUsers(c *fiber.Ctx) error {
	// Expected query parameters
	type queryParams struct {
		Page    uint   `query:"page"`
		PerPage uint   `query:"per_page"`
		Query   string `query:"query"`
	}
	qParams := &queryParams{}

	// Parse
	if err := c.QueryParser(qParams); err != nil {
		return fiber.NewError(400, "Error parsing query parameters: "+err.Error())
	}

	// Set simpler var names
	page := qParams.Page
	perPage := qParams.PerPage
	query := qParams.Query

	// Get row rowCount
	rowCount, err := db.GetRowCount("users")
	if err != nil {
		return fiber.NewError(500, "Error getting row count: "+err.Error())
	}
	if rowCount == 0 {
		return fiber.NewError(400, "No records were found in the database")
	}

	// Validate
	if page < 1 {
		return fiber.NewError(400, "Page number must be 1 or greater")
	}
	floatPageCount := math.Ceil(float64(rowCount) / float64(perPage))
	pageCount := uint(floatPageCount)
	if page > pageCount {
		return fiber.NewError(400, "The page number requested is larger than the total number of pages")
	}

	// Get users
	users, sql, err := db.GetUsers(page, perPage, query)
	if err != nil {
		return fiber.NewError(500, "Error getting users from database: "+err.Error())
	}

	// Type assert user that is requesting access (should be in c.Locals() from Authn() middleware)
	inquirer, ok := c.Locals("inquirer").(*models.User)
	if !ok {
		return fiber.NewError(500, `Type assertion failed for c.Locals("inquirer")`)
	}

	// Determine if the inquirer is an admin (for use when serializing)
	isAdmin := (inquirer.Role == userrole.ADMIN)

	// Serialize users for response
	var serializedUsers []*models.UserResponse
	if isAdmin {
		for _, user := range users {
			serializedUsers = append(serializedUsers, serialization.UserResponse(&user))
		}
	} else {
		for _, user := range users {
			isOwner := (inquirer.Id == user.Id)
			if !isOwner {
				user.Email = "********"
			}
			serializedUsers = append(serializedUsers, serialization.UserResponse(&user))
		}
	}

	// Create pagination data for response
	pre := "/api/v1"
	selfLink := fmt.Sprintf("%s/users?page=%d&per_page=%d&query=%s", pre, page, perPage, query)
	firstLink := fmt.Sprintf("%s/users?page=%d&per_page=%d&query=%s", pre, 1, perPage, query)
	previousLink := fmt.Sprintf("%s/users?page=%d&per_page=%d&query=%s", pre, page-1, perPage, query)
	if page == 1 {
		previousLink = ""
	}
	nextLink := fmt.Sprintf("%s/users?page=%d&per_page=%d&query=%s", pre, page+1, perPage, query)
	if page == pageCount {
		nextLink = ""
	}
	lastLink := fmt.Sprintf("%s/users?page=%d&per_page=%d&query=%s", pre, pageCount, perPage, query)
	pageData := &models.Pagination{
		Page:         page,
		PerPage:      perPage,
		TotalPages:   pageCount,
		TotalCount:   rowCount,
		SelfLink:     selfLink,
		FirstLink:    firstLink,
		PreviousLink: previousLink,
		NextLink:     nextLink,
		LastLink:     lastLink,
	}

	// Respond
	return utils.SendPaginationJSON(c, serializedUsers, pageData, sql)
}

func GetUser(c *fiber.Ctx) error {
	// Type assert user (the user should be in c.Locals() from AttachUser() middleware)
	user, ok := c.Locals("user").(*models.User)
	if !ok {
		return fiber.NewError(500, `Type assertion failed for c.Locals("user")`)
	}

	// Type assert user that is requesting access (should be in c.Locals() from Authn() middleware)
	inquirer, ok := c.Locals("inquirer").(*models.User)
	if !ok {
		return fiber.NewError(500, `Type assertion failed for c.Locals("inquirer")`)
	}

	// Hide email address if not admin or owner
	isAdmin := (inquirer.Role == userrole.ADMIN)
	isOwner := (inquirer.Id == user.Id)
	if !isAdmin && !isOwner {
		user.Email = "********"
	}

	// Serialize user
	userResponse := serialization.UserResponse(user)

	// Send user
	return utils.SendSuccessJSON(c, 200, userResponse, "Retrieved user from database")
}

func IsEmailAvailable(c *fiber.Ctx) error {
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
		return fiber.NewError(400, "Email is invalid")
	}

	// Check availability and respond
	if err := db.CheckEmailAvailability(body.Email); err != nil {
		if err.Error() == "email address is already in use by another user" {
			return utils.SendSuccessJSON(c, 200, false, err.Error())
		}
		return fiber.NewError(500, "Error looking up email in database: "+err.Error())
	}
	return utils.SendSuccessJSON(c, 200, true, "email address is available")
}

func IsUsernameAvailable(c *fiber.Ctx) error {
	// Shape of request body
	type reqBody struct {
		Username string `json:"username" form:"username"`
	}
	body := &reqBody{}

	// Parse
	if err := c.BodyParser(body); err != nil {
		return fiber.NewError(400, "Error parsing request body: "+err.Error())
	}

	// Validate
	if !validation.IsUsernameValid(body.Username) {
		return fiber.NewError(400, "Username is invalid")
	}

	// Check availability and respond
	if err := db.CheckUsernameAvailability(body.Username); err != nil {
		if err.Error() == "username is already in use by another user" {
			return utils.SendSuccessJSON(c, 200, false, err.Error())
		}
		return fiber.NewError(500, "Error looking up username in database: "+err.Error())
	}
	return utils.SendSuccessJSON(c, 200, true, "username is available")
}

func UpdateUser(c *fiber.Ctx) error {
	// Parse
	userUpdate := &models.UserUpdate{}
	if err := c.BodyParser(userUpdate); err != nil {
		return fiber.NewError(400, `Error parsing request body: `+err.Error())
	}

	// Validate
	if warnings := validation.ValidateUserUpdate(userUpdate); warnings != nil {
		return fiber.NewError(400, strings.Join(warnings, " "))
	}

	// Type assert user (the user should be in c.Locals() from AttachUser() middleware)
	user, ok := c.Locals("user").(*models.User)
	if !ok {
		return fiber.NewError(500, `Type assertion failed for c.Locals("user")`)
	}

	// Check email availability, if new email
	if userUpdate.Email != user.Email {
		err := db.CheckEmailAvailability(userUpdate.Email)
		if err != nil {
			return fiber.NewError(400, err.Error())
		}
	}

	// Check username availability, if new username
	if userUpdate.Username != user.Username {
		err := db.CheckUsernameAvailability(userUpdate.Username)
		if err != nil {
			return fiber.NewError(400, err.Error())
		}
	}

	// Save to db
	updatedUser, err := db.UpdateUser(user.Id, userUpdate)
	if err != nil {
		return fiber.NewError(500, "Error updating user in database: "+err.Error())
	}

	// Serialize user
	userResponse := serialization.UserResponse(&updatedUser)

	// Send response
	return utils.SendSuccessJSON(c, 200, userResponse, "Updated user in database")
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
		return fiber.NewError(400, "Error parsing request body: "+err.Error())
	}

	// Validate inputs
	warnings := make([]string, 0, 3)
	if !validation.IsPasswordValid(body.CurrentPassword) {
		warnings = append(warnings, "Current password is invalid.")
	}
	if !validation.IsPasswordValid(body.NewPassword) {
		warnings = append(warnings, "New password is invalid.")
	}
	if body.NewPassword != body.RepeatNewPassword {
		warnings = append(warnings, "NewPassword and RepeatNewPassword do not match.")
	}
	if len(warnings) > 0 {
		return fiber.NewError(400, strings.Join(warnings, " "))
	}

	// Type assert user (the user should be in c.Locals() from AttachUser() middleware)
	user, ok := c.Locals("user").(*models.User)
	if !ok {
		return fiber.NewError(500, `Type assertion failed for c.Locals("user")`)
	}

	// Type assert user that is requesting access (should be in c.Locals() from Authn() middleware)
	inquirer, ok := c.Locals("inquirer").(*models.User)
	if !ok {
		return fiber.NewError(500, `Type assertion failed for c.Locals("inquirer")`)
	}

	// Check password, if not admin (this is an extra security check)
	if inquirer.Role != userrole.ADMIN {
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.CurrentPassword)); err != nil {
			return fiber.NewError(500, "CurrentPassword input doesn't match saved password")
		}
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

	// Serialize user
	userResponse := serialization.UserResponse(&updatedUser)

	// Send response
	return utils.SendSuccessJSON(c, 200, userResponse, "Saved password to database")
}

func UpdateEmail(c *fiber.Ctx) error {
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
		return fiber.NewError(500, "Email input is invalid")
	}

	// Check if email is already taken
	err := db.CheckEmailAvailability(body.Email)
	if err != nil {
		return fiber.NewError(400, err.Error())
	}

	// Parse userId from path
	id, err := c.ParamsInt("userId")
	if err != nil || id == 0 {
		return fiber.NewError(400, "Error parsing path parameter: "+err.Error())
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
	jwtString, err := token.SignedString([]byte(config.API_SECRET))
	if err != nil {
		return fiber.NewError(500, "Error creating JWT: "+err.Error())
	}

	// Create email verification link
	link := fmt.Sprintf("%s/api/v1/auth/verify-email/?token=%s", config.API_BASE_URL, jwtString)

	// Send verification email
	err = mail.SendEmailVerification(body.Email, link)
	if err != nil {
		return fiber.NewError(500, "Error sending email: "+err.Error())
	}

	// Send response
	return utils.SendSuccessJSON(c, 200, nil, "A verification link has been emailed to "+body.Email)
}

func UpdateUsername(c *fiber.Ctx) error {
	// Shape of data in request body
	type reqBody struct {
		Username string `json:"username" form:"username"`
	}
	body := &reqBody{}

	// Parse
	if err := c.BodyParser(body); err != nil {
		return fiber.NewError(400, "Error parsing request body: "+err.Error())
	}

	// Validate
	if !validation.IsUsernameValid(body.Username) {
		return fiber.NewError(400, "Username is invalid")
	}

	// Check if username is already taken
	err := db.CheckUsernameAvailability(body.Username)
	if err != nil {
		return fiber.NewError(400, err.Error())
	}

	// Parse userId from path
	id, err := c.ParamsInt("userId")
	if err != nil || id == 0 {
		return fiber.NewError(400, "Error parsing path parameter: "+err.Error())
	}
	userId := uint(id)

	// Save to db
	updatedUser, err := db.UpdateUsername(userId, body.Username)
	if err != nil {
		return fiber.NewError(400, "Error updating username in database: "+err.Error())
	}

	// Serialize user
	userResponse := serialization.UserResponse(&updatedUser)

	// Send response
	return utils.SendSuccessJSON(c, 200, userResponse, "Saved username to database")
}

func SoftDeleteUser(c *fiber.Ctx) error {
	// Parse userId from path
	id, err := c.ParamsInt("userId")
	if err != nil || id == 0 {
		return fiber.NewError(400, "Error parsing path parameter: "+err.Error())
	}
	userId := uint(id)

	// Soft delete the user
	row, err := db.SoftDeleteUser(userId)
	if err != nil {
		return fiber.NewError(500, "Error deleting user from database: "+err.Error())
	}

	// Send response
	data := map[string]any{"id": row.Id}
	return utils.SendSuccessJSON(c, 200, data, "Deleted user from database")
}
