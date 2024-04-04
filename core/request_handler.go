package core

import (
	"github.com/gin-gonic/gin"
)

// RequestHandler function
type RequestHandler struct {
	Gin *gin.Engine
}

// NewRequestHandler creates a new request handler
func NewRequestHandler(logger *Logger) *RequestHandler {
	gin.DefaultWriter = logger.GetGinLogger()
	engine := gin.New()
	engine.Use(gin.Logger())
	return &RequestHandler{Gin: engine}
}
