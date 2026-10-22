package cmd

import (
	"github.com/spf13/cobra"

	"go-starter-kit/internal/di"
	"go-starter-kit/internal/handler"
	"go-starter-kit/internal/router"
	"go-starter-kit/internal/server"
	"go-starter-kit/internal/shared"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "run the HTTP server",
	RunE: func(_ *cobra.Command, _ []string) error {
		container := di.GetContainer()
		if err := di.RegisterDependencies(container); err != nil {
			return err
		}

		return container.Invoke(func(dep shared.Dependency, health *handler.HealthHandler, example *handler.ExampleHandler) error {
			defer func() { _ = dep.Close() }()

			e := server.New(dep.Logger)

			if err := router.Build(e, health, example); err != nil {
				return err
			}

			addr := dep.Config.AppConfig.Host + ":" + dep.Config.AppConfig.Port
			dep.Logger.Info("starting server on " + addr)
			return server.Start(e, addr, dep.Logger)
		})
	},
}
