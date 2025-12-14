package core

import (
	"sync"
)

// Cache is a single-shard high-performance in-memory cache.
//
// Design goals:
//   - no map
//   - predictable memory usage
//   - O(1) get/set
//   - minimal GC pressure
//   - byte-level storage
//
// Cache is NOT intended to be used directly by end users.
// It is a low-level core primitive.
type Cache struct {
	mu      sync.Mutex // protects all fields below
	size    uint32     // number of slots (power of two)
	mask    uint32     // size-1, used for fast modulo
	clock   uint32     // clock hand for eviction
	entries []entry    // hash table entries
	arena   arena      // byte arena for keys and values
	stopped bool       // disables cache operations
}

// New creates a new single-shard cache.
//
// capacity   - number of entries (rounded to power of two)
// arenaSize  - total memory size for keys and values (bytes)
func New(capacity int, arenaSize int) *Cache {
	size := uint32(1)
	for size < uint32(capacity) {
		size <<= 1
	}

	var a arena
	if arenaSize > 0 {
		a = newArena(arenaSize)
	} else {
		a = newArenaLazy(0)
	}

	return &Cache{
		size:    size,
		mask:    size - 1,
		entries: make([]entry, size),
		arena:   a,
	}
}

// NewFromConfig creates a single-shard cache from validated config.
func NewFromConfig(cfg Config) *Cache {
	return New(cfg.Capacity, cfg.ArenaSize)
}
