package cache

import (
	"context"
	"time"
)

type (
	Cache interface {
		Set(ctx context.Context, key string, value interface{}) error
		SetNX(ctx context.Context, key string, value interface{}, exp time.Duration) (bool, error)
		SetWithExp(ctx context.Context, key string, value interface{}, exp time.Duration) error
		Get(ctx context.Context, key string, object interface{}) error
		GetBytes(ctx context.Context, key string) ([]byte, error)
		Delete(ctx context.Context, keys ...string) (int64, error)

		IncrByFloat(ctx context.Context, key string, value float64) (float64, error)
		DecrByFloat(ctx context.Context, key string, value float64) (float64, error)
		Expire(ctx context.Context, key string, exp time.Duration) error

		Hash

		Ping(ctx context.Context) error
		Close() error
	}

	Hash interface {
		HSet(ctx context.Context, key string, field string, value interface{}, ttl ...time.Duration) error
		HGet(ctx context.Context, key string, field string, value interface{}) error
		HDel(ctx context.Context, key string, field string) error
	}
)
