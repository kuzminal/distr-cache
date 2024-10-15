package cache

import (
	"container/list"
	"sync"
	"time"
)

type CacheItem struct {
	Value      string
	ExpiryTime time.Time
}
type entry struct {
	key   string
	value CacheItem
}
type Cache struct {
	mu       sync.RWMutex
	items    map[string]*list.Element
	eviction *list.List
	capacity int
}

func NewCache(capacity int) *Cache {
	return &Cache{
		items:    make(map[string]*list.Element),
		eviction: list.New(),
		capacity: capacity,
	}
}
func (c *Cache) Set(key, value string, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if elem, found := c.items[key]; found {
		c.eviction.Remove(elem)
		delete(c.items, key)
	}
	if c.eviction.Len() >= c.capacity {
		c.evictLRU()
	}
	item := CacheItem{
		Value:      value,
		ExpiryTime: time.Now().Add(ttl),
	}
	elem := c.eviction.PushFront(&entry{key, item})
	c.items[key] = elem
}
func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	elem, found := c.items[key]
	if !found || time.Now().After(elem.Value.(*entry).value.ExpiryTime) {
		if found {
			c.eviction.Remove(elem)
			delete(c.items, key)
		}
		return "", false
	}
	c.eviction.MoveToFront(elem)
	return elem.Value.(*entry).value.Value, true
}

func (c *Cache) evictLRU() {
	elem := c.eviction.Back()
	if elem != nil {
		c.eviction.Remove(elem)
		kv := elem.Value.(*entry)
		delete(c.items, kv.key)
	}
}

func (c *Cache) StartEvictionTicker(d time.Duration) {
	ticker := time.NewTicker(d)
	go func() {
		for range ticker.C {
			c.evictExpiredItems()
		}
	}()
}
func (c *Cache) evictExpiredItems() {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	for key, elem := range c.items {
		if now.After(elem.Value.(*entry).value.ExpiryTime) {
			delete(c.items, key)
		}
	}
}
