package pokecache

import "time"
import "sync"

type Cache struct {
	cache map[string]cacheEntry
	mu    sync.Mutex
}

type cacheEntry struct {
	createdTime time.Time
	val         []byte
}

func NewCache(interval time.Duration) *Cache {
	newCache := Cache{
		cache: make(map[string]cacheEntry),
	}
	go newCache.reapLoop(interval)
	return &newCache
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	c.cache[key] = cacheEntry{createdTime: time.Now(), val: val}
	c.mu.Unlock()
}

func (c *Cache) Get(key string) ([]byte, bool) {
	defer c.mu.Unlock()
	c.mu.Lock()
	if val, ok := c.cache[key]; ok {
		return val.val, true
	}

	return []byte{}, false
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		c.mu.Lock()
		age := time.Now().Add(-interval)
		for cache := range c.cache {
			if age.After(c.cache[cache].createdTime) {
				delete(c.cache, cache)
			}
		}
		c.mu.Unlock()
	}
}
