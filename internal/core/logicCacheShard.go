package core

import "time"

// ShardedCache is a cache composed of multiple independent cache shards.
//
// Sharding reduces lock contention and improves scalability
// under high concurrency by distributing keys across shards.
type ShardedCache struct {
	shards []*Cache
	mask   uint64
}

// shard returns the cache shard responsible for the given key.
func (s *ShardedCache) shard(key []byte) *Cache {
	h := hash(key)
	return s.shards[h&s.mask]
}

//
// Cache interface implementation
//

// Set inserts or overwrites a value in the appropriate shard.
func (s *ShardedCache) Set(key, val []byte, ttl time.Duration) bool {
	return s.shard(key).Set(key, val, ttl)
}

// Get returns a value from the appropriate shard.
//
// The returned byte slice MUST be treated as read-only.
func (s *ShardedCache) Get(key []byte) ([]byte, bool) {
	return s.shard(key).Get(key)
}

// Delete removes a key from the appropriate shard.
func (s *ShardedCache) Delete(key []byte) {
	s.shard(key).Delete(key)
}

// LoadBulk loads multiple key-value pairs into the cache.
//
// Returns the number of successfully inserted entries.
func (s *ShardedCache) LoadBulk(
	keys [][]byte,
	values [][]byte,
	ttl time.Duration,
) int {
	n := 0
	for i := 0; i < len(keys); i++ {
		if s.Set(keys[i], values[i], ttl) {
			n++
		}
	}
	return n
}

// Reset clears all shards and drops all cached data.
func (s *ShardedCache) Reset() {
	for _, sh := range s.shards {
		sh.Reset()
	}
}

// Stop permanently disables all shards.
func (s *ShardedCache) Stop() {
	for _, sh := range s.shards {
		sh.Stop()
	}
}

//
// ShardedCache interface extensions
//

// Shards returns the total number of shards.
func (s *ShardedCache) Shards() int {
	return len(s.shards)
}

// Shard returns the shard responsible for the given key.
//
// The returned cache is shard-local.
func (s *ShardedCache) Shard(key []byte) *Cache {
	return s.shard(key)
}

// ForEachShard iterates over all shards.
//
// Useful for maintenance tasks, statistics collection,
// warm-up, or graceful shutdown.
func (s *ShardedCache) ForEachShard(fn func(idx int, shard *Cache)) {
	for i, sh := range s.shards {
		fn(i, sh)
	}
}

//
// Constructor
//

// NewShardedFromConfig creates a sharded cache from a validated Config.
//
// Config MUST be validated and defaulted before calling this function.
func NewShardedFromConfig(cfg Config) *ShardedCache {
	shards := cfg.Shards
	if shards <= 1 {
		panic("sharded cache requires shards > 1")
	}

	perShardCap := cfg.Capacity / shards
	perShardArena := 0
	if cfg.ArenaSize > 0 {
		perShardArena = cfg.ArenaSize / shards
	}

	s := &ShardedCache{
		shards: make([]*Cache, shards),
		mask:   uint64(shards - 1),
	}

	for i := 0; i < shards; i++ {
		s.shards[i] = New(perShardCap, perShardArena)
	}

	return s
}
