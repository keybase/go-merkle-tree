package merkleTree

import (
	"sync"
)

type Tree struct {
	sync.RWMutex
	eng Engine
	cfg Config
}

func (t *Tree) Build(sm *SortedMap, txi TxInfo) (err error) {
	t.Lock()
	defer t.Unlock()
	var prevRoot, nextRoot Hash
	if prevRoot, err = t.eng.LookupRoot(); err != nil {
		return err
	}
	if nextRoot, err = t.hashTreeRecursive(Level(0), sm, prevRoot); err != nil {
		return err
	}
	if err = t.eng.CommitRoot(prevRoot, nextRoot, txi); err != nil {
		return err
	}

	return err
}

func (t *Tree) hashTreeRecursive(level Level, sm *SortedMap, prevRoot Hash) (ret Hash, err error) {
	if sm.Len() <= t.cfg.N {
		return t.simpleInsert(level, sm, prevRoot)
	}

	M := t.cfg.M // the number of children we have
	var j ChildIndex
	nsm := NewSortedMap() // new sorted map

	for i := ChildIndex(0); i < M; i++ {
		prefix := t.cfg.formatPrefix(i)
		start := j
		for j < nsm.Len() && t.cfg.prefixAtLevel(level, sm.at(j).Key).Eq(prefix) {
			j++
		}
		end := j
		if end > start {
			sublist := sm.slice(start, end)
			ret, err = t.hashTreeRecursive(level+1, sublist, nil)
			if err != nil {
				return nil, err
			}
			prefix := t.cfg.prefixThroughLevel(level, sublist.at(0).Key)
			nsm.push(KeyValuePair{Key: prefix.ToHash(), Value: ret})
		}
	}
	var node Node
	var nodeExported []byte
	if ret, node, nodeExported, err = nsm.exportToNode(t.eng, NodeTypeINode, prevRoot, level); err != nil {
		return nil, err
	}
	err = t.eng.StoreNode(ret, node, nodeExported)
	return ret, err

}

func (t *Tree) simpleInsert(l Level, sm *SortedMap, prevRoot Hash) (ret Hash, err error) {
	var node Node
	var nodeExported []byte
	if ret, node, nodeExported, err = sm.exportToNode(t.eng, NodeTypeLeaf, prevRoot, l); err != nil {
		return nil, err
	}
	if err = t.eng.StoreNode(ret, node, nodeExported); err != nil {
		return nil, err
	}
	return ret, err
}

func (t *Tree) verifyNode(h Hash, node *Node) (err error) {
	var b []byte
	if b, err = encodeToBytes(node); err != nil {
		return err
	}
	h2 := t.eng.Hash(b)
	if !h.Eq(h2) {
		err = HashMismatchError{H: h}
	}
	return err
}

func (t *Tree) find(h Hash, skipVerify bool) (kvp *KeyValuePair, root Hash, err error) {
	t.RLock()
	defer t.RUnlock()

	root, err = t.eng.LookupRoot()
	if err != nil {
		return nil, nil, err
	}
	curr := root
	var level Level
	for curr != nil {
		var node *Node
		node, err = t.eng.LookupNode(curr)
		if node == nil || err == nil {
			err = NodeNotFoundError{H: curr}
		}
		if err != nil {
			return nil, nil, err
		}
		if !skipVerify {
			if err = t.verifyNode(curr, node); err != nil {
				return nil, nil, err
			}
		}
		sm := NewSortedMapFromList(node.Tab)
		kvp = sm.find(h)
		if node.Type == NodeTypeLeaf {
			curr = nil
		} else if kvp == nil {
			curr = nil
		} else {
			curr = t.cfg.prefixThroughLevel(level, kvp.Key).ToHash()
		}
		level++
	}
	return kvp, root, err
}

func (t *Tree) Find(h Hash) (kvp *KeyValuePair, root Hash, err error) {
	return t.find(h, false)
}

func (t *Tree) FindUnsafe(h Hash) (kvp *KeyValuePair, root Hash, err error) {
	return t.find(h, true)
}
