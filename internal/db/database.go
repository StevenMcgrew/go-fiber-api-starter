package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

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
	fmt.Println("Database is ready.")
}
