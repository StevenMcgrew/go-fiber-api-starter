package handlers

import (
	"fmt"
	"go-fiber-api-starter/internal/db"
	"go-fiber-api-starter/internal/models"

	"github.com/jackc/pgx/v5"
)

func getUserByEmail(email string) ([]models.User, error) {
	rows, err := db.Pool.Query(db.Ctx,
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
