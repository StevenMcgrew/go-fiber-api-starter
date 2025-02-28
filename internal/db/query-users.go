package db

import (
	"fmt"
	"go-fiber-api-starter/internal/enums/userstatus"
	"go-fiber-api-starter/internal/models"

	"github.com/jackc/pgx/v5"
)

func InsertUser(user *models.User) (models.User, error) {
	row, err := One(`INSERT INTO users (email, username, password, otp, role, status, image_url)
					 VALUES (@email,
					 		 @username,
							 @password,
							 @otp,
							 @role,
							 @status,
							 @imageUrl)
					 RETURNING *;`,
		pgx.NamedArgs{
			"email":    user.Email,
			"username": user.Username,
			"password": user.Password,
			"otp":      user.Otp,
			"role":     user.Role,
			"status":   user.Status,
			"imageUrl": user.ImageUrl},
		&models.User{})
	return row, err
}

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

func GetUsers(page uint, perPage uint, query string) ([]models.User, string, error) {
	// Query builder
	qb := NewQueryBuilder(
		page,
		perPage,
		query,
		"users",
		[]string{
			"id",
			"email",
			"username",
			"password",
			"role",
			"status",
			"image_url",
			"created_at",
			"updated_at",
			"deleted_at",
		},
	)

	// Build the query string
	queryString, err := qb.Build()
	if err != nil {
		return nil, "", err
	}

	// Run the query
	rows, err := Many(queryString, pgx.NamedArgs{}, &models.User{})
	return rows, queryString, err
}

func UpdateUser(id uint, userUpdate *models.UserUpdate) (models.User, error) {
	row, err := One(`UPDATE users
					 SET email = @email,
						 username = @username,
						 otp = @otp,
						 role = @role,
						 status = @status,
						 image_url = @imageUrl
					 WHERE id = @id RETURNING *;`,
		pgx.NamedArgs{
			"email":    userUpdate.Email,
			"username": userUpdate.Username,
			"otp":      userUpdate.Otp,
			"role":     userUpdate.Role,
			"status":   userUpdate.Status,
			"imageUrl": userUpdate.ImageUrl,
			"id":       id},
		&models.User{})
	return row, err
}

func UpdateImageUrl(id uint, imageUrl string) (models.User, error) {
	row, err := One(`UPDATE users SET image_url = @imageUrl WHERE id = @id RETURNING *;`,
		pgx.NamedArgs{"imageUrl": imageUrl, "id": id},
		&models.User{})
	return row, err
}

func UpdateUsername(id uint, username string) (models.User, error) {
	row, err := One(`UPDATE users SET username = @username WHERE id = @id RETURNING *;`,
		pgx.NamedArgs{"username": username, "id": id},
		&models.User{})
	return row, err
}

func UpdateEmail(id uint, email string) (models.User, error) {
	row, err := One(`UPDATE users SET email = @email WHERE id = @id RETURNING *;`,
		pgx.NamedArgs{"email": email, "id": id},
		&models.User{})
	return row, err
}

func UpdatePassword(id uint, password string) (models.User, error) {
	row, err := One(`UPDATE users SET password = @password WHERE id = @id RETURNING *;`,
		pgx.NamedArgs{"password": password, "id": id},
		&models.User{})
	return row, err
}

func UpdateOtp(id uint, otp string) (models.User, error) {
	row, err := One(`UPDATE users SET otp = @otp WHERE id = @id RETURNING *;`,
		pgx.NamedArgs{"otp": otp, "id": id},
		&models.User{})
	return row, err
}

func SoftDeleteUser(id uint) (models.User, error) {
	row, err := One(`UPDATE users
					SET status = @status, deleted_at = CURRENT_TIMESTAMP
					WHERE id = @id
					RETURNING *;`,
		pgx.NamedArgs{
			"status": userstatus.DELETED,
			"id":     id},
		&models.User{})
	return row, err
}

func HardDeleteUser(id uint) (models.User, error) {
	row, err := One(`DELETE FROM users WHERE id = @id RETURNING *;`,
		pgx.NamedArgs{"id": id},
		&models.User{})
	return row, err
}

func CheckEmailAvailability(email string) error {
	_, err := GetUserByEmail(email)
	if err != nil {
		if err != pgx.ErrNoRows {
			// some error other than ErrNoRows
			return err
		}
		// user not found (email is available)
		return nil
	} else {
		// user found (email is NOT available)
		return fmt.Errorf("email address is already in use by another user")
	}
}

func CheckUsernameAvailability(username string) error {
	_, err := GetUserByUsername(username)
	if err != nil {
		if err != pgx.ErrNoRows {
			// some error other than ErrNoRows
			return err
		}
		// user not found (username is available)
		return nil
	} else {
		// user found (username is NOT available)
		return fmt.Errorf("username is already in use by another user")
	}
}
