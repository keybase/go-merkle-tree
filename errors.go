package merkleTree

import (
	"fmt"
)

type HashMismatchError struct {
	H Hash
}

func (h HashMismatchError) Error() string {
	return fmt.Sprintf("Hash mismatch at %x", h.H)
}

type NodeNotFoundError struct {
	H Hash
}

func (n NodeNotFoundError) Error() string {
	return fmt.Sprintf("Node with hash %x not found", n.H)
}
