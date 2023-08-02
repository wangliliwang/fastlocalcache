
# fastlocalcache

## 考虑要素

### 接口

```golang
type Cache interface {
    Get(key string, value any) error
    Set(key string, value any, expiration *time.Duration) error
    Del(key string)
    Len() int64
}
```

### 功能

1. 不支持空间限制，也就是不能根据内存用量、存储条目数量限制空间用量。期望是支持的。
2. 淘汰机制。目前只支持按照时间淘汰，最好支持一种更好的淘汰机制，以在有限的空间内满足业务需求。淘汰机制对命中率有很大影响。

### 性能

1. 需要保证并发安全。使用`sync.Map`包存储数据，底层使用读写锁保证并发安全。支持并发读写可以提升吞吐量。
2. 使用锁的基础上，减少锁等待时间。使用256个`sync.Map`存储数据，使用key的hash值选取Map。
3. 需要考虑GC。除了256个`sync.Map`是引用类型，map中的value是`[]byte`也是引用类型，所以gc扫描会影响性能。
4. 目前支持按照时间淘汰，每隔1分钟执行一次全部扫描，遇到过期的key会执行删除操作，所以扫描时会加读写锁，造成性能损耗。
