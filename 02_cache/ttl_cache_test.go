package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTTLCache_OK(t *testing.T) {
	cache := NewTTLCache()

	cache.Put("foo", "bar", time.Now().Add(time.Hour))
	value, ok := cache.Get("foo")
	assert.Equal(t, "bar", value)
	assert.True(t, ok)
}

func TestTTLCache_KeyMismatch(t *testing.T) {
	cache := NewTTLCache()

	cache.Put("foo", "bar", time.Now().Add(time.Hour))
	value, ok := cache.Get("foo2")
	assert.Nil(t, value)
	assert.False(t, ok)
}

func TestTTLCache_Expired(t *testing.T) {
	cache := NewTTLCache()

	cache.Put("foo", "bar", time.Now())
	cache.Put("foo2", "bar", time.Now().Add(time.Hour))
	time.Sleep(time.Nanosecond)

	value, ok := cache.Get("foo")
	assert.Nil(t, value)
	assert.False(t, ok)

	value, ok = cache.Get("foo2")
	assert.Equal(t, "bar", value)
	assert.True(t, ok)
}

func TestTTLCache_DeleteExpired(t *testing.T) {
	cache := NewTTLCache()

	cache.Put("foo", "bar", time.Now())
	cache.Put("foo2", "bar", time.Now().Add(time.Hour))
	time.Sleep(time.Nanosecond)

	assert.Equal(t, cache.DeleteExpired(), 1)
}
