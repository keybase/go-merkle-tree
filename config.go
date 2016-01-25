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
