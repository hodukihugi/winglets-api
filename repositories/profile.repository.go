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

	// Check if a soft-deleted record with the same ID exists
	var existingProfile models.Profile
	err := db.Unscoped().First(&existingProfile, "id = ?", profile.ID).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		r.logger.Debug(err)
		return nil, err
	}

	// If a soft-deleted record exists, restore it
	if existingProfile.ID != "" {
		err = db.Unscoped().Model(&existingProfile).Update("deleted_at", nil).Error
		if err != nil {
			r.logger.Debug(err)
			return nil, err
		}

		// Update the restored record with the new data from the `profile` parameter
		err = db.Model(&existingProfile).Updates(&profile).Error
		if err != nil {
			r.logger.Debug(err)
			return nil, err
		}

		return &existingProfile, nil
	}

	// Create a new record if no soft-deleted record exists
	if err = db.Create(&profile).Error; err != nil {
		return nil, err
	}

	return &profile, nil
}

func (r *ProfileRepository) GetProfileById(id string) (*models.Profile, error) {
	var profile models.Profile
	db := r.Database.Model(&profile)
	err := db.First(&profile, "id = ?", id).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		r.logger.Debug(err)
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

	var existingProfile models.Profile
	err := db.First(&existingProfile, "id = ?", id).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		r.logger.Debug(err)
		return nil, err
	}

	if err = db.Model(&existingProfile).Updates(&profile).Error; err != nil {
		r.logger.Debug(err)
		return nil, err
	}

	return &profile, nil
}

func (r *ProfileRepository) DeleteProfileById(id string) error {
	var profile models.Profile
	db := r.Database.Model(&profile)
	checkProfilesExist := db.Select("*").Where("id = ?", id).Find(&profile)
	if checkProfilesExist.RowsAffected <= 0 {
		r.logger.Debug("profile not found")
		return errors.New("profile not found")
	}

	if err := db.Delete(&profile).Error; err != nil {
		r.logger.Debug("can't delete profile")
		r.logger.Error(err)
		return err
	}
	return nil
}
