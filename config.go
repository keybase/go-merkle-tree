package merkleTree

import (
	"fmt"
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
	l := level.ToChildIndex()
	return bitslice(h, l*c.c, (l+1)*c.c)
}

func div8roundUp(i ChildIndex) ChildIndex {
	return ((i + 7) >> 3)
}

func bitslice(b []byte, start ChildIndex, end ChildIndex) Prefix {
	tot := end - start
	bitrem := tot & 0x7

	// Shift off the first start/8 bytes
	b = b[(start >> 3):]
	// Clamp off all bytes after tot/8
	b = b[0:div8roundUp(tot)]

	out := make([]byte, len(b))

	if bitrem > 0 {
		fmt.Printf("pree -> %d %v\n", bitrem, b)
		carry := byte(0)
		for i := len(b) - 1; i >= 0; i-- {
			nxtCarry := (b[i] & (0xff << (8 - bitrem)))
			out[i] = (b[i] >> (8 - bitrem)) | carry
			carry = nxtCarry
		}
		fmt.Printf("post -> %v\n", out)
	} else {
		copy(out, b)
	}

	return Prefix(out)

}

func (c Config) formatPrefix(index ChildIndex) Prefix {
	ret := make([]byte, 4)
	binary.BigEndian.PutUint32(ret, uint32(index))
	return Prefix(ret[(4 - c.c):])
}
