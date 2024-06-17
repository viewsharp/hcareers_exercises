package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLRUCache_OK(t *testing.T) {
	cache, _ := NewLRUCache(1)

	cache.Put("foo", "bar")
	value, ok := cache.Get("foo")
	assert.Equal(t, "bar", value)
	assert.True(t, ok)
}

func TestLRUCache_Reset(t *testing.T) {
	cache, _ := NewLRUCache(1)

	cache.Put("foo", "bar")
	value, ok := cache.Get("foo")
	assert.Equal(t, "bar", value)
	assert.True(t, ok)

	cache.Put("foo", "bar1")
	value, ok = cache.Get("foo")
	assert.Equal(t, "bar1", value)
	assert.True(t, ok)
}

func TestLRUCache_KeyMismatch(t *testing.T) {
	cache, _ := NewLRUCache(1)

	cache.Put("foo", "bar")
	value, ok := cache.Get("foo2")
	assert.Nil(t, value)
	assert.False(t, ok)
}

func TestLRUCache_Eviction(t *testing.T) {
	cache, _ := NewLRUCache(2)
	cache.Put("foo1", "bar1")
	cache.Put("foo2", "bar2")
	cache.Put("foo3", "bar3")

	value, ok := cache.Get("foo1")
	assert.Nil(t, value)
	assert.False(t, ok)

	value, ok = cache.Get("foo2")
	assert.Equal(t, "bar2", value)
	assert.True(t, ok)

	value, ok = cache.Get("foo3")
	assert.Equal(t, "bar3", value)
	assert.True(t, ok)
}

func TestLRUCache_ResetAndEviction(t *testing.T) {
	cache, _ := NewLRUCache(2)
	cache.Put("foo1", "bar1")
	cache.Put("foo2", "bar2")
	cache.Put("foo1", "bar1")
	cache.Put("foo3", "bar3")

	value, ok := cache.Get("foo1")
	assert.Equal(t, "bar1", value)
	assert.True(t, ok)

	value, ok = cache.Get("foo2")
	assert.Nil(t, value)
	assert.False(t, ok)

	value, ok = cache.Get("foo3")
	assert.Equal(t, "bar3", value)
	assert.True(t, ok)
}

func TestLRUCache_GetAndEviction(t *testing.T) {
	cache, _ := NewLRUCache(2)
	cache.Put("foo1", "bar1")
	cache.Put("foo2", "bar2")
	cache.Get("foo1")
	cache.Put("foo3", "bar3")

	value, ok := cache.Get("foo1")
	assert.Equal(t, "bar1", value)
	assert.True(t, ok)

	value, ok = cache.Get("foo2")
	assert.Nil(t, value)
	assert.False(t, ok)

	value, ok = cache.Get("foo3")
	assert.Equal(t, "bar3", value)
	assert.True(t, ok)
}
