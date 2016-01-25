package merkleTree

func (h Hash) Len() int {
	ret := len(h)
	for _, c := range h {
		if c == 0 {
			ret--
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

func (h Hash) Less(h2 Hash) bool {
	return h.cmp(h2) < 0
}

func (h Hash) Eq(h2 Hash) bool {
	return h.cmp(h2) == 0
}
