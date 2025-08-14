package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisCache struct {
	redisClient *redis.Client
}

func NewRedisCache(client *redis.Client) Cache {
	return &redisCache{
		redisClient: client,
	}
}
func (r *redisCache) Get(ctx context.Context, key string) (interface{}, error) {
	val, err := r.redisClient.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return val, nil
}
func (r *redisCache) Set(ctx context.Context, key string, value interface{}) error {
	return r.redisClient.Set(ctx, key, value, 0).Err()
}

func (r *redisCache) SetTTL(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return r.redisClient.Set(ctx, key, value, ttl).Err()
}

func (r *redisCache) Delete(ctx context.Context, key string) error {
	return r.redisClient.Del(ctx, key).Err()
}
func (r *redisCache) Clear(ctx context.Context) error {
	return r.redisClient.FlushDB(ctx).Err()
}
func (r *redisCache) Close() error {
	return r.redisClient.Close()
}
