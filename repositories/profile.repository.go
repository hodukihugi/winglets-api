package repositories

import (
	"errors"
	"fmt"
	"github.com/hodukihugi/winglets-api/core"
	"github.com/hodukihugi/winglets-api/models"
	"github.com/hodukihugi/winglets-api/utils"
	"gorm.io/gorm"
	"time"
)

type IProfileRepository interface {
	CreateProfile(models.Profile) (*models.Profile, error)
	GetProfileById(string) (*models.Profile, error)
	GetListProfile(models.ProfileFilter) ([]models.Profile, error)
	UpdateProfileById(string, models.Profile) (*models.Profile, error)
	UpdateProfileImageById(string, []int) (*models.Profile, error)
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
	db := r.Database.Model(&models.Profile{})

	// Check if a record with the same ID exists
	var existingProfile models.Profile

	// Check if a record with the same ID exists, including soft-deleted records
	if err := db.Unscoped().First(&existingProfile, "id = ?", profile.ID).Error; err == nil {
		// If a soft-deleted record exists, restore it
		err = db.Model(&existingProfile).Update("deleted_at", nil).Error
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
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		r.logger.Debug(err)
		return nil, err
	}

	// Create a new record if no record exists
	db = r.Database.Model(&models.Profile{})
	r.logger.Info("Creating new profile")
	if err := db.Create(&profile).Error; err != nil {
		r.logger.Debug(err)
		return nil, err
	}
	return &profile, nil
}

func (r *ProfileRepository) GetProfileById(id string) (*models.Profile, error) {
	db := r.Database.Model(&models.Profile{})
	var profile models.Profile
	err := db.First(&profile, "id = ?", id).Error
	if err != nil {
		r.logger.Debug("Profile not found")
		return nil, err
	}
	return &profile, nil
}

func (r *ProfileRepository) GetListProfile(filter models.ProfileFilter) ([]models.Profile, error) {
	if filter.MinAge > filter.MaxAge {
		return nil, errors.New("min age is higher than max age")
	}

	if filter.MinDistance > filter.MaxDistance {
		return nil, errors.New("min distance is higher than max distance")
	}

	db := r.Database.Model(&models.Profile{})
	var profiles, results []models.Profile
	r.logger.Info(fmt.Sprintf("Filter: %+v", filter))
	if filter.MinAge > 0 && filter.MinDistance >= 0 && filter.Longitude != 0 && filter.Latitude != 0 {
		minimum := time.Now().AddDate(-filter.MaxAge, 0, 0).UTC()
		maximum := time.Now().AddDate(-filter.MinAge, 0, 0).UTC()
		r.logger.Debugf("Min birthday: %v, Max birthday: %v", minimum, maximum)
		db.
			Where("gender = ? "+
				"AND birthday >= ? AND birthday <= ? "+
				"AND id NOT IN (SELECT recommended_user_id FROM recommendation_bins WHERE user_id = ?) "+
				"AND id <> ?",
				filter.Gender, minimum, maximum, filter.ExcludedUserId, filter.ExcludedUserId).
			Limit(20).
			Find(&profiles)
		for _, profile := range profiles {
			longitude, latitude, err := utils.CoordinatesStringToPairFloat64(profile.Coordinates)

			if err != nil {
				return nil, err
			}

			distance := utils.CalculateDistance(filter.Longitude, filter.Latitude, longitude, latitude)
			if distance >= filter.MinDistance && distance <= filter.MaxDistance {
				results = append(results, profile)
			}
		}

	} else {
		r.logger.Debug("Finding all profiles")
		db.Find(&profiles)
		results = profiles
	}

	return results, nil
}

func (r *ProfileRepository) UpdateProfileById(id string, profile models.Profile) (*models.Profile, error) {
	db := r.Database.Model(&models.Profile{})

	var existingProfile models.Profile
	err := db.First(&existingProfile, "id = ?", id).Error
	if err != nil {
		r.logger.Debug(err)
		return nil, err
	}

	if err = db.Model(&existingProfile).Updates(&profile).Error; err != nil {
		r.logger.Debug(err)
		return nil, err
	}

	return &profile, nil
}

func (r *ProfileRepository) UpdateProfileImageById(id string, slots []int) (*models.Profile, error) {
	db := r.Database.Model(&models.Profile{})

	var existingProfile models.Profile
	err := db.First(&existingProfile, "id = ?", id).Error
	if err != nil {
		r.logger.Debug(err)
		return nil, err
	}

	updateMap := make(map[string]interface{})

	for _, slotId := range slots {
		updateMap[fmt.Sprintf("ImageId%d", slotId)] = nil
		updateMap[fmt.Sprintf("ImageUrl%d", slotId)] = nil
	}
	if err = db.Model(&existingProfile).Updates(updateMap).Error; err != nil {
		r.logger.Debug(err)
		return nil, err
	}

	return &existingProfile, nil
}

func (r *ProfileRepository) DeleteProfileById(id string) error {
	db := r.Database.Model(&models.Profile{})
	var profile models.Profile

	var existingProfile models.Profile
	err := db.First(&existingProfile, "id = ?", id).Error
	if err != nil {
		r.logger.Debug("Not found profile")
		return err
	}

	if err := db.Delete(&profile).Error; err != nil {
		r.logger.Debug("can't delete profile")
		r.logger.Error(err)
		return err
	}
	return nil
}
