package merkleTree

import (
	"sort"
)

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

func NewSortedMapFromKeyAndValue(kp KeyValuePair) *SortedMap {
	return NewSortedMapFromList([]KeyValuePair{kp})
}

type byKey []KeyValuePair

func (b byKey) Len() int           { return len(b) }
func (b byKey) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b byKey) Less(i, j int) bool { return b[i].Key.Less(b[j].Key) }

func (s *SortedMap) sort() {
	sort.Sort(byKey(s.list))
}

func (s *SortedMap) push(kp KeyValuePair) {
	s.list = append(s.list, kp)
}

func (s *SortedMap) exportToNode(h Hasher, typ NodeType, prevRoot Hash, level Level) (hash Hash, node Node, objExported []byte, err error) {
	if prevRoot != nil && level == Level(0) {
		node.PrevRoot = prevRoot
	}
	node.Type = typ
	node.Tab = s.list
	objExported, err = encodeToBytes(node)
	hash = h.Hash(objExported)
	return hash, node, objExported, err
}

func (s *SortedMap) binarySearch(k Hash) (ret int, eq bool) {
	beg := 0
	end := len(s.list) - 1

	for beg < end {
		mid := (end + beg) >> 1
		if s.list[mid].Key.Less(k) {
			beg = mid + 1
		} else {
			end = mid
		}
	}

	ret = beg
	if c := k.cmp(s.list[beg].Key); c > 0 {
		ret = beg + 1
	} else if c == 0 {
		eq = true
	}
	return ret, eq
}

func (s *SortedMap) replace(kvp KeyValuePair) *SortedMap {
	if len(s.list) > 0 {
		i, eq := s.binarySearch(kvp.Key)
		j := i
		if eq {
			j++
		}
		lst := append(s.list[0:i], kvp)
		lst = append(lst, s.list[j:]...)
		s.list = lst

	} else {
		s.list = []KeyValuePair{kvp}
	}
	return s
}
