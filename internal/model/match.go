package model

import "time"

type MatchStatus string

const (
	MatchOngoing  MatchStatus = "ongoing"
	MatchFinished MatchStatus = "finished"
	MatchAborted  MatchStatus = "aborted"
)

type Match struct {
	ID        uint64      `gorm:"primaryKey;autoIncrement" json:"id"`
	Player1ID uint64      `json:"player_1_id"`
	Player2ID *uint64     `json:"player_2_id"`
	WinnerID  *uint64     `json:"winner_id"`
	SubtestID uint64      `json:"subtest_id"`
	Status    MatchStatus `gorm:"type:varchar(20)" json:"status"`
	CreatedAt time.Time   `json:"created_at"`
	Player1   User        `gorm:"foreignKey:Player1ID" json:"player_1"`
	Player2   *User       `gorm:"foreignKey:Player2ID" json:"player_2"`
	Winner    *User       `gorm:"foreignKey:WinnerID" json:"winner"`
	Subtest   Subtest     `gorm:"foreignKey:SubtestID" json:"subtest"`
}

type MatchDetail struct {
	ID            uint64                 `gorm:"primaryKey;autoIncrement" json:"id"`
	MatchID       uint64                 `json:"match_id"`
	QuestionID    uint64                 `json:"question_id"`
	Player1Answer *QuestionCorrectAnswer `gorm:"type:varchar(5)" json:"player_1_answer"`
	Player2Answer *QuestionCorrectAnswer `gorm:"type:varchar(5)" json:"player_2_answer"`
	IsCorrectP1   bool                   `json:"is_correct_p1"`
	IsCorrectP2   bool                   `json:"is_correct_p2"`
	Match         Match                  `gorm:"foreignKey:MatchID" json:"match"`
	Question      Question               `gorm:"foreignKey:QuestionID" json:"question"`
}
