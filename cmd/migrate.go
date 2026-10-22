package cmd

import (
	"fmt"

	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"

	"go-starter-kit/configs/db"
	"go-starter-kit/internal/di"
)

// migrationsDir is resolved relative to the working directory, so migrate is
// run from the project root.
const migrationsDir = "migrations"

var migrateCmd = &cobra.Command{
	Use:       "migrate [up|down|status]",
	Short:     "run goose database migrations",
	ValidArgs: []string{"up", "down", "status"},
	Args:      cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		container := di.GetContainer()
		if err := di.RegisterDependencies(container); err != nil {
			return err
		}

		return container.Invoke(func(database *db.Database) error {
			sqlDB, err := database.DB.DB()
			if err != nil {
				return err
			}
			if err := goose.SetDialect("postgres"); err != nil {
				return err
			}

			switch args[0] {
			case "up":
				return goose.Up(sqlDB, migrationsDir)
			case "down":
				return goose.Down(sqlDB, migrationsDir)
			case "status":
				return goose.Status(sqlDB, migrationsDir)
			default:
				return fmt.Errorf("unknown migrate subcommand %q (use up|down|status)", args[0])
			}
		})
	},
}
