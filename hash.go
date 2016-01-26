package merkleTree

// Len returns the number of bytes in the hash, but after shifting off
// leading 0s from the length size of the hash
func (h Hash) Len() int {
	ret := len(h)
	for _, c := range h {
		if c == 0 {
			ret--
		} else {
			break
		}
	}
	return ret
}

func (h Hash) cmp(h2 Hash) int {
	if h.Len() < h2.Len() {
		return -1
	}
	if h.Len() > h2.Len() {
		return 1
	}
	for i, b := range h {
		b2 := h2[i]
		if b < b2 {
			return -1
		}
		if b > b2 {
			return 1
		}
	}
	// Equal in this case
	return 0
}

// Less determines if the receiver is less than the arg, after shifting off all
// leading 0 bytes and using big-endian byte ordering.
func (h Hash) Less(h2 Hash) bool {
	return h.cmp(h2) < 0
}

// Eq determines if the two hashes are equal.
func (h Hash) Eq(h2 Hash) bool {
	return h.cmp(h2) == 0
}
