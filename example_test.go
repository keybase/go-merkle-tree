package merkletree_test

import (
	"context"

	merkletree "github.com/keybase/go-merkle-tree"
)

func ExampleTree_Build() {
	// factory is an "object factory" that makes a whole bunch
	// of phony objects. Importantly, it fits the 'ValueConstructor'
	// interface, so that it can tell the MerkleTree class how
	// to pull type values out of the tree.
	factory := merkletree.NewTestObjFactory()

	// Make a whole bunch of phony objects in our Object Factory.
	objs := factory.Mproduce(1024)

	// Collect and sort the objects into a "sorted map"
	sm := merkletree.NewSortedMapFromList(objs)

	// Make a test storage engine
	eng := merkletree.NewMemEngine()

	// 256 children per node; once there are 512 entries in a leaf,
	// then split the leaf by adding more parents.
	config := merkletree.NewConfig(merkletree.SHA512Hasher{}, 256, 512, factory)

	// Make a new tree object with this engine and these config
	// values
	tree := merkletree.NewTree(eng, config)

	// Make an empty Tranaction info for now
	var txInfo merkletree.TxInfo

	// Build the tree
	err := tree.Build(context.TODO(), sm, txInfo)
	if err != nil {
		panic(err.Error())
	}
}
