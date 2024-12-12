package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func ConnectDB() {

	var err error
	Pool, err = pgxpool.New(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	// defer Pool.Close()

	// connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", username, password, host, port, database, schema)
	// DB, err := sql.Open("pgx", connStr)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	fmt.Println("Connection pool created for the database.")

	// Create database tables, etc., if not exists

	fmt.Println("Database is ready.")

	// var err error // this is required, otherwise we get a panic elsewhere
	// DB, err = gorm.Open(sqlite.Open(os.Getenv("DB_URL")), &gorm.Config{})
	// if err != nil {
	// 	panic("failed to connect database")
	// }

	// fmt.Println("Connection Opened to Database")
	// DB.AutoMigrate(&model.Something{}, &model.User{})
	// fmt.Println("Database Migrated")
}
