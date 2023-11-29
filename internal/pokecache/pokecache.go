package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
	permanent bool
}

type Cache struct {
	cache map[string]cacheEntry
	mutex *sync.Mutex
}

func NewCache(interval time.Duration) Cache {
	c := Cache{
		cache: make(map[string]cacheEntry),
		mutex: &sync.Mutex{},
	}
	go c.reapLoop(interval)
	return c
}

func (c *Cache) Add(key string, val []byte, permanent bool) {
	c.cache[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
		permanent: permanent,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	for k, v := range c.cache {
		if k == key {
			return v.val, true
		}
	}
	return nil, false
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		c.reap(time.Now().UTC(), interval)
	}
}

func (c *Cache) reap(t time.Time, interval time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for key, entry := range c.cache {
		if !entry.permanent && t.Sub(entry.createdAt) > interval {
			delete(c.cache, key)
		}
	}
}
