package adapters

import (
	"context"
	"github.com/swanden/storage/pkg/cache"
	"time"
)

type CacheAdapter struct {
	cache *cache.Cache
}

func NewCacheAdapter(cache *cache.Cache) *CacheAdapter {
	return &CacheAdapter{
		cache: cache,
	}
}

func (ca *CacheAdapter) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	ca.cache.Set(key, value, ttl)

	return nil
}

func (ca *CacheAdapter) Get(ctx context.Context, key string) (string, error) {
	value, ok := ca.cache.Get(key)
	if !ok {
		return value, ErrNotFound
	}

	return value, nil
}

func (ca *CacheAdapter) Delete(ctx context.Context, key string) error {
	ca.cache.Delete(key)

	return nil
}

func (ca *CacheAdapter) Close() {
}
