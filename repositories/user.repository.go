package repositories

import (
	"github.com/hodukihugi/winglets-api/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(input *models.UserModel) (*models.UserModel, string)
	GetUserById(id string) (*models.UserModel, string)
	GetListUser() ([]models.UserModel, string)
	UpdateUserById(id string, input *models.UserModel) (*models.UserModel, string)
	DeleteUserById(id string) string
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(input *models.UserModel) (*models.UserModel, string) {
	var user models.UserModel
	db := r.db.Model(&user)
	user.Clone(input)
	if error := db.Create(&user).Error; error != nil {
		return &user, "DB_ERROR_CREATE_USER_FAILED"
	}

	return &user, "nil"
}

func (r *userRepository) GetUserById(id string) (*models.UserModel, string) {
	var user models.UserModel
	db := r.db.Model(&user)
	checkUsersExist := db.Select("*").Where("id = ?", id).Find(&user)
	db.Commit()
	if checkUsersExist.RowsAffected <= 0 {
		return &user, "USER_NOT_FOUND_404"
	}

	return &user, "nil"
}

func (r *userRepository) GetListUser() ([]models.UserModel, string) {
	var users []models.UserModel
	db := r.db.Model(&models.UserModel{})
	db.Find(&users)
	db.Commit()
	return users, "nil"
}

func (r *userRepository) UpdateUserById(id string, input *models.UserModel) (*models.UserModel, string) {
	var user models.UserModel
	db := r.db.Model(&user)
	user.Clone(input)
	if err := db.Where("id = ?", id).Updates(&user).Error; err != nil {
		return &user, "DB_ERROR_UPDATE_USER_FAILED"
	}
	return &user, "nil"
}

func (r *userRepository) DeleteUserById(id string) string {
	var user models.UserModel
	db := r.db.Model(&user)
	checkUsersExist := db.Select("*").Where("id = ?", id).Find(&user)
	if checkUsersExist.RowsAffected <= 0 {
		return "USER_NOT_FOUND_404"
	}

	db.Delete(&user)
	db.Commit()
	return "nil"
}
