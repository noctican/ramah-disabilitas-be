package service

import (
	"errors"
	"ramah-disabilitas-be/internal/model"
	"ramah-disabilitas-be/internal/repository"
	"time"
)

type AssignmentInput struct {
	Title       string    `json:"title" form:"title" binding:"required"`
	Instruction string    `json:"instruction" form:"instruction" binding:"required"`
	ModuleID    *uint64   `json:"module_id" form:"module_id"`
	MaxPoints   int       `json:"max_points" form:"max_points" binding:"required"`
	Deadline    time.Time `json:"deadline" form:"deadline" binding:"required" time_format:"2006-01-02T15:04:05Z07:00"`
	AllowFile   bool      `json:"allow_file" form:"allow_file"`
	AllowText   bool      `json:"allow_text" form:"allow_text"`
	AllowLate   bool      `json:"allow_late" form:"allow_late"`
}

func CreateAssignment(courseID uint64, input AssignmentInput, teacherID uint64) (*model.Assignment, error) {
	// Verify Course Ownership
	course, err := repository.GetCourseByID(courseID)
	if err != nil {
		return nil, errors.New("kelas tidak ditemukan")
	}

	if course.TeacherID != teacherID {
		return nil, errors.New("unauthorized: anda tidak memiliki akses ke kelas ini")
	}

	// Verify Module if provided
	if input.ModuleID != nil {
		module, err := repository.GetModuleByID(*input.ModuleID)
		if err != nil {
			return nil, errors.New("modul tidak ditemukan")
		}
		if module.CourseID != courseID {
			return nil, errors.New("modul tidak valid untuk kelas ini")
		}
	}

	assignment := &model.Assignment{
		CourseID:    courseID,
		ModuleID:    input.ModuleID,
		Title:       input.Title,
		Instruction: input.Instruction,
		MaxPoints:   input.MaxPoints,
		Deadline:    input.Deadline,
		AllowFile:   input.AllowFile,
		AllowText:   input.AllowText,
		AllowLate:   input.AllowLate,
		// Defaulting AllowVoice based on text/file logic or keeping it false for now as not in UI
		AllowVoice: false,
	}

	if err := repository.CreateAssignment(assignment); err != nil {
		return nil, err
	}

	return assignment, nil
}

func GetAssignmentsByCourse(courseID uint64, teacherID uint64) ([]model.Assignment, error) {
	course, err := repository.GetCourseByID(courseID)
	if err != nil {
		return nil, errors.New("kelas tidak ditemukan")
	}

	if course.TeacherID != teacherID {
		return nil, errors.New("unauthorized: anda tidak memiliki akses ke kelas ini")
	}

	return repository.GetAssignmentsByCourseID(courseID)
}

func GetStudentAssignments(studentID uint64) ([]model.Assignment, error) {
	return repository.GetAssignmentsByStudentID(studentID)
}
