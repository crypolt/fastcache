package tests

import (
	"testing"
	"time"

	"github.com/crypolt/fastcache"
)

func TestSingleCache_SetGetDelete(t *testing.T) {

	var c fastcache.Cache = fastcache.New()

	key := []byte("key")
	val := []byte("value")

	ok := c.Set(key, val, time.Second)
	if !ok {
		t.Fatal("Set failed")
	}

	got, found := c.Get(key)
	if !found {
		t.Fatal("Get failed")
	}

	if string(got) != "value" {
		t.Fatalf("unexpected value: %s", got)
	}

	c.Delete(key)

	if _, found := c.Get(key); found {
		t.Fatal("expected key to be deleted")
	}
}
