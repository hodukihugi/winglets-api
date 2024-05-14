package routers

import (
	"github.com/hodukihugi/winglets-api/api/controllers"
	"github.com/hodukihugi/winglets-api/api/middlewares"
	"github.com/hodukihugi/winglets-api/core"
)

// AuthRouter struct
type AuthRouter struct {
	logger         *core.Logger
	handler        *core.RequestHandler
	authController *controllers.AuthController
	authMiddleware *middlewares.JWTMiddleware
	corsMiddleware *middlewares.CorsMiddleware
}

// Setup user routes
func (s *AuthRouter) Setup() {
	s.logger.Info("Setting up routes")
	auth := s.handler.Gin.Group("/api/auth").Use()
	{
		auth.POST("/login", s.authController.SignIn)
		auth.POST("/register", s.authController.Register)
		auth.POST("/verify-email", s.authController.VerifyEmail)
		auth.POST("/send-verification-email", s.authController.SendVerificationEmail)
		auth.POST("/refresh", s.authMiddleware.Handler(), s.authController.Refresh)
	}
}

// NewAuthRouter creates new user controller
func NewAuthRouter(
	handler *core.RequestHandler,
	authController *controllers.AuthController,
	authMiddleware *middlewares.JWTMiddleware,
	corsMiddleware *middlewares.CorsMiddleware,
	logger *core.Logger,
) *AuthRouter {
	return &AuthRouter{
		handler:        handler,
		logger:         logger,
		authController: authController,
		authMiddleware: authMiddleware,
		corsMiddleware: corsMiddleware,
	}
}
