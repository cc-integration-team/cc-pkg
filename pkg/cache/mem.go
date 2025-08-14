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

func (m *memoryCache) Get(key string) (interface{}, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	item, exists := m.store[key]
	if !exists {
		return nil, nil // Key does not exist
	}
	if time.Now().After(item.expireTime) {
		delete(m.store, key) // Remove expired item
		return nil, nil
	}
	return item.value, nil
}

func (m *memoryCache) GetWithContext(ctx context.Context, key string) (interface{}, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	item, exists := m.store[key]
	if !exists {
		return nil, nil // Key does not exist
	}
	if time.Now().After(item.expireTime) {
		delete(m.store, key) // Remove expired item
		return nil, nil
	}
	return item.value, nil
}

func (m *memoryCache) Set(key string, value interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.store[key] = memItem{
		value:      value,
		expireTime: time.Time{}, // No expiration
	}
	return nil
}

func (m *memoryCache) SetWithContext(ctx context.Context, key string, value interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.store[key] = memItem{
		value:      value,
		expireTime: time.Time{}, // No expiration
	}
	return nil
}

func (m *memoryCache) SetTTL(key string, value interface{}, ttl time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if ttl <= 0 {
		return errors.New("TTL must be greater than zero")
	}

	m.store[key] = memItem{
		value:      value,
		expireTime: time.Now().Add(ttl),
	}
	return nil
}

func (m *memoryCache) SetTTLWithContext(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if ttl <= 0 {
		return errors.New("TTL must be greater than zero")
	}

	m.store[key] = memItem{
		value:      value,
		expireTime: time.Now().Add(ttl),
	}
	return nil
}

func (m *memoryCache) Delete(key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.store, key)
	return nil
}
func (m *memoryCache) DeleteWithContext(ctx context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.store, key)
	return nil
}
func (m *memoryCache) Clear(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.store = make(map[string]memItem) // Clear the store
	return nil
}
func (m *memoryCache) ClearWithContext(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.store = make(map[string]memItem) // Clear the store
	return nil
}
func (m *memoryCache) Close() error {
	// No resources to close for in-memory cache
	return nil
}
