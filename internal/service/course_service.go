package service

import (
	"context"
	"errors"
	"math/rand"
	"ramah-disabilitas-be/internal/model"
	"ramah-disabilitas-be/internal/repository"
	"ramah-disabilitas-be/pkg/ai"
	"ramah-disabilitas-be/pkg/utils"
	"time"
)

type MaterialInput struct {
	ID          uint64             `json:"id"`
	Title       string             `json:"title" binding:"required"`
	Type        model.MaterialType `json:"type" binding:"required"`
	SourceURL   string             `json:"source_url"`
	RawContent  string             `json:"raw_content"`
	DurationMin int                `json:"duration_min"`
	HasCaptions bool               `json:"has_captions"`
}

type ModuleInput struct {
	ID        uint64          `json:"id"`
	Title     string          `json:"title" binding:"required"`
	Order     int             `json:"order"`
	Materials []MaterialInput `json:"materials,omitempty"`
}

type CourseInput struct {
	Title       string        `json:"title" binding:"required"`
	Description string        `json:"description"`
	Thumbnail   string        `json:"thumbnail"`
	ClassCode   string        `json:"class_code"`
	Status      string        `json:"status"` // published, draft
	Modules     []ModuleInput `json:"modules,omitempty"`
}

func CreateCourse(input CourseInput, teacherID uint64) (*model.Course, error) {
	if input.ClassCode == "" {
		input.ClassCode = generateClassCode()
	}

	status := "draft"
	if input.Status != "" {
		status = input.Status
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
		Status:      status,
		Modules:     modules,
	}

	if err := repository.CreateCourse(course); err != nil {
		return nil, err
	}

	return course, nil
}

func GetCoursesByTeacher(teacherID uint64, search string, status string, sort string) ([]model.Course, error) {
	return repository.GetCoursesByTeacherID(teacherID, search, status, sort)
}

func GetCourseDetail(id uint64) (*model.Course, error) {
	return repository.GetCourseByID(id)
}

func UpdateCourse(id uint64, input CourseInput, teacherID uint64) (*model.Course, error) {
	existingCourse, err := repository.GetCourseByID(id)
	if err != nil {
		return nil, err
	}

	if existingCourse.TeacherID != teacherID {
		return nil, errors.New("unauthorized: you do not own this course")
	}

	// Prepare the new state
	course := &model.Course{
		ID:          id,
		TeacherID:   teacherID,
		Title:       input.Title,
		Description: input.Description,
		Thumbnail:   existingCourse.Thumbnail, // Default to existing
		ClassCode:   existingCourse.ClassCode, // Default to existing
		Status:      existingCourse.Status,    // Default to existing
		CreatedAt:   existingCourse.CreatedAt,
	}

	if input.Thumbnail != "" {
		course.Thumbnail = input.Thumbnail
	}
	if input.ClassCode != "" {
		course.ClassCode = input.ClassCode
	}
	if input.Status != "" {
		course.Status = input.Status
	}

	// Map Modules from Input
	var modules []model.Module
	for _, m := range input.Modules {
		var materials []model.Material
		for _, mat := range m.Materials {
			materials = append(materials, model.Material{
				ID:          mat.ID,
				ModuleID:    m.ID, // Will be 0 if new module
				Title:       mat.Title,
				Type:        mat.Type,
				SourceURL:   mat.SourceURL,
				RawContent:  mat.RawContent,
				DurationMin: mat.DurationMin,
				HasCaptions: mat.HasCaptions,
			})
		}
		modules = append(modules, model.Module{
			ID:        m.ID,
			CourseID:  id,
			Title:     m.Title,
			Order:     m.Order,
			Materials: materials,
		})
	}
	course.Modules = modules

	if err := repository.UpdateCourse(course); err != nil {
		return nil, err
	}

	// Return updated course
	updatedCourse, err := repository.GetCourseByID(id)
	return updatedCourse, err
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

func DeleteModule(moduleID uint64, teacherID uint64) error {
	module, err := repository.GetModuleByID(moduleID)
	if err != nil {
		return errors.New("modul tidak ditemukan")
	}

	course, err := repository.GetCourseByID(module.CourseID)
	if err != nil {
		return errors.New("kelas tidak ditemukan")
	}

	if course.TeacherID != teacherID {
		return errors.New("unauthorized: anda tidak memiliki akses ke modul ini")
	}

	return repository.DeleteModule(moduleID)
}

func DeleteMaterial(materialID uint64, teacherID uint64) error {
	material, err := repository.GetMaterialByID(materialID)
	if err != nil {
		return errors.New("materi tidak ditemukan")
	}

	module, err := repository.GetModuleByID(material.ModuleID)
	if err != nil {
		return errors.New("modul tidak ditemukan (data inkonsisten)")
	}

	course, err := repository.GetCourseByID(module.CourseID)
	if err != nil {
		return errors.New("kelas tidak ditemukan (data inkonsisten)")
	}

	if course.TeacherID != teacherID {
		return errors.New("unauthorized: anda tidak memiliki akses ke materi ini")
	}

	return repository.DeleteMaterial(materialID)
}

func CreateMaterial(moduleID uint64, input MaterialInput, teacherID uint64) (*model.Material, error) {
	module, err := repository.GetModuleByID(moduleID)
	if err != nil {
		return nil, errors.New("modul tidak ditemukan")
	}

	course, err := repository.GetCourseByID(module.CourseID)
	if err != nil {
		return nil, errors.New("kelas tidak ditemukan")
	}

	if course.TeacherID != teacherID {
		return nil, errors.New("unauthorized: anda tidak memiliki akses ke modul ini")
	}

	material := &model.Material{
		ModuleID:    moduleID,
		Title:       input.Title,
		Type:        input.Type,
		SourceURL:   input.SourceURL,
		RawContent:  input.RawContent,
		DurationMin: input.DurationMin,
		HasCaptions: input.HasCaptions,
	}

	if err := repository.CreateMaterial(material); err != nil {
		return nil, err
	}

	return material, nil
}

func UpdateMaterial(materialID uint64, input MaterialInput, teacherID uint64) (*model.Material, error) {
	material, err := repository.GetMaterialByID(materialID)
	if err != nil {
		return nil, errors.New("materi tidak ditemukan")
	}

	module, err := repository.GetModuleByID(material.ModuleID)
	if err != nil {
		return nil, errors.New("modul tidak ditemukan")
	}

	course, err := repository.GetCourseByID(module.CourseID)
	if err != nil {
		return nil, errors.New("kelas tidak ditemukan")
	}

	if course.TeacherID != teacherID {
		return nil, errors.New("unauthorized: anda tidak memiliki akses ke materi ini")
	}

	material.Title = input.Title
	material.Type = input.Type

	// Only update source url if not empty (or if we want to allow clearing it, we need logic. Assuming update meant to set new value)
	if input.SourceURL != "" {
		material.SourceURL = input.SourceURL
	}
	material.RawContent = input.RawContent
	material.DurationMin = input.DurationMin
	material.HasCaptions = input.HasCaptions

	if err := repository.UpdateMaterial(material); err != nil {
		return nil, err
	}

	return material, nil
}

func ToggleMaterialCompletion(userID, materialID uint64) (bool, error) {
	return repository.ToggleMaterialCompletion(userID, materialID)
}

func GetStudentCourseDetail(courseID, studentID uint64) (*model.Course, error) {
	// 1. Check if student is enrolled
	inCourse, err := repository.IsStudentInCourse(courseID, studentID)
	if err != nil {
		return nil, err
	}
	if !inCourse {
		return nil, errors.New("unauthorized: anda belum bergabung di kelas ini")
	}

	// 2. Get Course Detail
	course, err := repository.GetCourseByID(courseID)
	if err != nil {
		return nil, errors.New("kelas tidak ditemukan")
	}

	// 3. Get User's Completed Materials
	completedMap, err := repository.GetCompletedMaterialsMap(courseID, studentID)
	if err != nil {
		return nil, err
	}

	// 4. Map completion status
	for i := range course.Modules {
		for j := range course.Modules[i].Materials {
			if completedMap[course.Modules[i].Materials[j].ID] {
				course.Modules[i].Materials[j].IsCompleted = true
			}
		}
	}

	return course, nil
}

func GetMaterialDetailWithStatus(materialID, userID uint64) (*model.Material, error) {
	material, err := repository.GetMaterialByID(materialID)
	if err != nil {
		return nil, errors.New("materi tidak ditemukan")
	}

	module, err := repository.GetModuleByID(material.ModuleID)
	if err != nil {
		return nil, errors.New("module not found")
	}

	// Verify enrollment (or ownership)
	// If teacher is the same as userID?
	// But usually this view is for students.
	// For teachers, they can view it too?
	// Let's assume Student context for now since "WithStatus" implies student progress.
	// We check enrollment.
	inCourse, err := repository.IsStudentInCourse(module.CourseID, userID)
	if err != nil {
		return nil, err
	}

	// If not student, check if teacher
	if !inCourse {
		course, err := repository.GetCourseByID(module.CourseID)
		if err == nil && course.TeacherID == userID {
			// Is teacher, allowed. Status logic might be different or irrelevant (always false or true?)
			// Let's just return material without IsCompleted (default false)
			return material, nil
		}
		// Not teacher either
		return nil, errors.New("unauthorized: anda tidak memiliki akses ke materi ini")
	}

	// Is Student, check completion
	isCompleted := repository.GetMaterialCompletionStatus(userID, materialID)
	material.IsCompleted = isCompleted

	return material, nil
}

func CreateStudentAndEnroll(courseID uint64, input CreateStudentInput, teacherID uint64) (*model.User, error) {
	// 1. Verify Teacher owns the course
	course, err := repository.GetCourseByID(courseID)
	if err != nil {
		return nil, errors.New("kelas tidak ditemukan")
	}
	if course.TeacherID != teacherID {
		return nil, errors.New("unauthorized: anda tidak memiliki akses ke kelas ini")
	}

	// 2. Create Student Account
	user, err := CreateStudent(input)
	if err != nil {
		return nil, err
	}

	// 3. Enroll Student to Course
	if err := repository.AddStudentToCourse(courseID, user.ID); err != nil {
		// Ideally we should rollback user creation here if transaction support was passed down
		// But for now let's just error out. The user exists but is not enrolled.
		// Retrying enroll might be manual.
		return nil, errors.New("user created but failed to enroll: " + err.Error())
	}

	return user, nil
}

func GenerateMaterialSummary(materialID uint64) (*model.SmartFeature, error) {
	material, err := repository.GetMaterialByID(materialID)
	if err != nil {
		return nil, errors.New("materi tidak ditemukan")
	}

	// Check if summary already exists (simple caching)
	// You might want to allow re-generation, but strictly sticking to 'generate if needed' for cost efficiency first.
	// If the user wants to regenerate, we can clear this field or add a force flag.
	if material.SmartFeature != nil && material.SmartFeature.Summary != "" {
		return material.SmartFeature, nil
	}

	var textContent string
	if material.Type == model.TypePDF {
		// Extract from source URL
		// Ensure SourceURL is not empty
		if material.SourceURL == "" {
			return nil, errors.New("file PDF tidak ditemukan (URL kosong)")
		}

		extracted, err := utils.ExtractTextFromPDF(material.SourceURL)
		if err != nil {
			return nil, errors.New("gagal membaca PDF: " + err.Error())
		}
		if len(extracted) > 200000 {
			// Truncate if too huge? Gemini Flash has 1M token context, so 200k chars is fine (~50k tokens).
			// But let's be safe against abuse.
			extracted = extracted[:200000]
		}
		textContent = extracted
	} else if material.Type == model.TypeText {
		textContent = material.RawContent
	} else {
		return nil, errors.New("tipe materi ini belum didukung untuk ringkasan otomatis")
	}

	if textContent == "" {
		return nil, errors.New("konten materi kosong, tidak bisa diringkas")
	}

	// Call AI
	// Using a context with timeout is good practice
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Prompt refined for rich formatting
	prompt := "Buatkan ringkasan yang komprehensif dari materi berikut. Gunakan format Markdown untuk struktur yang rapi: gunakan **bold** untuk istilah penting atau heading, list bullet points untuk rincian, dan _italic_ untuk penekanan. Pastikan ringkasan mudah dipahami oleh mahasiswa:\n\n" + textContent
	summary, err := ai.GenerateContent(ctx, prompt)
	if err != nil {
		return nil, errors.New("gagal menghasilkan ringkasan AI: " + err.Error())
	}

	// Return ephemeral result (not saved yet)
	return &model.SmartFeature{
		MaterialID:  materialID,
		Summary:     summary,
		IsGenerated: true, // Marked as AI generated
	}, nil
}

func SaveMaterialSummary(materialID uint64, summary string) (*model.SmartFeature, error) {
	material, err := repository.GetMaterialByID(materialID)
	if err != nil {
		return nil, errors.New("materi tidak ditemukan")
	}

	smartFeature := material.SmartFeature
	if smartFeature == nil {
		smartFeature = &model.SmartFeature{
			MaterialID: materialID,
		}
	}

	smartFeature.Summary = summary
	smartFeature.IsGenerated = true

	if err := repository.SaveSmartFeature(smartFeature); err != nil {
		return nil, errors.New("gagal menyimpan ringkasan")
	}

	return smartFeature, nil
}
