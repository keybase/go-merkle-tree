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

	// If we have M children per node, how many bits does it take to represent it?
	c ChildIndex
}

func log2(y ChildIndex) ChildIndex {
	ret := ChildIndex(0)
	for y > 1 {
		y = (y >> 1)
		ret++
	}
	return ret
}

func NewConfig(h Hasher, m ChildIndex, n ChildIndex) Config {
	return Config{hasher: h, m: m, n: n, c: log2(m)}
}

func (c Config) prefixAtLevel(level Level, h Hash) Prefix {
	return bitslice(h, c.c, level)
}

func div8roundUp(i ChildIndex) ChildIndex {
	return ((i + 7) >> 3)
}

func (c Config) formatPrefix(index ChildIndex) Prefix {
	ret := make([]byte, 4)
	binary.BigEndian.PutUint32(ret, uint32(index))
	return Prefix(ret[(4 - c.c):])
}
