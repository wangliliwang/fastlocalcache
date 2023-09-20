package fastlocalcache

import (
	"errors"
	"log"
	"sync"
)

type setState uint64

const (
	setStateSet     setState = 1
	setStateReplace setState = 2
)

type deleteState uint64

const (
	deleteStateDoDelete  deleteState = 1
	deleteStateDoNothing deleteState = 2
)

type cacheShard struct {
	mu    sync.RWMutex
	clock Clocker

	ringIndex map[uint64]uint32
	ring      *Ring

	logger log.Logger
}

// 需要搞一个state. FirstSet, Replace
// 是否删除了元素
func (cs *cacheShard) Set(key string, keyHash uint64, ety entry) (setState, bool, error) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	// 均摊算法：尝试从head删除非正常元素
	deleteOnSet := false
	{
		headEty, peekErr := cs.ring.Peek()
		if peekErr != nil {
			cs.logger.Println("peek error")
		} else {
			if !headEty.isNormal() || !headEty.isExpired(cs.clock.CurrentTimestamp()) {
				// delete data
				popedEty, popErr := cs.ring.Pop()
				if popErr != nil {
					cs.logger.Println("pop error")
				} else {
					deleteOnSet = true
					// delete index
					delete(cs.ringIndex, popedEty.getKeyHash())
				}
			}
		}
	}

	// todo(扩容机制)

	// get index
	keyIndex, ok := cs.ringIndex[keyHash]
	if !ok { // ring.pushNew
		enqueueErr := cs.ring.Push(ety)
		if enqueueErr != nil {
			return 0, deleteOnSet, enqueueErr
		}
		return setStateSet, deleteOnSet, nil
	}

	// 检测是否有hash冲突. 有则报错.
	oldEty, getErr := cs.ring.Get(keyIndex)
	if getErr != nil {
		return 0, deleteOnSet, getErr
	}
	if oldEty.getKey() != key {
		return 0, deleteOnSet, errors.New("hash collision!")
	}

	// 如果这个key存在，那么标记之前的key为删除状态
	setStateErr := cs.ring.SetState(keyIndex, entryStateDeleted)
	if setStateErr != nil {
		return 0, deleteOnSet, setStateErr
	}

	// 插入新元素
	enqueueErr := cs.ring.Push(ety)
	if enqueueErr != nil {
		return 0, deleteOnSet, enqueueErr
	}
	return setStateReplace, deleteOnSet, nil
}

func (cs *cacheShard) Get(key string, keyHash uint64) (entry, error) {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	// get index
	keyIndex, ok := cs.ringIndex[keyHash]
	if !ok { // ring.pushNew
		return nil, errors.New("not found")
	}

	// 只需要取出来keyHash对应的即可。
	ety, getErr := cs.ring.Get(keyIndex)
	if getErr != nil {
		return nil, getErr
	}

	// 比对key是否一致
	if ety.getKey() != key {
		return nil, errors.New("key not equal")
	}

	return ety, nil
}

func (cs *cacheShard) Delete(key string, keyHash uint64) (deleteState, error) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	// get index
	keyIndex, ok := cs.ringIndex[keyHash]
	if !ok { // ring.pushNew
		return deleteStateDoNothing, nil
	}

	// exec delete
	setStateErr := cs.ring.SetState(keyIndex, entryStateDeleted)
	if setStateErr != nil {
		return 0, setStateErr
	}
	return deleteStateDoDelete, nil
}
