package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/hodukihugi/winglets-api/models"
	"github.com/hodukihugi/winglets-api/services"
	"github.com/hodukihugi/winglets-api/utils"
	"net/http"
)

type userController struct {
	service services.UserService
}

func NewUserController(service services.UserService) *userController {
	return &userController{service: service}
}

func (c *userController) CreateUser(ctx *gin.Context) {
	var input models.CreateUserEntity
	if err := ctx.ShouldBind(&input); err != nil {
		utils.APIResponse(ctx, "Bad request", http.StatusBadRequest, http.MethodPost, nil)
		return
	}

	switch errCreateUser := c.service.CreateUser(&input); errCreateUser {
	case "DB_ERROR_CREATE_USER_FAILED":
		utils.APIResponse(ctx, "Create new user account failed", http.StatusForbidden, http.MethodPost, nil)
	default:
		utils.APIResponse(ctx, "Create new user account successfully", http.StatusCreated, http.MethodPost, nil)
	}
}

func (c *userController) GetUserById(ctx *gin.Context) {
	data, errCreateUser := c.service.GetUserById(ctx.Param("id"))
	switch errCreateUser {
	case "USER_NOT_FOUND_404":
		utils.APIResponse(ctx, "User not exist", http.StatusNotFound, http.MethodGet, nil)
	default:
		utils.APIResponse(ctx, "Get user account by id successfully", http.StatusCreated, http.MethodPost, data)
	}
}

func (c *userController) GetListUser(ctx *gin.Context) {
	resultGetListUser, _ := c.service.GetListUser()
	utils.APIResponse(ctx, "Get list user successfully", http.StatusOK, http.MethodGet, resultGetListUser)
}

func (c *userController) UpdateUserById(ctx *gin.Context) {
	var input models.UpdateUserEntity
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.APIResponse(ctx, "Bad request", http.StatusBadRequest, http.MethodPost, nil)
		return
	}

	switch errUpdateUser := c.service.UpdateUserById(ctx.Param("id"), &input); errUpdateUser {
	case "DB_ERROR_UPDATE_USER_FAILED":
		utils.APIResponse(ctx, "Update user account failed", http.StatusConflict, http.MethodPut, nil)
	default:
		utils.APIResponse(ctx, "Update user account by id successfully", http.StatusOK, http.MethodPut, nil)
	}
}

func (c *userController) DeleteUserById(ctx *gin.Context) {
	errDeleteUser := c.service.DeleteUserById(ctx.Param("id"))
	switch errDeleteUser {
	case "USER_NOT_FOUND_404":
		utils.APIResponse(ctx, "User not exist", http.StatusNotFound, http.MethodDelete, nil)
	default:
		utils.APIResponse(ctx, "Delete user successfully", http.StatusOK, http.MethodGet, nil)
	}
}
