package service

import (
	"errors"
	"math/rand"
	"ramah-disabilitas-be/internal/model"
	"ramah-disabilitas-be/internal/repository"
	"time"
)

type MaterialInput struct {
	Title       string             `json:"title" binding:"required"`
	Type        model.MaterialType `json:"type" binding:"required"`
	SourceURL   string             `json:"source_url"`
	RawContent  string             `json:"raw_content"`
	DurationMin int                `json:"duration_min"`
	HasCaptions bool               `json:"has_captions"`
}

type ModuleInput struct {
	Title     string          `json:"title" binding:"required"`
	Order     int             `json:"order"`
	Materials []MaterialInput `json:"materials,omitempty"`
}

type CourseInput struct {
	Title       string        `json:"title" binding:"required"`
	Description string        `json:"description"`
	Thumbnail   string        `json:"thumbnail"`
	ClassCode   string        `json:"class_code"`
	Modules     []ModuleInput `json:"modules,omitempty"`
}

func CreateCourse(input CourseInput, teacherID uint64) (*model.Course, error) {
	if input.ClassCode == "" {
		input.ClassCode = generateClassCode()
	}

	var modules []model.Module
	for _, m := range input.Modules {
		var materials []model.Material
		for _, mat := range m.Materials {
			materials = append(materials, model.Material{
				Title:       mat.Title,
				Type:        mat.Type,
				SourceURL:   mat.SourceURL,
				RawContent:  mat.RawContent,
				DurationMin: mat.DurationMin,
				HasCaptions: mat.HasCaptions,
			})
		}
		modules = append(modules, model.Module{
			Title:     m.Title,
			Order:     m.Order,
			Materials: materials,
		})
	}

	course := &model.Course{
		TeacherID:   teacherID,
		Title:       input.Title,
		Description: input.Description,
		Thumbnail:   input.Thumbnail,
		ClassCode:   input.ClassCode,
		Modules:     modules,
	}

	if err := repository.CreateCourse(course); err != nil {
		return nil, err
	}

	return course, nil
}

func GetCoursesByTeacher(teacherID uint64) ([]model.Course, error) {
	return repository.GetCoursesByTeacherID(teacherID)
}

func UpdateCourse(id uint64, input CourseInput, teacherID uint64) (*model.Course, error) {
	course, err := repository.GetCourseByID(id)
	if err != nil {
		return nil, err
	}

	if course.TeacherID != teacherID {
		return nil, errors.New("unauthorized: you do not own this course")
	}

	course.Title = input.Title
	course.Description = input.Description
	course.Thumbnail = input.Thumbnail
	if input.ClassCode != "" {
		course.ClassCode = input.ClassCode
	}

	if err := repository.UpdateCourse(course); err != nil {
		return nil, err
	}

	return course, nil
}

func DeleteCourse(id uint64, teacherID uint64) error {
	course, err := repository.GetCourseByID(id)
	if err != nil {
		return err
	}

	if course.TeacherID != teacherID {
		return errors.New("unauthorized: you do not own this course")
	}

	return repository.DeleteCourse(id)
}

func generateClassCode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length := 6

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
}

func JoinCourse(classCode string, studentID uint64) error {
	course, err := repository.GetCourseByClassCode(classCode)
	if err != nil {
		return errors.New("kelas tidak ditemukan")
	}

	exists, err := repository.IsStudentInCourse(course.ID, studentID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("anda sudah bergabung di kelas ini")
	}

	if course.TeacherID == studentID {
		return errors.New("anda adalah pengajar di kelas ini")
	}

	return repository.AddStudentToCourse(course.ID, studentID)
}

func GetStudentCourses(studentID uint64) ([]model.Course, error) {
	return repository.GetCoursesByStudentID(studentID)
}
