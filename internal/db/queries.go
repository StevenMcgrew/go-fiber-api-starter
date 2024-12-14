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
	parsedRows, parsedErr := pgx.CollectRows(rows, pgx.RowToStructByName[T])
	if parsedErr != nil {
		fmt.Println("Error parsing database rows:", parsedErr)
		return nil, parsedErr
	}
	return parsedRows, nil
}

func getUserByEmail(email string) ([]models.User, error) {
	rows, err := Pool.Query(Ctx,
		"SELECT * FROM users WHERE email = $1 LIMIT 1;",
		email)
	if err != nil {
		fmt.Println("Error getting user by email address:", err)
		return nil, err
	}
	var parsedRows []models.User
	parsedRows, err = pgx.CollectRows(rows, pgx.RowToStructByName[models.User])
	if err != nil {
		fmt.Println("Error parsing db rows to User:", err)
		return nil, err
	}
	return parsedRows, nil
}
