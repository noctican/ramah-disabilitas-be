package service

import (
	"encoding/csv"
	"io"
	"ramah-disabilitas-be/internal/model"
	"ramah-disabilitas-be/internal/repository"
	"ramah-disabilitas-be/pkg/utils"
	"strings"
)

type AccessibilityInput struct {
	Categories []string `json:"categories" binding:"required"`
}

type CreateStudentInput struct {
	Name       string   `json:"name" binding:"required"`
	Email      string   `json:"email" binding:"required,email"`
	Password   string   `json:"password" binding:"required,min=6"`
	Categories []string `json:"disabilities"`
}

func CreateStudent(input CreateStudentInput) (*model.User, error) {
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: hashedPassword,
		Role:     model.RoleStudent,
	}

	if err := repository.CreateUser(user); err != nil {
		return nil, err
	}

	// Create Accessibility Profile
	if len(input.Categories) > 0 {
		_, err := UpdateAccessibilityProfile(user.ID, AccessibilityInput{Categories: input.Categories})
		if err != nil {
			return nil, err
		}
	} else {
		// Create default empty profile
		profile := &model.AccessibilityProfile{UserID: user.ID}
		repository.SaveAccessibilityProfile(profile)
	}

	return user, nil
}

func UpdateAccessibilityProfile(userID uint64, input AccessibilityInput) (*model.AccessibilityProfile, error) {
	profile := &model.AccessibilityProfile{
		UserID: userID,
	}

	for _, category := range input.Categories {
		category = strings.ToLower(strings.TrimSpace(category))
		switch category {
		case "a", "vision", "penglihatan", "tuna netra", "tuna_netra":
			profile.VisionImpaired = true
			profile.ScreenReaderCompatible = true
			profile.AudioDescription = true
		case "b", "hearing", "pendengaran", "tuna rungu", "tuna_rungu":
			profile.HearingImpaired = true
			profile.SubtitlesRequired = true
			profile.VisualNotifications = true
		case "c", "physical", "motorik", "daksa", "tuna daksa", "tuna_daksa":
			profile.PhysicalImpaired = true
			profile.KeyboardNavigation = true
			profile.VoiceCommand = true
		case "d", "cognitive", "fokus", "adhd", "disleksia", "kesulitan kognitif", "kesulitan_kognitif", "tuna grahita":
			profile.CognitiveImpaired = true
			profile.AISummary = true
			profile.FocusMode = true
		case "e", "speech", "wicara", "bisu", "tuna wicara", "tuna_wicara":
			profile.SpeechImpaired = true
			profile.TextBasedSubmission = true
		}
	}

	if err := repository.SaveAccessibilityProfile(profile); err != nil {
		return nil, err
	}

	return profile, nil
}

func GetAccessibilityProfile(userID uint64) (*model.AccessibilityProfile, error) {
	return repository.FindAccessibilityProfileByUserID(userID)
}

func ImportStudentsFromCSV(reader io.Reader) ([]*model.User, error) {
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	var createdUsers []*model.User

	// Skip header if it exists (assuming first row is header)
	if len(records) > 0 {
		firstCell := strings.ToLower(records[0][0])
		if firstCell == "name" || firstCell == "nama" {
			records = records[1:]
		}
	}

	for _, record := range records {
		if len(record) < 3 {
			continue // Skip invalid rows
		}

		name := record[0]
		email := record[1]
		password := record[2]

		var disabilities []string
		if len(record) > 3 && record[3] != "" {
			// Split by semicolon or comma (if strictly CSV, comma is delimiter, so inside a field it usually avoids comma unless quoted. Let's assume semicolon for multiple disabilities)
			disabilities = strings.Split(record[3], ";")
			for i := range disabilities {
				disabilities[i] = strings.TrimSpace(disabilities[i])
			}
		}

		input := CreateStudentInput{
			Name:       name,
			Email:      email,
			Password:   password,
			Categories: disabilities,
		}

		user, err := CreateStudent(input)
		if err != nil {
			// If error (e.g. duplicate email), skip or return error?
			// Ideally collect errors, but for now let's just log or skip?
			// Let's return error to be safe, or continue?
			// User request "import", usually expects all or nothing or partial success.
			// Let's stop on error for simplicity, or maybe just skip duplicates.
			if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
				continue
			}
			return nil, err
		}
		createdUsers = append(createdUsers, user)
	}

	return createdUsers, nil
}
