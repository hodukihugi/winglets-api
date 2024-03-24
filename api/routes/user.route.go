package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hodukihugi/winglets-api/api/controllers"
	"github.com/hodukihugi/winglets-api/repositories"
	"github.com/hodukihugi/winglets-api/services"
	"gorm.io/gorm"
)

func InitUserRoutes(db *gorm.DB, route *gin.Engine) {
	/**
	@description Init user controller
	*/
	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepository)
	userController := controllers.NewUserController(userService)

	/**
	@description All User Route
	*/
	groupRoute := route.Group("/api/v1")
	groupRoute.POST("/users", userController.CreateUser)
	groupRoute.GET("/users", userController.GetListUser)
	groupRoute.GET("/users/:id", userController.GetUserById)
	groupRoute.PUT("/users/:id", userController.UpdateUserById)
	groupRoute.DELETE("/users/:id", userController.DeleteUserById)
}
