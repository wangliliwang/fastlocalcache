package fastlocalcache

import (
	"errors"
	"math"
	"sync/atomic"
	"time"
)

const (
	shardsCount int64         = 256
	neverExpire time.Duration = -1
)

type DeleteReason int64

const (
	DeleteReasonByUser DeleteReason = iota
)

type Cacher interface {
	Get(key string, value any) error
	Set(key string, value any, expiration *time.Duration) error
	Contains(key string) bool
	Len() int64
	Delete(key string)
	//Iterator()
}

type OnDelete func()

type Config struct {
	Hasher     Hasher
	Serializer Serializer

	CleanUpInterval time.Duration
	OnDelete        OnDelete // 删除callback

	Verbose bool // 是否打印详情
}

type Cache struct {
	len int64

	serializer  Serializer
	cacheShards [shardsCount]*cacheShard
	hasher      Hasher
	clocker     Clocker

	expiration     time.Duration
	expirationRing []byte // 方便失效的圆环
}

func (c *Cache) getShardInfo(key string) (uint64, int) {
	hash := c.hasher.Sum64(key)
	return hash, int(hash % uint64(shardsCount))
}

func (c *Cache) hasExpired(expireTimestampInseconds uint64) bool {
	return expireTimestampInseconds < uint64(c.clocker.CurrentTime().Unix())
}

func (c *Cache) incrLen() {
	atomic.AddInt64(&c.len, 1)
}

func (c *Cache) descLen() {
	atomic.AddInt64(&c.len, -1)
}

func (c *Cache) Get(key string, value any) error {
	// hash
	keyHash, shardIndex := c.getShardInfo(key)

	ety, getErr := c.cacheShards[shardIndex].Get(key, keyHash)
	if getErr != nil {
		return getErr
	}

	// unpack
	if c.hasExpired(ety.getTimestamp()) {
		return errors.New("expired")
	}

	// dese
	unmarshalErr := c.serializer.Unmarshal(ety.getValue(), value)
	if unmarshalErr != nil {
		return unmarshalErr
	}

	return nil
}

func (c *Cache) Set(key string, value any, expiration *time.Duration) error {
	// key 长度检查
	if len(key) > 1<<17 {
		return errors.New("key too long")
	}

	// hash
	keyHash, shardIndex := c.getShardInfo(key)

	// serialize
	valueBytes, marsharErr := c.serializer.Marshal(value)
	if marsharErr != nil {
		return marsharErr
	}
	if len(valueBytes) > 1<<33 {
		return errors.New("value too long")
	}

	// 计算过期时间.
	// 不传过期时间，按照默认过期时间来
	// 传-1的过期时间，表示用不失效
	var expireTimestampInSeconds uint64
	if expiration == nil { // == nil, 表示无过期时间
		expiration = &c.expiration
	}
	if *expiration == neverExpire {
		expireTimestampInSeconds = math.MaxUint64
	} else {
		expireTimestampInSeconds = uint64(c.clocker.CurrentTime().Add(*expiration).Unix())
	}

	ety := packEntry(expireTimestampInSeconds, keyHash, key, valueBytes)

	ss, setErr := c.cacheShards[shardIndex].Set(key, keyHash, ety)
	if setErr != nil {
		return setErr
	}

	// counter
	if ss == setStateSet {
		c.incrLen()
	}
	return nil
}

func (c *Cache) Contains(key string) bool {
	// hash
	keyHash, shardIndex := c.getShardInfo(key)

	ety, getErr := c.cacheShards[shardIndex].Get(key, keyHash)
	if getErr != nil {
		return false
	}

	// unpack
	if c.hasExpired(ety.getTimestamp()) {
		return false
	}

	return true
}

func (c *Cache) Len() int64 {
	return atomic.LoadInt64(&c.len)
}

func (c *Cache) Delete(key string) {
	// hash
	keyHash, shardIndex := c.getShardInfo(key)

	ds := c.cacheShards[shardIndex].Delete(key, keyHash)
	if ds == deleteStateDoDelete {
		c.descLen()
	}
}

func NewCache() *Cache {
	c := &Cache{
		serializer: JSONSerializer{},
		hasher:     newDefaultHasher(),
	}

	// init shards

	// init clean up

	return c
}
