package database

import (
	"fmt"
	"os"

	"go-fiber-api-starter/internal/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectDB connect to db
func ConnectDB() {
	var err error // this is required, otherwise we get a panic elsewhere
	DB, err = gorm.Open(sqlite.Open(os.Getenv("DB_URL")), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println("Connection Opened to Database")
	DB.AutoMigrate(&model.Something{}, &model.User{})
	fmt.Println("Database Migrated")
}
