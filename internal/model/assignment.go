package model

import (
	"time"
)

type Assignment struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	CourseID    uint64    `json:"course_id"`
	ModuleID    *uint64   `json:"module_id"` // Bisa link ke modul tertentu
	Title       string    `gorm:"type:varchar(255)" json:"title"`
	Instruction string    `gorm:"type:text" json:"instruction"`
	Deadline    time.Time `json:"deadline"`

	MaxPoints int `json:"max_points"`

	AllowText  bool `json:"allow_text"`
	AllowFile  bool `json:"allow_file"`
	AllowVoice bool `json:"allow_voice"`
	AllowLate  bool `json:"allow_late"`

	Submissions []Submission `gorm:"foreignKey:AssignmentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"submissions,omitempty"`

	MySubmission *Submission `gorm:"-" json:"my_submission,omitempty"`
}

type Submission struct {
	ID           uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	AssignmentID uint64 `json:"assignment_id"`
	StudentID    uint64 `json:"student_id"`

	TextAnswer   string `gorm:"type:text" json:"text_answer"`
	FileURL      string `gorm:"type:text" json:"file_url"`
	VoiceNoteURL string `gorm:"type:text" json:"voice_note_url"`

	Grade       float64   `json:"grade"`
	Feedback    string    `gorm:"type:text" json:"feedback"`
	SubmittedAt time.Time `json:"submitted_at"`

	Student User `gorm:"foreignKey:StudentID" json:"student,omitempty"`
}
