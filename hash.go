package merkletree

import "crypto/sha512"

// Len returns the number of bytes in the hash, but after shifting off
// leading 0s from the length size of the hash
func (h Hash) Len() int {
	for idx, x := range h {
		if x != 0 {
			return len(h) - idx
		}
	}
	return 0
}

func (h Hash) cmp(h2 Hash) int {
	if h.Len() < h2.Len() {
		return -1
	}
	if h.Len() > h2.Len() {
		return 1
	}

	// Preserve the historical byte ordering after the significant-length
	// comparison. Only compare the common prefix so that inputs which used to
	// run past the end of h2 are handled without changing any previously
	// defined comparison result.
	n := min(len(h2), len(h))
	for i := 0; i < n; i++ {
		if h[i] < h2[i] {
			return -1
		}
		if h[i] > h2[i] {
			return 1
		}
	}
	return 0
}

// Less determines if the receiver precedes the argument in the tree's
// historical ordering: significant length first, then the original bytes.
func (h Hash) Less(h2 Hash) bool {
	return h.cmp(h2) < 0
}

// Eq determines if the two hashes are equal.
func (h Hash) Eq(h2 Hash) bool {
	return h.cmp(h2) == 0
}

// SHA512Hasher is a simple SHA512 hash function application
type SHA512Hasher struct{}

// Hash the data
func (s SHA512Hasher) Hash(b []byte) Hash {
	tmp := sha512.Sum512(b)
	return Hash(tmp[:])
}
