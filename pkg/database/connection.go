package database

import (
	"fmt"
	"log"
	"os"

	"ramah-disabilitas-be/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connection established")
}

func Migrate() {
	err := DB.AutoMigrate(
		&model.User{},
		&model.Friendship{},
		&model.Subtest{},
		&model.Material{},
		&model.Question{},
		&model.Match{},
		&model.MatchDetail{},
		&model.PracticeSession{},
		&model.QuestionReport{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Println("Database migration completed")
}
