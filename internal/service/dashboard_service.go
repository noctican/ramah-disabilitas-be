package service

import (
	"ramah-disabilitas-be/internal/repository"
	"strconv"
	"time"
)

type DashboardSummaryStats struct {
	TotalClasses      int64   `json:"total_classes"`
	TotalStudents     int64   `json:"total_students"`
	PendingGrades     int64   `json:"pending_grades"`
	AverageCompletion float64 `json:"average_completion"`
}

func GetDashboardSummary(teacherID uint64) (*DashboardSummaryStats, error) {
	courseCount, err := repository.GetCourseCountByTeacherID(teacherID)
	if err != nil {
		return nil, err
	}

	studentCount, err := repository.GetStudentCountByTeacherID(teacherID)
	if err != nil {
		return nil, err
	}

	ungradedCount, err := repository.GetUngradedAssignmentCountByTeacherID(teacherID)
	if err != nil {
		return nil, err
	}

	avgProgress, err := repository.GetAverageProgressByTeacherID(teacherID)
	if err != nil {
		return nil, err
	}

	return &DashboardSummaryStats{
		TotalClasses:      courseCount,
		TotalStudents:     studentCount,
		PendingGrades:     ungradedCount,
		AverageCompletion: avgProgress,
	}, nil
}

type ActiveClassResponse struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Semester     string  `json:"semester"`
	StudentCount int64   `json:"student_count"`
	Progress     float64 `json:"progress"`
}

func GetActiveClasses(teacherID uint64) ([]ActiveClassResponse, error) {
	classes, err := repository.GetActiveClassesByTeacherID(teacherID)
	if err != nil {
		return nil, err
	}

	var response []ActiveClassResponse
	for _, c := range classes {
		response = append(response, ActiveClassResponse{
			ID:           "cls_" + strconv.FormatUint(c.ID, 10), // Adding prefix as per example
			Name:         c.Name,
			Semester:     c.AvailableTime,
			StudentCount: c.StudentCount,
			Progress:     c.Progress,
		})
	}
	// Return empty array instead of null
	if response == nil {
		response = []ActiveClassResponse{}
	}
	return response, nil
}

type ActivityResponse struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CourseID    string    `json:"course_id"`
	CreatedAt   time.Time `json:"created_at"`
}

func GetRecentActivities(teacherID uint64, limit int) ([]ActivityResponse, error) {
	activities, err := repository.GetRecentActivitiesByTeacherID(teacherID, limit)
	if err != nil {
		return nil, err
	}

	var response []ActivityResponse
	for _, a := range activities {
		response = append(response, ActivityResponse{
			ID:          "act_" + strconv.FormatUint(a.ID, 10),
			Type:        a.Type,
			Title:       a.Title,
			Description: a.Description,
			CourseID:    "cls_" + strconv.FormatUint(a.CourseID, 10),
			CreatedAt:   a.CreatedAt,
		})
	}
	if response == nil {
		response = []ActivityResponse{}
	}
	return response, nil
}

type PendingAssignmentResponse struct {
	AssignmentID   string `json:"assignment_id"`
	Title          string `json:"title"`
	Course         string `json:"course"`
	SubmittedCount int64  `json:"submitted_count"`
}

func GetPendingAssignments(teacherID uint64) ([]PendingAssignmentResponse, error) {
	assignments, err := repository.GetPendingAssignmentsByTeacherID(teacherID)
	if err != nil {
		return nil, err
	}

	var response []PendingAssignmentResponse
	for _, a := range assignments {
		response = append(response, PendingAssignmentResponse{
			AssignmentID:   "asg_" + strconv.FormatUint(a.AssignmentID, 10),
			Title:          a.Title,
			Course:         a.CourseName,
			SubmittedCount: a.SubmittedCount,
		})
	}
	if response == nil {
		response = []PendingAssignmentResponse{}
	}
	return response, nil
}

type ClassProgressResponse struct {
	ClassID    string  `json:"class_id"`
	Completion float64 `json:"completion"`
}

type ProgressSummaryResponse struct {
	AverageCompletion float64                 `json:"average_completion"`
	ByClass           []ClassProgressResponse `json:"by_class"`
}

func GetProgressSummary(teacherID uint64) (*ProgressSummaryResponse, error) {
	avg, err := repository.GetAverageProgressByTeacherID(teacherID)
	if err != nil {
		return nil, err
	}

	classes, err := repository.GetActiveClassesByTeacherID(teacherID)
	if err != nil {
		return nil, err
	}

	var byClass []ClassProgressResponse
	for _, c := range classes {
		byClass = append(byClass, ClassProgressResponse{
			ClassID:    "cls_" + strconv.FormatUint(c.ID, 10),
			Completion: c.Progress,
		})
	}
	if byClass == nil {
		byClass = []ClassProgressResponse{}
	}

	return &ProgressSummaryResponse{
		AverageCompletion: avg,
		ByClass:           byClass,
	}, nil
}

// Keep the old function for backward compatibility if used elsewhere (though not used in updated handler below)
func GetLecturerDashboardStats(teacherID uint64) (*DashboardSummaryStats, error) {
	return GetDashboardSummary(teacherID)
}
