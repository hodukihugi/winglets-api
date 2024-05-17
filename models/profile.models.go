package models

import (
	"gorm.io/gorm"
	"strings"
	"time"
)

// ---------- DAO ----------------

// Profile model
type Profile struct {
	gorm.Model
	ID          string    `gorm:"primaryKey;column:id"`
	Name        string    `gorm:"column:name"`
	Gender      string    `gorm:"column:gender"`
	Birthday    time.Time `gorm:"column:birthday"`
	Height      int       `gorm:"column:height"`
	Horoscope   string    `gorm:"column:horoscope"`
	Hobby       string    `gorm:"column:hobby"`
	Language    string    `gorm:"column:language"`
	Education   string    `gorm:"column:education"`
	HomeTown    string    `gorm:"column:home_town"`
	Coordinates string    `gorm:"column:coordinates"`
	ImageId1    string    `gorm:"column:image_id_1"`
	ImageId2    string    `gorm:"column:image_id_2"`
	ImageId3    string    `gorm:"column:image_id_3"`
	ImageId4    string    `gorm:"column:image_id_4"`
	ImageId5    string    `gorm:"column:image_id_5"`
	ImageUrl1   string    `gorm:"column:image_url_1"`
	ImageUrl2   string    `gorm:"column:image_url_2"`
	ImageUrl3   string    `gorm:"column:image_url_3"`
	ImageUrl4   string    `gorm:"column:image_url_4"`
	ImageUrl5   string    `gorm:"column:image_url_5"`
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
		Hobby:     strings.Split(p.Hobby, ","),
		Language:  strings.Split(p.Language, ","),
		Education: p.Education,
		HomeTown:  p.HomeTown,
		Image1:    p.ImageUrl1,
		Image2:    p.ImageUrl2,
		Image3:    p.ImageUrl3,
		Image4:    p.ImageUrl4,
		Image5:    p.ImageUrl5,
	}
}

func (p *Profile) ConvertToMatchProfile() *MatchProfile {
	if p == nil {
		return nil
	}

	return &MatchProfile{
		ID:        p.ID,
		Name:      p.Name,
		Gender:    p.Gender,
		Birthday:  p.Birthday,
		Height:    p.Height,
		Horoscope: p.Horoscope,
		Hobby:     strings.Split(p.Hobby, ","),
		Language:  strings.Split(p.Language, ","),
		Education: p.Education,
		HomeTown:  p.HomeTown,
		Image1:    p.ImageUrl1,
		Image2:    p.ImageUrl2,
		Image3:    p.ImageUrl3,
		Image4:    p.ImageUrl4,
		Image5:    p.ImageUrl5,
	}
}

type SerializableProfile struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Gender    string    `json:"gender"`
	Birthday  time.Time `json:"birthday"`
	Height    int       `json:"height"`
	Horoscope string    `json:"horoscope"`
	Hobby     []string  `json:"hobby"`
	Language  []string  `json:"language"`
	Education string    `json:"education"`
	Location  string    `json:"location"`
	HomeTown  string    `json:"home_town"`
	Image1    string    `json:"image_1"`
	Image2    string    `json:"image_2"`
	Image3    string    `json:"image_3"`
	Image4    string    `json:"image_4"`
	Image5    string    `json:"image_5"`
}

type MatchProfile struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Gender          string    `json:"gender"`
	Birthday        time.Time `json:"birthday"`
	Height          int       `json:"height"`
	Horoscope       string    `json:"horoscope"`
	Hobby           []string  `json:"hobby"`
	Language        []string  `json:"language"`
	Education       string    `json:"education"`
	Location        string    `json:"location"`
	HomeTown        string    `json:"home_town"`
	Distance        float64   `json:"distance"`
	MatchPercentage float64   `json:"match_percentage"`
	Image1          string    `json:"image_1"`
	Image2          string    `json:"image_2"`
	Image3          string    `json:"image_3"`
	Image4          string    `json:"image_4"`
	Image5          string    `json:"image_5"`
}

type ProfileCreateRequest struct {
	Name              string   `json:"name" validate:"required"`
	Gender            string   `json:"gender" validate:"oneof=male female,required"`
	BirthdayInSeconds int64    `json:"birthday_in_seconds" validate:"required"`
	Height            int      `json:"height"`
	Horoscope         string   `json:"horoscope"`
	Hobby             []string `json:"hobby"`
	Language          []string `json:"language"`
	Education         string   `json:"education"`
	Location          string   `json:"location"`
	HomeTown          string   `json:"home_town"`
	Coordinates       struct {
		Longitude float64 `json:"longitude"`
		Latitude  float64 `json:"latitude"`
	} `json:"coordinates"`
}

type ProfileUpdateRequest struct {
	Name              string   `json:"name"`
	Gender            string   `json:"gender" validate:"oneof=male female"`
	BirthdayInSeconds int64    `json:"birthday_in_seconds"`
	Height            int      `json:"height"`
	Horoscope         string   `json:"horoscope"`
	Hobby             []string `json:"hobby"`
	Language          []string `json:"language"`
	Education         string   `json:"education"`
	Location          string   `json:"location"`
	HomeTown          string   `json:"home_town"`
	Coordinates       struct {
		Longitude float64 `json:"longitude"`
		Latitude  float64 `json:"latitude"`
	} `json:"coordinates"`
	ImageId1  string `json:"image_id_1"`
	ImageId2  string `json:"image_id_2"`
	ImageId3  string `json:"image_id_3"`
	ImageId4  string `json:"image_id_4"`
	ImageId5  string `json:"image_id_5"`
	ImageUrl1 string `json:"image_url_1"`
	ImageUrl2 string `json:"image_url_2"`
	ImageUrl3 string `json:"image_url_3"`
	ImageUrl4 string `json:"image_url_4"`
	ImageUrl5 string `json:"image_url_5"`
}

type ImageUploadResult struct {
	SlotId  int
	FileId  string
	FileUrl string
}

type ProfileImageDeleteRequest struct {
	Slots []int `json:"slots"`
}

type ProfileFilter struct {
	ExcludedUserId string
	Gender         string
	MinAge         int
	MaxAge         int
	MinDistance    float64
	MaxDistance    float64
	Longitude      float64
	Latitude       float64
}
