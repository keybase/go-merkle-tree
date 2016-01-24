package merkleTree

import (
	"sort"
)

type HexKey []byte

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

type KeyValuePair struct {
	Key   HexKey
	Value interface{}
}

type SortedMap struct {
	list []KeyValuePair
}

func NewSortedMapFromSortedList(l []KeyValuePair) *SortedMap {
	return &SortedMap{list: l}
}

func NewSortedMapFromList(l []KeyValuePair) *SortedMap {
	ret := NewSortedMapFromSortedList(l)
	ret.sort()
	return ret
}

type byKey []KeyValuePair

func (b byKey) Len() int           { return len(b) }
func (b byKey) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b byKey) Less(i, j int) bool { return b[i].Key.Less(b[j].Key) }

func (s *SortedMap) sort() {
	sort.Sort(byKey(s.list))
}
