package repository

import (
	"ramah-disabilitas-be/internal/model"
	"ramah-disabilitas-be/pkg/database"
)

func CreateAssignment(assignment *model.Assignment) error {
	return database.DB.Create(assignment).Error
}

func GetAssignmentsByCourseID(courseID uint64) ([]model.Assignment, error) {
	var assignments []model.Assignment
	err := database.DB.Where("course_id = ?", courseID).Find(&assignments).Error
	return assignments, err
}

func GetAssignmentByID(id uint64) (*model.Assignment, error) {
	var assignment model.Assignment
	err := database.DB.Preload("Submissions").First(&assignment, id).Error
	return &assignment, err
}

func GetAssignmentsByStudentID(studentID uint64) ([]model.Assignment, error) {
	var assignments []model.Assignment
	// Join with course_students to find courses this student joined
	// Then get assignments from those courses
	err := database.DB.Table("assignments").
		Joins("JOIN courses ON assignments.course_id = courses.id").
		Joins("JOIN course_students ON courses.id = course_students.course_id").
		Where("course_students.user_id = ?", studentID).
		Order("assignments.id DESC"). // Terbaru based on ID (creation order)
		Find(&assignments).Error
	return assignments, err
}
