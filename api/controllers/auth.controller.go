package controllers

import (
	"crypto/rand"
	"encoding/base32"
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
	"sync"
	"time"
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

	if user.VerificationStatus == 0 {
		ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
			Message: "email is not verified",
		})
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
	ctx.JSON(http.StatusCreated, models.HTTPResponse{
		Message: "success",
		Data: map[string]interface{}{
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

	result, err := c.service.Register(payload)
	if err != nil {
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

	//send verification token to user's email
	var wg sync.WaitGroup
	ch := make(chan error, 1)
	email := []string{result.Email}
	wg.Add(1)
	go utils.SendVerificationEmailAsync(&wg, ch, c.env, result.VerificationCode, email)
	wg.Wait()
	if err = <-ch; err != nil {
		ctx.JSON(http.StatusInternalServerError, models.HTTPResponse{
			Message: "server error",
		})
		c.logger.Debug(err)
		return
	}

	ctx.JSON(http.StatusOK, models.HTTPResponse{
		Message: "success",
	})
	return
}

// VerifyEmail verify user email
func (c *AuthController) VerifyEmail(ctx *gin.Context) {
	var verifyEmailRequest models.VerifyEmailRequest
	if err := ctx.ShouldBindJSON(&verifyEmailRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
			Message: "fail to parse request body",
		})
		return
	}

	email, verifyToken := verifyEmailRequest.Email, verifyEmailRequest.VerificationCode
	verificationRequest := models.VerifyEmailRequest{
		Email:            email,
		VerificationCode: verifyToken,
	}

	if errs := c.validator.Validate.Struct(&verificationRequest); errs != nil {
		var invalidFields []string
		for _, err := range errs.(validator.ValidationErrors) {
			invalidFields = append(invalidFields, utils.PascalToSnake(err.Field()))
		}
		ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
			Message:       "invalid request",
			InvalidFields: invalidFields,
		})
		return
	}

	filterUser := models.OneUserFilter{
		Email: email,
	}

	result, err := c.userService.First(filterUser)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.HTTPResponse{
			Message: "server error",
		})
		return
	}

	if result.VerificationStatus == 1 {
		ctx.JSON(http.StatusConflict, models.HTTPResponse{
			Message: "email has already been verified",
		})
		return
	}

	expireTime := result.VerificationTime.Add(c.env.EmailVerificationExpiresIn)

	isValid := (result.VerificationCode == verifyToken) &&
		(time.Now().UTC().Before(expireTime))

	if !isValid {
		ctx.JSON(http.StatusRequestTimeout, models.HTTPResponse{
			Message: "email verification timeout",
		})
		return
	}

	userUpdateRequest := models.UserUpdateRequest{
		VerificationStatus: 1,
	}

	if err = c.userService.UpdateById(result.ID, userUpdateRequest); err != nil {
		c.logger.Debug(err)
		ctx.JSON(http.StatusInternalServerError, models.HTTPResponse{
			Message: "server error",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.HTTPResponse{
		Message: "success",
	})
}

// SendVerificationEmail send user verification email
func (c *AuthController) SendVerificationEmail(ctx *gin.Context) {
	var request models.SendVerificationEmailRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
			Message: "fail to parse request body",
		})
		return
	}

	if errs := c.validator.Validate.Struct(&request); errs != nil {
		var invalidFields []string
		for _, err := range errs.(validator.ValidationErrors) {
			invalidFields = append(invalidFields, utils.PascalToSnake(err.Field()))
		}
		ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
			Message:       "invalid request",
			InvalidFields: invalidFields,
		})
		return
	}

	result, err := c.userService.First(models.OneUserFilter{
		Email: request.Email,
	})

	if err != nil {
		ctx.JSON(http.StatusNotFound, models.HTTPResponse{
			Message: "user not exist",
		})
		return
	}

	if result.VerificationStatus == 1 {
		ctx.JSON(http.StatusConflict, models.HTTPResponse{
			Message: "user has already been verified",
		})
		return
	}

	//generate verification token
	randomBytes := make([]byte, 10)
	_, err = rand.Read(randomBytes)
	if err != nil {
		c.logger.Debug(err)
		ctx.JSON(http.StatusInternalServerError, models.HTTPResponse{
			Message: "server error",
		})
		return
	}

	verificationCode := base32.StdEncoding.EncodeToString(randomBytes)[:10]

	err = c.userService.UpdateById(result.ID, models.UserUpdateRequest{
		VerificationCode: verificationCode,
		VerificationTime: time.Now().UTC(),
	})

	if err != nil {
		c.logger.Debug(err)
		ctx.JSON(http.StatusInternalServerError, models.HTTPResponse{
			Message: "server error",
		})
		return
	}

	//send verification token to user's email
	var wg sync.WaitGroup
	ch := make(chan error, 1)
	email := []string{request.Email}
	wg.Add(1)
	go utils.SendVerificationEmailAsync(&wg, ch, c.env, verificationCode, email)
	wg.Wait()
	if err = <-ch; err != nil {
		ctx.JSON(http.StatusInternalServerError, models.HTTPResponse{
			Message: "server error",
		})
		c.logger.Debug(err)
		return
	}

	ctx.JSON(http.StatusOK, models.HTTPResponse{
		Message: "success",
	})
	return
}
