package repositories

import (
	"errors"
	"github.com/hodukihugi/winglets-api/core"
	"github.com/hodukihugi/winglets-api/models"
	"gorm.io/gorm"
)

type IProfileRepository interface {
	CreateProfile(models.Profile) (*models.Profile, error)
	GetProfileById(string) (*models.Profile, error)
	GetListProfile() ([]models.Profile, error)
	UpdateProfileById(string, models.Profile) (*models.Profile, error)
	DeleteProfileById(string) error
}

type ProfileRepository struct {
	*core.Database
	logger *core.Logger
}

func NewProfileRepository(db *core.Database, logger *core.Logger) IProfileRepository {
	return &ProfileRepository{
		Database: db,
		logger:   logger,
	}
}

func (r *ProfileRepository) CreateProfile(profile models.Profile) (*models.Profile, error) {
	db := r.Database.Model(&profile)

	if error := db.Create(&profile).Error; error != nil {
		return nil, error
	}

	return &profile, nil
}

func (r *ProfileRepository) GetProfileById(id string) (*models.Profile, error) {
	var profile models.Profile
	db := r.Database.Model(&profile)
	err := db.First(&profile, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.Info("profile not found")
			return nil, errors.New("profile not found")
		}
		return nil, err
	}
	return &profile, nil
}

func (r *ProfileRepository) GetListProfile() ([]models.Profile, error) {
	var profiles []models.Profile
	db := r.Database.Model(&models.Profile{})
	db.Find(&profiles)
	return profiles, nil
}

func (r *ProfileRepository) UpdateProfileById(id string, profile models.Profile) (*models.Profile, error) {
	db := r.Database.Model(&profile)

	checkProfilesExist := db.Select("*").Where("id = ?", id).Find(&profile)
	if checkProfilesExist.RowsAffected <= 0 {
		return nil, errors.New("profile not found")
	}

	if err := db.Where("id = ?", id).Updates(&profile).Error; err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *ProfileRepository) DeleteProfileById(id string) error {
	var profile models.Profile
	db := r.Database.Model(&profile)
	checkProfilesExist := db.Select("*").Where("id = ?", id).Find(&profile)
	if checkProfilesExist.RowsAffected <= 0 {
		r.logger.Info("profile not found")
		return errors.New("profile not found")
	}

	if err := db.Delete(&profile).Error; err != nil {
		r.logger.Info("can't delete profile")
		r.logger.Error(err)
		return err
	}
	return nil
}
