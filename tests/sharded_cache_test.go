package tests

import (
	"testing"
	"time"

	"github.com/crypolt/fastcache"
)

func TestShardedCache_SetGet(t *testing.T) {

	var c fastcache.Cache = fastcache.New()

	keys := [][]byte{
		[]byte("a"),
		[]byte("b"),
		[]byte("c"),
		[]byte("d"),
	}

	for _, k := range keys {
		if !c.Set(k, k, time.Second) {
			t.Fatalf("Set failed for key %s", k)
		}
	}

	for _, k := range keys {
		v, ok := c.Get(k)
		if !ok {
			t.Fatalf("Get failed for key %s", k)
		}
		if string(v) != string(k) {
			t.Fatalf("unexpected value for key %s", k)
		}
	}
}
