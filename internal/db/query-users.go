package db

import (
	"fmt"
	"go-fiber-api-starter/internal/enums/userstatus"
	"go-fiber-api-starter/internal/models"

	"github.com/jackc/pgx/v5"
)

func InsertUser(user *models.User) (models.User, error) {
	row, err := One(`INSERT INTO users (email, username, password, role, status, image_url)
					 VALUES (@email,
					 		 @username,
							 @password,
							 @role,
							 @status,
							 @imageUrl)
					 RETURNING *;`,
		pgx.NamedArgs{
			"email":    user.Email,
			"username": user.Username,
			"password": user.Password,
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

func UpdateUser(id uint, userUpdate *models.UserUpdate) (models.User, error) {
	row, err := One(`UPDATE users
					 SET email = @email,
						 username = @username,
						 password = @password,
						 role = @role,
						 status = @status,
						 image_url = @imageUrl
					 WHERE id = @id RETURNING *;`,
		pgx.NamedArgs{
			"email":    userUpdate.Email,
			"username": userUpdate.Username,
			"password": userUpdate.Password,
			"role":     userUpdate.Role,
			"status":   userUpdate.Status,
			"imageUrl": userUpdate.ImageUrl,
			"id":       id},
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

func IsEmailAvailable(email string) (bool, error) {
	_, err := GetUserByEmail(email)
	if err != nil {
		if err != pgx.ErrNoRows {
			// some error other than ErrNoRows
			return false, err
		}
		// user not found (email is available)
		return true, nil
	} else {
		// user found (email is NOT available)
		return false, fmt.Errorf("email address is already in use by another user")
	}
}

func IsUsernameAvailable(username string) (bool, error) {
	_, err := GetUserByUsername(username)
	if err != nil {
		if err != pgx.ErrNoRows {
			// some error other than ErrNoRows
			return false, err
		}
		// user not found (username is available)
		return true, nil
	} else {
		// user found (username is NOT available)
		return false, fmt.Errorf("username is already in use by another user")
	}
}
