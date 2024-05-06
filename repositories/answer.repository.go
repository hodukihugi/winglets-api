package repositories

import (
	"errors"
	"github.com/hodukihugi/winglets-api/core"
	"github.com/hodukihugi/winglets-api/models"
)

type IAnswerRepository interface {
	CreateAnswer(models.Answer) (*models.Answer, error)
	FindAnswer(models.Answer) (*models.Answer, error)
	FindListAnswerByUserId(userId string) ([]models.Answer, error)
	DeleteAnswer(models.Answer) error
}

type AnswerRepository struct {
	*core.Database
	logger *core.Logger
}

func NewAnswerRepository(database *core.Database, logger *core.Logger) IAnswerRepository {
	return &AnswerRepository{
		Database: database,
		logger:   logger,
	}
}

func (r *AnswerRepository) CreateAnswer(answer models.Answer) (*models.Answer, error) {
	db := r.Database.Model(&answer)
	if err := db.Create(&answer).Error; err != nil {
		return nil, err
	}
	return &answer, nil
}

func (r *AnswerRepository) DeleteAnswer(answer models.Answer) error {
	db := r.Database.Model(&models.Answer{})

	if err := db.Unscoped().Delete(&models.Answer{},
		"user_id = ? AND question_id = ?",
		answer.UserID,
		answer.QuestionID,
	).Error; err != nil {
		return err
	}
	return nil
}

func (r *AnswerRepository) FindAnswer(filter models.Answer) (*models.Answer, error) {

	if filter.UserID == "" || filter.QuestionID == 0 {
		return nil, errors.New("you need to specify a userID or questionID")
	}

	var answer models.Answer

	db := r.Database.Model(&models.Answer{})
	if err := db.
		Where("user_id = ? AND question_id = ?", filter.UserID, filter.QuestionID).
		First(&answer).Error; err != nil {
		return nil, err
	}
	return &answer, nil
}

func (r *AnswerRepository) FindListAnswerByUserId(userId string) ([]models.Answer, error) {
	if userId == "" {
		return nil, errors.New("you need to specify a userId")
	}

	var result []models.Answer
	db := r.Database.Model(&models.Answer{})
	if err := db.Where("user_id = ?", userId).Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}
