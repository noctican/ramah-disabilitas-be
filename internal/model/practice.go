package model

import "time"

type PracticeSession struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       uint64    `json:"user_id"`
	SubtestID    uint64    `json:"subtest_id"`
	Score        int       `json:"score"`
	AiEvaluation string    `gorm:"type:text" json:"ai_evaluation"`
	CreatedAt    time.Time `json:"created_at"`
	User         User      `gorm:"foreignKey:UserID" json:"user"`
	Subtest      Subtest   `gorm:"foreignKey:SubtestID" json:"subtest"`
}

type ReportStatus string

const (
	ReportPending  ReportStatus = "pending"
	ReportResolved ReportStatus = "resolved"
	ReportIgnored  ReportStatus = "ignored"
)

type QuestionReport struct {
	ID         uint64       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     uint64       `json:"user_id"`
	QuestionID uint64       `json:"question_id"`
	Reason     string       `gorm:"type:text" json:"reason"`
	Status     ReportStatus `gorm:"type:varchar(20)" json:"status"`
	User       User         `gorm:"foreignKey:UserID" json:"user"`
	Question   Question     `gorm:"foreignKey:QuestionID" json:"question"`
}
