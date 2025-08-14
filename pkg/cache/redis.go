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
func (r *redisCache) Get(key string) (interface{}, error) {
	ctx := context.Background()
	val, err := r.redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // Key does not exist
	} else if err != nil {
		return nil, err // Other error
	}
	return val, nil
}
func (r *redisCache) GetWithContext(ctx context.Context, key string) (interface{}, error) {
	val, err := r.redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // Key does not exist
	} else if err != nil {
		return nil, err // Other error
	}
	return val, nil
}

func (r *redisCache) Set(key string, value interface{}) error {
	ctx := context.Background()
	return r.redisClient.Set(ctx, key, value, 0).Err()
}

func (r *redisCache) SetWithContext(ctx context.Context, key string, value interface{}) error {
	return r.redisClient.Set(ctx, key, value, 0).Err()
}

func (r *redisCache) SetTTL(key string, value interface{}, ttl time.Duration) error {
	ctx := context.Background()
	return r.redisClient.Set(ctx, key, value, ttl).Err()
}

func (r *redisCache) SetTTLWithContext(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return r.redisClient.Set(ctx, key, value, ttl).Err()
}

func (r *redisCache) Delete(key string) error {
	ctx := context.Background()
	return r.redisClient.Del(ctx, key).Err()
}

func (r *redisCache) DeleteWithContext(ctx context.Context, key string) error {
	return r.redisClient.Del(ctx, key).Err()
}

func (r *redisCache) Clear(ctx context.Context) error {
	return r.redisClient.FlushDB(ctx).Err()
}

func (r *redisCache) ClearWithContext(ctx context.Context) error {
	return r.redisClient.FlushDB(ctx).Err()
}

func (r *redisCache) Close() error {
	return r.redisClient.Close()
}
