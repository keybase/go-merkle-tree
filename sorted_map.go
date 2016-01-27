package merkleTree

import (
	"sort"
)


// SortedMap is a list of KeyValuePairs, kept in sorted order so
// that binary search can work.
type SortedMap struct {
	list []KeyValuePair
}

// NewSortedMap makes an empty sorted map.
func NewSortedMap() *SortedMap {
	return &SortedMap{}
}

// NewSortedMapFromSortedList just wraps the given sorted list and
// doesn't check that it's sorted.  So don't pass it an unsorted list.
func NewSortedMapFromSortedList(l []KeyValuePair) *SortedMap {
	return &SortedMap{list: l}
}

func newSortedMapFromNode(n *Node) *SortedMap {
	return NewSortedMapFromSortedList(n.Tab)
}

// NewSortedMapFromList makes a sorted map from an unsorted list
// of KeyValuePairs
func NewSortedMapFromList(l []KeyValuePair) *SortedMap {
	ret := NewSortedMapFromSortedList(l)
	ret.sort()
	return ret
}

func newSortedMapFromKeyAndValue(kp KeyValuePair) *SortedMap {
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

func (s *SortedMap) find(k Hash) *KeyValuePair {
	i, eq := s.binarySearch(k)
	if !eq {
		return nil
	}
	ret := s.at(ChildIndex(i))
	return &ret
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

func (s *SortedMap) Len() ChildIndex {
	return ChildIndex(len(s.list))
}

func (s *SortedMap) at(i ChildIndex) KeyValuePair {
	return s.list[i]
}

func (s *SortedMap) slice(begin, end ChildIndex) *SortedMap {
	return NewSortedMapFromList(s.list[begin:end])
}

func (n *Node) findValueInLeaf(h Hash) interface{} {
	kvp := newSortedMapFromNode(n).find(h)
	if kvp == nil {
		return nil
	}
	return kvp.Value
}

func (n *Node) findChildByPrefix(p Prefix) (Hash, error) {
	kvp := newSortedMapFromNode(n).find(Hash(p))
	if kvp == nil {
		return nil, nil
	}
	b, ok := (kvp.Value).(Hash)
	if !ok {
		return nil, BadChildPointerError{kvp.Value}
	}
	return Hash(b), nil
}
