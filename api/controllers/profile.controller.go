package controllers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hodukihugi/winglets-api/core"
	"github.com/hodukihugi/winglets-api/models"
	"github.com/hodukihugi/winglets-api/services"
	"github.com/hodukihugi/winglets-api/utils"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ProfileController struct {
	service          services.IProfileService
	recommendService services.IRecommendService
	ik               *core.ImageKit
	logger           *core.Logger
}

func NewProfileController(
	service services.IProfileService,
	recommendService services.IRecommendService,
	ik *core.ImageKit,
	logger *core.Logger,
) *ProfileController {
	return &ProfileController{
		service:          service,
		recommendService: recommendService,
		ik:               ik,
		logger:           logger,
	}
}

func (c *ProfileController) CreateProfile(ctx *gin.Context) {
	var request models.ProfileCreateRequest
	if err := ctx.ShouldBind(&request); err != nil {
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

	result, _ := c.service.GetProfileById(userID)
	if result != nil {
		ctx.JSON(http.StatusConflict, models.HTTPResponse{
			Message: "profile exists",
		})
		return
	}

	if err = c.service.CreateProfile(userID, request); err != nil {
		ctx.JSON(http.StatusForbidden, models.HTTPResponse{
			Message: err.Error(),
		})
		return
	} else {
		ctx.JSON(http.StatusCreated, models.HTTPResponse{
			Message: "success",
		})
	}
}

func (c *ProfileController) GetProfileById(ctx *gin.Context) {
	data, err := c.service.GetProfileById(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusOK, models.HTTPResponse{
			Message: err.Error(),
		})
		return
	} else {
		ctx.JSON(http.StatusOK, models.HTTPResponse{
			Message: "success",
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
		return
	}

	result, err := c.service.GetProfileById(userID)

	if err != nil {
		ctx.JSON(http.StatusOK, models.HTTPResponse{
			Message: "profile not found",
		})
		return
	}

	answered := 1
	_, err = c.recommendService.GetAnswersByUserId(userID)
	if err != nil {
		if err.Error() == "user answers not found" {
			answered = 0
		} else {
			ctx.JSON(http.StatusInternalServerError, models.HTTPResponse{
				Message: "server error",
			})
			return
		}
	}

	serializeProfile := result.Serialize()
	serializeProfile.Answered = answered

	ctx.JSON(http.StatusOK, models.HTTPResponse{
		Message: "success",
		Data:    serializeProfile,
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
		return
	}

	if err := c.service.UpdateProfileById(userID, request); err != nil {
		ctx.JSON(http.StatusConflict, models.HTTPResponse{
			Message: err.Error(),
		})
		return
	} else {
		ctx.JSON(http.StatusOK, models.HTTPResponse{
			Message: "success",
		})
	}
}

func (c *ProfileController) UploadImage(ctx *gin.Context) {
	context, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	userID, err := utils.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.HTTPResponse{
			Message: "server error",
		})
		return
	}

	var wg sync.WaitGroup
	ch := make(chan models.ImageUploadResult, 5)

	// Parse the multipart form
	err = ctx.Request.ParseMultipartForm(10 << 20) // Max file size of 10MB
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
			Message: "bad request",
		})
	}

	// Get the files from the form
	form, _ := ctx.MultipartForm()
	files := form.File

	var invalidFields []string

	for key, images := range files {
		if !strings.HasPrefix(key, "image_") {
			continue
		}

		if len(images) <= 0 {
			invalidFields = append(invalidFields, fmt.Sprintf("Empty image: %s", key))
			break
		}

		slotId, err := strconv.Atoi(strings.Trim(key, "image_"))
		if err != nil {
			invalidFields = append(invalidFields, fmt.Sprintf("Can't convert image slot id: %s", key))
			c.logger.Debug(err)
			break
		}

		// Open the uploaded file
		fileHeader := images[0]
		src, err := fileHeader.Open()
		defer src.Close()

		if err != nil {
			invalidFields = append(invalidFields, fmt.Sprintf("Can't open image: %s", key))
			c.logger.Debug(err)
			break
		}

		// Upload the file to ImageKit
		wg.Add(1)
		imageName := fmt.Sprintf("user_profile_image_%d", slotId)
		go utils.UploadProfileImageAsync(&wg, context, ch, c.ik.ImageKit, src, slotId, imageName, userID)
	}

	wg.Wait()
	close(ch)

	var request models.ProfileUpdateRequest
	for imageUploadResult := range ch {
		v := reflect.ValueOf(&request)
		f_id := v.Elem().FieldByName(fmt.Sprintf("ImageId%d", imageUploadResult.SlotId))
		f_url := v.Elem().FieldByName(fmt.Sprintf("ImageUrl%d", imageUploadResult.SlotId))
		if f_id.IsValid() {
			f_id.Set(reflect.ValueOf(imageUploadResult.FileId))
		}

		if f_url.IsValid() {
			f_url.Set(reflect.ValueOf(imageUploadResult.FileUrl))
		}
	}

	if err = c.service.UpdateProfileById(userID, request); err != nil {
		ctx.JSON(http.StatusConflict, models.HTTPResponse{
			Message: err.Error(),
		})
		return
	} else {

		if len(invalidFields) > 0 {
			ctx.JSON(http.StatusOK, models.HTTPResponse{
				Message:       "Upload profile image fail",
				InvalidFields: invalidFields,
			})
		} else {
			ctx.JSON(http.StatusOK, models.HTTPResponse{
				Message: "success",
			})
		}
	}
}

func (c *ProfileController) RemoveImage(ctx *gin.Context) {
	var deleteRequest models.ProfileImageDeleteRequest
	if err := ctx.ShouldBindJSON(&deleteRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
			Message: err.Error(),
		})
		return
	}

	if len(deleteRequest.Slots) <= 0 {
		ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
			Message: "image slot list empty",
		})
		return
	}

	context, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	userID, err := utils.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.HTTPResponse{
			Message: "server error",
		})
		return
	}

	profile, err := c.service.GetProfileById(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.HTTPResponse{
			Message: err.Error(),
		})
	}

	var invalidFields []string

	for _, slotId := range deleteRequest.Slots {
		v := reflect.ValueOf(profile)
		f := v.Elem().FieldByName(fmt.Sprintf("ImageId%d", slotId))
		if !f.IsValid() || f.String() == "" {
			invalidFields = append(invalidFields, fmt.Sprintf("Image slot %d is empty", slotId))
			c.logger.Debugf("Invalid ImageId%d: %s", slotId, f.String())
		} else {
			c.logger.Debugf("Valid ImageId%d: %s", slotId, f.String())
		}
	}

	if len(invalidFields) > 0 {
		ctx.JSON(http.StatusOK, models.HTTPResponse{
			Message:       "Remove profile image fail",
			InvalidFields: invalidFields,
		})
		return
	}

	var wg sync.WaitGroup
	for _, slotId := range deleteRequest.Slots {
		wg.Add(1)
		v := reflect.ValueOf(profile)
		f := v.Elem().FieldByName(fmt.Sprintf("ImageId%d", slotId))
		go utils.RemoveProfileImageAsync(&wg, context, c.ik.ImageKit, f.String())
		c.logger.Debugf("Run remove profile image!")
	}
	wg.Wait()

	if err = c.service.UpdateProfileImageById(userID, deleteRequest.Slots); err != nil {
		ctx.JSON(http.StatusConflict, models.HTTPResponse{
			Message: err.Error(),
		})
		return
	} else {
		ctx.JSON(http.StatusOK, models.HTTPResponse{
			Message: "success",
		})
	}
}

func (c *ProfileController) DeleteProfile(ctx *gin.Context) {
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.HTTPResponse{
			Message: "server error",
		})
		return
	}

	err = c.service.DeleteProfileById(userID)
	if err != nil {
		ctx.JSON(http.StatusOK, models.HTTPResponse{
			Message: err.Error(),
		})
		return
	} else {
		ctx.JSON(http.StatusOK, models.HTTPResponse{
			Message: "success",
		})
	}
}
