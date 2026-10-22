package shared

import (
	"go-starter-kit/configs"
	"go-starter-kit/configs/db"
	"go-starter-kit/internal/shared/cache"
	"go-starter-kit/internal/shared/log"

	"go.uber.org/dig"
)

type Dependency struct {
	dig.In

	Logger log.Logger
	Config *configs.Config
	ORM    *db.Database
	Cache  cache.Cache
}

func (d Dependency) Close() error {
	if d.ORM != nil {
		if sql, err := d.ORM.DB.DB(); err == nil {
			_ = sql.Close()
		}
	}
	if d.Cache != nil {
		_ = d.Cache.Close()
	}
	return nil
}
