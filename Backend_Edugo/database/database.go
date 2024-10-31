package database

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DatabaseInit() {
	var err error
	err = godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	DB_USER := os.Getenv("DB_USER")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_HOST := os.Getenv("DB_HOST")
	DB_PORT := os.Getenv("DB_PORT")
	DB_NAME := os.Getenv("DB_NAME")
	DB_OPTION := os.Getenv("DB_OPTION")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME, DB_OPTION)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}
}
