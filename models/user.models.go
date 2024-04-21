package models

import (
	"gorm.io/gorm"
	"time"
)

// ---------------- DAO ----------------

// User model
type User struct {
	gorm.Model
	ID                 string    `gorm:"primaryKey"`
	Email              string    `gorm:"column:email"`
	Password           string    `gorm:"column:password"`
	VerificationCode   string    `gorm:"column:verification_code"`
	VerificationStatus int       `gorm:"column:verification_status"`
	VerificationTime   time.Time `gorm:"column:verification_time"`
}

// TableName gives table name of model
func (u *User) TableName() string {
	return "users"
}

// ---------------- DTO ----------------

func (u *User) Serialize() *SerializableUser {
	if u == nil {
		return nil
	}
	return &SerializableUser{
		ID:    u.ID,
		Email: u.Email,
	}
}

type SerializableUser struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

func (u *User) SerializeNested() *SerializableNestedUser {
	return &SerializableNestedUser{
		User: u.Serialize(),
	}
}

type SerializableNestedUser struct {
	User *SerializableUser `json:"user,omitempty"`
}

type OneUserFilter struct {
	ID     string               `form:"id"`
	Email  string               `form:"email"`
	Joins  *ArrStringFilterType `form:"joins"`
	Fields *ArrStringFilterType `form:"fields"`
}

type UserUpdateRequest struct {
	VerificationCode   string `json:"verification_code"`
	VerificationStatus int    `json:"verification_status"`
}
