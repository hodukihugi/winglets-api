package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/hodukihugi/winglets-api/constants"
	"github.com/hodukihugi/winglets-api/core"
	"github.com/hodukihugi/winglets-api/models"
)

func AttachCookiesToResponse(env *core.Env, accessTokenJWT, refreshTokenJWT string, c *gin.Context) {
	// secure is set to false if development is local
	isSecure := true
	if env.Environment == "development" {
		isSecure = false
	}

	accessCookie := accessTokenJWT
	c.SetCookie("accessCookie", accessCookie, int(env.AccessTokenExpiresIn), "/", "localhost", isSecure, true)

	refreshCookie := refreshTokenJWT
	c.SetCookie("refreshCookie", refreshCookie, int(env.RefreshTokenExpiresIn), "/", "localhost", isSecure, true)
}

func GetUserID(ctx *gin.Context) (string, error) {
	payload, ok := ctx.Get(constants.CtxKey_JWTClaim)
	if !ok {
		return "", nil
	}

	var jwtClaim = payload.(*models.JWTClaim)
	return jwtClaim.UserID, nil
}
