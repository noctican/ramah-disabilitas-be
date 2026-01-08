package model

import (
	"time"

	"gorm.io/datatypes"
)

type Course struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	TeacherID   uint64    `json:"teacher_id"`
	Title       string    `gorm:"type:varchar(255)" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	Thumbnail   string    `gorm:"type:varchar(255)" json:"thumbnail"`
	ClassCode   string    `gorm:"uniqueIndex;type:varchar(20)" json:"class_code"`
	Status      string    `gorm:"type:varchar(20);default:'draft'" json:"status"` // published, draft, archived
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Modules     []Module     `gorm:"foreignKey:CourseID" json:"modules,omitempty"`
	Assignments []Assignment `gorm:"foreignKey:CourseID" json:"assignments,omitempty"`
	Students    []User       `gorm:"many2many:course_students;" json:"students,omitempty"`
}

type Module struct {
	ID       uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	CourseID uint64 `json:"course_id"`
	Title    string `gorm:"type:varchar(255)" json:"title"`
	Order    int    `json:"order"`

	Materials []Material `gorm:"foreignKey:ModuleID" json:"materials,omitempty"`
}

type MaterialType string

const (
	TypeYoutube MaterialType = "youtube"
	TypePDF     MaterialType = "pdf"
	TypeText    MaterialType = "text"
)

type Material struct {
	ID       uint64       `gorm:"primaryKey;autoIncrement" json:"id"`
	ModuleID uint64       `json:"module_id"`
	Title    string       `gorm:"type:varchar(255)" json:"title"`
	Type     MaterialType `gorm:"type:varchar(20)" json:"type"`

	SourceURL  string `gorm:"type:text" json:"source_url"`
	RawContent string `gorm:"type:text" json:"raw_content,omitempty"`

	DurationMin int  `json:"duration_min"`
	HasCaptions bool `json:"has_captions"`

	SmartFeature *SmartFeature `gorm:"foreignKey:MaterialID" json:"smart_feature,omitempty"`
}

type SmartFeature struct {
	ID         uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	MaterialID uint64 `gorm:"uniqueIndex" json:"material_id"`

	Summary string `gorm:"type:text" json:"summary"`

	Simplified string `gorm:"type:text" json:"simplified_content"`

	QuizData datatypes.JSON `json:"quiz_data"`

	IsGenerated bool `gorm:"default:false" json:"is_generated"`
}
