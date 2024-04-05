package murmur

import (
	"encoding/binary"
	"math/big"
)

// Hash computes the MurmurHash3 of a string and returns a URL-safe
// base64-encoded string
func Hash(key string) string {
	return toBase62(toBytes(Hash32(key)))
}

// toBase62 converts a UUID to a base62 string.
// Based on: https://ucarion.com/go-base62
func toBase62(key []byte) string {
	var i big.Int
	i.SetBytes(key)
	return i.Text(62)
}

// Hash32 computes the MurmurHash3 of a string and returns a uint32. Mostly
// used for testing.
func Hash32(key string) uint32 {
	return murmur32([]byte(key), 0)
}

// murmur32 computes the MurmurHash3 of a byte slice with the given seed.
// Based on: https://en.wikipedia.org/wiki/MurmurHash#Algorithm
func murmur32(key []byte, seed uint32) uint32 {
	var h, k uint32
	h = seed
	length := len(key)

	// Read in groups of 4
	for len(key) >= 4 {
		k = binary.LittleEndian.Uint32(key[:4])
		key = key[4:]
		h ^= scramble(k)
		h = (h << 13) | (h >> 19)
		h = h*5 + 0xe6546b64
	}

	// Read the rest
	k = 0
	for i, v := range key {
		k |= uint32(v) << (8 * uint(i))
	}
	h ^= scramble(k)

	// Finalize
	h ^= uint32(length)
	h ^= h >> 16
	h *= 0x85ebca6b
	h ^= h >> 13
	h *= 0xc2b2ae35
	h ^= h >> 16

	return h
}

// scramble scrambles a 32-bit integer.
func scramble(k uint32) uint32 {
	k *= 0xcc9e2d51
	k = (k << 15) | (k >> 17)
	k *= 0x1b873593
	return k
}

func toBytes(value uint32) []byte {
	buf := make([]byte, 4)                    // Create a byte slice of length 4, as uint32 is 4 bytes
	binary.LittleEndian.PutUint32(buf, value) // Convert uint32 to byte slice
	return buf
}
