package core

func (c *Cache) evict() {
	for {
		e := &c.entries[c.clock]
		if e.used == 0 {
			e.keyHash = 0
			return
		}
		e.used = 0
		c.clock = (c.clock + 1) & c.mask
	}
}
