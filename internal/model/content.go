package model

import "time"

type Subtest struct {
	ID   uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"type:varchar(100)" json:"name"`
	Slug string `gorm:"type:varchar(100);unique" json:"slug"`
	Icon string `gorm:"type:varchar(255)" json:"icon"`
}

type Material struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	SubtestID uint64    `json:"subtest_id"`
	Title     string    `gorm:"type:varchar(255)" json:"title"`
	Content   string    `gorm:"type:text" json:"content"`
	VideoURL  string    `gorm:"type:varchar(255)" json:"video_url"`
	CreatedAt time.Time `json:"created_at"`
	Subtest   Subtest   `gorm:"foreignKey:SubtestID" json:"subtest"`
}

type QuestionCorrectAnswer string

const (
	AnswerA QuestionCorrectAnswer = "a"
	AnswerB QuestionCorrectAnswer = "b"
	AnswerC QuestionCorrectAnswer = "c"
	AnswerD QuestionCorrectAnswer = "d"
	AnswerE QuestionCorrectAnswer = "e"
)

type QuestionDifficulty string

const (
	DifficultyEasy   QuestionDifficulty = "easy"
	DifficultyMedium QuestionDifficulty = "medium"
	DifficultyHard   QuestionDifficulty = "hard"
)

type Question struct {
	ID            uint64                `gorm:"primaryKey;autoIncrement" json:"id"`
	SubtestID     uint64                `json:"subtest_id"`
	QuestionText  string                `gorm:"type:text" json:"question_text"`
	ImageURL      string                `gorm:"type:varchar(255)" json:"image_url"`
	OptionA       string                `gorm:"type:text" json:"option_a"`
	OptionB       string                `gorm:"type:text" json:"option_b"`
	OptionC       string                `gorm:"type:text" json:"option_c"`
	OptionD       string                `gorm:"type:text" json:"option_d"`
	OptionE       string                `gorm:"type:text" json:"option_e"`
	CorrectAnswer QuestionCorrectAnswer `gorm:"type:varchar(5)" json:"correct_answer"`
	Explanation   string                `gorm:"type:text" json:"explanation"`
	Difficulty    QuestionDifficulty    `gorm:"type:varchar(20)" json:"difficulty"`
	Subtest       Subtest               `gorm:"foreignKey:SubtestID" json:"subtest"`
}
