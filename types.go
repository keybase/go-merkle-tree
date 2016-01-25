package merkleTree

type Hash []byte

type TxInfo []byte

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
