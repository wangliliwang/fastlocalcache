package fastlocalcache

import (
	"strings"
	"testing"
	"time"
)

func TestEntryPackUnpack(t *testing.T) {
	state := entryStateNormal
	timestamp := uint64(time.Now().Unix())
	key := "this is a key"
	keyHash := newDefaultHasher().Sum64(key)
	value := []byte(key)

	// pack
	ety, packErr := packEntry(state, timestamp, keyHash, key, value)
	noError(t, packErr)

	// unpack
	assertEqual(t, state, ety.getState())
	assertEqual(t, timestamp, ety.getTimestamp())
	assertEqual(t, keyHash, ety.getKeyHash())
	assertEqual(t, uint16(len(key)), ety.getKeyLen())
	assertEqual(t, key, ety.getKey())
	assertEqual(t, uint32(len(value)), ety.getValueLen())
	assertEqual(t, value, ety.getValue())

	// empty key
	emptyKey := ""
	emptyKeyHash := newDefaultHasher().Sum64(emptyKey)
	_, emptyKeyPackErr := packEntry(state, timestamp, emptyKeyHash, emptyKey, value)
	hasError(t, emptyKeyPackErr)

	// too long key
	tooLongKey := strings.Repeat("a", maxKeyLen+1)
	tooLongKeyHash := newDefaultHasher().Sum64(tooLongKey)
	_, tooLongKeyPackErr := packEntry(state, timestamp, tooLongKeyHash, tooLongKey, value)
	hasError(t, tooLongKeyPackErr)

	// empty value
	emptyValue := make([]byte, 0)
	_, emptyValuePackErr := packEntry(state, timestamp, keyHash, key, emptyValue)
	hasError(t, emptyValuePackErr)

	// too long value
	tooLongValue := make([]byte, maxValueLen+1)
	_, tooLongValuePackErr := packEntry(state, timestamp, keyHash, key, tooLongValue)
	hasError(t, tooLongValuePackErr)
}
