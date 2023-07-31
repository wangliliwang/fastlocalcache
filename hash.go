package fastlocalcache

import "hash/fnv"

func KeyToHash(key string) uint64 {
	h := fnv.New64()
	h.Write([]byte(key))
	return h.Sum64()
}
