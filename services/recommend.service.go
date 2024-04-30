package services

import (
	"errors"
	"fmt"
	"github.com/hodukihugi/winglets-api/core"
	"github.com/hodukihugi/winglets-api/models"
	"github.com/hodukihugi/winglets-api/repositories"
	"gorm.io/gorm"
)

type IRecommendService interface {
	CreateUserAnswer(models.SerializableAnswer) error
	GetMatchesById(string) error
	GetRecommendationById(string) error
	SmashById(string) error
	PassById(string) error
}

type RecommendService struct {
	answerRepository repositories.IAnswerRepository
	matchRepository  repositories.IMatchRepository
	logger           *core.Logger
}

func NewRecommendService(
	answerRepository repositories.IAnswerRepository,
	matchRepository repositories.IMatchRepository,
	logger *core.Logger,
) IRecommendService {
	return &RecommendService{
		answerRepository: answerRepository,
		matchRepository:  matchRepository,
		logger:           logger,
	}
}

func (s *RecommendService) CreateUserAnswer(answer models.SerializableAnswer) error {
	existingAnswer, err := s.answerRepository.FindAnswer(models.Answer{
		UserID:     answer.UserID,
		QuestionID: answer.QuestionID,
	})

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if existingAnswer != nil {
		err := s.answerRepository.DeleteAnswer(models.Answer{
			UserID:     answer.UserID,
			QuestionID: answer.QuestionID,
		})

		if err != nil {
			return err
		}
		s.logger.Debug(fmt.Sprintf("exsistingAnswer: %v", existingAnswer))
	}

	_, err = s.answerRepository.CreateAnswer(models.Answer{
		UserID:       answer.UserID,
		QuestionID:   answer.QuestionID,
		UserAnswer:   answer.UserAnswer,
		PreferAnswer: answer.PreferAnswer,
		Importance:   answer.Importance,
	})
	return err
}

func (s *RecommendService) GetMatchesById(id string) error {
	return errors.New("not implemented")
}

func (s *RecommendService) GetRecommendationById(id string) error {
	return errors.New("not implemented")
}

func (s *RecommendService) SmashById(id string) error {
	return errors.New("not implemented")
}

func (s *RecommendService) PassById(id string) error {
	return errors.New("not implemented")
}
