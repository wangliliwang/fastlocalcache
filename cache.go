package fastlocalcache

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	shardsCount int64 = 256
	neverExpire int64 = -1
)

type Cache struct {
	serializer Serializer
	shardedMap *shardedMap
}

func NewCache() *Cache {
	c := &Cache{
		serializer: JSONSerializer{},
		shardedMap: newShardedMap(),
	}
	go c.scanAndExpire()
	return c
}

func hasExpired(now, expireAt int64) bool {
	return expireAt != neverExpire && now > expireAt
}

func (c *Cache) scanAndExpire() {
	quitChannel := make(chan os.Signal)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case <-time.After(time.Minute):
			c.shardedMap.scanAndExpire()
		case <-quitChannel:
			break
		}
	}
}

func (c *Cache) Get(key string, value any) error {
	// get from store
	ci, ok := c.shardedMap.Get(key)
	if !ok {
		return errors.New("missing key")
	}

	// delete expired key
	if hasExpired(time.Now().Unix(), ci.expireAt) {
		c.shardedMap.Del(key)
		return errors.New("missing key")
	}

	// unmarshal
	err := c.serializer.Unmarshal(ci.value, value)
	if err != nil {
		return fmt.Errorf("unmarshal error: %w", err)
	}

	return nil
}

func (c *Cache) Set(key string, value any, expiration *time.Duration) error {
	// marshal
	bs, err := c.serializer.Marshal(value)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	// set to store
	ci := &cacheItem{
		value: bs,
	}
	if expiration == nil {
		ci.expireAt = neverExpire
	} else {
		ci.expireAt = time.Now().Add(*expiration).Unix()
	}
	c.shardedMap.Set(key, ci)

	return nil
}

func (c *Cache) Len() int64 {
	return c.shardedMap.len
}

func (c *Cache) Del(key string) {
	c.shardedMap.Del(key)
}

type cacheItem struct {
	value    []byte
	expireAt int64 // unix timestamp, in seconds
}

func newShardedMap() *shardedMap {
	shards := make([]*sync.Map, shardsCount)
	for i := 0; i < int(shardsCount); i++ {
		shards[i] = &sync.Map{}
	}
	return &shardedMap{
		shards:      shards,
		shardsCount: shardsCount,
		keyToHash:   KeyToHash,
	}
}

type shardedMap struct {
	shards      []*sync.Map
	shardsCount int64
	len         int64 // 实际存储的key的数量，包括失效的
	keyToHash   func(key string) uint64
}

func (m *shardedMap) getShard(key string) *sync.Map {
	hash := m.keyToHash(key)
	return m.shards[hash%uint64(m.shardsCount)]
}

func (m *shardedMap) Get(key string) (*cacheItem, bool) {
	value, ok := m.getShard(key).Load(key)
	if !ok {
		return nil, false
	}
	ci, ok := value.(*cacheItem)
	if !ok {
		panic("unsupported value")
	}
	return ci, true
}

func (m *shardedMap) Set(key string, value *cacheItem) {
	_, loaded := m.getShard(key).Swap(key, value)
	if !loaded {
		m.len++
	}
}

func (m *shardedMap) Del(key string) {
	m.del(key)
}

func (m *shardedMap) del(key string) {
	_, loaded := m.getShard(key).LoadAndDelete(key)
	if loaded {
		m.len--
	}
}

func (m *shardedMap) scanAndExpire() {
	now := time.Now().Unix()
	for _, shard := range m.shards {
		shard.Range(func(key, value any) bool {
			ci, ok := value.(*cacheItem)
			if !ok {
				panic("unsupported value")
			}
			if hasExpired(now, ci.expireAt) {
				keyStr, ok := key.(string)
				if !ok {
					panic("unsupported key type")
				}
				m.del(keyStr)
			}

			return true
		})
	}
}
