package fastcache_lru

import (
	"sync"
)

type Cache[T any] struct {
	mp        map[string]*linkedListNode[T]
	ll        *linkedList[T]
	len       int
	cap       int
	m         sync.RWMutex
	s         *stats
	zeroValue T
}

func New[T any](cap int) *Cache[T] {
	var zeroValue T
	s := &stats{}
	runStatsServer(s)
	return &Cache[T]{
		mp:        make(map[string]*linkedListNode[T]),
		ll:        newLinkedList[T](),
		len:       0,
		cap:       cap,
		m:         sync.RWMutex{},
		s:         s,
		zeroValue: zeroValue,
	}
}

func (c *Cache[T]) Get(key string) (T, bool) {
	c.m.Lock()
	defer c.m.Unlock()

	c.s.GetCalls++

	nd, ok := c.mp[key]
	if !ok {
		c.s.Misses++
		return c.zeroValue, false
	}

	// move the nd to head
	c.ll.moveToHead(nd)

	return nd.value, true
}

func (c *Cache[T]) Set(key string, t T) {
	c.m.Lock()
	defer c.m.Unlock()

	c.s.SetCalls++

	if nd, ok := c.mp[key]; ok {
		nd.value = t
		// move to head
		c.ll.moveToHead(nd)
	} else {
		newNd := &linkedListNode[T]{
			value: t,
			key:   key,
		}
		c.mp[key] = newNd
		c.ll.moveToHead(newNd)
		c.len++

		// check over len
		if c.len > c.cap {
			delNd := c.ll.delTail()
			delete(c.mp, delNd.key)
			c.len--
		}
	}
}

func (c *Cache[T]) Del(key string) {
	c.m.Lock()
	defer c.m.Unlock()

	c.s.DelCalls++

	if nd, ok := c.mp[key]; ok {
		delete(c.mp, key)
		c.ll.del(nd)
		c.len--
	}
}

func (c *Cache[T]) String() string {
	return ""
}
