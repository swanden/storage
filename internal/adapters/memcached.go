package adapters

import (
	"context"
	"github.com/pkg/errors"
	"github.com/swanden/storage/pkg/memcached"
	"time"
)

type MemcachedAdapter struct {
	client *memcached.Client
}

func NewMemcachedAdapter(client *memcached.Client) *MemcachedAdapter {
	return &MemcachedAdapter{
		client: client,
	}
}

func (ma *MemcachedAdapter) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	return ma.client.Set(ctx, key, value, int(ttl.Seconds()))
}

func (ma *MemcachedAdapter) Get(ctx context.Context, key string) (string, error) {
	value, err := ma.client.Get(ctx, key)
	if errors.Is(err, memcached.ErrNotFound) {
		return value, ErrNotFound
	}

	return value, err
}

func (ma *MemcachedAdapter) Delete(ctx context.Context, key string) error {
	return ma.client.Delete(ctx, key)
}

func (ma *MemcachedAdapter) Close() {
	ma.client.Close()
}
