package core

func hash(b []byte) uint64 {
	var h uint64 = 1469598103934665603
	for i := 0; i < len(b); i++ {
		h ^= uint64(b[i])
		h *= 1099511628211
	}
	return h
}
