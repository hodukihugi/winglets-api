package services

import (
	"github.com/hodukihugi/winglets-api/core"
	"github.com/hodukihugi/winglets-api/models"
	"github.com/hodukihugi/winglets-api/repositories"
	"time"
)

type IProfileService interface {
	CreateProfile(string, models.ProfileCreateRequest) error
	GetProfileById(string) (*models.Profile, error)
	UpdateProfileById(string, models.ProfileUpdateRequest) error
	DeleteProfileById(string) error
}

type ProfileService struct {
	repository repositories.IProfileRepository
	logger     *core.Logger
}

func NewProfileService(repository repositories.IProfileRepository, logger *core.Logger) IProfileService {
	return &ProfileService{
		repository: repository,
		logger:     logger,
	}
}

func (s *ProfileService) CreateProfile(userID string, request models.ProfileCreateRequest) error {
	_, err := s.repository.CreateProfile(models.Profile{
		ID:        userID,
		Name:      request.Name,
		Gender:    request.Gender,
		Birthday:  time.Unix(request.BirthdayInSeconds, 0).UTC(),
		Height:    request.Height,
		Horoscope: request.Horoscope,
		Hobby:     request.Hobby,
		Language:  request.Language,
		Education: request.Education,
	})

	return err
}

func (s *ProfileService) GetProfileById(id string) (*models.Profile, error) {
	result, err := s.repository.GetProfileById(id)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *ProfileService) UpdateProfileById(id string, request models.ProfileUpdateRequest) error {
	_, err := s.repository.UpdateProfileById(id, models.Profile{
		Name:      request.Name,
		Gender:    request.Gender,
		Birthday:  time.Unix(request.BirthdayInSeconds, 0).UTC(),
		Height:    request.Height,
		Horoscope: request.Horoscope,
		Hobby:     request.Hobby,
		Language:  request.Language,
		Education: request.Education,
	})

	return err
}
func (s *ProfileService) DeleteProfileById(id string) error {
	err := s.repository.DeleteProfileById(id)
	return err
}
