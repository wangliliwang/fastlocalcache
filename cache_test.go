package fastlocalcache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCache(t *testing.T) {
	// test set
	cache := NewCache()
	setErr := cache.Set("t-1", "a", nil)
	assert.Nil(t, setErr)
	assert.Equal(t, int64(1), cache.Len())
	setErr = cache.Set("t-1", "b", nil)
	assert.Nil(t, setErr)
	assert.Equal(t, int64(1), cache.Len())

	// test get
	var str string
	getErr := cache.Get("t-1", &str)
	assert.Nil(t, getErr)
	assert.Equal(t, "b", str)

	// test del
	cache.Del("t-1")
	assert.Equal(t, int64(0), cache.Len())
	getErr = cache.Get("t-1", str)
	assert.NotNil(t, getErr)
}
