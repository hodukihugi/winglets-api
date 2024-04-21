package services

import (
	"github.com/hodukihugi/winglets-api/core"
	"github.com/hodukihugi/winglets-api/models"
	"github.com/hodukihugi/winglets-api/repositories"
)

// IUserService interface
type IUserService interface {
	First(models.OneUserFilter) (*models.User, error)
	UpdateById(string, models.UserUpdateRequest) error
}

// UserService service layer
type UserService struct {
	logger     *core.Logger
	repository repositories.IUserRepository
}

// NewUserService creates a new user service
func NewUserService(logger *core.Logger, repository repositories.IUserRepository) IUserService {
	return &UserService{
		logger:     logger,
		repository: repository,
	}
}

func (u *UserService) First(filter models.OneUserFilter) (user *models.User, err error) {
	return u.repository.First(filter)
}

func (u *UserService) UpdateById(id string, request models.UserUpdateRequest) error {
	return u.repository.UpdateById(id, models.User{
		VerificationCode:   request.VerificationCode,
		VerificationStatus: request.VerificationStatus,
	})
}
