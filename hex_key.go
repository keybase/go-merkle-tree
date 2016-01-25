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

func (h HexKey) Less(h2 HexKey) bool {
	if h.Len() < h2.Len() {
		return true
	}
	if h.Len() > h2.Len() {
		return false
	}
	for i, b := range h {
		b2 := h2[i]
		if b < b2 {
			return true
		}
		if b > b2 {
			return false
		}
	}
	// Equal in this case
	return false
}
