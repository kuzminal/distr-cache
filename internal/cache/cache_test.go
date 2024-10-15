package cache

import (
	"testing"
)

func TestCacheInitialization(t *testing.T) {
	cache := NewCache()
	if cache == nil {
		t.Errorf("NewCache() = %v, want non-nil", cache)
	}
}

func TestCacheSetAndGetBehavior(t *testing.T) {
	cache := NewCache()
	cache.Set("key1", "value1")

	value, found := cache.Get("key1")
	if !found || value != "value1" {
		t.Errorf("Get() = %v, %v, want %v, %v", value, found, "value1", true)
	}
}

func TestCacheGetNonExistentKeyBehavior(t *testing.T) {
	cache := NewCache()

	_, found := cache.Get("nonExistentKey")
	if found {
		t.Errorf("Get() = %v, want %v", found, false)
	}
}

func TestCacheSetOverwritesValueBehavior(t *testing.T) {
	cache := NewCache()
	cache.Set("key1", "value1")
	cache.Set("key1", "value2")

	value, _ := cache.Get("key1")
	if value != "value2" {
		t.Errorf("Get() = %v, want %v", value, "value2")
	}
}
