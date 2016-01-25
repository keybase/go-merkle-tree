package merkleTree

import (
	"encoding/binary"
)

type Config struct {
	// The number of children per node
	M ChildIndex

	// The maximum number of leaves before we split
	N ChildIndex

	// If we have M children per node, how many binary chars does it
	// take to represent it?
	C ChildIndex
}

func log256(y ChildIndex) ChildIndex {
	ret := ChildIndex(0)
	for y > 1 {
		y = (y >> 8)
		ret++
	}
	return ret
}

func NewConfig(M, N ChildIndex) Config {
	return Config{M: M, N: N, C: log256(M)}
}

func (c Config) prefixAtLevel(level Level, h Hash) Prefix {
	l := level.ToChildIndex()
	return Prefix(h[(l * c.C):((l + 1) * c.C)])
}

func (c Config) prefixThroughLevel(level Level, h Hash) Prefix {
	l := level.ToChildIndex()
	return Prefix(h[0:((l + 1) * c.C)])
}

func (c Config) formatPrefix(index ChildIndex) Prefix {
	ret := make([]byte, 4)
	binary.BigEndian.PutUint32(ret, uint32(index))
	return Prefix(ret)
}
