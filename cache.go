package fastcache

import (
	"log"

	"github.com/crypolt/fastcache/internal/core"
)

// New creates a cache with default configuration.
func New() Cache {
	return NewWithConfig(core.Config{})
}

// NewWithConfig creates a cache with custom configuration.
// Panic-safe.
func NewWithConfig(cfg core.Config) Cache {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[fastcache] panic recovered in New(): %v", r)
		}
	}()

	cfg = cfg.ApplyDefaults()
	if err := cfg.Validate(); err != nil {
		panic(err)
	}

	if cfg.Shards <= 1 {
		return core.NewFromConfig(cfg)
	}

	return core.NewShardedFromConfig(cfg)
}
