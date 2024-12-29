package db

import (
	"fmt"
	"go-fiber-api-starter/internal/enums/userstatus"
	"go-fiber-api-starter/internal/models"

	"github.com/jackc/pgx/v5"
)

func GetUserById(id uint) (models.User, error) {
	row, err := One("SELECT * FROM users WHERE id = @id LIMIT 1;",
		pgx.NamedArgs{"id": id},
		&models.User{})
	return row, err
}

func GetUserByEmail(email string) (models.User, error) {
	row, err := One("SELECT * FROM users WHERE email = @email LIMIT 1;",
		pgx.NamedArgs{"email": email},
		&models.User{})
	return row, err
}

func GetUserByUsername(username string) (models.User, error) {
	row, err := One("SELECT * FROM users WHERE username = @username LIMIT 1;",
		pgx.NamedArgs{"username": username},
		&models.User{})
	return row, err
}
func InsertUser(user *models.User) (models.User, error) {
	row, err := One(`INSERT INTO users (email, username, password, otp, role, status, image_url, deleted_at)
					 VALUES (@email, @username, @password, @otp, @role, @status, @imageUrl, @deletedAt)
					 RETURNING *;`,
		pgx.NamedArgs{"email": user.Email, "username": user.Username, "password": user.Password, "otp": user.OTP,
			"role": user.Role, "status": user.Status, "imageUrl": user.ImageUrl, "deletedAt": user.DeletedAt},
		&models.User{})
	return row, err
}

func UpdateUser(id uint, userUpdate *models.UserUpdate) (models.User, error) {
	row, err := One(`UPDATE users
					 SET email = @email,
						 username = @username,
						 password = @password,
						 otp = @otp,
						 role = @role,
						 status = @status,
						 image_url = @imageUrl
					 WHERE id = @id RETURNING *;`,
		pgx.NamedArgs{
			"email":    userUpdate.Email,
			"username": userUpdate.Username,
			"password": userUpdate.Password,
			"otp":      userUpdate.OTP,
			"role":     userUpdate.Role,
			"status":   userUpdate.Status,
			"imageUrl": userUpdate.ImageUrl,
			"id":       id},
		&models.User{})
	return row, err
}

func UpdatePassword(userId uint, password string) (models.User, error) {
	row, err := One(`UPDATE users SET password = @password WHERE id = @userId RETURNING *;`,
		pgx.NamedArgs{"password": password, "userId": userId},
		&models.User{})
	return row, err
}

func SoftDeleteUser(id uint) error {
	if err := None(`UPDATE users
					SET status = @status, deleted_at = CURRENT_TIMESTAMP
					WHERE id = @id;`,
		pgx.NamedArgs{"status": userstatus.DELETED, "id": id}); err != nil {
		return err
	}
	return nil
}

func HardDeleteUser(id uint) error {
	if err := None(`DELETE FROM users WHERE id = @id;`,
		pgx.NamedArgs{"id": id}); err != nil {
		return err
	}
	return nil
}

func CheckEmailAvailability(email string) (string, error) {
	_, err := GetUserByEmail(email)
	if err != nil {
		if err != pgx.ErrNoRows {
			// some error other than ErrNoRows
			return "Server error on email lookup", err
		}
		// user not found (email is available)
		return "Email is available", nil
	} else {
		// user found (email is NOT available)
		return "Email address is already in use by another user", fmt.Errorf("Email address is already in use by another user")
	}
}

func CheckUsernameAvailability(username string) (string, error) {
	_, err := GetUserByUsername(username)
	if err != nil {
		if err != pgx.ErrNoRows {
			// some error other than ErrNoRows
			return "Server error on username lookup", err
		}
		// user not found (username is available)
		return "Username is available", nil
	} else {
		// user found (username is NOT available)
		return "Username is already in use by another user", fmt.Errorf("username is already in use by another user")
	}
}
