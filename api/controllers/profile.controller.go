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
	service services.IProfileService
	ik      *core.ImageKit
	logger  *core.Logger
}

func NewProfileController(service services.IProfileService, ik *core.ImageKit, logger *core.Logger) *ProfileController {
	return &ProfileController{
		service: service,
		ik:      ik,
		logger:  logger,
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
		return
	}

	result, err := c.service.GetProfileById(userID)

	if err != nil {
		ctx.JSON(http.StatusNotFound, models.HTTPResponse{
			Message: "profile not found",
		})
		return
	}

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
		return
	}

	if err := c.service.UpdateProfileById(userID, request); err != nil {
		ctx.JSON(http.StatusConflict, models.HTTPResponse{
			Message: err.Error(),
		})
		return
	} else {
		ctx.JSON(http.StatusOK, models.HTTPResponse{
			Message: "Update profile by id successfully",
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
			continue
		}

		// Upload the file to ImageKit
		wg.Add(1)
		imageName := fmt.Sprintf("user_profile_image_%d", slotId)
		go utils.UploadProfileImageAsync(&wg, context, ch, c.ik.ImageKit, src, slotId, imageName, userID)
	}

	wg.Wait()
	c.logger.Debug("Done upload image!")
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
				Message: "Upload profile image successfully",
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

	var wg sync.WaitGroup
	var invalidFields []string

	for slotId := range deleteRequest.Slots {
		v := reflect.ValueOf(&profile)
		f := v.Elem().FieldByName(fmt.Sprintf("Image%d", slotId))
		if f.IsValid() && f.String() != "" {
			wg.Add(1)
			go utils.RemoveProfileImageAsync(&wg, context, c.ik.ImageKit, userID)
		} else {
			invalidFields = append(invalidFields, fmt.Sprintf("Image slot is empty: %d", slotId))
		}
	}

	wg.Wait()
	c.logger.Debug("Done remove images!")

	if len(invalidFields) > 0 {
		ctx.JSON(http.StatusOK, models.HTTPResponse{
			Message:       "Remove profile image fail",
			InvalidFields: invalidFields,
		})
	}

	var request models.ProfileUpdateRequest
	for slotId := range deleteRequest.Slots {
		v := reflect.ValueOf(&request)
		f := v.Elem().FieldByName(fmt.Sprintf("Image%d", slotId))
		if f.IsValid() {
			f.Set(reflect.ValueOf(nil))
		}
	}

	if err = c.service.UpdateProfileById(userID, request); err != nil {
		ctx.JSON(http.StatusConflict, models.HTTPResponse{
			Message: err.Error(),
		})
		return
	} else {
		ctx.JSON(http.StatusOK, models.HTTPResponse{
			Message: "Upload profile image successfully",
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
		ctx.JSON(http.StatusNotFound, models.HTTPResponse{
			Message: err.Error(),
		})
		return
	} else {
		ctx.JSON(http.StatusOK, models.HTTPResponse{
			Message: "Delete profile successfully",
		})
	}
}

// ====================== private ======================

func (c *ProfileController) getAvailableSlotId(profile models.Profile) int {
	if profile.ID == "" {
		return -1
	}

	// Define an array of image field names
	imageFields := []string{"Image1", "Image2", "Image3", "Image4", "Image5"}

	// Iterate over image fields
	for i, field := range imageFields {
		// Check if the image field is empty
		if reflect.ValueOf(profile).FieldByName(field).String() == "" {
			return i + 1
		}
	}

	return -1
}

func (c *ProfileController) getListAvailableSlotId(profile models.Profile) []int {
	if profile.ID == "" {
		return nil
	}

	// Define an array of image field names
	var result []int
	imageFields := []string{"Image1", "Image2", "Image3", "Image4", "Image5"}

	// Iterate over image fields
	for i, field := range imageFields {
		// Check if the image field is empty
		if reflect.ValueOf(profile).FieldByName(field).String() == "" {
			result = append(result, i+1)
		}
	}

	if len(result) == 0 {
		return nil
	}
	return result
}
