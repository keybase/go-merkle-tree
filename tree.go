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

func (t *Tree) lookupNode(h Hash) (*Node, error) {
	node, err := t.eng.LookupNode(h)
	if node == nil || err == nil {
		err = NodeNotFoundError{H: h}
	}
	if err != nil {
		return nil, err
	}
	return node, err
}

func (t *Tree) find(h Hash, skipVerify bool) (ret interface{}, root Hash, err error) {
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
		node, err = t.lookupNode(curr)
		if err != nil {
			return nil, nil, err
		}
		if !skipVerify {
			if err = t.verifyNode(curr, node); err != nil {
				return nil, nil, err
			}
		}

		if node.Type == NodeTypeLeaf {
			ret = node.findValueInLeaf(h)
			break
		}
		prfx := t.cfg.prefixThroughLevel(level, h)
		curr, err = node.findChildByPrefix(prfx)
		if err != nil {
			return nil, nil, err
		}
		level++
	}
	return ret, root, err
}

func (t *Tree) Find(h Hash) (ret interface{}, root Hash, err error) {
	return t.find(h, false)
}

func (t *Tree) FindUnsafe(h Hash) (ret interface{}, root Hash, err error) {
	return t.find(h, true)
}

type step struct {
	p Prefix
	n *Node
	l Level
}

type path struct{ steps []step }

func (p *path) push(s step) { p.steps = append(p.steps, s) }
func (p path) len() Level   { return Level(len(p.steps)) }

func (p *path) reverse() {
	j := len(p.steps) - 1
	for i := 0; i < j; i++ {
		p.steps[i], p.steps[j] = p.steps[j], p.steps[i]
		j--
	}
}

func (t *Tree) Upsert(kvp KeyValuePair, txinfo TxInfo) (err error) {
	t.Lock()
	defer t.Unlock()

	root, err := t.eng.LookupRoot()
	if err != nil {
		return err
	}

	prevRoot := root
	var last *Node

	var path path
	var level Level

	curr, err := t.lookupNode(root)
	if err != nil {
		return err
	}

	// Find the path from the key up to the root;
	// find by walking down from the root.
	for curr != nil {
		prefix := t.cfg.prefixThroughLevel(level, kvp.Key)
		path.push(step{p: prefix, n: curr, l: level})
		level++
		last = curr
		if curr.Type == NodeTypeLeaf {
			break
		}
		nxt, err := curr.findChildByPrefix(prefix)
		if err != nil {
			return err
		}
		if nxt == nil {
			break
		}
		curr, err = t.lookupNode(nxt)
		if err != nil {
			return err
		}
	}

	// Figure out what to store at the node where we stopped going down the path.
	var sm *SortedMap
	if last == nil || last.Type == NodeTypeINode {
		sm = NewSortedMapFromKeyAndValue(kvp)
		level = 0
	} else if val2 := last.findValueInLeaf(kvp.Key); val2 == nil || !deepEqual(val2, kvp.Value) {
		sm = newSortedMapFromNode(last).replace(kvp)
		level = path.len() - 1
	} else {
		return nil
	}

	// Make a new subtree out of our new node.
	_, err = t.hashTreeRecursive(level, sm, prevRoot)
	if err != nil {
		return err
	}

	path.reverse()

	return nil
}
