package models

import (
	"gorm.io/gorm"
	"time"
)

// ============= DAO ================

type Match struct {
	MatcherId   string `gorm:"primaryKey;column:matcher_id"`
	MatcheeId   string `gorm:"primaryKey;column:matchee_id"`
	MatchStatus int    `gorm:"column:match_status"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (m *Match) TableName() string {
	return "matches"
}

// ============= DTO ================

type MatchCalculationResult struct {
	MatchPercentage float64
	MatchedProfile  Profile
}

type GetRecommendationRequest struct {
	MinAge      int     `json:"min_age"`
	MaxAge      int     `json:"max_age"`
	MinDistance float64 `json:"min_distance"`
	MaxDistance float64 `json:"max_distance"`
}

type SmashRequest struct {
	UserId string `json:"user_id"`
}

type PassRequest struct {
	UserId string `json:"user_id"`
}
