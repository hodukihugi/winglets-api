package models

import (
	"gorm.io/gorm"
)

// ---------- DAO ----------------

// Answer model
type Answer struct {
	gorm.Model
	UserID       string `gorm:"primaryKey;column:user_id"`
	QuestionID   int    `gorm:"primaryKey;column:question_id"`
	UserAnswer   int    `gorm:"column:user_answer"`
	PreferAnswer int    `gorm:"column:prefer_answer"`
	Importance   int    `gorm:"column:importance"`
}

// TableName gives table name of model
func (p *Answer) TableName() string {
	return "answers"
}

// ---------- DTO ----------------

func (p *Answer) Serialize() *SerializableAnswer {
	if p == nil {
		return nil
	}
	return &SerializableAnswer{
		UserID:       p.UserID,
		QuestionID:   p.QuestionID,
		UserAnswer:   p.UserAnswer,
		PreferAnswer: p.PreferAnswer,
		Importance:   p.Importance,
	}
}

type SerializableAnswer struct {
	UserID       string `json:"user_id"`
	QuestionID   int    `json:"question_id"`
	UserAnswer   int    `json:"user_answer"`
	PreferAnswer int    `json:"prefer_answer"`
	Importance   int    `json:"importance"`
}

type AnswerRequest struct {
	Answers []struct {
		QuestionID   int `json:"question_id"`
		UserAnswer   int `json:"user_answer"`
		PreferAnswer int `json:"prefer_answer"`
		Importance   int `json:"importance"`
	} `json:"answers"`
}
