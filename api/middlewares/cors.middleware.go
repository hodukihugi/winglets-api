package middlewares

import (
	"github.com/hodukihugi/winglets-api/core"
	cors "github.com/rs/cors/wrapper/gin"
)

// CorsMiddleware middleware for cors
type CorsMiddleware struct {
	handler *core.RequestHandler
	logger  *core.Logger
	env     *core.Env
}

// NewCorsMiddleware creates new cors middleware
func NewCorsMiddleware(handler *core.RequestHandler, logger *core.Logger, env *core.Env) *CorsMiddleware {
	return &CorsMiddleware{
		handler: handler,
		logger:  logger,
		env:     env,
	}
}

// Setup sets up cors middleware
func (m *CorsMiddleware) Setup() {
	m.logger.Info("Setting up cors middleware")

	if m.env.Environment == "development" {
		return
	}

	debug := m.env.Environment == "development"
	m.handler.Gin.Use(cors.New(cors.Options{
		AllowCredentials: true,
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "HEAD", "OPTIONS", "PATCH", "DELETE"},
		Debug:            debug,
	}))
}
