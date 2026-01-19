package repository

import (
	"ramah-disabilitas-be/pkg/database"
	"time"
)

type ActiveClassResult struct {
	ID            uint64
	Name          string
	AvailableTime string // Not really semester, but maybe CreatedAt or a field? Using CreatedAt for now or static.
	StudentCount  int64
	Progress      float64
}

func GetActiveClassesByTeacherID(teacherID uint64) ([]ActiveClassResult, error) {
	var results []ActiveClassResult

	// Get courses
	type CourseData struct {
		ID        uint64
		Title     string
		CreatedAt time.Time
	}
	var courses []CourseData

	err := database.DB.Table("courses").
		Select("id, title, created_at").
		Where("teacher_id = ? AND status = ?", teacherID, "published").
		Scan(&courses).Error
	if err != nil {
		return nil, err
	}

	for _, c := range courses {
		// Get Student Count
		var studentCount int64
		database.DB.Table("course_students").Where("course_id = ?", c.ID).Count(&studentCount)

		// Get Progress (Average completion of students in this course)
		// Logic: (Sum(Completed Items per Student) / (Total Items * Student Count)) * 100

		// Total Items in Course
		var materialCount int64
		database.DB.Table("materials").
			Joins("JOIN modules ON materials.module_id = modules.id").
			Where("modules.course_id = ?", c.ID).
			Count(&materialCount)

		var assignmentCount int64
		database.DB.Table("assignments").Where("course_id = ?", c.ID).Count(&assignmentCount)

		totalItems := materialCount + assignmentCount

		var progress float64 = 0
		if totalItems > 0 && studentCount > 0 {
			// Count completed materials by students in this course
			var completedMaterials int64
			database.DB.Table("material_completions").
				Joins("JOIN materials ON material_completions.material_id = materials.id").
				Joins("JOIN modules ON materials.module_id = modules.id").
				Where("modules.course_id = ? AND material_completions.completed = ?", c.ID, true).
				Count(&completedMaterials)

			// Count submitted assignments by students in this course
			var submittedAssignments int64
			database.DB.Table("submissions").
				Joins("JOIN assignments ON submissions.assignment_id = assignments.id").
				Where("assignments.course_id = ?", c.ID).
				Count(&submittedAssignments)

			totalCompleted := completedMaterials + submittedAssignments
			totalPossible := totalItems * studentCount

			progress = (float64(totalCompleted) / float64(totalPossible)) * 100
		}

		results = append(results, ActiveClassResult{
			ID:            c.ID,
			Name:          c.Title,
			AvailableTime: "Genap 2025", // Hardcoded as per request example/context not having semester
			StudentCount:  studentCount,
			Progress:      progress,
		})
	}

	return results, nil
}

type ActivityResult struct {
	ID          uint64
	Type        string
	Title       string
	Description string
	CourseID    uint64
	CreatedAt   time.Time
}

func GetRecentActivitiesByTeacherID(teacherID uint64, limit int) ([]ActivityResult, error) {
	var results []ActivityResult

	// We only assume assignment_submission for now
	type SubmissionActivity struct {
		SubmissionID    uint64
		StudentName     string
		AssignmentTitle string
		CourseID        uint64
		SubmittedAt     time.Time
	}

	var subs []SubmissionActivity
	err := database.DB.Table("submissions").
		Select("submissions.id as submission_id, users.name as student_name, assignments.title as assignment_title, courses.id as course_id, submissions.submitted_at").
		Joins("JOIN users ON submissions.student_id = users.id").
		Joins("JOIN assignments ON submissions.assignment_id = assignments.id").
		Joins("JOIN courses ON assignments.course_id = courses.id").
		Where("courses.teacher_id = ?", teacherID).
		Order("submissions.submitted_at desc").
		Limit(limit).
		Scan(&subs).Error

	if err != nil {
		return nil, err
	}

	for _, s := range subs {
		results = append(results, ActivityResult{
			ID:          s.SubmissionID,
			Type:        "assignment_submission",
			Title:       "Pengumpulan Tugas: " + s.AssignmentTitle,
			Description: s.StudentName + " mengumpulkan tugas",
			CourseID:    s.CourseID,
			CreatedAt:   s.SubmittedAt,
		})
	}

	return results, nil
}

type PendingAssignmentResult struct {
	AssignmentID   uint64
	Title          string
	CourseName     string
	SubmittedCount int64
}

func GetPendingAssignmentsByTeacherID(teacherID uint64) ([]PendingAssignmentResult, error) {
	var results []PendingAssignmentResult

	rows, err := database.DB.Table("assignments").
		Select("assignments.id, assignments.title, courses.title as course_name, COUNT(submissions.id) as submitted_count").
		Joins("JOIN courses ON assignments.course_id = courses.id").
		Joins("JOIN submissions ON assignments.id = submissions.assignment_id").
		Where("courses.teacher_id = ? AND submissions.grade = 0", teacherID).
		Group("assignments.id, assignments.title, courses.title").
		Having("COUNT(submissions.id) > 0").
		Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var r PendingAssignmentResult
		if err := rows.Scan(&r.AssignmentID, &r.Title, &r.CourseName, &r.SubmittedCount); err != nil {
			continue
		}
		results = append(results, r)
	}

	return results, nil
}

// Reuse existing functions for simpler stats or keeping them compatible

func GetCourseCountByTeacherID(teacherID uint64) (int64, error) {
	var count int64
	err := database.DB.Table("courses").Where("teacher_id = ? AND status = ?", teacherID, "published").Count(&count).Error
	return count, err
}

func GetStudentCountByTeacherID(teacherID uint64) (int64, error) {
	var count int64
	err := database.DB.Table("course_students").
		Joins("JOIN courses ON course_students.course_id = courses.id").
		Where("courses.teacher_id = ? AND courses.status = ?", teacherID, "published").
		Distinct("course_students.user_id").
		Count(&count).Error
	return count, err
}

func GetUngradedAssignmentCountByTeacherID(teacherID uint64) (int64, error) {
	var count int64
	err := database.DB.Table("submissions").
		Joins("JOIN assignments ON submissions.assignment_id = assignments.id").
		Joins("JOIN courses ON assignments.course_id = courses.id").
		Where("courses.teacher_id = ? AND submissions.grade = 0", teacherID).
		Count(&count).Error
	return count, err
}

func GetAverageProgressByTeacherID(teacherID uint64) (float64, error) {
	// Reusing logic but maybe simplified or duplicated for robustness
	activeClasses, err := GetActiveClassesByTeacherID(teacherID)
	if err != nil {
		return 0, err
	}

	if len(activeClasses) == 0 {
		return 0, nil
	}

	var totalProgress float64
	for _, c := range activeClasses {
		totalProgress += c.Progress
	}
	return totalProgress / float64(len(activeClasses)), nil
}
