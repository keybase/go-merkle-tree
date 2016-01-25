package merkleTree

type Hasher func([]byte) Hash

type MerkleStorer interface {
	GetHasher() Hasher
	StoreNode(Hash, Node, []byte) error
	CommitRoot(curr Hash, prev Hash, txinfo TxInfo) error
}
