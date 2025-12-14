# fastcache

`fastcache` is a high-performance in-memory cache library for Go, designed for
low latency, predictable memory usage, and minimal garbage collection overhead.

The package targets latency-sensitive backend paths where map-based caches or
external systems introduce unnecessary overhead.

The API is intentionally small and explicit, allowing the cache to be used as a
low-level building block in performance-critical code.

The package is intended to be published and consumed via pkg.go.dev and follows
standard Go documentation and API stability conventions.

---

## Installation

```bash
go get github.com/crypolt/fastcache
```

---

## Quick start (no configuration)

In most cases, no configuration is required. A cache can be created using
`fastcache.New()`, which applies safe and well-tested defaults.

```go
package main

import (
	"fmt"
	"time"

	"github.com/crypolt/fastcache"
)

func main() {
	cache := fastcache.New()

	cache.Set([]byte("key"), []byte("value"), time.Minute)

	if v, ok := cache.Get([]byte("key")); ok {
		fmt.Println(string(v))
	}
}
```

---

## Default configuration

When `fastcache.New()` is used, the cache is created with the following default
settings:

- Capacity: `65536` entries
- Arena size: `32 MB`
- Shards: `1` (single-shard mode)
- TTL support: enabled
- Default TTL: `0` (entries do not expire unless TTL is provided)
- Eviction: enabled (clock-based)
- Reset on full: disabled
- Read-only mode: disabled

These defaults are applied internally via `ApplyDefaults()` and validated before
initialization.

Equivalent configuration:

```go
Config{
	Capacity:       65536,
	ArenaSize:      32 << 20,
	Shards:         1,
	DefaultTTL:     0,
	EnableTTL:      true,
	EnableEviction: true,
	ResetOnFull:    false,
	ReadOnly:       false,
	Debug:          false,
}
```

---

## Public API

### Cache interface

```go
type Cache interface {
	Set(key, value []byte, ttl time.Duration) bool
	Get(key []byte) ([]byte, bool)
	Delete(key []byte)
	LoadBulk(keys [][]byte, values [][]byte, ttl time.Duration) int
	Reset()
	Stop()
}
```

---

## Method semantics

### Set

```go
Set(key, value []byte, ttl time.Duration) bool
```

Stores or overwrites a value associated with the given key.

- `ttl` defines the time-to-live for the entry
- if `ttl == 0`, `DefaultTTL` from configuration is applied
- if both are zero, the entry does not expire
- returns `false` if the cache is stopped or in read-only mode

---

### Get

```go
Get(key []byte) ([]byte, bool)
```

Returns the value associated with the given key.

- the returned byte slice is read-only
- the slice may reference internal cache memory
- callers must copy the data if long-term ownership is required
- returns `false` if the key does not exist or is expired

---

### Delete

```go
Delete(key []byte)
```

Removes the key from the cache.

- no-op if the key does not exist

---

### LoadBulk

```go
LoadBulk(keys [][]byte, values [][]byte, ttl time.Duration) int
```

Loads multiple key-value pairs into the cache.

- `keys` and `values` must have equal length
- `ttl` is applied to all entries
- returns the number of successfully stored entries

---

### Reset

```go
Reset()
```

Immediately drops all cached data and resets internal state.

---

### Stop

```go
Stop()
```

Permanently disables the cache.

- all subsequent operations fail
- intended for controlled shutdown scenarios

---

## Sharded cache

For high-concurrency workloads, `fastcache` supports sharded operation.

```go
type ShardedCache interface {
	Cache
	Shards() int
	Shard(key []byte) Cache
	ForEachShard(fn func(idx int, shard Cache))
}
```

Notes:

- sharding reduces lock contention
- the number of shards must be a power of two
- shard-level cache instances are not goroutine-safe
- shard-level access is intended for advanced or internal use

---

## Architecture

```text
fastcache
├── fastcache.go        // public constructors
├── interface.go        // Cache and ShardedCache interfaces
├── internal/
│   └── core/
│       ├── cache.go    // single-shard cache implementation
│       ├── shard.go    // sharded cache wrapper
│       ├── arena.go    // byte arena allocator
│       ├── entry.go    // hash table entries
│       └── config.go   // configuration and validation
```

Design principles:

- no Go maps
- no per-entry allocations
- contiguous memory storage
- predictable memory behavior
- minimal GC interaction

---

## Limitations

- not a distributed cache
- no built-in persistence
- not a Redis replacement
- shard-level instances are not goroutine-safe

---

## License

MIT

---

## Contact

Vladislav Poltavskiy  
Telegram: @crypolt
