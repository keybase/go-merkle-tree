package merkleTree

import (
	"encoding/binary"
	"fmt"
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

func rightShift(b []byte, n uint) {
	carry := byte(0)
	for i := 0; i < len(b); i++ {
		nxtCarry := b[i] & ((1 << n) - 1)
		b[i] = (b[i] >> n) | (carry << (8 - n))
		carry = nxtCarry
	}
}

func bitslice(b []byte, start ChildIndex, end ChildIndex) Prefix {
	tot := end - start

	// Shift off the first start/8 bytes
	b = b[(start >> 3):]
	// Clamp off all bytes after tot/8
	b = b[0:div8roundUp(tot)]

	// Output vector
	out := make([]byte, len(b))
	copy(out, b)

	startrem := start & 0x7
	if startrem > 0 && len(out) > 0 {
		fmt.Printf("pre LSH %v %d %d\n", out, (8 - startrem), (1 << (8 - startrem)))
		out[0] = out[0] & ((1 << (8 - startrem)) - 1)
		fmt.Printf("post LSH %v\n", out)
	}

	// We might have to shift the whole buffer over rsh bytes
	if rsh := (end & 0x7); rsh > 0 {
		rightShift(out, uint(8-rsh))
	}

	return Prefix(out)

}

func (c Config) formatPrefix(index ChildIndex) Prefix {
	ret := make([]byte, 4)
	binary.BigEndian.PutUint32(ret, uint32(index))
	return Prefix(ret[(4 - c.c):])
}
