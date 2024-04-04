package models

import (
	"gorm.io/gorm"
	"time"
)

// ---------------- DAO ----------------

// User model
type User struct {
	gorm.Model
	Name     string
	Email    string
	Password string
	Birthday time.Time
	//Gender    string
	//Height    int
	//Horoscope string
	//Hobby     string
	//Language  string
	//Education string
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
		ID:       u.ID,
		Name:     u.Name,
		Email:    u.Email,
		Birthday: u.Birthday,
		//Gender:    u.Gender,
		//Height:    u.Height,
		//Horoscope: u.Horoscope,
		//Hobby:     u.Hobby,
		//Language:  u.Language,
		//Education: u.Education,
	}
}

type SerializableUser struct {
	ID       uint      `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Birthday time.Time `json:"birthday"`
	//Gender    string    `json:"gender"`
	//Height    int       `json:"height"`
	//Horoscope string    `json:"horoscope"`
	//Hobby     string    `json:"hobby"`
	//Language  string    `json:"language"`
	//Education string    `json:"education"`
}

//func (u *User) SerializeNested() *SerializableNestedUser {
//	return &SerializableNestedUser{
//		User: u.Serialize(),
//	}
//}
//
//type SerializableNestedUser struct {
//	User *SerializableUser `json:"user,omitempty"`
//}

type OneUserFilter struct {
	ID     uint                 `form:"id"`
	Email  string               `form:"email"`
	Joins  *ArrStringFilterType `form:"joins"`
	Fields *ArrStringFilterType `form:"fields"`
}
