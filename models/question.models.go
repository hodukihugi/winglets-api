package models

import (
	"gorm.io/gorm"
	"strings"
)

// ==================== DAO ==============

// Question model
type Question struct {
	gorm.Model
	QuestionID      int    `gorm:"primaryKey;column:question_id"`
	QuestionContent string `gorm:"column:content"`
	QuestionAnswer  string `gorm:"column:answers"`
}

func (q *Question) TableName() string {
	return "questions"
}

// ==================== DTO ==============

func (q *Question) Serialize() *SerializableQuestion {
	return &SerializableQuestion{
		QuestionID:      q.QuestionID,
		QuestionContent: q.QuestionContent,
		QuestionAnswer:  strings.Split(q.QuestionAnswer, ","),
	}
}

type SerializableQuestion struct {
	QuestionID      int      `json:"question_id"`
	QuestionContent string   `json:"question_content"`
	QuestionAnswer  []string `json:"question_answer"`
}
