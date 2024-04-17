package models

import (
	"gorm.io/gorm"
)

// ---------------- DAO ----------------

// User model
type User struct {
	gorm.Model
	ID       string `gorm:"primaryKey"`
	Email    string `gorm:"column:email"`
	Password string `gorm:"column:password"`
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
