# fastlocalcache

## 考虑要素

### 设计目标

1. 高性能
2. 支持并发
3. 简便易用

### 接口

```golang
type Cache interface {
    Get(key string, value any) error
    Set(key string, value any, expiration *time.Duration) error
    Del(key string)
    Len() int64
}
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
   + 减少指针使用。使用index+length+环形缓冲区索引数据
   + 减少序列化消耗。直接使用T来存储如何？也就是只能存储一种数据。这样也挺好。
   + 使用锁保证并发安全
   + 增加锁，减少锁冲突
   + 选择合适的hash值计算函数
3. 淘汰
   + 肯定得支持淘汰机制。我觉得只支持按照时间淘汰就行。方便。
   + 容量不做限制。

### 模块拆分

1. key-value 引擎
2. 接口层

## 参考

1. https://www.luozhiyun.com/archives/531

## 性能调优

1. benchmark
2. 感知

