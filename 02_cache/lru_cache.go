package cache

import (
	"container/list"
	"errors"
	"sync"
)

type LRUCache struct {
	keyToElement map[string]*list.Element
	list         *list.List
	capacity     int
	mutex        sync.Mutex
}

type lruCacheItem struct {
	key   string
	value interface{}
}

func NewLRUCache(capacity int) (*LRUCache, error) {
	if capacity <= 0 {
		return nil, errors.New("must provide a positive capacity")
	}

	return &LRUCache{
		keyToElement: make(map[string]*list.Element),
		list:         list.New(),
		capacity:     capacity,
	}, nil
}

func (c *LRUCache) Put(key string, value any) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if elem, ok := c.keyToElement[key]; ok {
		c.list.MoveToFront(elem)
		elem.Value = lruCacheItem{
			key:   key,
			value: value,
		}
		return
	}

	item := lruCacheItem{
		key:   key,
		value: value,
	}

	if c.list.Len() < c.capacity {
		c.keyToElement[key] = c.list.PushFront(item)
		return
	}

	elem := c.list.Back()
	delete(c.keyToElement, elem.Value.(lruCacheItem).key)

	elem.Value = item
	c.keyToElement[key] = elem
	c.list.MoveToFront(elem)
}

func (c *LRUCache) Get(key string) (any, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if elem, ok := c.keyToElement[key]; ok {
		c.list.MoveToFront(elem)
		return elem.Value.(lruCacheItem).value, true
	}

	return nil, false
}
