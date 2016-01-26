package merkleTree

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
)

// MemEngine is an in-memory MerkleTree engine, used now mainly for testing
type MemEngine struct {
	root  Hash
	nodes map[string](*Node)
}

func NewMemEngine() *MemEngine {
	return &MemEngine{
		nodes: make(map[string](*Node)),
	}
}

var _ Engine = (*MemEngine)(nil)

// CommitRoot "commits" the root ot the blessed memory slot
func (m *MemEngine) CommitRoot(prev Hash, curr Hash, txinfo TxInfo) error {
	m.root = curr
	return nil
}

// Hash runs SHA512
func (m *MemEngine) Hash(d []byte) Hash {
	sum := sha512.Sum512(d)
	return Hash(sum[:])
}

// LookupNode looks up a MerkleTree node by hash
func (m *MemEngine) LookupNode(h Hash) (*Node, error) {
	ret := m.nodes[hex.EncodeToString(h)]
	return ret, nil
}

// LookupRoot fetches the root of the in-memory tree back out
func (m *MemEngine) LookupRoot() (Hash, error) {
	return m.root, nil
}

// StoreNode stores the given node under the given key.
func (m *MemEngine) StoreNode(key Hash, val Node, _ []byte) error {
	fmt.Printf("Store %x -> %+v\n", key, val)
	m.nodes[hex.EncodeToString(key)] = &val
	return nil
}