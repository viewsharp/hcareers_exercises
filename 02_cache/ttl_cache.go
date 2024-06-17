package cache

import (
	"container/heap"
	"time"
)

type TTLCache struct {
	keyToItem       map[string]*ttlCacheItem
	expirationQueue expirationQueue
}

func NewTTLCache() *TTLCache {
	return &TTLCache{
		keyToItem: make(map[string]*ttlCacheItem),
	}
}

func (c *TTLCache) Put(key string, value interface{}, expirationTime time.Time) {
	if item, ok := c.keyToItem[key]; ok {
		item.value = value
		item.expirationTime = expirationTime
		heap.Fix(&c.expirationQueue, item.index)
		return
	}

	item := &ttlCacheItem{
		key:            key,
		value:          value,
		expirationTime: expirationTime,
	}
	heap.Push(&c.expirationQueue, item)
	c.keyToItem[key] = item
}

func (c *TTLCache) Get(key string) (any, bool) {
	if item, ok := c.keyToItem[key]; ok && time.Now().Before(item.expirationTime) {
		return item.value, true
	}
	return nil, false
}

func (c *TTLCache) DeleteExpired() int {
	count := 0

	for len(c.expirationQueue) > 0 && c.expirationQueue[0].expirationTime.Before(time.Now()) {
		item := c.expirationQueue[0]
		delete(c.keyToItem, item.key)
		heap.Pop(&c.expirationQueue)
		count += 1
	}

	return count
}

type ttlCacheItem struct {
	key            string
	value          any
	index          int
	expirationTime time.Time
}

type expirationQueue []*ttlCacheItem

func (eq expirationQueue) Len() int { return len(eq) }

func (eq expirationQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return eq[i].expirationTime.Before(eq[j].expirationTime)
}

func (eq expirationQueue) Swap(i, j int) {
	eq[i], eq[j] = eq[j], eq[i]
	eq[i].index = i
	eq[j].index = j
}

func (eq *expirationQueue) Push(x any) {
	n := len(*eq)
	item := x.(*ttlCacheItem)
	item.index = n
	*eq = append(*eq, item)
}

func (eq *expirationQueue) Pop() any {
	old := *eq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*eq = old[0 : n-1]
	return item
}
