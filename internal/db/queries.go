package db

import (
	"fmt"
	"go-fiber-api-starter/internal/models"

	"github.com/jackc/pgx/v5"
)

func Query[T any](sql string, args pgx.NamedArgs, ptrModel *T) ([]T, error) {
	// Run the query
	rows, err := Pool.Query(Ctx, sql, args)
	if err != nil {
		fmt.Println("Error querying the database:", err)
		return nil, err
	}
	// Parse the rows
	parsedRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[T])
	if err != nil {
		fmt.Println("Error parsing database rows:", err)
		return nil, err
	}
	return parsedRows, nil
}

func GetUserById(id uint) ([]models.User, error) {
	rows, err := Query("SELECT * FROM users WHERE id = @id LIMIT 1;",
		pgx.NamedArgs{"id": id},
		&models.User{})
	return rows, err
}

func GetUserByEmail(email string) ([]models.User, error) {
	rows, err := Query("SELECT * FROM users WHERE email = @email LIMIT 1;",
		pgx.NamedArgs{"email": email},
		&models.User{})
	return rows, err
}

func GetUserByUserName(userName string) ([]models.User, error) {
	rows, err := Query("SELECT * FROM users WHERE username = @userName LIMIT 1;",
		pgx.NamedArgs{"userName": userName},
		&models.User{})
	return rows, err
}
func InsertUser(user *models.User) ([]models.User, error) {
	rows, err := Query(`INSERT INTO users (email, username, password, user_type, user_status, image_url)
						VALUES (@email, @userName, @password, @userType, @userStatus, @imageUrl)
						RETURNING *;`,
		pgx.NamedArgs{"email": user.Email, "userName": user.Username, "password": user.Password,
			"userType": user.Role, "userStatus": user.Status, "imageUrl": user.ImageUrl},
		&models.User{})
	return rows, err
}

func UpdateUser(user *models.User) ([]models.User, error) {
	rows, err := Query(`UPDATE users
						SET email = @email,
							username = @userName,
							password = @password,
							user_type = @userType,
							user_status = @userStatus,
							image_url = @imageUrl
						WHERE id = @userId RETURNING *;`,
		pgx.NamedArgs{
			"userId":     user.Id,
			"email":      user.Email,
			"userName":   user.Username,
			"password":   user.Password,
			"userType":   user.Role,
			"userStatus": user.Status,
			"imageUrl":   user.ImageUrl},
		&models.User{})
	return rows, err
}
