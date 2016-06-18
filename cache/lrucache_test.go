package cache

import (
	"testing"
)

func TestLRUCache(t *testing.T) {
	// illegal size lru cache test
	cache := NewLRUCache(-1)
	if cache != nil {
		t.Error("illegal new LRU cache assigned, cache size less than 0")
	}

	cache = NewLRUCache(0)
	if cache != nil {
		t.Error("illegal new LRU cache assighid, cache size equal 0")
	}

	// legal size lru cache test
	cache = NewLRUCache(3)
	cache.Put(1, "1")
	if len(cache.hmap) != 1 || cache.freeIdx != 1 {
		t.Error("Put new data(key=1) in a initial LRL cache in error")
	}

	cache.Put(2, "2")
	if len(cache.hmap) != 2 || cache.freeIdx != 2 {
		t.Error("Put new data(key=2) in LRL cache in error")
	}

	cache.Put(1, "new 1")
	if len(cache.hmap) != 2 || cache.freeIdx != 2 {
		t.Error("Put a new data(exist key, key=1) in LRU cache in error")
	}
	if cache.head.next.key != 1 {
		t.Error("position of attaching new data(exist key, key=1) to LRU cache in error")
	}

	cache.Put(3, "3")
	cache.Put(4, "4")
	if len(cache.hmap) != 3 || cache.freeIdx != 3 {
		t.Error("cache size control in error")
	}
	if cache.head.next.key != 4 {
		t.Error("position of attaching new data to LRU cache in error")
	}
	b := false
	for i := cache.head.next; i != cache.tail; i = i.next {
		if i.key == 2 {
			b = true
			break
		}
	}
	if b {
		t.Error("detach old data for new data in error")
	}

	// Get func test
	v := cache.Get(1)
	if v == nil {
		t.Error("Get a exist data return nil")
	}
	if v != "new 1" {
		t.Error("Get a updated data value incorrectly")
	}
}
