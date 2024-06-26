package services

import (
	"fmt"
	"github.com/hodukihugi/winglets-api/core"
	"github.com/hodukihugi/winglets-api/models"
	"github.com/hodukihugi/winglets-api/repositories"
	"strings"
	"time"
)

type IProfileService interface {
	CreateProfile(string, models.ProfileCreateRequest) error
	GetProfileById(string) (*models.Profile, error)
	UpdateProfileById(string, models.ProfileUpdateRequest) error
	UpdateProfileImageById(string, []int) error
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
	var coordinates = []string{
		fmt.Sprintf("%.6f", request.Coordinates.Longitude),
		fmt.Sprintf("%.6f", request.Coordinates.Latitude),
	}

	_, err := s.repository.CreateProfile(models.Profile{
		ID:          userID,
		Name:        request.Name,
		Gender:      request.Gender,
		Birthday:    time.Unix(request.BirthdayInSeconds, 0).UTC(),
		Height:      request.Height,
		Horoscope:   request.Horoscope,
		Hobby:       strings.Join(request.Hobby, ","),
		Language:    strings.Join(request.Language, ","),
		Education:   request.Education,
		HomeTown:    request.HomeTown,
		Coordinates: strings.Join(coordinates, ","),
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
	var coordinates = []string{
		fmt.Sprintf("%.6f", request.Coordinates.Longitude),
		fmt.Sprintf("%.6f", request.Coordinates.Latitude),
	}

	_, err := s.repository.UpdateProfileById(id, models.Profile{
		Name:        request.Name,
		Gender:      request.Gender,
		Birthday:    time.Unix(request.BirthdayInSeconds, 0).UTC(),
		Height:      request.Height,
		Horoscope:   request.Horoscope,
		Hobby:       strings.Join(request.Hobby, ","),
		Language:    strings.Join(request.Language, ","),
		Education:   request.Education,
		HomeTown:    request.HomeTown,
		Coordinates: strings.Join(coordinates, ","),
		ImageId1:    request.ImageId1,
		ImageId2:    request.ImageId2,
		ImageId3:    request.ImageId3,
		ImageId4:    request.ImageId4,
		ImageId5:    request.ImageId5,
		ImageUrl1:   request.ImageUrl1,
		ImageUrl2:   request.ImageUrl2,
		ImageUrl3:   request.ImageUrl3,
		ImageUrl4:   request.ImageUrl4,
		ImageUrl5:   request.ImageUrl5,
	})

	return err
}

func (s *ProfileService) UpdateProfileImageById(id string, slots []int) error {
	_, err := s.repository.UpdateProfileImageById(id, slots)
	return err
}
func (s *ProfileService) DeleteProfileById(id string) error {
	err := s.repository.DeleteProfileById(id)
	return err
}
