package fastlocalcache

import (
	"encoding/binary"
	"errors"
)

// entry pack struct:
// state(1) - timestamp(8) - keyHash(8) - keyLen(2) - valueLen(4) - key - value
// key最大64K; value最大4G

type entry []byte

type entryState byte

const (
	entryStateNormal  entryState = 0
	entryStateDeleted entryState = 1

	maxKeyLen   = 1<<16 - 1
	maxValueLen = 1<<32 - 1

	entryStateSizeInBytes     = 1
	entryTimestampSizeInBytes = 8
	entryKeyHashSizeInBytes   = 8
	entryKeyLenSizeInBytes    = 2
	entryValueLenSizeInBytes  = 4
	entryHeaderSizeInBytes    = entryStateSizeInBytes + entryTimestampSizeInBytes +
		entryKeyHashSizeInBytes + entryKeyLenSizeInBytes + entryValueLenSizeInBytes
)

func packEntry(state entryState, timestamp uint64, keyHash uint64, key string, value []byte) (entry, error) {
	keyLen, valueLen := len(key), len(value)
	if keyLen == 0 || keyLen > maxKeyLen || valueLen == 0 || valueLen > maxValueLen {
		return nil, errors.New("invalid key or value")
	}
	entrySizeInBytes := entryHeaderSizeInBytes + keyLen + valueLen
	ety := make([]byte, entrySizeInBytes)
	ety[0] = byte(state)
	binary.LittleEndian.PutUint64(ety[entryStateSizeInBytes:], timestamp)
	binary.LittleEndian.PutUint64(ety[entryStateSizeInBytes+entryTimestampSizeInBytes:], keyHash)
	binary.LittleEndian.PutUint16(
		ety[entryStateSizeInBytes+entryTimestampSizeInBytes+entryKeyHashSizeInBytes:],
		uint16(keyLen),
	)
	binary.LittleEndian.PutUint32(
		ety[entryStateSizeInBytes+entryTimestampSizeInBytes+entryKeyHashSizeInBytes+entryKeyLenSizeInBytes:],
		uint32(valueLen),
	)
	copy(ety[entryHeaderSizeInBytes:], key)
	copy(ety[entryHeaderSizeInBytes+keyLen:], value)
	return ety, nil
}

func (e entry) getState() entryState {
	return entryState(e[0])
}

func (e entry) getTimestamp() uint64 {
	return binary.LittleEndian.Uint64(e[entryStateSizeInBytes:])
}

func (e entry) getKeyHash() uint64 {
	return binary.LittleEndian.Uint64(e[entryStateSizeInBytes+entryTimestampSizeInBytes:])
}

func (e entry) getKeyLen() uint16 {
	return binary.LittleEndian.Uint16(e[entryStateSizeInBytes+entryTimestampSizeInBytes+entryKeyHashSizeInBytes:])
}

func (e entry) getValueLen() uint32 {
	return binary.LittleEndian.Uint32(e[entryStateSizeInBytes+entryTimestampSizeInBytes+entryKeyHashSizeInBytes+entryKeyLenSizeInBytes:])
}

func (e entry) getKey() string {
	keyLen := int(e.getKeyLen())
	dst := make([]byte, keyLen)
	copy(dst, e[entryHeaderSizeInBytes:entryHeaderSizeInBytes+keyLen])
	return string(dst)
}

func (e entry) getValue() []byte {
	keyLen := int(e.getKeyLen())
	valueLen := int(e.getValueLen())
	dst := make([]byte, valueLen)
	copy(dst, e[entryHeaderSizeInBytes+keyLen:entryHeaderSizeInBytes+keyLen+valueLen])
	return dst
}

func (e entry) isNormal() bool {
	return e.getState() == entryStateNormal
}

func (e entry) isExpired(currentTimestamp uint64) bool {
	return e.getTimestamp() <= currentTimestamp
}
