package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	Pool *pgxpool.Pool
	Ctx  context.Context = context.Background()
)

func Connect(connString string) {
	var err error
	Pool, err = pgxpool.New(Ctx, connString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Connection pool created for the database.")
}
