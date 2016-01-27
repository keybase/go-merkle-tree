package merkleTree

import (
	"crypto/sha512"
)

type sha512Hasher struct{}

func (s sha512Hasher) Hash(b []byte) Hash {
	tmp := sha512.Sum512(b)
	return Hash(tmp[:])
}

func ExampleBuild() {

	// Make a whole bunch of phony objects in our Object Factory.
	var objs []KeyValuePair
	objs = newObjFactory().mproduce(1024)

	// Collect and sort the objects into a "sorted map"
	var sm *SortedMap
	sm = NewSortedMapFromList(objs)

	// Make a test storage engine
	var eng StorageEngine
	eng = NewMemEngine()

	// 256 children per node; once there are 512 entries in a leaf,
	// then split the leaf by adding more parents.
	var config Config
	config = NewConfig(sha512Hasher{}, 256, 512)

	// Make a new tree object with this engine and these config
	// values
	var tree *Tree
	tree = NewTree(eng, config)

	// Make an empty Tranaction info for now
	var txInfo TxInfo

	// Build the tree
	tree.Build(sm, txInfo)
}
