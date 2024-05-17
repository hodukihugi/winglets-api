package repositories

import (
	"github.com/hodukihugi/winglets-api/core"
	"github.com/hodukihugi/winglets-api/models"
)

type IRecommendationBinRepository interface {
	GetRecommendedUserByUserId(string) ([]models.RecommendationBin, error)
	Create(models.RecommendationBin) error
}

// RecommendationBinRepository database structure
type RecommendationBinRepository struct {
	*core.Database
	logger *core.Logger
}

// NewRecommendationBinRepository creates a new user repository
func NewRecommendationBinRepository(db *core.Database, logger *core.Logger) IRecommendationBinRepository {
	return &RecommendationBinRepository{
		Database: db,
		logger:   logger,
	}
}

func (r *RecommendationBinRepository) GetRecommendedUserByUserId(userId string) ([]models.RecommendationBin, error) {
	var result []models.RecommendationBin
	db := r.Database.Model(models.RecommendationBin{})

	if err := db.Where("user_id = ?", userId).Find(&result).Error; err != nil {
		r.logger.Error(err)
		return nil, err
	}

	return result, nil
}

func (r *RecommendationBinRepository) Create(recommendedUser models.RecommendationBin) error {
	db := r.Database.Model(models.RecommendationBin{})
	if err := db.Create(&recommendedUser).Error; err != nil {
		r.logger.Error(err)
		return err
	}
	return nil
}
