package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/hodukihugi/winglets-api/models"
	"github.com/hodukihugi/winglets-api/services"
	"github.com/hodukihugi/winglets-api/utils"
	"net/http"
)

type profileController struct {
	service services.ProfileService
}

func NewProfileController(service services.ProfileService) *profileController {
	return &profileController{service: service}
}

func (c *profileController) CreateProfile(ctx *gin.Context) {
	var request models.CreateProfileRequest
	if err := ctx.ShouldBind(&request); err != nil {
		utils.APIResponseSimple(ctx, err.Error(), http.StatusBadRequest)
		return
	}

	if err := c.service.CreateProfile(&request); err != nil {
		utils.APIResponse(ctx, err.Error(), http.StatusForbidden, nil)
	} else {
		utils.APIResponse(ctx, "Create new profile successfully", http.StatusCreated, nil)
	}
}

func (c *profileController) GetProfileById(ctx *gin.Context) {
	data, err := c.service.GetProfileById(ctx.Param("id"))
	if err != nil {
		utils.APIResponse(ctx, err.Error(), http.StatusNotFound, nil)
	} else {
		utils.APIResponse(ctx, "Get profile by id successfully", http.StatusCreated, data)
	}
}

func (c *profileController) GetListProfile(ctx *gin.Context) {
	result, _ := c.service.GetListProfile()
	utils.APIResponse(ctx, "Get list profile successfully", http.StatusOK, result)
}

func (c *profileController) UpdateProfileById(ctx *gin.Context) {
	var input models.UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.APIResponseSimple(ctx, err.Error(), http.StatusBadRequest)
		return
	}

	if err := c.service.UpdateProfileById(ctx.Param("id"), &input); err != nil {
		utils.APIResponse(ctx, err.Error(), http.StatusConflict, nil)
	} else {
		utils.APIResponse(ctx, "Update profile by id successfully", http.StatusOK, nil)
	}
}

func (c *profileController) DeleteProfileById(ctx *gin.Context) {
	err := c.service.DeleteProfileById(ctx.Param("id"))
	if err != nil {
		utils.APIResponse(ctx, err.Error(), http.StatusNotFound, nil)
	} else {
		utils.APIResponse(ctx, "Delete profile successfully", http.StatusOK, nil)
	}
}
