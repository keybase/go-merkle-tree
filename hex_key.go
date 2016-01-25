package merkleTree

func (h HexKey) Len() int {
	ret := len(h)
	for _, c := range h {
		if c == 0 {
			ret--
		}
	}
	return ret
}

func (h HexKey) cmp(h2 HexKey) int {
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

func (h HexKey) Less(h2 HexKey) bool {
	return h.cmp(h2) < 0
}

func (h HexKey) Eq(h2 HexKey) bool {
	return h.cmp(h2) == 0
}
