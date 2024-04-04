package repositories

import (
	"errors"
	"github.com/hodukihugi/winglets-api/models"
	"gorm.io/gorm"
)

type ProfileRepository interface {
	CreateProfile(input *models.Profile) (*models.Profile, error)
	GetProfileById(id string) (*models.Profile, error)
	GetListProfile() ([]models.Profile, error)
	UpdateProfileById(id string, input *models.Profile) (*models.Profile, error)
	DeleteProfileById(id string) error
}

type profileRepository struct {
	db *gorm.DB
}

func NewProfileRepository(db *gorm.DB) *profileRepository {
	return &profileRepository{db: db}
}

func (r *profileRepository) CreateProfile(input *models.Profile) (*models.Profile, error) {
	var profile models.Profile
	db := r.db.Model(&profile)
	profile.Clone(input)
	if error := db.Create(&profile).Error; error != nil {
		return nil, error
	}

	return &profile, nil
}

func (r *profileRepository) GetProfileById(id string) (*models.Profile, error) {
	var profile models.Profile
	db := r.db.Model(&profile)
	checkProfilesExist := db.Select("*").Where("id = ?", id).Find(&profile)
	if checkProfilesExist.RowsAffected <= 0 {
		return nil, errors.New("profile not found")
	}

	return &profile, nil
}

func (r *profileRepository) GetListProfile() ([]models.Profile, error) {
	var profiles []models.Profile
	db := r.db.Model(&models.Profile{})
	db.Find(&profiles)
	db.Commit()
	return profiles, nil
}

func (r *profileRepository) UpdateProfileById(id string, input *models.Profile) (*models.Profile, error) {
	var profile models.Profile
	db := r.db.Model(&profile)
	profile.Clone(input)

	checkProfilesExist := db.Select("*").Where("id = ?", id).Find(&profile)
	if checkProfilesExist.RowsAffected <= 0 {
		return nil, errors.New("profile not found")
	}

	if err := db.Where("id = ?", id).Updates(&profile).Error; err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *profileRepository) DeleteProfileById(id string) error {
	var profile models.Profile
	db := r.db.Model(&profile)
	checkProfilesExist := db.Select("*").Where("id = ?", id).Find(&profile)
	if checkProfilesExist.RowsAffected <= 0 {
		return errors.New("profile not found")
	}

	db.Delete(&profile)
	return nil
}
