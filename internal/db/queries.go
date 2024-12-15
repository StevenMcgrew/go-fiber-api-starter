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
