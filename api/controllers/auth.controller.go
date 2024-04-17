package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hodukihugi/winglets-api/constants"
	"github.com/hodukihugi/winglets-api/core"
	"github.com/hodukihugi/winglets-api/models"
	"github.com/hodukihugi/winglets-api/services"
	"github.com/hodukihugi/winglets-api/utils"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

// AuthController struct
type AuthController struct {
	logger      *core.Logger
	service     services.IAuthService
	userService services.IUserService
	validator   *core.Validator
	env         *core.Env
}

// NewAuthController creates new controller
func NewAuthController(
	logger *core.Logger,
	service services.IAuthService,
	userService services.IUserService,
	validator *core.Validator,
	env *core.Env,
) *AuthController {
	return &AuthController{
		logger:      logger,
		service:     service,
		userService: userService,
		validator:   validator,
		env:         env,
	}
}

// SignIn signs in user
func (c *AuthController) SignIn(ctx *gin.Context) {
	var payload models.SignInRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
			Message: "fail to parse request body",
		})
		return
	}

	if validationErrors := c.validator.Validate.Struct(&payload); validationErrors != nil {
		var invalidFields []string
		for _, validationErr := range validationErrors.(validator.ValidationErrors) {
			invalidFields = append(invalidFields, utils.PascalToSnake(validationErr.Field()))
		}
		ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
			Message:       "invalid request body",
			InvalidFields: invalidFields,
		})
		return
	}

	user, err := c.userService.First(models.OneUserFilter{Email: payload.Email})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, models.HTTPResponse{
				Message: "user not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, models.HTTPResponse{
			Message: "server error",
		})
		c.logger.Errorf("fail to sign in, payload [%v], error [%v]", payload, err)
		return
	}

	if err = utils.VerifyPassword(user.Password, payload.Password); err != nil {
		ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
			Message: "wrong password",
		})
		return
	}

	accessToken, refreshToken, accessExpired, refreshExpired, err := c.service.GenerateJWTTokens(*user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.HTTPResponse{
			Message: "server error",
		})
		c.logger.Errorf("fail to gen jwt tokens, payload [%v], error [%v]", payload, err)
		return
	}

	utils.AttachCookiesToResponse(c.env, accessToken, refreshToken, ctx)
	ctx.JSON(http.StatusCreated, models.HTTPResponse{Message: "success", Data: map[string]interface{}{
		"access_token":    accessToken,
		"refresh_token":   refreshToken,
		"access_expired":  accessExpired,
		"refresh_expired": refreshExpired,
	}})
	return
}

// Refresh renew session
func (c *AuthController) Refresh(ctx *gin.Context) {
	claim, ok := ctx.Get(constants.CtxKey_JWTClaim)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, models.HTTPResponse{
			Message: "server error",
		})
		c.logger.Error("fail to get jwt claim")
		return
	}

	var jwtClaim = claim.(*models.JWTClaim)
	accessToken, accessExpired, err := c.service.Refresh(models.User{
		Email: jwtClaim.UserEmail,
		ID:    jwtClaim.UserID,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.HTTPResponse{
			Message: "server error",
		})
		c.logger.Errorf("fail to refresh, error [%v]", err)
		return
	}

	ctx.JSON(http.StatusCreated, models.HTTPResponse{
		Message: "success",
		Data: map[string]interface{}{
			"access_token":   accessToken,
			"access_expired": accessExpired,
		}})
	return
}

// Register registers user
func (c *AuthController) Register(ctx *gin.Context) {
	var payload models.RegisterRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, models.HTTPResponse{Message: "fail to parse request body"})
		return
	}

	if errs := c.validator.Validate.Struct(&payload); errs != nil {
		var invalidFields []string
		for _, err := range errs.(validator.ValidationErrors) {
			invalidFields = append(invalidFields, utils.PascalToSnake(err.Field()))
		}
		ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
			Message:       "invalid request body",
			InvalidFields: invalidFields,
		})
		return
	}

	if err := c.service.Register(payload); err != nil {
		emailDuplicated :=
			strings.Contains(err.Error(), "users.email_unique") &&
				strings.Contains(err.Error(), "Duplicate")
		if emailDuplicated {
			ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
				Message: "duplicate email",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, models.HTTPResponse{Message: "server error"})
		c.logger.Errorf("fail to register new user, payload [%v], error [%v]", payload, err)
		return
	}
	ctx.JSON(http.StatusCreated, models.HTTPResponse{Message: "success"})
	return
}
