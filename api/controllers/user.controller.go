package controllers

import (
	"github.com/hodukihugi/winglets-api/models"
	"github.com/hodukihugi/winglets-api/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hodukihugi/winglets-api/core"
)

// UserController data type
type UserController struct {
	service services.IUserService
	logger  *core.Logger
}

// NewUserController creates new user controller
func NewUserController(userService services.IUserService, logger *core.Logger) *UserController {
	return &UserController{
		service: userService,
		logger:  logger,
	}
}

// GetByID gets one user
func (u *UserController) GetByID(c *gin.Context) {
	paramID := c.Param("id")

	id, err := strconv.ParseUint(paramID, 10, 32)
	if err != nil {
		u.logger.Error(err)
		c.JSON(http.StatusBadRequest, models.HTTPResponse{
			Message:       err.Error(),
			InvalidFields: []string{"id"},
		})
		return
	}

	var uintID = uint(id)
	user, err := u.service.First(models.OneUserFilter{ID: uintID})

	if err != nil {
		u.logger.Errorf("fail to get user by id, id: [%v], err: [%v]", id, err)
		c.JSON(http.StatusInternalServerError, models.HTTPResponse{
			Message: "server error",
		})
		return
	}

	c.JSON(http.StatusOK, models.HTTPResponse{Data: map[string]interface{}{"user": user.Serialize()}})
}
