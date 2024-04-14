package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/hodukihugi/winglets-api/core"
	"github.com/hodukihugi/winglets-api/models"
	"github.com/hodukihugi/winglets-api/services"
	"github.com/hodukihugi/winglets-api/utils"
	"net/http"
)

type ProfileController struct {
	service services.IProfileService
	logger  *core.Logger
}

func NewProfileController(service services.IProfileService, logger *core.Logger) *ProfileController {
	return &ProfileController{
		service: service,
		logger:  logger,
	}
}

func (c *ProfileController) CreateProfile(ctx *gin.Context) {
	var request models.ProfileCreateRequest
	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
			Message: err.Error(),
		})
		ctx.Abort()
		return
	}

	userID, err := utils.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.HTTPResponse{
			Message: "server error",
		})
		ctx.Abort()
		return
	}

	_, err = c.service.GetProfileById(userID)
	if err == nil {
		ctx.JSON(http.StatusConflict, models.HTTPResponse{
			Message: "profile exists",
		})
		ctx.Abort()
		return
	}

	if err := c.service.CreateProfile(userID, request); err != nil {
		ctx.JSON(http.StatusForbidden, models.HTTPResponse{
			Message: err.Error(),
		})
		ctx.Abort()
		return
	} else {
		ctx.JSON(http.StatusCreated, models.HTTPResponse{
			Message: "Create new profile successfully",
		})
	}
}

func (c *ProfileController) GetProfileById(ctx *gin.Context) {
	data, err := c.service.GetProfileById(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, models.HTTPResponse{
			Message: err.Error(),
		})
		ctx.Abort()
		return
	} else {
		ctx.JSON(http.StatusNotFound, models.HTTPResponse{
			Message: "Get profile by id successfully",
			Data:    data,
		})
	}
}

func (c *ProfileController) GetMyProfile(ctx *gin.Context) {
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.HTTPResponse{
			Message: "server error",
		})
		ctx.Abort()
		return
	}

	result, _ := c.service.GetProfileById(userID)
	ctx.JSON(http.StatusOK, models.HTTPResponse{
		Message: "Get user profile successfully",
		Data:    result.Serialize(),
	})
}

func (c *ProfileController) UpdateProfile(ctx *gin.Context) {
	var request models.ProfileUpdateRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
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
		ctx.Abort()
		return
	}

	if err := c.service.UpdateProfileById(userID, request); err != nil {
		ctx.JSON(http.StatusConflict, models.HTTPResponse{
			Message: err.Error(),
		})
		ctx.Abort()
		return
	} else {
		ctx.JSON(http.StatusOK, models.HTTPResponse{
			Message: "Update profile by id successfully",
		})
	}
}

func (c *ProfileController) DeleteProfile(ctx *gin.Context) {
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.HTTPResponse{
			Message: "server error",
		})
		ctx.Abort()
		return
	}

	err = c.service.DeleteProfileById(userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, models.HTTPResponse{
			Message: err.Error(),
		})
		ctx.Abort()
		return
	} else {
		ctx.JSON(http.StatusOK, models.HTTPResponse{
			Message: "Delete profile successfully",
		})
	}
}
