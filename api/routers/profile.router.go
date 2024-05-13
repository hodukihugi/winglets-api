package routers

import (
	"github.com/hodukihugi/winglets-api/api/controllers"
	"github.com/hodukihugi/winglets-api/api/middlewares"
	"github.com/hodukihugi/winglets-api/core"
)

// ProfileRouter struct
type ProfileRouter struct {
	handler           *core.RequestHandler
	profileController *controllers.ProfileController
	authMiddleware    *middlewares.JWTMiddleware
}

func (r *ProfileRouter) Setup() {
	api := r.handler.Gin.Group("/api").Use(r.authMiddleware.Handler())
	{
		api.GET("/profile/:id", r.profileController.GetProfileById)
		api.GET("/profile", r.profileController.GetMyProfile)
		api.POST("/profile", r.profileController.CreateProfile)
		api.POST("/profile/upload-image", r.profileController.UploadImage)
		api.DELETE("/profile/remove-image", r.profileController.RemoveImage)
		api.PUT("/profile", r.profileController.UpdateProfile)
		api.DELETE("/profile", r.profileController.DeleteProfile)
	}
}

func NewProfileRouter(
	handler *core.RequestHandler,
	profileController *controllers.ProfileController,
	authMiddleware *middlewares.JWTMiddleware,
) *ProfileRouter {
	return &ProfileRouter{
		handler:           handler,
		profileController: profileController,
		authMiddleware:    authMiddleware,
	}
}
