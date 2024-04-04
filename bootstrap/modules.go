package bootstrap

import (
	"github.com/hodukihugi/winglets-api/api/controllers"
	"github.com/hodukihugi/winglets-api/api/middlewares"
	"github.com/hodukihugi/winglets-api/api/routers"
	"github.com/hodukihugi/winglets-api/core"
	"github.com/hodukihugi/winglets-api/repositories"
	"github.com/hodukihugi/winglets-api/services"
	"go.uber.org/fx"
)

var CommonModules = fx.Options(
	controllers.Module,
	routers.Module,
	core.Module,
	services.Module,
	middlewares.Module,
	repositories.Module,
)
