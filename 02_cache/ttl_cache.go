package cache

import (
	"sync"
	"time"
)

type TTLCache struct {
	keyToItem sync.Map
}

func NewTTLCache() *TTLCache {
	return &TTLCache{}
}

func (c *TTLCache) Put(key string, value interface{}, expirationTime time.Time) {
	c.keyToItem.Store(key, ttlCacheItem{
		value:          value,
		expirationTime: expirationTime,
	})
}

func (c *TTLCache) Get(key string) (any, bool) {
	value, ok := c.keyToItem.Load(key)
	if !ok {
		return nil, false
	}

	item := value.(ttlCacheItem)
	if item.expirationTime.Before(time.Now()) {
		return nil, false
	}
	return item.value, true
}

func (c *TTLCache) DeleteExpired() int {
	count := 0

	c.keyToItem.Range(func(key, value any) bool {
		if value.(ttlCacheItem).expirationTime.Before(time.Now()) {
			c.keyToItem.Delete(key)
			count += 1
		}
		return true
	})

	return count
}

type ttlCacheItem struct {
	value          any
	expirationTime time.Time
}
