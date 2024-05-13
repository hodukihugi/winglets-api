package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/hodukihugi/winglets-api/core"
	"github.com/hodukihugi/winglets-api/models"
	"github.com/hodukihugi/winglets-api/services"
	"github.com/hodukihugi/winglets-api/utils"
	"net/http"
	"strconv"
)

// RecommendController data type
type RecommendController struct {
	service services.IRecommendService
	logger  *core.Logger
}

// NewRecommendController creates new match controller
func NewRecommendController(recommendService services.IRecommendService, logger *core.Logger) *RecommendController {
	return &RecommendController{
		service: recommendService,
		logger:  logger,
	}
}

func (c *RecommendController) Answer(ctx *gin.Context) {
	var request models.AnswerRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		c.logger.Error(err)
		ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
			Message: err.Error(),
		})
		return
	}

	userID, err := utils.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.HTTPResponse{
			Message: "server error",
		})
		return
	}

	var invalidFields []string

	for i := 0; i < len(request.Answers); i++ {
		err := c.service.CreateUserAnswer(models.SerializableAnswer{
			UserID:       userID,
			QuestionID:   request.Answers[i].QuestionID,
			UserAnswer:   request.Answers[i].UserAnswer,
			PreferAnswer: request.Answers[i].PreferAnswer,
			Importance:   request.Answers[i].Importance,
		})

		if err != nil {
			invalidFields = append(invalidFields, err.Error())
		}
	}

	if len(invalidFields) > 0 {
		ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
			Message:       "invalid fields",
			InvalidFields: invalidFields,
		})
		return
	}

	ctx.JSON(http.StatusOK, models.HTTPResponse{
		Message: "success",
	})

}

func (c *RecommendController) GetUserMatches(ctx *gin.Context) {

}

func (c *RecommendController) GetRecommendations(ctx *gin.Context) {
	var minAgeInt, maxAgeInt int
	var minDistanceFloat, maxDistanceFloat float64
	var err error

	// Get query params
	minAge, ok := ctx.GetQuery("min_age")
	if !ok {
		minAgeInt = 18
	} else {
		minAgeInt, err = strconv.Atoi(minAge)
		if err != nil {
			c.logger.Error(err)
			ctx.JSON(http.StatusInternalServerError, models.HTTPResponse{
				Message: "server error",
			})
			return
		}
	}

	maxAge, ok := ctx.GetQuery("max_age")
	if !ok {
		maxAgeInt = 99
	} else {
		maxAgeInt, err = strconv.Atoi(maxAge)
		if err != nil {
			c.logger.Error(err)
			ctx.JSON(http.StatusInternalServerError, models.HTTPResponse{
				Message: "server error",
			})
			return
		}
	}

	minDistance, ok := ctx.GetQuery("min_distance")
	if !ok {
		minDistanceFloat = 0
	} else {
		minDistanceFloat, err = strconv.ParseFloat(minDistance, 64)
		if err != nil {
			c.logger.Error(err)
			ctx.JSON(http.StatusInternalServerError, models.HTTPResponse{
				Message: "server error",
			})
			return
		}
	}

	maxDistance, ok := ctx.GetQuery("max_distance")
	if !ok {
		maxDistanceFloat = 100
	} else {
		maxDistanceFloat, err = strconv.ParseFloat(maxDistance, 64)
		if err != nil {
			c.logger.Error(err)
			ctx.JSON(http.StatusInternalServerError, models.HTTPResponse{
				Message: "server error",
			})
			return
		}
	}

	userID, err := utils.GetUserID(ctx)
	if err != nil {
		c.logger.Error(err)
		ctx.JSON(http.StatusInternalServerError, models.HTTPResponse{
			Message: "server error",
		})
		return
	}

	profiles, err := c.service.GetRecommendationById(userID, minAgeInt, maxAgeInt, minDistanceFloat, maxDistanceFloat)
	if err != nil {
		c.logger.Error(err)
		ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.HTTPResponse{
		Message: "get recommendations success",
		Data:    profiles,
	})

}

func (c *RecommendController) Smash(ctx *gin.Context) {

}

func (c *RecommendController) Pass(ctx *gin.Context) {

}