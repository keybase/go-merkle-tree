package merkleTree

import (
	"encoding/binary"
)

type Hasher interface {
	Hash([]byte) Hash
}

type Config struct {
	// A hasher is used to compute hashes in this configuration
	hasher Hasher

	// The number of children per node
	m ChildIndex

	// The maximum number of leaves before we split
	n ChildIndex

	// If we have M children per node, how many binary chars does it
	// take to represent it?
	c ChildIndex
}

func log256(y ChildIndex) ChildIndex {
	ret := ChildIndex(0)
	for y > 1 {
		y = (y >> 8)
		ret++
	}
	return ret
}

func NewConfig(h Hasher, m ChildIndex, n ChildIndex) Config {
	return Config{hasher: h, m: m, n: n, c: log256(m)}
}

func (c Config) prefixAtLevel(level Level, h Hash) Prefix {
	l := level.ToChildIndex()
	return Prefix(h[(l * c.c):((l + 1) * c.c)])
}

func (c Config) prefixThroughLevel(level Level, h Hash) Prefix {
	l := level.ToChildIndex()
	return Prefix(h[0:((l + 1) * c.c)])
}

func (c Config) formatPrefix(index ChildIndex) Prefix {
	ret := make([]byte, 4)
	binary.BigEndian.PutUint32(ret, uint32(index))
	return Prefix(ret[(4 - c.c):])
}
