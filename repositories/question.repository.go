package repositories

import (
	"github.com/hodukihugi/winglets-api/core"
	"github.com/hodukihugi/winglets-api/models"
)

type IQuestionRepository interface {
	GetListQuestions() ([]*models.Question, error)
}

// QuestionRepository database structure
type QuestionRepository struct {
	*core.Database
	logger *core.Logger
}

// NewQuestionRepository creates a new user repository
func NewQuestionRepository(db *core.Database, logger *core.Logger) IQuestionRepository {
	return &QuestionRepository{
		Database: db,
		logger:   logger,
	}
}

func (r *QuestionRepository) GetListQuestions() ([]*models.Question, error) {
	var questions []*models.Question
	db := r.Database.Model(&models.Question{})
	db.Find(&questions)
	return questions, nil
}
