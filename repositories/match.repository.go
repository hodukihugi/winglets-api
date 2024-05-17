package repositories

import (
	"errors"
	"github.com/hodukihugi/winglets-api/core"
	"github.com/hodukihugi/winglets-api/models"
)

type IMatchRepository interface {
	First(string, string) (*models.Match, error)
	Create(models.Match) error
	Match(models.Match) error
	Pass(models.Match) error
}

type MatchRepository struct {
	*core.Database
	logger *core.Logger
}

func NewMatchRepository(db *core.Database, logger *core.Logger) IMatchRepository {
	return &MatchRepository{
		Database: db,
		logger:   logger,
	}
}

func (r *MatchRepository) First(matcherId, matcheeId string) (*models.Match, error) {
	var match models.Match
	db := r.Database.Model(models.Match{})
	if matcherId == "" || matcheeId == "" {
		return nil, errors.New("matcher id or matchee is empty")
	}

	if err := db.Where("matcher_id = ? AND matchee_id = ?", matcherId, matcheeId).First(&match).Error; err != nil {
		return nil, err
	}

	return &match, nil
}

func (r *MatchRepository) Create(match models.Match) error {
	db := r.Database.Model(models.Match{})
	if err := db.Create(&match).Error; err != nil {
		return err
	}

	return nil
}

func (r *MatchRepository) Match(match models.Match) error {
	db := r.Database.Model(models.Match{})

	// Kiểm tra xem đã có match record chưa, nếu chưa có thì trở về
	if err := db.Where("matcher_id = ? AND matchee_id = ?", match.MatcherId, match.MatcheeId).First(&match).Error; err != nil {
		return err
	}

	// Nếu đã có thì update match status lên 1
	if err := db.Where("matcher_id = ? AND matchee_id = ?", match.MatcherId, match.MatcheeId).
		Update("match_status", 1).Error; err != nil {
		return err
	}
	return nil

}

func (r *MatchRepository) Pass(match models.Match) error {
	return errors.New("not implemented")
}
