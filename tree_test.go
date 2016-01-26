package merkleTree

import (
	"crypto/sha512"
	"testing"
)

type sha512Hasher struct{}

func (s sha512Hasher) Hash(b []byte) Hash {
	tmp := sha512.Sum512(b)
	return Hash(tmp[:])
}

func makeConfig() Config {
	return NewConfig(sha512Hasher{}, 4, 16)
}

func newMemTree() (t *Tree, m *MemEngine) {
	m = NewMemEngine()
	t = NewTree(m, makeConfig())
	return t, m
}

func TestSimpleBuild(t *testing.T) {
	of := newObjFactory()
	objs := of.mproduce(1024)
	sm := NewSortedMapFromList(objs)
	tree, _ := newMemTree()
	if err := tree.Build(sm, nil); err != nil {
		t.Fatalf("Error in build: %v", err)
	}
	findAll(t, tree, objs)
}

func findAll(t *testing.T, tree *Tree, objs []KeyValuePair) {
	for i, kvp := range objs {
		v, _, err := tree.Find(kvp.Key)
		if err != nil {
			t.Fatalf("Find for obj %d yielded an error: %v", i, err)
		}
		if v == nil {
			t.Fatalf("Find for obj %i with key %v returned no results", i, kvp.Key)
		}
		if !deepEqual(v, kvp.Value) {
			t.Fatalf("Didn't get object equality for %d: %+v != %+v", v, kvp.Value)
		}
	}
}
