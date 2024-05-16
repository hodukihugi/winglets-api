package services

import (
	"errors"
	"fmt"
	"github.com/hodukihugi/winglets-api/core"
	"github.com/hodukihugi/winglets-api/models"
	"github.com/hodukihugi/winglets-api/repositories"
	"github.com/hodukihugi/winglets-api/utils"
	"gorm.io/gorm"
	"sort"
	"sync"
)

type IRecommendService interface {
	CreateUserAnswer(models.SerializableAnswer) error
	GetMatchesByUserId(string) error
	GetAnswersByUserId(string) ([]models.SerializableAnswer, error)
	GetListQuestions() ([]models.SerializableQuestion, error)
	GetRecommendationByUserId(string, int, int, float64, float64) ([]models.MatchProfile, error)
	SmashById(string) error
	PassById(string) error
}

type RecommendService struct {
	profileRepository  repositories.IProfileRepository
	answerRepository   repositories.IAnswerRepository
	matchRepository    repositories.IMatchRepository
	questionRepository repositories.IQuestionRepository
	logger             *core.Logger
}

func NewRecommendService(
	profileRepository repositories.IProfileRepository,
	answerRepository repositories.IAnswerRepository,
	matchRepository repositories.IMatchRepository,
	questionRepository repositories.IQuestionRepository,
	logger *core.Logger,
) IRecommendService {
	return &RecommendService{
		profileRepository:  profileRepository,
		answerRepository:   answerRepository,
		matchRepository:    matchRepository,
		questionRepository: questionRepository,
		logger:             logger,
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

func (s *RecommendService) GetMatchesByUserId(id string) error {
	return errors.New("not implemented")
}

func (s *RecommendService) GetAnswersByUserId(id string) ([]models.SerializableAnswer, error) {
	var result []models.SerializableAnswer
	answers, err := s.answerRepository.FindListAnswerByUserId(id)
	if err != nil {
		return nil, err
	}

	for _, answer := range answers {
		serializableAnswer := answer.Serialize()
		result = append(result, *serializableAnswer)
	}

	if result == nil || len(result) == 0 {
		return nil, errors.New("user answers not found")
	}

	return result, nil
}

func (s *RecommendService) GetListQuestions() ([]models.SerializableQuestion, error) {
	var result []models.SerializableQuestion
	questions, err := s.questionRepository.GetListQuestions()
	if err != nil {
		return nil, err
	}

	for _, question := range questions {
		serializeQuestion := question.Serialize()
		result = append(result, *serializeQuestion)
	}

	return result, nil
}

func (s *RecommendService) GetRecommendationByUserId(
	userId string,
	minAge int,
	maxAge int,
	minDistance float64,
	maxDistance float64,
) ([]models.MatchProfile, error) {
	var recommendedProfiles []models.MatchProfile

	userProfile, err := s.profileRepository.GetProfileById(userId)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	longitude, latitude, err := utils.CoordinatesStringToPairFloat64(userProfile.Coordinates)

	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	satisfiedProfiles, err := s.profileRepository.GetListProfile(models.ProfileFilter{
		ExcludedUserId: userId,
		MinAge:         minAge,
		MaxAge:         maxAge,
		MinDistance:    minDistance,
		MaxDistance:    maxDistance,
		Longitude:      longitude,
		Latitude:       latitude,
	})

	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	userAnswers, err := s.answerRepository.FindListAnswerByUserId(userId)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	mapUserAnswers := make(map[int]*models.Answer)
	for _, userAnswer := range userAnswers {
		mapUserAnswers[userAnswer.QuestionID] = &userAnswer
	}

	var wg sync.WaitGroup
	var matchCalculationResultChan = make(chan models.MatchCalculationResult, len(satisfiedProfiles))
	for _, otherProfile := range satisfiedProfiles {
		otherAnswers, err := s.answerRepository.FindListAnswerByUserId(otherProfile.ID)
		mapOtherAnswers := make(map[int]*models.Answer)
		for _, otherAnswer := range otherAnswers {
			mapOtherAnswers[otherAnswer.QuestionID] = &otherAnswer
		}
		if err != nil {
			s.logger.Error(err)
			return nil, err
		}
		wg.Add(1)
		go utils.CalculateMatchPercentage(&wg, matchCalculationResultChan, mapUserAnswers, mapOtherAnswers, otherProfile)
	}
	wg.Wait()
	close(matchCalculationResultChan)

	var matchResults []models.MatchCalculationResult

	for result := range matchCalculationResultChan {
		matchResults = append(matchResults, result)
	}

	sort.Slice(matchResults, func(i, j int) bool {
		return matchResults[i].MatchPercentage > matchResults[j].MatchPercentage
	})

	for _, result := range matchResults {
		matchProfile := result.MatchedProfile.ConvertToMatchProfile()
		lon, lat, err := utils.CoordinatesStringToPairFloat64(result.MatchedProfile.Coordinates)
		if err != nil {
			s.logger.Error(err)
			return nil, err
		}
		matchProfile.Distance = utils.CalculateDistance(longitude, latitude, lon, lat)
		matchProfile.MatchPercentage = result.MatchPercentage
		recommendedProfiles = append(recommendedProfiles, *matchProfile)
	}
	// TODO: Lưu những profile đã được gợi í
	return recommendedProfiles, nil
}

func (s *RecommendService) SmashById(id string) error {
	return errors.New("not implemented")
}

func (s *RecommendService) PassById(id string) error {
	return errors.New("not implemented")
}
