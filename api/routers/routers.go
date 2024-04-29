package routers

import (
	"go.uber.org/fx"
)

// Module exports dependency to container
var Module = fx.Options(
	fx.Provide(NewUserRouter),
	fx.Provide(NewAuthRouter),
	fx.Provide(NewProfileRouter),
	fx.Provide(NewRecommendRouter),
	fx.Provide(NewRouters),
)

// Routers contains multiple routes
type Routers []Router

// Router interface
type Router interface {
	Setup()
}

// NewRouters sets up routes
func NewRouters(
	userRouter *UserRouter,
	authRouter *AuthRouter,
	profileRouter *ProfileRouter,
	recommendRouter *RecommendRouter,
) Routers {
	return Routers{
		userRouter,
		authRouter,
		profileRouter,
		recommendRouter,
	}
}

// Setup all the route
func (r Routers) Setup() {
	for _, route := range r {
		route.Setup()
	}
}
