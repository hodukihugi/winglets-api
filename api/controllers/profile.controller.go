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

	// Parse the multipart form
	err = ctx.Request.ParseMultipartForm(10 << 20) // Max file size of 10MB
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
			Message: "bad request",
		})
	}

	// Get the files from the form
	images := ctx.Request.MultipartForm.File["images"]

	// Get profile
	result, err := c.service.GetProfileById(userID)
	if err != nil {
		c.logger.Debug(err)
		ctx.JSON(http.StatusInternalServerError, models.HTTPResponse{
			Message: "server error",
		})
		return
	}

	// Get profile's available image slot
	availableSlotIds := c.getListAvailableSlotId(*result)
	c.logger.Debug(availableSlotIds)
	slotIndex := 0

	if availableSlotIds == nil {
		ctx.JSON(http.StatusConflict, models.HTTPResponse{
			Message: "user has full image",
		})
		return
	}

	var wg sync.WaitGroup
	ch := make(chan string, len(availableSlotIds))

	for _, fileHeader := range images {
		// No more available slot
		if slotIndex >= len(availableSlotIds) {
			break
		}

		// Open the uploaded file
		src, err := fileHeader.Open()
		defer src.Close()

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, models.HTTPResponse{
				Message: "Failed to open file",
			})
			continue
		}

		// Upload the file to ImageKit
		wg.Add(1)
		imageName := fmt.Sprintf("user_profile_image_%d", availableSlotIds[slotIndex])
		go utils.UploadImageAsync(&wg, context, ch, c.ik.ImageKit, src, imageName, userID)
		slotIndex++
		c.logger.Debugf("\nSlot index: %d", slotIndex)
	}
	wg.Wait()
	c.logger.Debug("Done upload image!")
	close(ch)

	var request models.ProfileUpdateRequest
	slotIndex = 0
	for url := range ch {
		// No more available slot
		if slotIndex >= len(availableSlotIds) {
			break
		}

		v := reflect.ValueOf(&request)
		f := v.Elem().FieldByName(fmt.Sprintf("Image%d", availableSlotIds[slotIndex]))
		if f.IsValid() {
			f.Set(reflect.ValueOf(url))
		}
		slotIndex++
	}
	c.logger.Debugf("Image request value %v: ", request)

	if err = c.service.UpdateProfileById(userID, request); err != nil {
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
