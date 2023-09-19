# fastlocalcache

## 考虑要素

### 设计目标

1. 高性能
2. 支持并发
3. 简便易用

### 接口

```golang
type Cache interface {
    Get(key string, value []byte) ([]byte, error)
    Set(key string, value []byte, expiration *time.Duration) error
    Delete(key string)
    Len() int64
	Iterator() any // range over
}
```

bigcache interface(6.8k stars)

```golang
type Stats struct {
   // Hits is a number of successfully found keys
   Hits int64 `json:"hits"`
   // Misses is a number of not found keys
   Misses int64 `json:"misses"`
   // DelHits is a number of successfully deleted keys
   DelHits int64 `json:"delete_hits"`
   // DelMisses is a number of not deleted keys
   DelMisses int64 `json:"delete_misses"`
   // Collisions is a number of happened key-collisions
   Collisions int64 `json:"collisions"`
}

type Metadata struct {
    RequestCount uint32
}

type BigCache interface {
   Close() error
   Get(key string) ([]byte, error)
   GetWithInfo(key string) ([]byte, Response, error)
   Set(key string, entry []byte) error
   Append(key string, entry []byte) error // 追加多个entry到key上. 使用场景是啥？
   Delete(key string) error
   Reset() error
   ResetStats() error
   Len() int
   Capacity() int
   Stats() Stats
   KeyMetadata(key string) Metadata
   Iterator() *EntryInfoIterator
}
```

go-cache(https://github.com/patrickmn/go-cache)

```golang
// 结论：没啥学习价值。过期逻辑太烂，整个期间都Lock；没有对GC指针做优化。

// 过期原理：开启goroutine，定期运行DeleteExpired()

// 构造函数
New(defaultExpiration, cleanupInterval time.Duration) *Cache
NewFrom(defaultExpiration, cleanupInterval time.Duration, items map[string]Item) *Cache

// 同时有shard版本的
type GoCache interface {
    Set(k string, x interface{}, d time.Duration)
    Add(k string, x interface{}, d time.Duration) error // set on not exist
    Replace(k string, x interface{}, d time.Duration) error // set on exist
    Get(k string) (interface{}, bool)
    GetWithExpiration(k string) (interface{}, time.Time, bool)
    Delete(k string)
    DeleteExpired() // 删除过期的
    Items() map[string]Item
    Flush()

    Increment(k string, n int64) error
    Decrement(k string, n int64) error
    
    OnEvicted(f func(string, interface{}))
	
	// 持久化
    SaveFile(fname string) error
    Load(r io.Reader) error
}

```

fastcache (1.8k)

fastcache interface(6.8k stars)
```golang
type Cache
   func LoadFromFile(filePath string) (*Cache, error)
   func LoadFromFileOrNew(filePath string, maxBytes int) *Cache
   func New(maxBytes int) *Cache
   func (c *Cache) Del(k []byte)
   func (c *Cache) Get(dst, k []byte) []byte
   func (c *Cache) GetBig(dst, k []byte) (r []byte)
   func (c *Cache) Has(k []byte) bool
   func (c *Cache) HasGet(dst, k []byte) ([]byte, bool)
   func (c *Cache) Reset()
   func (c *Cache) SaveToFile(filePath string) error
   func (c *Cache) SaveToFileConcurrent(filePath string, concurrency int) error
   func (c *Cache) Set(k, v []byte)
   func (c *Cache) SetBig(k, v []byte)
   func (c *Cache) UpdateStats(s *Stats)
type Stats
func (s *Stats) Reset()
```

### 功能方面

1. 提供的接口
2. 淘汰机制。
3. 并发安全。一般通过RWMutex实现

### 性能

1. 读写性能。考虑存储机制、GC（减少指针使用）、序列化、内存分配（减少频繁扩容、缩容）
2. 并发性能。考虑加锁对性能的影响

### 方案要素

1. 接口。参考往上的实现。map, sync.Map, bigcache等

2. 存储
   + 大容量。因为小容量使用map、sync.Map也不会有什么问题
   + 减少指针使用。使用index+length+环形缓冲区索引数据
   + 减少序列化消耗。直接使用T来存储如何？也就是只能存储一种数据。这样也挺好。但是与GC优化冲突。
   + 使用锁保证并发安全
   + 增加锁，减少锁冲突
   + 选择合适的hash值计算函数
3. 淘汰
   + 肯定得支持淘汰机制。我觉得只支持按照时间淘汰就行。方便。
   + 容量不做限制。

其他：
1. GC优化，可以考虑堆外内存。这块内存需要自己进行管理。syscall.Mmap
2. 

### 功能实现

如何实现按时间的淘汰机制？
1. 定期扫描
2. 某种均摊算法
3. Set的时候，注册过期时间，到某个数据结构上。（靠谱）
   a. 实现思路：ring-buffer，时间取模放入；这样每个gc周期，只需要扫描很少一部分key；放入hash比较好。
4. Set的时候，注册一个回调，time.Sleep多久后，删除（不行，性能太低）

使用环状缓冲区存储，如何利用已经被删除的空间？可以参考内存分配算法。
1. 线性分配。最简单，但是比较低效。用链表来维护可以使用的内存区域。
2. 多级链表。像tcmalloc一样，预先划分不同的segment，每个segment若干资源。用的时候只需要向上取整就行。（靠谱）

如何解决hash冲突？
1. 参考map的实现。但是map的实现，bucket分配的空间是预先知道的。
2. 不考虑删除，开放寻址可以
3. 外挂法。这样会导致一些指针。这样也有一个前提：hash冲突不严重。（靠谱）

扩容时机？
1. 装载系数扩容。有效容量 / 占用的容量。占坑的容量一般都是内存碎片。同时，一部分kv外挂在ring外边。
2. 找不到分配的空间时，扩容。

扩容方式？
1. 另外一套相似的存储。index需要rehash。kv需要搬迁。

hash算法如何选取？完全不懂。

其他：
1. 是否要对key、value大小做限制？
2. 

### 模块拆分

1. key-value 引擎
2. 接口层

## 参考

1. 如何打造高性能的 Go 缓存库 https://www.luozhiyun.com/archives/531
   a. 值得借鉴的地方：堆外内存，减少GC压力
   b. 未解决的问题：删除、按照时间失效、hash冲突

2. bigcache分析 https://pandaychen.github.io/2020/03/03/BIGCACHE-ANALYSIS/
3. localcache选型 https://jonasx.com/archives/%E6%B5%85%E6%9E%90golang%E9%87%8C%E7%9A%84%E4%BC%97%E5%A4%9A%E7%9A%84localcache
4. writing a fast cache https://blog.allegro.tech/2016/03/writing-fast-cache-service-in-go.html
5. localcache比较 https://zhuanlan.zhihu.com/p/487455942

如何被pkg.go.dev收录？

## 性能调优

1. benchmark
2. 感知

