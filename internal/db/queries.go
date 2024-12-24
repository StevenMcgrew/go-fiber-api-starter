package db

import (
	"go-fiber-api-starter/internal/models"
	"log"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5"
)

func ExecuteSqlFile(path string) {
	path, pathErr := filepath.Abs(path)
	if pathErr != nil {
		log.Fatal("Error getting absolute path to sql file:", pathErr)
	}
	bytes, ioErr := os.ReadFile(path)
	if ioErr != nil {
		log.Fatal("Error reading sql file: ", ioErr)
	}
	sql := string(bytes)
	_, execErr := Pool.Exec(Ctx, sql)
	if execErr != nil {
		log.Fatal("Error executing sql from file: ", execErr)
	}
}

func Many[T any](sql string, args pgx.NamedArgs, ptrModel *T) ([]T, error) {
	// Run the query
	rows, err := Pool.Query(Ctx, sql, args)
	if err != nil {
		return nil, err
	}
	// Parse the rows
	parsedRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[T])
	if err != nil {
		return nil, err
	}
	return parsedRows, nil
}

func One[T any](sql string, args pgx.NamedArgs, ptrModel *T) (T, error) {
	// Run the query
	rows, err := Pool.Query(Ctx, sql, args)
	if err != nil {
		return *ptrModel, err
	}
	// Parse the rows
	parsedRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[T])
	if err != nil {
		return *ptrModel, err
	}
	if len(parsedRows) == 0 {
		return *ptrModel, pgx.ErrNoRows
	}
	return parsedRows[0], nil
}

func None(sql string, args pgx.NamedArgs) error {
	_, err := Pool.Exec(Ctx, sql, args)
	if err != nil {
		return err
	}
	return nil
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

func GetUserByUserName(username string) (models.User, error) {
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

func UpdateUser(user *models.User) (models.User, error) {
	row, err := One(`UPDATE users
						SET email = @email,
							username = @username,
							password = @password,
							otp = @otp,
							role = @role,
							status = @status,
							image_url = @imageUrl,
							deleted_at = @deletedAt
						WHERE id = @id RETURNING *;`,
		pgx.NamedArgs{
			"id":        user.Id,
			"email":     user.Email,
			"username":  user.Username,
			"password":  user.Password,
			"otp":       user.OTP,
			"role":      user.Role,
			"status":    user.Status,
			"imageUrl":  user.ImageUrl,
			"deletedAt": user.DeletedAt},
		&models.User{})
	return row, err
}

func InsertNotification(n *models.Notification) (models.Notification, error) {
	row, err := One(`INSERT INTO notifications (text_content, has_viewed, user_id)
						VALUES (@textContent, @hasViewed, @userId)
						RETURNING *;`,
		pgx.NamedArgs{"textContent": n.TextContent, "hasViewed": n.HasViewed, "userId": n.UserId},
		&models.Notification{})
	return row, err
}
