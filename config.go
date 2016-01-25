package merkleTree

type Config struct {
	// The number of children per node
	M uint

	// The maximum number of leaves before we split
	N uint

	// If we have M children per node, how many binary chars does it
	// take to represent it?
	C uint
}

func log256(y uint) uint {
	ret := uint(0)
	for y > 1 {
		y = (y >> 8)
		ret++
	}
	return ret
}

func NewConfig(M,N uint) Config{
	return Config{M : M, N : N, C : log256(M) }
}

func (c Config) prefixAtLevel(l Level, h Hash) []byte {
	return h[(uint(l)*c.C):((uint(l)+1)*c.C)]
}

func (c Config) prefixThroughLevel(l Level, h Hash) []byte {
	return h[0:((uint(l)+1)*c.C)]
}