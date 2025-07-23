package database

import (
	"fmt"
	"log"
	"order-api/models"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading env file")
	}

	dsn := os.Getenv("DATABASE_URL")

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to connect database", err)
	}

	fmt.Println("Database is successfully connected")
	DB.AutoMigrate(&models.User{}, &models.Product{}, &models.Order{})
}
