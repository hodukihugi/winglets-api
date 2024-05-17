package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hodukihugi/winglets-api/core"
	"github.com/hodukihugi/winglets-api/models"
	"github.com/hodukihugi/winglets-api/services"
	"github.com/hodukihugi/winglets-api/utils"
	"net/http"
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

func (c *RecommendController) GetUserAnswers(ctx *gin.Context) {
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		c.logger.Error(err)
		ctx.JSON(http.StatusInternalServerError, models.HTTPResponse{
			Message: "server error",
		})
		return
	}

	userAnswers, err := c.service.GetAnswersByUserId(userID)
	if err != nil {
		c.logger.Error(err)
		if err.Error() == "user answers not found" {
			ctx.JSON(http.StatusConflict, models.HTTPResponse{
				Message: "user answers not found",
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, models.HTTPResponse{
				Message: "server error",
			})
		}
		return
	}

	ctx.JSON(http.StatusOK, models.HTTPResponse{
		Message: "success",
		Data:    map[string]interface{}{"answers": userAnswers},
	})
}

func (c *RecommendController) GetQuestions(ctx *gin.Context) {
	questions, err := c.service.GetListQuestions()
	if err != nil {
		c.logger.Error(err)
		ctx.JSON(http.StatusInternalServerError, models.HTTPResponse{
			Message: "server error",
		})
	}

	ctx.JSON(http.StatusOK, models.HTTPResponse{
		Message: "success",
		Data:    map[string]interface{}{"questions": questions},
	})
}

func (c *RecommendController) GetRecommendations(ctx *gin.Context) {
	var err error
	var request models.GetRecommendationRequest
	if err = ctx.ShouldBindJSON(&request); err != nil {
		c.logger.Error(err)
		ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
			Message: err.Error(),
		})
		return
	}

	minAgeInt, maxAgeInt := request.MinAge, request.MaxAge
	minDistanceFloat, maxDistanceFloat := request.MinDistance, request.MaxDistance
	c.logger.Info(fmt.Sprintf("Request: %+v", request))
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		c.logger.Error(err)
		ctx.JSON(http.StatusInternalServerError, models.HTTPResponse{
			Message: "server error",
		})
		return
	}

	profiles, err := c.service.GetRecommendationByUserId(userID, minAgeInt, maxAgeInt, minDistanceFloat, maxDistanceFloat)
	if err != nil {
		c.logger.Error(err)
		ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.HTTPResponse{
		Message: "success",
		Data:    map[string]interface{}{"profiles": profiles},
	})

}

func (c *RecommendController) Smash(ctx *gin.Context) {
	var request models.SmashRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		c.logger.Error(err)
		ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
			Message: "fail to parse request",
		})
	}

	userID, err := utils.GetUserID(ctx)
	if err != nil {
		c.logger.Error(err)
		ctx.JSON(http.StatusInternalServerError, models.HTTPResponse{
			Message: "server error",
		})
	}

	message, profile, err := c.service.SmashById(userID, request.UserId)
	if err != nil {
		c.logger.Error(err)
		ctx.JSON(http.StatusInternalServerError, models.HTTPResponse{
			Message: "server error",
		})
		return
	}

	if message != "match wait" && message != "match finish" {
		c.logger.Error(message)
		ctx.JSON(http.StatusInternalServerError, models.HTTPResponse{
			Message: "server error",
		})
		return
	}

	if message == "match wait" {
		ctx.JSON(http.StatusOK, models.HTTPResponse{
			Message: "success, match wait",
		})
	} else {
		ctx.JSON(http.StatusOK, models.HTTPResponse{
			Message: "success, match finish",
			Data: map[string]*models.Profile{
				"profile": profile,
			},
		})
	}
}

func (c *RecommendController) Pass(ctx *gin.Context) {
	var request models.PassRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		c.logger.Error(err)
		ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
			Message: "fail to parse request",
		})
	}

	userID, err := utils.GetUserID(ctx)
	if err != nil {
		c.logger.Error(err)
		ctx.JSON(http.StatusInternalServerError, models.HTTPResponse{
			Message: "server error",
		})
	}

	if err := c.service.PassById(userID, request.UserId); err != nil {
		c.logger.Error(err)
		ctx.JSON(http.StatusInternalServerError, models.HTTPResponse{
			Message: "server error",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.HTTPResponse{
		Message: "success",
	})
}
