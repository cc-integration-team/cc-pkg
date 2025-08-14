package cache

import (
	"context"
	"errors"
	"sync"
	"time"
)

type memItem struct {
	value      interface{}
	expireTime time.Time
}

type memoryCache struct {
	store map[string]memItem
	mu    sync.RWMutex
}

func NewMemoryCache() Cache {
	return &memoryCache{
		store: make(map[string]memItem),
	}
}

func (m *memoryCache) Get(ctx context.Context, key string) (interface{}, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	item, ok := m.store[key]
	if !ok {
		return nil, errors.New("key not found")
	}

	// Check TTL
	if !item.expireTime.IsZero() && time.Now().After(item.expireTime) {
		return nil, errors.New("key expired")
	}

	return item.value, nil
}

func (m *memoryCache) Set(ctx context.Context, key string, value interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.store[key] = memItem{
		value:      value,
		expireTime: time.Time{},
	}
	return nil
}

func (m *memoryCache) SetTTL(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.store[key] = memItem{
		value:      value,
		expireTime: time.Now().Add(ttl),
	}
	return nil
}

func (m *memoryCache) Delete(ctx context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.store, key)
	return nil
}

func (m *memoryCache) Clear(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.store = make(map[string]memItem)
	return nil
}

func (m *memoryCache) Close() error {
	return nil
}
