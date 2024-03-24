package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// Data Model

type UserModel struct {
	ID          string    `json:"id" gorm:"column:id"`
	Name        string    `json:"name" gorm:"column:name"`
	DateOfBirth string    `json:"date_of_birth" gorm:"column:date_of_birth"`
	Gender      string    `json:"gender" gorm:"column:gender"`
	Height      *int      `json:"height" gorm:"column:height"`
	Horoscope   *string   `json:"horoscope" gorm:"column:horoscope"`
	Hobby       *string   `json:"hobby" gorm:"column:hobby"`
	Language    *string   `json:"language" gorm:"column:language"`
	Education   *string   `json:"education" gorm:"column:education"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (UserModel) TableName() string {
	return "users"
}

func (entity *UserModel) BeforeCreate(*gorm.DB) error {
	entity.ID = uuid.New().String()
	entity.CreatedAt = time.Now().Local()
	return nil
}

func (entity *UserModel) BeforeUpdate(*gorm.DB) error {
	entity.UpdatedAt = time.Now().Local()
	return nil
}

func (entity *UserModel) Clone(clone *UserModel) {
	entity.Name = clone.Name
	entity.DateOfBirth = clone.DateOfBirth
	entity.Gender = clone.Gender
	entity.Height = clone.Height
	entity.Horoscope = clone.Horoscope
	entity.Hobby = clone.Hobby
	entity.Language = clone.Language
	entity.Education = clone.Education
}

// CreateUserEntity

type CreateUserEntity struct {
	Name        string `json:"name" gorm:"column:name"`
	DateOfBirth string `json:"date_of_birth" gorm:"column:date_of_birth"`
	Gender      string `json:"gender" gorm:"column:gender"`
}

func (CreateUserEntity) TableName() string {
	return UserModel{}.TableName()
}

// UpdateUserEntity

type UpdateUserEntity struct {
	Name        string  `json:"name" gorm:"column:name"`
	DateOfBirth string  `json:"date_of_birth" gorm:"column:date_of_birth"`
	Gender      string  `json:"gender" gorm:"column:gender"`
	Height      *int    `json:"height" gorm:"column:height"`
	Horoscope   *string `json:"horoscope" gorm:"column:horoscope"`
	Hobby       *string `json:"hobby" gorm:"column:hobby"`
	Language    *string `json:"language" gorm:"column:language"`
	Education   *string `json:"education" gorm:"column:education"`
}

func (UpdateUserEntity) TableName() string {
	return UserModel{}.TableName()
}
