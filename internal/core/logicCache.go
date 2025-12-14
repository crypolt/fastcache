package core

import "time"

// Set inserts or overwrites a value in the cache.
//
// Returns false if the cache is stopped or if memory allocation failed.
func (c *Cache) Set(key, val []byte, ttl time.Duration) bool {
	if c.stopped {
		return false
	}

	h := hash(key)

	c.mu.Lock()
	defer c.mu.Unlock()

	// Calculate expiration timestamp under lock
	var exp int64
	if ttl > 0 {
		exp = time.Now().Add(ttl).UnixNano()
	} else {
		exp = 0 // no TTL
	}

	i := uint32(h) & c.mask

	for {
		e := &c.entries[i]

		// Empty slot or overwrite existing key
		if e.keyHash == 0 || e.keyHash == h {
			koff, ok := c.arena.alloc(uint32(len(key)))
			if !ok {
				c.evict()
				c.arena.reset()
				return false
			}

			voff, ok := c.arena.alloc(uint32(len(val)))
			if !ok {
				c.evict()
				c.arena.reset()
				return false
			}

			// Copy key/value into arena
			copy(c.arena.buf[koff:koff+uint32(len(key))], key)
			copy(c.arena.buf[voff:voff+uint32(len(val))], val)

			*e = entry{
				keyHash: h,
				keyLen:  uint16(len(key)),
				valLen:  uint32(len(val)),
				keyOff:  koff,
				valOff:  voff,
				expire:  exp,
				used:    1,
			}

			return true
		}

		// Linear probing
		i = (i + 1) & c.mask
	}
}

// Get returns a value by key.
//
// The returned byte slice points directly into the internal arena memory.
// The caller MUST treat it as read-only.
func (c *Cache) Get(key []byte) ([]byte, bool) {
	if c.stopped {
		return nil, false
	}

	h := hash(key)

	c.mu.Lock()
	defer c.mu.Unlock()

	i := uint32(h) & c.mask

	for {
		e := &c.entries[i]

		// Empty slot means key not found.
		if e.keyHash == 0 {
			return nil, false
		}

		if e.keyHash == h {
			// Check TTL expiration.
			if e.expire > 0 && time.Now().UnixNano() > e.expire {
				e.keyHash = 0
				return nil, false
			}

			// Mark entry as recently used.
			e.used = 1

			return c.arena.buf[e.valOff : e.valOff+e.valLen], true
		}

		i = (i + 1) & c.mask
	}
}

// Delete removes a key from the cache.
// It is a no-op if the key does not exist.
func (c *Cache) Delete(key []byte) {
	h := hash(key)

	c.mu.Lock()
	defer c.mu.Unlock()

	i := uint32(h) & c.mask
	for {
		e := &c.entries[i]
		if e.keyHash == 0 {
			return
		}
		if e.keyHash == h {
			e.keyHash = 0
			return
		}
		i = (i + 1) & c.mask
	}
}

// Reset completely clears the cache and its memory arena.
// All stored data is lost.
func (c *Cache) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for i := range c.entries {
		c.entries[i] = entry{}
	}
	c.arena.reset()
	c.clock = 0
}

// Stop permanently disables the cache.
// All subsequent operations will fail.
func (c *Cache) Stop() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.stopped = true
}

// LoadBulk loads multiple key-value pairs into the cache.
//
// keys and values slices must have equal length.
// Returns the number of successfully inserted entries.
func (c *Cache) LoadBulk(
	keys [][]byte,
	values [][]byte,
	ttl time.Duration,
) int {
	n := 0

	limit := len(keys)
	if len(values) < limit {
		limit = len(values)
	}

	for i := 0; i < limit; i++ {
		if c.Set(keys[i], values[i], ttl) {
			n++
		}
	}

	return n
}
