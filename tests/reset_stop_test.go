package tests

import (
	"testing"
	"time"

	"github.com/crypolt/fastcache"
)

func TestReset(t *testing.T) {

	c := fastcache.New()

	c.Set([]byte("a"), []byte("b"), time.Second)
	c.Reset()

	if _, ok := c.Get([]byte("a")); ok {
		t.Fatal("expected cache to be empty after reset")
	}
}

func TestStop(t *testing.T) {

	c := fastcache.New()

	c.Stop()

	if ok := c.Set([]byte("x"), []byte("y"), time.Second); ok {
		t.Fatal("Set should fail after Stop")
	}

	if _, ok := c.Get([]byte("x")); ok {
		t.Fatal("Get should fail after Stop")
	}
}
