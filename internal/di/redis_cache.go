package di

import (
	"fmt"

	"go-starter-kit/configs"
	"go-starter-kit/internal/shared/cache"
	redisCache "go-starter-kit/internal/shared/cache/redis"

	"github.com/redis/go-redis/v9"
)

func NewRedisCache(cfg *configs.Config) cache.Cache {
	addr := fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port)
	redisClient := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    []string{addr},
		Username: cfg.Redis.User,
		Password: cfg.Redis.Pass,
		DB:       cfg.Redis.Index,
	})
	return redisCache.NewCache(redisClient)
}
