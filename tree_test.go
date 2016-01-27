package merkleTree

import (
	"testing"
)

func makeTestConfig() Config {
	return NewConfig(sha512Hasher{}, 4, 16)
}

func newTestMemTree() (t *Tree, m *MemEngine) {
	m = NewMemEngine()
	t = NewTree(m, makeTestConfig())
	return t, m
}

func testSimpleBuild(t *testing.T) {
	of := newObjFactory()
	objs := of.mproduce(1024)
	sm := NewSortedMapFromList(objs)
	tree, _ := newTestMemTree()
	if err := tree.Build(sm, nil); err != nil {
		t.Fatalf("Error in build: %v", err)
	}
	findAll(t, tree, objs)
}

// This example shows how to build the tree.
func TestSimpleBuild(t *testing.T) {
	for i := 0; i < 32; i++ {
		testSimpleBuild(t)
	}
}

func findAll(t *testing.T, tree *Tree, objs []KeyValuePair) {
	for i, kvp := range objs {
		v, _, err := tree.Find(kvp.Key)
		if err != nil {
			t.Fatalf("Find for obj %d yielded an error: %v", i, err)
		}
		if v == nil {
			t.Fatalf("Find for obj %d with key %v returned no results", i, kvp.Key)
		}
		if !deepEqual(v, kvp.Value) {
			t.Fatalf("Didn't get object equality for %d: %+v != %+v", i, v, kvp.Value)
		}
	}
}
