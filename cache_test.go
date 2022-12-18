package fastcache_lru

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	c := New[int](5)

	v1, ok1 := c.Get("k1")
	fmt.Println(v1, ok1)

	// set
	c.Set("k1", 1)
	c.Set("k2", 2)
	c.Set("k3", 3)
	c.Set("k4", 4)
	c.Set("k5", 5)
	fmt.Println(c.ll)
	c.Set("k1", 1)
	fmt.Println(c.ll)
	c.Set("k6", 6)
	fmt.Println(c.ll)

	c.Del("k1")
	fmt.Println("del k1: ", c.ll, c.mp)

	c.Del("k2")
	fmt.Println("del k2: ", c.ll)
}

func BenchmarkCache_Get(b *testing.B) {
	c := New[int](100)
	// set
	for i := 0; i < 101; i++ {
		key := fmt.Sprintf("k%d", i)
		c.Set(key, 10)
	}
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("k%d", rand.Intn(1000))
		c.Get(key)
	}
}

func BenchmarkCache_Set(b *testing.B) {
	c := New[int](100)
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("k%d", rand.Intn(10000))
		c.Set(key, 10)
	}
}

func BenchmarkCache_Del(b *testing.B) {
	c := New[int](100)
	// set
	for i := 0; i < 101; i++ {
		key := fmt.Sprintf("k%d", i)
		c.Set(key, 10)
	}
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("k%d", rand.Intn(1000))
		c.Del(key)
	}
}

func TestCache_Stats(t *testing.T) {
	c := New[int](1000)

	for i := 0; i < 101; i++ {
		key := fmt.Sprintf("k%d", i)
		c.Set(key, 10)
	}
	time.Sleep(10 * time.Second)

	// get some
	for i := 0; i < 10000; i++ {
		key := fmt.Sprintf("k%d", rand.Intn(1000))
		c.Get(key)
	}
	time.Sleep(10 * time.Second)

	// del some
	for i := 0; i < 100000; i++ {
		key := fmt.Sprintf("k%d", rand.Intn(50))
		c.Del(key)
	}
	time.Sleep(10 * time.Second)

	// sleep
	time.Sleep(100 * time.Second)
}
