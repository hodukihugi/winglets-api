package models

import (
	"gorm.io/gorm"
	"time"
)

// ---------- DAO ----------------

// Profile model
type Profile struct {
	gorm.Model
	ID        string    `gorm:"primaryKey"`
	Name      string    `gorm:"column:name"`
	Gender    string    `gorm:"column:gender"`
	Birthday  time.Time `gorm:"column:birthday"`
	Height    int       `gorm:"column:height"`
	Horoscope string    `gorm:"column:horoscope"`
	Hobby     string    `gorm:"column:hobby"`
	Language  string    `gorm:"column:language"`
	Education string    `gorm:"column:education"`
}

// TableName gives table name of model
func (p *Profile) TableName() string {
	return "profiles"
}

// ---------- DTO ----------------

func (p *Profile) Serialize() *SerializableProfile {
	if p == nil {
		return nil
	}
	return &SerializableProfile{
		ID:        p.ID,
		Name:      p.Name,
		Gender:    p.Gender,
		Birthday:  p.Birthday,
		Height:    p.Height,
		Horoscope: p.Horoscope,
		Hobby:     p.Hobby,
		Language:  p.Language,
		Education: p.Education,
	}
}

type SerializableProfile struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Gender    string    `json:"gender"`
	Birthday  time.Time `json:"birthday"`
	Height    int       `json:"height"`
	Horoscope string    `json:"horoscope"`
	Hobby     string    `json:"hobby"`
	Language  string    `json:"language"`
	Education string    `json:"education"`
}

type ProfileCreateRequest struct {
	Name              string `json:"name"`
	Gender            string `json:"gender"`
	BirthdayInSeconds int64  `json:"birthday_in_seconds" validate:"required"`
	Height            int    `json:"height"`
	Horoscope         string `json:"horoscope"`
	Hobby             string `json:"hobby"`
	Language          string `json:"language"`
	Education         string `json:"education"`
}

type ProfileUpdateRequest struct {
	Name              string `json:"name"`
	Gender            string `json:"gender"`
	BirthdayInSeconds int64  `json:"birthday_in_seconds" validate:"required"`
	Height            int    `json:"height"`
	Horoscope         string `json:"horoscope"`
	Hobby             string `json:"hobby"`
	Language          string `json:"language"`
	Education         string `json:"education"`
}
