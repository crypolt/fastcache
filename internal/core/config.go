package core

import (
	"errors"
	"time"
)

// Config defines cache configuration parameters.
//
// This config is the single source of truth for:
//   - default values
//   - validation rules
//   - internal behavior switches
//
// The fastcache package MUST NOT modify config values
// except calling ApplyDefaults() and Validate().
type Config struct {
	// ========================
	// Core memory configuration
	// ========================

	// Capacity is the total number of cache entries.
	// Rounded up to the nearest power of two.
	Capacity int

	// ArenaSize is the total memory size (in bytes)
	// used to store keys and values.
	//
	// If ArenaSize == 0, arena is allocated lazily
	// and grows on demand.
	ArenaSize int

	// Shards defines the number of cache shards.
	// 0 or 1 means single-shard mode.
	// Must be a power of two.
	Shards int

	// ========================
	// TTL & expiration
	// ========================

	// DefaultTTL is applied when Set is called with ttl == 0.
	// Zero means entries never expire by default.
	DefaultTTL time.Duration

	// EnableTTL enables TTL checks.
	EnableTTL bool

	// ========================
	// Eviction & memory behavior
	// ========================

	// EnableEviction enables clock-based eviction.
	EnableEviction bool

	// ResetOnFull resets the entire cache when arena memory is exhausted.
	// Faster than eviction, but drops all data.
	ResetOnFull bool

	// PreallocArena forces arena memory preallocation.
	PreallocArena bool

	// ========================
	// Read/write modes
	// ========================

	// ReadOnly disables all write operations.
	ReadOnly bool

	// ========================
	// Observability & debug
	// ========================

	EnableStats bool
	Debug       bool

	// ========================
	// Distributed / persistence (future)
	// ========================

	NodeID            string
	EnableSync        bool
	EnablePersistence bool
	DataDir           string
	SnapshotInterval  time.Duration
}

// ApplyDefaults fills zero-value fields with sane defaults.
func (c Config) ApplyDefaults() Config {
	if c.Capacity == 0 {
		c.Capacity = 65536
	}
	if c.ArenaSize == 0 {
		c.ArenaSize = 32 << 20 // 32 MB
	}
	if c.Shards == 0 {
		c.Shards = 1
	}
	if c.DefaultTTL < 0 {
		c.DefaultTTL = 0
	}
	if !c.EnableTTL {
		c.EnableTTL = true
	}
	if !c.EnableEviction && !c.ResetOnFull {
		c.EnableEviction = true
	}
	if c.SnapshotInterval == 0 {
		c.SnapshotInterval = 30 * time.Second
	}
	return c
}

// Validate checks configuration correctness.
// Returns error instead of panicking.
func (c Config) Validate() error {
	if c.Capacity <= 0 {
		return errors.New("capacity must be > 0")
	}
	if c.ArenaSize < 0 {
		return errors.New("arenaSize must be >= 0")
	}
	if c.Shards < 0 {
		return errors.New("shards must be >= 0")
	}
	if c.Shards > 1 && (c.Shards&(c.Shards-1)) != 0 {
		return errors.New("shards must be a power of two")
	}
	if c.EnablePersistence && c.DataDir == "" {
		return errors.New("dataDir must be set when persistence is enabled")
	}
	return nil
}
