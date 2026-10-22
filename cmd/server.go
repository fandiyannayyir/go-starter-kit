package cmd

import (
	"github.com/spf13/cobra"

	"go-starter-kit/internal/di"
	"go-starter-kit/internal/handler"
	"go-starter-kit/internal/repository"
	"go-starter-kit/internal/server"
	"go-starter-kit/internal/service"
	"go-starter-kit/internal/shared"
	"go-starter-kit/pkg/route"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "run the HTTP server",
	RunE: func(_ *cobra.Command, _ []string) error {
		container := di.GetContainer()
		if err := di.RegisterDependencies(container); err != nil {
			return err
		}

		return container.Invoke(func(dep shared.Dependency) error {
			defer dep.Close()

			e := server.New(dep.Logger)

			exampleRepo := repository.NewExampleRepository(dep.ORM.DB)
			exampleSvc := service.NewExampleService(exampleRepo)

			reg := route.NewRegistry()
			reg.Add(handler.NewHealthHandler(dep.ORM, dep.Cache).Routes()...)
			reg.Add(handler.NewExampleHandler(exampleSvc).Routes()...)
			if err := reg.Register(e); err != nil {
				return err
			}

			addr := dep.Config.AppConfig.Host + ":" + dep.Config.AppConfig.Port
			dep.Logger.Info("starting server on " + addr)
			return server.Start(e, addr, dep.Logger)
		})
	},
}
