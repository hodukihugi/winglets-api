package models

import "gorm.io/gorm"

// ---------- DAO ----------------

// RecommendationBin model
type RecommendationBin struct {
	gorm.Model
	UserID            string `gorm:"primaryKey;column:user_id"`
	RecommendedUserID string `gorm:"primaryKey;column:recommended_user_id"`
}

// TableName gives table name of model
func (p *RecommendationBin) TableName() string {
	return "recommendation_bins"
}
