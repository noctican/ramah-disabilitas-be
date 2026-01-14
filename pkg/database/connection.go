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
	// DB.Migrator().DropTable(
	// 	&model.Submission{},
	// 	&model.QuestionReport{},
	// 	&model.PracticeSession{},
	// 	&model.MatchDetail{},
	// 	&model.Match{},
	// 	&model.SmartFeature{},
	// 	&model.Assignment{},
	// 	&model.Question{},
	// 	&model.Material{},
	// 	&model.Module{},
	// 	&model.Course{},
	// 	&model.AccessibilityProfile{},
	// 	&model.Subtest{},
	// 	&model.Friendship{},
	// 	&model.User{},
	// 	&model.Course{},
	// 	&model.MaterialCompletion{},
	// )
	if os.Getenv("APP_ENV") != "production" {
		log.Println("Running AutoMigrate...")

		err := DB.AutoMigrate(
			&model.User{},
			&model.Friendship{},
			&model.Subtest{},
			&model.AccessibilityProfile{},
		)
		if err != nil {
			log.Fatal("Failed to migrate Step 1 (Users):", err)
		}

		err = DB.AutoMigrate(
			&model.Course{},
			&model.Module{},
		)
		if err != nil {
			log.Fatal("Failed to migrate Step 2 (Courses):", err)
		}

		err = DB.AutoMigrate(
			&model.Material{},
			&model.Question{},
			&model.Assignment{},
		)
		if err != nil {
			log.Fatal("Failed to migrate Step 3 (Materials):", err)
		}

		err = DB.AutoMigrate(
			&model.SmartFeature{},
			&model.Match{},
			&model.MatchDetail{},
			&model.PracticeSession{},
			&model.QuestionReport{},
			&model.Submission{},
			&model.MaterialCompletion{},
		)
		if err != nil {
			log.Fatal("Failed to migrate Step 4 (Features):", err)
		}

		log.Println("Database migration completed successfully")
	} else {
		log.Println("Production mode: Skipping AutoMigrate to save startup time.")
	}
}
