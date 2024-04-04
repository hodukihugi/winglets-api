package commands

import (
	"github.com/hodukihugi/winglets-api/api/middlewares"
	"github.com/hodukihugi/winglets-api/api/routers"
	"github.com/hodukihugi/winglets-api/core"
	"github.com/spf13/cobra"
)

// ServeCommand test command
type ServeCommand struct{}

func (s *ServeCommand) Short() string {
	return "serve application"
}

func (s *ServeCommand) Setup(cmd *cobra.Command) {}

func (s *ServeCommand) Run() core.CommandRunner {
	return func(
		middleware middlewares.Middlewares,
		env *core.Env,
		router *core.RequestHandler,
		route routers.Routers,
		logger *core.Logger,
		database *core.Database,
	) {
		middleware.Setup()
		route.Setup()

		logger.Info("Running server")
		if env.ServerPort == "" {
			_ = router.Gin.Run()
		} else {
			_ = router.Gin.Run(":" + env.ServerPort)
		}
	}
}

func NewServeCommand() *ServeCommand {
	return &ServeCommand{}
}
