package tests

import (
	"testing"
	"time"

	"github.com/crypolt/fastcache"
)

func TestTTLExpiration(t *testing.T) {
	c := fastcache.New()

	key := []byte("ttl-key")
	val := []byte("ttl-val")

	c.Set(key, val, 10*time.Millisecond)

	time.Sleep(15 * time.Millisecond)

	if _, ok := c.Get(key); ok {
		t.Fatal("expected key to expire")
	}
}
