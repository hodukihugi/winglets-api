package routers

import (
	"github.com/hodukihugi/winglets-api/api/controllers"
	"github.com/hodukihugi/winglets-api/api/middlewares"
	"github.com/hodukihugi/winglets-api/core"
)

// UserRouter struct
type UserRouter struct {
	handler        *core.RequestHandler
	userController *controllers.UserController
	authMiddleware *middlewares.JWTMiddleware
}

// Setup user routes
func (s *UserRouter) Setup() {
	api := s.handler.Gin.Group("/api").Use(s.authMiddleware.AuthorizationWithCookie())
	{
		api.GET("/users/:id", s.userController.GetByID)
	}
}

// NewUserRouter creates new user controller
func NewUserRouter(
	handler *core.RequestHandler,
	userController *controllers.UserController,
	authMiddleware *middlewares.JWTMiddleware,
) *UserRouter {
	return &UserRouter{
		handler:        handler,
		userController: userController,
		authMiddleware: authMiddleware,
	}
}
