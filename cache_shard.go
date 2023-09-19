package fastlocalcache

import "sync"

// key最大64K; value最大4G
// timestamp(8) - keyHash(8) - keyLen(2) - valueLen(4) - key - value
type entry []byte

func packEntry(timestamp uint64, keyHash uint64, key string, value []byte) entry {

}

func (e entry) getTimestamp() uint64 {
	panic("impl me")
}

func (e entry) getValue() []byte {
	panic("impl me")
}

type setState uint64

const (
	setStateSet     setState = 0
	setStateReplace setState = 1
)

type deleteState uint64

const (
	deleteStateDoDelete  deleteState = 0
	deleteStateDoNothing deleteState = 1
)

type cacheShard struct {
	mu sync.RWMutex

	ringIndex map[uint64]uint32
	data      []byte
}

// 需要搞一个state. FirstSet, Replace,
func (cs *cacheShard) Set(key string, keyHash uint64, ety entry) (setState, error) {

}

func (cs *cacheShard) Get(key string, keyHash uint64) (entry, error) {
	// 只需要取出来keyHash对应的即可。
	// 需要匹配key吗？需要
	// 需要校验过期时间吗？	不需要
}

func (cs *cacheShard) Delete(key string, keyHash uint64) deleteState {

}
