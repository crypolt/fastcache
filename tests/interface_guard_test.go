package tests

import (
	"github.com/crypolt/fastcache"
	"github.com/crypolt/fastcache/internal/core"
)

// Compile-time interface checks.
var _ fastcache.Cache = (*core.Cache)(nil)
var _ fastcache.Cache = (*core.ShardedCache)(nil)
