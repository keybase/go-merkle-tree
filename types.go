package merkleTree

type HexKey []byte

type KeyValuePair struct {
	_struct bool        `codec:",toarray"`
	Key     HexKey      `codec:"key"`
	Value   interface{} `codec:"value"`
}

type NodeType int

const (
	NodeTypeNone  NodeType = 0
	NodeTypeINode NodeType = 1
	NodeTypeLeaf  NodeType = 2
)

type Node struct {
	PrevRoot []byte         `codec:"prev_root,omitempty"`
	Tab      []KeyValuePair `codec:"tab"`
	Type     NodeType       `codec:"type"`
}
