package merkletree

import (
	"encoding/binary"
	"fmt"
)

// Hasher is an interface for hashing MerkleTree data structures into their
// cryptographic hashes.
type Hasher interface {
	Hash([]byte) Hash
}

// ValueConstructor is an interface for constructing values, so that typed
// values can be pulled out of the Merkle Tree. We are of course assuming
// that there is only one type of Value at the leaves, which makes sense.
type ValueConstructor interface {
	// Construct a new template empty value for the leaf, so that the
	// Unmarshalling routine has the correct type template.
	Construct() any
}

// Config defines the shape of the MerkleTree.
type Config struct {
	// A hasher is used to compute hashes in this configuration
	hasher Hasher

	// The number of children per node
	m ChildIndex

	// The maximum number of leaves before we split
	n ChildIndex

	// If we have M children per node, how many bits does it take to represent it?
	c ChildIndex

	// Construct a new object to unmarshal values into
	v ValueConstructor
}

func log2(y ChildIndex) ChildIndex {
	ret := ChildIndex(0)
	for y > 1 {
		y = (y >> 1)
		ret++
	}
	return ret
}

// NewConfig makes a new config object. Pass it a Hasher
// (though we suggest sha512.Sum512); a parameter `m` which
// is the number of children per interior node (we recommend 256),
// and `n`, the maximum number of entries in a leaf before a
// new level of the tree is introduced.
//
// Panics if m cannot branch or n cannot hold any entries.
func NewConfig(h Hasher, m ChildIndex, n ChildIndex, v ValueConstructor) Config {
	// These degenerate configurations cannot reduce a non-empty subtree and
	// would recurse forever once a leaf needs splitting.
	if m < 2 {
		panic(fmt.Sprintf("invalid config: m must be at least 2, got %d", m))
	}
	if n == 0 {
		panic("invalid config: n must be positive")
	}
	return Config{hasher: h, m: m, n: n, c: log2(m), v: v}
}

// maxKeyLevels returns the number of complete child-index chunks available in
// h. Tree traversal must not try to descend through an interior node once all
// of these chunks have been consumed.
func (c Config) maxKeyLevels(h Hash) Level {
	if c.c == 0 {
		return 0
	}
	return Level((uint64(len(h)) * 8) / uint64(c.c))
}

func (c Config) prefixAndIndexAtLevel(level Level, h Hash) (Prefix, ChildIndex, error) {
	maxLevel := c.maxKeyLevels(h)
	if level >= maxLevel {
		return nil, 0, fmt.Errorf("max tree depth %d exceeded", maxLevel)
	}
	prfx, ci := bitslice(h, int(c.c), int(level)) //nolint:gosec // G115: Level is bounded by config
	return prfx, ChildIndex(ci), nil
}

func (c Config) prefixAtLevel(level Level, h Hash) Prefix {
	ret, _, _ := c.prefixAndIndexAtLevel(level, h)
	return ret
}

func (c Config) PrefixAndIndexAtLevel(level Level, h Hash) (Prefix, ChildIndex) {
	prfx, ci, _ := c.prefixAndIndexAtLevel(level, h)
	return prfx, ci
}

func (c Config) PrefixAtLevel(level Level, h Hash) Prefix {
	return c.prefixAtLevel(level, h)
}

func (c Config) formatPrefix(index ChildIndex) Prefix {
	ret := make([]byte, 4)
	binary.BigEndian.PutUint32(ret, uint32(index))
	return Prefix(ret[(4 - (7+c.c)/8):])
}
