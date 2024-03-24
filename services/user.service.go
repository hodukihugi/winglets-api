package services

import (
	"github.com/hodukihugi/winglets-api/models"
	"github.com/hodukihugi/winglets-api/repositories"
)

type UserService interface {
	CreateUser(*models.CreateUserEntity) string
	GetUserById(string) (*models.UserModel, string)
	GetListUser() ([]models.UserModel, string)
	UpdateUserById(string, *models.UpdateUserEntity) string
	DeleteUserById(string) string
}

type userService struct {
	repository repositories.UserRepository
}

func NewUserService(repository repositories.UserRepository) *userService {
	return &userService{repository: repository}
}

func (s *userService) CreateUser(input *models.CreateUserEntity) string {
	user := models.UserModel{
		Name:        input.Name,
		DateOfBirth: input.DateOfBirth,
		Gender:      input.Gender,
	}

	_, errCreateUser := s.repository.CreateUser(&user)
	return errCreateUser
}

func (s *userService) GetUserById(id string) (*models.UserModel, string) {
	resultGetUserById, errGetUser := s.repository.GetUserById(id)
	return resultGetUserById, errGetUser
}

func (s *userService) GetListUser() ([]models.UserModel, string) {
	resultGetListUser, errGetUser := s.repository.GetListUser()
	return resultGetListUser, errGetUser
}

func (s *userService) UpdateUserById(id string, input *models.UpdateUserEntity) string {
	user := models.UserModel{
		Name:        input.Name,
		DateOfBirth: input.DateOfBirth,
		Gender:      input.Gender,
		Height:      input.Height,
		Horoscope:   input.Horoscope,
		Hobby:       input.Hobby,
		Language:    input.Language,
		Education:   input.Education,
	}

	_, errUpdateUser := s.repository.UpdateUserById(id, &user)
	return errUpdateUser
}

func (s *userService) DeleteUserById(id string) string {
	errUpdateUser := s.repository.DeleteUserById(id)
	return errUpdateUser
}
