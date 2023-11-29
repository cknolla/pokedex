package pokecache

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	waitTime := 5 * time.Millisecond
	c := NewCache(waitTime)
	c.Add("test_key", []byte{1, 2, 3}, false)
	c.Add("permanent_key", []byte{4, 5, 6}, true)
	assert.Equal(t, []string{"test_key", "permanent_key"}, c.GetKeys())
	value, found := c.Get("test_key")
	assert.True(t, found, "test_key not found in cache")
	assert.Equal(t, []byte{1, 2, 3}, value, "value retrieved from cache is %v", value)
	time.Sleep(waitTime * 2)
	value, found = c.Get("test_key")
	assert.False(t, found, "test_key still in cache after waitTime")
	assert.Nil(t, value, "value not returned nil if key not found")
	value, found = c.Get("permanent_key")
	assert.True(t, found)
	assert.Equal(t, []byte{4, 5, 6}, value)
}
