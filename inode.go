package merkleTree

import ()

type childPointerMap struct {
	tab []Hash
}

func newChildPointerMap(capacity ChildIndex) *childPointerMap {
	return &childPointerMap{
		tab: make([]Hash, capacity),
	}
}

func newChildPointerMapFromNode(n *Node) *childPointerMap {
	return &childPointerMap{
		tab: n.INodes,
	}
}

func (c *childPointerMap) exportToNode(h Hasher, prevRoot Hash, level Level) (hash Hash, node Node, objExported []byte, err error) {
	node.Type = nodeTypeINode
	node.INodes = c.tab
	return node.export(h, prevRoot, level)
}

func (n Node) export(h Hasher, prevRoot Hash, level Level) (hash Hash, node Node, objExported []byte, err error) {
	if prevRoot != nil && level == Level(0) {
		n.PrevRoot = prevRoot
	}
	objExported, err = encodeToBytes(n)
	if err == nil {
		hash = h.Hash(objExported)
	}
	return hash, n, objExported, err
}

func (n *Node) findChildByIndex(i ChildIndex) (Hash, error) {
	if n.INodes == nil || int(i) >= len(n.INodes) {
		return nil, ErrBadINode
	}
	return n.INodes[i], nil
}

func (c *childPointerMap) set(i ChildIndex, h Hash) *childPointerMap {
	c.tab[i] = h
	return c
}
