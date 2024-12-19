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
		pgx.NamedArgs{"email": user.Email, "userName": user.UserName, "password": user.Password,
			"userType": user.UserType, "userStatus": user.UserStatus, "imageUrl": user.ImageUrl},
		&models.User{})
	return rows, err
}

func UpdateUserStatus(userId uint, userStatus string) ([]models.User, error) {
	rows, err := Query(`UPDATE users SET user_status = @userStatus WHERE id = @userId RETURNING *;`,
		pgx.NamedArgs{"userStatus": userStatus, "userId": userId},
		&models.User{})
	return rows, err
}
