package tests

import (
	"testing"
	"time"

	"github.com/crypolt/fastcache"
)

func BenchmarkSingleCache_Get(b *testing.B) {

	c := fastcache.New()
	key := []byte("bench")

	c.Set(key, []byte("value"), time.Minute)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Get(key)
	}
}

func BenchmarkShardedCache_Get(b *testing.B) {

	c := fastcache.New()
	key := []byte("bench")

	c.Set(key, []byte("value"), time.Minute)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Get(key)
	}
}
