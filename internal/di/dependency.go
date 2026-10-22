package di

import (
	"sync"

	"go-starter-kit/configs"
	"go-starter-kit/configs/db"
	"go-starter-kit/internal/shared"

	"go.uber.org/dig"
	"gorm.io/gorm"
)

var (
	Container *dig.Container
	once      sync.Once
)

func GetContainer() *dig.Container {
	once.Do(func() {
		Container = dig.New()
	})
	return Container
}

func RegisterDependency(container *dig.Container) error {
	if err := container.Provide(configs.GetConfig); err != nil {
		return err
	}
	if err := container.Provide(NewLogger); err != nil {
		return err
	}
	if err := container.Provide(NewDatabase); err != nil {
		return err
	}
	if err := container.Provide(NewRedisCache); err != nil {
		return err
	}
	if err := container.Provide(func(database *db.Database) shared.TxManager {
		return shared.NewTxManager(database.DB)
	}); err != nil {
		return err
	}
	if err := container.Provide(func(database *db.Database) *gorm.DB {
		return database.DB
	}); err != nil {
		return err
	}
	return nil
}
