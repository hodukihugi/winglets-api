package middlewares

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/hodukihugi/winglets-api/constants"
	"github.com/hodukihugi/winglets-api/models"
	"github.com/hodukihugi/winglets-api/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/hodukihugi/winglets-api/core"
)

// JWTMiddleware middleware for jwt authentication
type JWTMiddleware struct {
	service services.IAuthService
	logger  *core.Logger
}

// NewJWTMiddleware creates new jwt auth middleware
func NewJWTMiddleware(
	logger *core.Logger,
	service services.IAuthService,
) *JWTMiddleware {
	return &JWTMiddleware{
		service: service,
		logger:  logger,
	}
}

// Setup sets up jwt auth middleware
func (m *JWTMiddleware) Setup() {}

// Handler handles middleware functionality
func (m *JWTMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		t := strings.Fields(authHeader)
		if len(t) == 2 && t[0] == "Bearer" {
			authToken := t[1]
			claim, err := m.service.Authorize(authToken)
			if err == nil {
				c.Set(constants.CtxKey_JWTClaim, claim)
				c.Next()
				return
			}
			var ve *jwt.ValidationError
			if errors.As(err, &ve) {
				if ve.Errors&jwt.ValidationErrorMalformed != 0 {
					c.JSON(http.StatusUnauthorized, models.HTTPResponse{
						Message: "token malformed",
					})
					c.Abort()
					return
				}
				if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
					c.JSON(http.StatusUnauthorized, models.HTTPResponse{
						Message: "token expired",
					})
					c.Abort()
					return
				}
			}
			c.JSON(http.StatusUnauthorized, models.HTTPResponse{
				Message: "fail to authorize",
			})
			m.logger.Errorf("fail to authorize: [%v]", err)
			c.Abort()
			return
		}
		c.JSON(http.StatusUnauthorized, models.HTTPResponse{
			Message: "you are not authorized",
		})
		c.Abort()
		return
	}
}
