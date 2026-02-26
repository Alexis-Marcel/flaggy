package engine

// MurmurHash3 (32-bit) — deterministic hash for consistent rollout bucketing.
// Given the same flag_key + entity_id, a user always lands in the same bucket.
func murmur3_32(data []byte, seed uint32) uint32 {
	const (
		c1 = 0xcc9e2d51
		c2 = 0x1b873593
	)

	h := seed
	nblocks := len(data) / 4

	// body
	for i := 0; i < nblocks; i++ {
		k := uint32(data[i*4]) |
			uint32(data[i*4+1])<<8 |
			uint32(data[i*4+2])<<16 |
			uint32(data[i*4+3])<<24

		k *= c1
		k = (k << 15) | (k >> 17)
		k *= c2

		h ^= k
		h = (h << 13) | (h >> 19)
		h = h*5 + 0xe6546b64
	}

	// tail
	tail := data[nblocks*4:]
	var k1 uint32
	switch len(tail) {
	case 3:
		k1 ^= uint32(tail[2]) << 16
		fallthrough
	case 2:
		k1 ^= uint32(tail[1]) << 8
		fallthrough
	case 1:
		k1 ^= uint32(tail[0])
		k1 *= c1
		k1 = (k1 << 15) | (k1 >> 17)
		k1 *= c2
		h ^= k1
	}

	// finalization
	h ^= uint32(len(data))
	h ^= h >> 16
	h *= 0x85ebca6b
	h ^= h >> 13
	h *= 0xc2b2ae35
	h ^= h >> 16

	return h
}

// RolloutBucket returns a bucket 0-99 for the given flag key and entity ID.
// Deterministic: same inputs always produce the same bucket.
func RolloutBucket(flagKey, entityID string) int {
	key := flagKey + ":" + entityID
	h := murmur3_32([]byte(key), 0)
	return int(h % 100)
}

// InRollout checks if the entity falls within the rollout percentage.
// percentage=0 means no rollout (always false), percentage=100 means always true.
func InRollout(flagKey, entityID string, percentage int) bool {
	if percentage <= 0 {
		return false
	}
	if percentage >= 100 {
		return true
	}
	return RolloutBucket(flagKey, entityID) < percentage
}
