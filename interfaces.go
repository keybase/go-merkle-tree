package merkleTree

type Hasher interface {
	Hash([]byte) Hash
}

type Engine interface {
	Hasher
	StoreNode(Hash, Node, []byte) error
	CommitRoot(prev Hash, curr Hash, txinfo TxInfo) error
	LookupNode(Hash) (*Node, error)
	LookupRoot() (Hash, error)
}
