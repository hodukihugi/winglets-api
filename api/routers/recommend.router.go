package routers

import (
	"github.com/hodukihugi/winglets-api/api/controllers"
	"github.com/hodukihugi/winglets-api/api/middlewares"
	"github.com/hodukihugi/winglets-api/core"
)

// RecommendRouter struct
type RecommendRouter struct {
	handler             *core.RequestHandler
	recommendController *controllers.RecommendController
	authMiddleware      *middlewares.JWTMiddleware
}

func (r *RecommendRouter) Setup() {
	api := r.handler.Gin.Group("/api").Use(r.authMiddleware.Handler())
	{
		api.POST("/answer", r.recommendController.Answer)
		api.GET("/get-user-matches", r.recommendController.GetUserMatches)
		api.GET("/get-recommendations", r.recommendController.GetRecommendations)
		api.POST("/smash/:id", r.recommendController.Smash)
		api.POST("/pass/:id", r.recommendController.Pass)
	}
}

func NewRecommendRouter(
	handler *core.RequestHandler,
	recommendController *controllers.RecommendController,
	authMiddleware *middlewares.JWTMiddleware,
) *RecommendRouter {
	return &RecommendRouter{
		handler:             handler,
		recommendController: recommendController,
		authMiddleware:      authMiddleware,
	}
}
