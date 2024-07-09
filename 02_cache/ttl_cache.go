package cache

import (
	"time"
)

type TTLCache struct {
	keyToItem map[string]ttlCacheItem
}

func NewTTLCache() *TTLCache {
	return &TTLCache{
		keyToItem: make(map[string]ttlCacheItem),
	}
}

func (c *TTLCache) Put(key string, value interface{}, expirationTime time.Time) {
	c.keyToItem[key] = ttlCacheItem{
		value:          value,
		expirationTime: expirationTime,
	}
}

func (c *TTLCache) Get(key string) (any, bool) {
	if item, ok := c.keyToItem[key]; ok && !item.expirationTime.Before(time.Now()) {
		return item.value, true
	}
	return nil, false
}

func (c *TTLCache) DeleteExpired() int {
	count := 0

	for key, item := range c.keyToItem {
		if item.expirationTime.Before(time.Now()) {
			delete(c.keyToItem, key)
			count += 1
		}
	}

	return count
}

type ttlCacheItem struct {
	value          any
	expirationTime time.Time
}
