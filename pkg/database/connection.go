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
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta default_query_exec_mode=simple_protocol",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: false,
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connection established")
}

func Migrate() {
	if os.Getenv("APP_ENV") != "production" {
		// DB.Migrator().DropTable("course_students") // Added this line
		// DB.Migrator().DropTable(
		// 	&model.User{},
		// 	&model.Friendship{},
		// 	&model.Subtest{},
		// 	&model.AccessibilityProfile{},
		// 	&model.Course{},
		// 	&model.Module{},
		// 	&model.Material{},
		// 	&model.Question{},
		// 	&model.Assignment{},
		// 	&model.SmartFeature{},
		// 	&model.Match{},
		// 	&model.MatchDetail{},
		// 	&model.PracticeSession{},
		// 	&model.QuestionReport{},
		// 	&model.Submission{},
		// 	&model.MaterialCompletion{},
		// )

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
