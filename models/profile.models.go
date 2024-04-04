package models

import "gorm.io/gorm"

// ---------- DAO ----------------

// profile model
type Profile struct {
	gorm.Model
	Gender    string
	Height    int
	Horoscope string
	Hobby     string
	Language  string
	Education string
}

// TableName gives table name of model
func (p *Profile) TableName() string {
	return "profile"
}

func (p *Profile) Clone(clone *Profile) {
	p.ID = clone.ID
	p.Gender = clone.Gender
	p.Height = clone.Height
	p.Hobby = clone.Hobby
	p.Education = clone.Education
	p.Language = clone.Language
}

// ---------- DTO ----------------
func (p *Profile) Serialize() *SerializableProfile {
	if p == nil {
		return nil
	}
	return &SerializableProfile{
		ID:        p.ID,
		Gender:    p.Gender,
		Height:    p.Height,
		Horoscope: p.Horoscope,
		Hobby:     p.Hobby,
		Language:  p.Language,
		Education: p.Education,
	}
}

type SerializableProfile struct {
	ID        uint   `json:"id"`
	Gender    string `json:"gender"`
	Height    int    `json:"height"`
	Horoscope string `json:"horoscope"`
	Hobby     string `json:"hobby"`
	Language  string `json:"language"`
	Education string `json:"education"`
}
