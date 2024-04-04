package bootstrap

import (
	"github.com/hodukihugi/winglets-api/commands"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "Innov8tor",
	Short: "Innov8tor",
	Long: `
This is a command runner or cli for api architecture in golang. 
Using this we can use underlying dependency injection container for running scripts. 
Main advantage is that, we can use same services, repositories, infrastructure present in the application itself`,
	TraverseChildren: true,
}

// App root of application
type App struct {
	*cobra.Command
}

func NewApp() *App {
	cmd := App{
		Command: rootCmd,
	}
	cmd.AddCommand(commands.GetSubCommands(CommonModules)...)
	return &cmd
}

var RootApp = NewApp()
