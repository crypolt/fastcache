// Package fastcache provides a high-performance in-memory cache for Go.
//
// The package is designed for low latency, predictable memory usage, and
// minimal garbage collection overhead. It is intended for latency-sensitive
// backend paths where map-based caches or external systems introduce
// unnecessary overhead.
//
// fastcache can be used without any configuration via fastcache.New(),
// which applies safe and well-tested defaults suitable for most workloads.
//
// The API is intentionally small and explicit, allowing the cache to be used
// as a low-level building block in performance-critical code.
//
// Returned values from Get must be treated as read-only. They may reference
// internal cache memory, and callers must copy the data if long-term ownership
// is required.
//
// The package is intended to be published and consumed via pkg.go.dev and
// follows standard Go documentation and API stability conventions.
package fastcache
