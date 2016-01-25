package merkleTree

import (
	"sort"
)

type sortedMap struct {
	list []KeyValuePair
}

func newSortedMapFromSortedList(l []KeyValuePair) *sortedMap {
	return &sortedMap{list: l}
}

func newSortedMapFromList(l []KeyValuePair) *sortedMap {
	ret := newSortedMapFromSortedList(l)
	ret.sort()
	return ret
}

func newSortedMapFromKeyAndValue(kp KeyValuePair) *sortedMap {
	return newSortedMapFromList([]KeyValuePair{kp})
}

type byKey []KeyValuePair

func (b byKey) Len() int           { return len(b) }
func (b byKey) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b byKey) Less(i, j int) bool { return b[i].Key.Less(b[j].Key) }

func (s *sortedMap) sort() {
	sort.Sort(byKey(s.list))
}

func (s *sortedMap) push(kp KeyValuePair) {
	s.list = append(s.list, kp)
}
