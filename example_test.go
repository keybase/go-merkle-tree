package merkleTree_test

import (
	merkleTree "github.com/keybase/go-merkle-tree"
)

func ExampleTree_Build() {

	// Make a whole bunch of phony objects in our Object Factory.
	var objs []merkleTree.KeyValuePair
	objs = merkleTree.NewTestObjFactory().Mproduce(1024)

	// Collect and sort the objects into a "sorted map"
	var sm *merkleTree.SortedMap
	sm = merkleTree.NewSortedMapFromList(objs)

	// Make a test storage engine
	var eng merkleTree.StorageEngine
	eng = merkleTree.NewMemEngine()

	// 256 children per node; once there are 512 entries in a leaf,
	// then split the leaf by adding more parents.
	var config merkleTree.Config
	config = merkleTree.NewConfig(merkleTree.SHA512Hasher{}, 256, 512)

	// Make a new tree object with this engine and these config
	// values
	var tree *merkleTree.Tree
	tree = merkleTree.NewTree(eng, config)

	// Make an empty Tranaction info for now
	var txInfo merkleTree.TxInfo

	// Build the tree
	tree.Build(sm, txInfo)
}
