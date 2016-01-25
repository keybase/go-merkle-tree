package merkleTree

type Hash []byte
type TxInfo []byte
type Prefix []byte

type KeyValuePair struct {
	_struct bool        `codec:",toarray"`
	Key     Hash        `codec:"key"`
	Value   interface{} `codec:"value"`
}

type NodeType int

const (
	NodeTypeNone  NodeType = 0
	NodeTypeINode NodeType = 1
	NodeTypeLeaf  NodeType = 2
)

type Node struct {
	PrevRoot Hash           `codec:"prev_root,omitempty"`
	Tab      []KeyValuePair `codec:"tab"`
	Type     NodeType       `codec:"type"`
}

type Level uint
type ChildIndex uint32

func (l Level) ToChildIndex() ChildIndex { return ChildIndex(l) }
func (p Prefix) Eq(p2 Prefix) bool       { return Hash(p).Eq(Hash(p2)) }
func (p Prefix) ToHash() Hash            { return Hash(p) }
