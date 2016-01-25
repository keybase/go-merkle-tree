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
