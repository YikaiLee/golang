package cache

import (
	"sync"
)

type entity struct {
	key  interface{}
	data interface{}
	prev *entity
	next *entity
}

type LRUCache struct {
	sync.RWMutex
	head     *entity
	tail     *entity
	hmap     map[interface{}]*entity
	entities []entity
	freeIdx  int // next free entity index in entities
}

func NewLRUCache(size int) (cache *LRUCache) {
	if size <= 0 {
		return nil
	}

	cache = &LRUCache{
		head:     &entity{nil, nil, nil, nil},
		tail:     &entity{nil, nil, nil, nil},
		hmap:     make(map[interface{}]*entity),
		entities: make([]entity, size, size),
		freeIdx:  0,
	}
	cache.head.next = cache.tail
	cache.tail.prev = cache.head

	return cache
}

func (cache *LRUCache) detach(en *entity) {
	en.prev.next = en.next
	en.next.prev = en.prev
}

// attach entity to head
func (cache *LRUCache) attach(en *entity) {
	en.prev = cache.head
	en.next = cache.head.next
	en.next.prev = en
	en.prev.next = en
}

func (cache *LRUCache) Put(key interface{}, data interface{}) {
	cache.Lock()
	defer cache.Unlock()
	en, ok := cache.hmap[key]
	if ok {
		// exist entity in hmap
		cache.detach(en)
		en.data = data
		cache.attach(en)
		return
	}

	// not exist entity in hmap
	if cache.freeIdx == len(cache.entities) {
		en = cache.tail.prev
		cache.detach(en)
		delete(cache.hmap, en.key)
	} else {
		en = &cache.entities[cache.freeIdx]
		cache.freeIdx++
	}

	// reset en
	en.key, en.data = key, data
	cache.attach(en)

	cache.hmap[key] = en
}

func (cache *LRUCache) Get(key interface{}) interface{} {
	cache.Lock()
	defer cache.Unlock()
	en, ok := cache.hmap[key]
	if !ok {
		return nil
	}

	cache.detach(en)
	cache.attach(en)
	return en.data
}
