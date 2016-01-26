package merkleTree

type Engine interface {
	StoreNode(Hash, Node, []byte) error
	CommitRoot(prev Hash, curr Hash, txinfo TxInfo) error
	LookupNode(Hash) (*Node, error)
	LookupRoot() (Hash, error)
}
