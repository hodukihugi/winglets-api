package utils

import (
	"github.com/hodukihugi/winglets-api/models"
	"math"
	"sync"
)

var (
	importancePoint = map[int]int{
		1: 0,
		2: 1,
		3: 10,
		4: 50,
		5: 250,
	}
)

func CalculateMatchPercentage(
	wg *sync.WaitGroup,
	resultChan chan models.MatchCalculationResult,
	userAnswers map[int]*models.Answer,
	otherAnswers map[int]*models.Answer,
	otherProfile models.Profile,
) {
	defer wg.Done()

	var userSatisfaction, otherSatisfaction float64
	var userTotalPoint, userMaximumPoint int
	var otherTotalPoint, otherMaximumPoint int
	var matchPercentage float64

	for questionID, userAnswer := range userAnswers {
		// Lấy ra câu hỏi giống user
		otherUserAnswer := otherAnswers[questionID]
		if otherUserAnswer == nil {
			continue
		}

		// So sánh câu trả lời của hai người
		userMaximumPoint += importancePoint[otherUserAnswer.Importance]
		otherMaximumPoint += importancePoint[userAnswer.Importance]

		if userAnswer.UserAnswer == otherUserAnswer.PreferAnswer {
			userTotalPoint += importancePoint[otherUserAnswer.Importance]
		}

		if otherUserAnswer.UserAnswer == userAnswer.PreferAnswer {
			otherTotalPoint += importancePoint[userAnswer.Importance]
		}
	}

	userSatisfaction = float64(userTotalPoint) / float64(userMaximumPoint)
	otherSatisfaction = float64(otherTotalPoint) / float64(otherMaximumPoint)

	matchPercentage = math.Pow(userSatisfaction*otherSatisfaction, 1.0/float64(len(userAnswers)))

	resultChan <- models.MatchCalculationResult{
		MatchPercentage: matchPercentage,
		MatchedProfile:  otherProfile,
	}
}
