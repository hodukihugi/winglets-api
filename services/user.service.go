package services

import (
	"github.com/hodukihugi/winglets-api/core"
	"github.com/hodukihugi/winglets-api/models"
	"github.com/hodukihugi/winglets-api/repositories"
)

// IUserService interface
type IUserService interface {
	First(models.OneUserFilter) (*models.User, error)
}

// UserService service layer
type UserService struct {
	logger     *core.Logger
	repository repositories.IUserRepository
}

func (u *UserService) First(filter models.OneUserFilter) (user *models.User, err error) {
	return u.repository.First(filter)
}

// NewUserService creates a new userservice
func NewUserService(logger *core.Logger, repository repositories.IUserRepository) IUserService {
	return &UserService{
		logger:     logger,
		repository: repository,
	}
}
