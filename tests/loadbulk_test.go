package tests

import (
	"testing"
	"time"

	"github.com/crypolt/fastcache"
)

func TestLoadBulk(t *testing.T) {

	c := fastcache.New()

	keys := [][]byte{
		[]byte("k1"),
		[]byte("k2"),
		[]byte("k3"),
	}
	values := [][]byte{
		[]byte("v1"),
		[]byte("v2"),
		[]byte("v3"),
	}

	n := c.LoadBulk(keys, values, time.Second)
	if n != 3 {
		t.Fatalf("expected 3 inserted entries, got %d", n)
	}

	for i := range keys {
		v, ok := c.Get(keys[i])
		if !ok {
			t.Fatalf("missing key %s", keys[i])
		}
		if string(v) != string(values[i]) {
			t.Fatalf("unexpected value for key %s", keys[i])
		}
	}
}
