package middlewares

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/hodukihugi/winglets-api/constants"
	"github.com/hodukihugi/winglets-api/models"
	"github.com/hodukihugi/winglets-api/services"
	"github.com/hodukihugi/winglets-api/utils"
	"net/http"
	"strings"

	"github.com/hodukihugi/winglets-api/core"
)

// JWTMiddleware middleware for jwt authentication
type JWTMiddleware struct {
	env     *core.Env
	service services.IAuthService
	logger  *core.Logger
}

// NewJWTMiddleware creates new jwt auth middleware
func NewJWTMiddleware(
	env *core.Env,
	logger *core.Logger,
	service services.IAuthService,
) *JWTMiddleware {
	return &JWTMiddleware{
		env:     env,
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

func (m *JWTMiddleware) AuthorizationWithCookie() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, err := c.Cookie("accessCookie")
		if err != nil {
			refreshToken, err := c.Cookie("refreshCookie")
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
				c.Abort()
				return
			}

			payload, err := m.service.Authorize(refreshToken)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token expired"})
				c.Abort()
				return
			}

			accessToken, _, err := m.service.Refresh(models.User{
				Email: payload.UserEmail,
				ID:    payload.UserID,
			})

			utils.AttachCookiesToResponse(m.env, accessToken, refreshToken, c)
			c.Set(constants.CtxKey_JWTClaim, payload)
			c.Next()
		}

		payload, err := m.service.Authorize(accessToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Access token expired"})
			c.Abort()
			return
		}

		c.Set(constants.CtxKey_JWTClaim, payload)
		c.Next()
	}
}
