package di

import (
	"go-starter-kit/configs"
	"go-starter-kit/configs/db"
)

func NewDatabase(cfg *configs.Config) *db.Database {
	return db.NewPostgreDB(cfg)
}
