package merkleTree

import (
	"testing"
)

func makeTestConfig(of *TestObjFactory, m ChildIndex, n ChildIndex) Config {
	return NewConfig(SHA512Hasher{}, m, n, of)
}

func newTestMemTree(of *TestObjFactory, m ChildIndex, n ChildIndex) (t *Tree, me *MemEngine) {
	me = NewMemEngine()
	t = NewTree(me, makeTestConfig(of, m, n))
	return t, me
}

func testSimpleBuild(t *testing.T, numElem int, m ChildIndex, n ChildIndex) {
	of := NewTestObjFactory()
	objs := of.Mproduce(numElem)
	sm := NewSortedMapFromList(objs)
	tree, _ := newTestMemTree(of, m, n)
	if err := tree.Build(sm, nil); err != nil {
		t.Fatalf("Error in build: %v", err)
	}
	findAll(t, tree, objs)
}

func TestSimpleBuild4by16(t *testing.T) {
	for i := 0; i < 10; i++ {
		testSimpleBuild(t, 1024, ChildIndex(4), ChildIndex(16))
	}
}

func TestSimpleBuild2by4(t *testing.T) {
	for i := 0; i < 10; i++ {
		testSimpleBuild(t, 1024, ChildIndex(2), ChildIndex(4))
	}
}

func TestSimpleBuild256by256(t *testing.T) {
	for i := 0; i < 2; i++ {
		testSimpleBuild(t, 8192, ChildIndex(256), ChildIndex(256))
	}
}
func TestSimpleBuildSmall(t *testing.T) {
	testSimpleBuild(t, 1, ChildIndex(2), ChildIndex(2))
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
