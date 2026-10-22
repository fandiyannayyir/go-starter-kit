package redis

import (
	"context"
	"encoding/json"
	"time"

	"go-starter-kit/internal/shared/cache"

	"github.com/redis/go-redis/v9"
)

type redisCache struct {
	client redis.UniversalClient
}

func NewCache(client redis.UniversalClient) cache.Cache {
	return &redisCache{client: client}
}

func (r *redisCache) Set(ctx context.Context, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, 0).Err()
}

func (r *redisCache) SetNX(ctx context.Context, key string, value interface{}, exp time.Duration) (bool, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return false, err
	}
	return r.client.SetNX(ctx, key, data, exp).Result()
}

func (r *redisCache) SetWithExp(ctx context.Context, key string, value interface{}, exp time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, exp).Err()
}

func (r *redisCache) Get(ctx context.Context, key string, object interface{}) error {
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, object)
}

func (r *redisCache) GetBytes(ctx context.Context, key string) ([]byte, error) {
	return r.client.Get(ctx, key).Bytes()
}

func (r *redisCache) Delete(ctx context.Context, keys ...string) (int64, error) {
	return r.client.Del(ctx, keys...).Result()
}

func (r *redisCache) HSet(ctx context.Context, key string, field string, value interface{}, ttl ...time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	if err := r.client.HSet(ctx, key, field, data).Err(); err != nil {
		return err
	}
	if len(ttl) > 0 && ttl[0] > 0 {
		r.client.Expire(ctx, key, ttl[0])
	}
	return nil
}

func (r *redisCache) HGet(ctx context.Context, key string, field string, value interface{}) error {
	data, err := r.client.HGet(ctx, key, field).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, value)
}

func (r *redisCache) HDel(ctx context.Context, key string, field string) error {
	return r.client.HDel(ctx, key, field).Err()
}

func (r *redisCache) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

func (r *redisCache) Close() error {
	return r.client.Close()
}

func (r *redisCache) IncrByFloat(ctx context.Context, key string, value float64) (float64, error) {
	return r.client.IncrByFloat(ctx, key, value).Result()
}

func (r *redisCache) DecrByFloat(ctx context.Context, key string, value float64) (float64, error) {
	return r.client.IncrByFloat(ctx, key, -value).Result()
}

func (r *redisCache) Expire(ctx context.Context, key string, exp time.Duration) error {
	return r.client.Expire(ctx, key, exp).Err()
}
