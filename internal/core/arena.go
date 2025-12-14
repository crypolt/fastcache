package core

const (
	// defaultChunkSize defines size of a single arena chunk.
	// 1 MB is a good balance between allocation cost and fragmentation.
	defaultChunkSize = 1 << 20 // 1 MB
)

type arena struct {
	buf       []byte
	pos       uint32
	chunkSize uint32
	curChunk  uint32
}

func newArena(size int) arena {
	return arena{
		buf: make([]byte, size),
	}
}

func newArenaLazy(chunkSize int) arena {
	if chunkSize <= 0 {
		chunkSize = defaultChunkSize
	}

	return arena{
		chunkSize: uint32(chunkSize),
	}
}

func (a *arena) alloc(n uint32) (uint32, bool) {
	if int(a.pos+n) > len(a.buf) {
		return 0, false
	}
	off := a.pos
	a.pos += n
	return off, true
}

func (a *arena) reset() {
	a.pos = 0
}
