package fastcache

import "time"

// Cache defines the public cache interface.
//
// Implementations may be single-shard or sharded,
// but must obey the same semantic contract.
type Cache interface {
	// Set inserts or overwrites a value associated with the key.
	//
	// ttl defines time-to-live for the entry.
	// If ttl is zero, DefaultTTL from configuration may be applied.
	//
	// Returns false if the value was not stored
	// (e.g. cache is read-only or stopped).
	Set(key, value []byte, ttl time.Duration) bool

	// Get returns a value associated with the key.
	//
	// The returned byte slice MUST be treated as read-only.
	// It may point directly to internal cache memory.
	Get(key []byte) ([]byte, bool)

	// Delete removes a key from the cache.
	// It is a no-op if the key does not exist.
	Delete(key []byte)

	// LoadBulk loads multiple key-value pairs into the cache.
	//
	// keys and values slices must have equal length.
	// Returns the number of successfully stored entries.
	LoadBulk(keys [][]byte, values [][]byte, ttl time.Duration) int

	// Reset completely clears the cache and drops all data.
	Reset()

	// Stop permanently disables the cache.
	// All subsequent operations will fail.
	Stop()
}

// ShardedCache extends Cache with shard-level access.
//
// This interface is optional and should only be used
// when shard-level operations are required.
type ShardedCache interface {
	Cache

	// Shards returns the total number of shards.
	Shards() int

	// Shard returns a cache instance responsible for the given key.
	//
	// Returned cache is shard-local and MUST NOT be shared
	// across goroutines without proper synchronization.
	Shard(key []byte) Cache

	// ForEachShard iterates over all shards.
	//
	// Useful for maintenance operations, metrics collection,
	// or bulk warm-up.
	ForEachShard(fn func(idx int, shard Cache))
}
