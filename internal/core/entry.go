package core

type entry struct {
	keyHash uint64

	keyLen uint16
	valLen uint32

	keyOff uint32
	valOff uint32

	expire int64
	used   uint8
}
