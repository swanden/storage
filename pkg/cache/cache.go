package cache

import (
	"sync"
	"time"
)

type Item struct {
	putTime time.Time
	ttl     time.Duration
	value   string
}

type Cache struct {
	mu   sync.RWMutex
	data map[string]Item
}

func New() *Cache {
	c := &Cache{}
	c.data = make(map[string]Item)

	return c
}

func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	item, ok := c.data[key]
	defer c.mu.RUnlock()

	now := time.Now()
	if item.ttl > 0 && now.Sub(item.putTime) > item.ttl {
		return "", false
	}

	return item.value, ok
}

// Set sets key-value pair
// ttl - expiration time, if 0 - no expire time
func (c *Cache) Set(key string, value string, ttl time.Duration) {
	c.mu.Lock()
	c.data[key] = Item{putTime: time.Now(), value: value, ttl: ttl}
	c.mu.Unlock()
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	delete(c.data, key)
	c.mu.Unlock()
}
