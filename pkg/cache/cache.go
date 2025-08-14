package cache

import (
	"context"
	"time"
)

type Cache interface {
	Get(key string) (interface{}, error)
	GetWithContext(ctx context.Context, key string) (interface{}, error)

	Set(key string, value interface{}) error
	SetWithContext(ctx context.Context, key string, value interface{}) error

	SetTTL(key string, value interface{}, ttl time.Duration) error
	SetTTLWithContext(ctx context.Context, key string, value interface{}, ttl time.Duration) error

	Delete(key string) error
	DeleteWithContext(ctx context.Context, key string) error

	Clear(ctx context.Context) error
	ClearWithContext(ctx context.Context) error

	Close() error
}
