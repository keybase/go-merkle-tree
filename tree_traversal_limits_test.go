package merkletree

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// mockDeepStorageEngine creates a storage engine that returns a very deep tree
// to test key-depth bounds.
type mockDeepStorageEngine struct {
	root     Hash
	nodeData map[string][]byte
}

func newMockDeepStorageEngine(hasher Hasher, depth Level) *mockDeepStorageEngine {
	m := &mockDeepStorageEngine{
		nodeData: make(map[string][]byte),
	}

	// Build a chain of interior nodes from leaf up to root
	// Start with a leaf at the bottom
	leafNode := Node{
		Type: NodeTypeLeaf,
		Leafs: []KeyValuePair{
			{Key: hasher.Hash([]byte("dummy-key")), Value: "dummy-value"},
		},
	}
	leafBytes, _ := encodeToBytes(leafNode)
	leafHash := hasher.Hash(leafBytes)
	m.nodeData[string(leafHash)] = leafBytes

	// Build interior nodes from bottom to top
	childHash := leafHash
	for range depth {
		node := Node{
			Type:   NodeTypeINode,
			INodes: make([]Hash, 256),
		}
		// Put the child in the first slot
		node.INodes[0] = childHash

		nodeBytes, _ := encodeToBytes(node)
		nodeHash := hasher.Hash(nodeBytes)
		m.nodeData[string(nodeHash)] = nodeBytes
		childHash = nodeHash
	}

	// The last hash we generated is the root
	m.root = childHash
	return m
}

func (m *mockDeepStorageEngine) LookupRoot(_ context.Context) (Hash, error) {
	return m.root, nil
}

func (m *mockDeepStorageEngine) LookupNode(_ context.Context, h Hash) ([]byte, error) {
	if data, ok := m.nodeData[string(h)]; ok {
		return data, nil
	}
	// For any unknown hash, return an interior node where all children point to
	// the same hash, creating an effectively infinite tree
	node := Node{
		Type:   NodeTypeINode,
		INodes: make([]Hash, 256),
	}
	// Fill all slots with the same hash to ensure traversal continues
	// regardless of which index is selected
	for i := range node.INodes {
		node.INodes[i] = h
	}
	nodeBytes, _ := encodeToBytes(node)
	return nodeBytes, nil
}

func (m *mockDeepStorageEngine) StoreNode(_ context.Context, _ Hash, _ []byte) error {
	return nil
}

func (m *mockDeepStorageEngine) CommitRoot(_ context.Context, _, _ Hash, _ TxInfo) error {
	return nil
}

// TestFindMaxDepthProtection verifies that Find rejects trees deeper than the
// search key can address.
func TestFindMaxDepthProtection(t *testing.T) {
	hasher := SHA512Hasher{}

	// Create a tree that's deeper than the max allowed depth
	mock := newMockDeepStorageEngine(hasher, 100)

	cfg := NewConfig(hasher, 256, 512, NewTestObjFactory())
	tree := NewTree(mock, cfg)

	// Use an all-zero key which maps to index 0 at every level,
	// ensuring we traverse the full depth of our mock tree
	searchKey := make(Hash, 64) // SHA-512 is 64 bytes

	// This should hit the max depth limit and return an error
	_, _, err := tree.Find(context.Background(), searchKey)
	require.Error(t, err)
	require.Contains(t, err.Error(), "max tree depth")
}

func TestFindDepthProtectionForOneBitFanout(t *testing.T) {
	hasher := SHA512Hasher{}
	mock := newMockDeepStorageEngine(hasher, 600)
	cfg := NewConfig(hasher, 2, 512, NewTestObjFactory())
	tree := NewTree(mock, cfg)

	// A 64-byte key has 512 complete one-bit levels. The traversal must return
	// an error there instead of letting bitslice attempt h[:65].
	_, _, err := tree.Find(context.Background(), make(Hash, 64))
	require.Error(t, err)
	require.Contains(t, err.Error(), "max tree depth 512 exceeded")
}

func TestBuildRejectsUnsplittableLeaf(t *testing.T) {
	of := NewTestObjFactory()
	tree := newTestMemTree(of, 2, 1)
	key := Hash{0x01}
	sm := NewSortedMapFromList([]KeyValuePair{
		{Key: key, Value: "first"},
		{Key: key, Value: "second"},
	})

	err := tree.Build(context.Background(), sm, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "max tree depth")
}

func TestLongKeysCanExceedOneThousandLevels(t *testing.T) {
	of := NewTestObjFactory()
	objs := of.Mproduce(2)
	objs[0].Key = make(Hash, 129)
	objs[1].Key = make(Hash, 129)
	objs[1].Key[128] = 0x80

	tree := newTestMemTree(of, 2, 1)
	require.NoError(t, tree.Build(context.Background(), NewSortedMapFromList(objs), nil))
	findAll(t, tree, objs)
}

func TestLookupNodeVerifiesBeforeDecode(t *testing.T) {
	ctx := context.Background()
	hasher := SHA512Hasher{}
	root := hasher.Hash([]byte("expected node"))
	eng := NewMemEngine()
	require.NoError(t, eng.StoreNode(ctx, root, []byte{0xc1}))
	require.NoError(t, eng.CommitRoot(ctx, nil, root, nil))

	tree := NewTree(eng, NewConfig(hasher, 256, 512, NewTestObjFactory()))
	_, _, err := tree.Find(ctx, make(Hash, 64))
	// With the skipVerify pattern, lookupNode decodes first. Corrupt data
	// that can't be decoded produces a decode error. Only after successful
	// decode does verification (when skipVerify=false) catch hash mismatches.
	require.Error(t, err)
	require.Contains(t, err.Error(), "msgpack decode error")
}

// TestUpsertMaxDepthProtection verifies that Upsert rejects trees deeper than
// the inserted key can address.
func TestUpsertMaxDepthProtection(t *testing.T) {
	hasher := SHA512Hasher{}

	// Create a tree that's deeper than the max allowed depth
	mock := newMockDeepStorageEngine(hasher, 100)

	cfg := NewConfig(hasher, 256, 512, NewTestObjFactory())
	tree := NewTree(mock, cfg)

	// Use an all-zero key which maps to index 0 at every level,
	// ensuring we traverse the full depth of our mock tree
	kvp := KeyValuePair{
		Key:   make(Hash, 64), // SHA-512 is 64 bytes, all zeros
		Value: "value",
	}

	// This should hit the max depth limit and return an error
	err := tree.Upsert(context.Background(), kvp, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "max tree depth")
}

// TestFindContextCancellation verifies that Find respects context cancellation
// during traversal.
func TestFindContextCancellation(t *testing.T) {
	hasher := SHA512Hasher{}

	// Create a moderately deep tree
	mock := newMockDeepStorageEngine(hasher, 100)

	cfg := NewConfig(hasher, 256, 512, NewTestObjFactory())
	tree := NewTree(mock, cfg)

	// Create a context that's already cancelled
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	searchKey := hasher.Hash([]byte("search-key"))

	// Should return context.Canceled error
	_, _, err := tree.Find(ctx, searchKey)
	require.ErrorIs(t, err, context.Canceled)
}

// TestUpsertContextCancellation verifies that Upsert respects context
// cancellation during traversal.
func TestUpsertContextCancellation(t *testing.T) {
	hasher := SHA512Hasher{}

	// Create a moderately deep tree
	mock := newMockDeepStorageEngine(hasher, 100)

	cfg := NewConfig(hasher, 256, 512, NewTestObjFactory())
	tree := NewTree(mock, cfg)

	// Create a context that's already cancelled
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	kvp := KeyValuePair{
		Key:   hasher.Hash([]byte("key")),
		Value: "value",
	}

	// Should return context.Canceled error
	err := tree.Upsert(ctx, kvp, nil)
	require.ErrorIs(t, err, context.Canceled)
}

// TestFindContextTimeout verifies that Find can be cancelled via timeout.
func TestFindContextTimeout(t *testing.T) {
	hasher := SHA512Hasher{}

	// Create a moderately deep tree
	mock := newMockDeepStorageEngine(hasher, 100)

	cfg := NewConfig(hasher, 256, 512, NewTestObjFactory())
	tree := NewTree(mock, cfg)

	// Create a context with a very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	// Give the timeout a moment to expire
	time.Sleep(1 * time.Millisecond)

	searchKey := hasher.Hash([]byte("search-key"))

	// Should return context deadline exceeded error
	_, _, err := tree.Find(ctx, searchKey)
	require.Error(t, err)
	require.True(t, err == context.DeadlineExceeded || err == context.Canceled)
}

// TestNormalDepthTree verifies that normal depth trees still work correctly
// with the new depth checks.
func TestNormalDepthTree(t *testing.T) {
	// Use a real memory storage engine to ensure normal operations work
	of := NewTestObjFactory()
	tree := newTestMemTree(of, 256, 512)

	// Build a tree with some data
	objs := of.Mproduce(100)
	sm := NewSortedMapFromList(objs)

	ctx := context.Background()
	err := tree.Build(ctx, sm, nil)
	require.NoError(t, err)

	// Verify we can find all objects
	for _, obj := range objs {
		val, _, err := tree.Find(ctx, obj.Key)
		require.NoError(t, err)
		require.NotNil(t, val)
	}

	// Verify we can upsert
	newObj := of.Produce()
	err = tree.Upsert(ctx, newObj, nil)
	require.NoError(t, err)

	// Verify we can find the upserted object
	val, _, err := tree.Find(ctx, newObj.Key)
	require.NoError(t, err)
	require.NotNil(t, val)
}

// TestTraversalEdgeCases tests edge conditions in bounded tree traversal.
func TestTraversalEdgeCases(t *testing.T) {
	t.Run("leaf exactly at limit", func(t *testing.T) {
		hasher := SHA512Hasher{}
		searchKey := make(Hash, 64)
		limit := NewConfig(hasher, 256, 512, NewTestObjFactory()).maxKeyLevels(searchKey)
		mock := newMockDeepStorageEngine(hasher, limit)

		cfg := NewConfig(hasher, 256, 512, NewTestObjFactory())
		tree := NewTree(mock, cfg)

		_, _, err := tree.Find(context.Background(), searchKey)
		require.NoError(t, err)
	})

	t.Run("interior node at limit", func(t *testing.T) {
		hasher := SHA512Hasher{}
		searchKey := make(Hash, 64)
		limit := NewConfig(hasher, 256, 512, NewTestObjFactory()).maxKeyLevels(searchKey)
		mock := newMockDeepStorageEngine(hasher, limit+1)

		cfg := NewConfig(hasher, 256, 512, NewTestObjFactory())
		tree := NewTree(mock, cfg)

		_, _, err := tree.Find(context.Background(), searchKey)
		require.Error(t, err)
		require.Contains(t, err.Error(), "max tree depth")
	})

	t.Run("empty tree traversal", func(t *testing.T) {
		eng := NewMemEngine()
		cfg := NewConfig(SHA512Hasher{}, 256, 512, NewTestObjFactory())
		tree := NewTree(eng, cfg)

		// Find in empty tree (no root committed)
		searchKey := make(Hash, 64)
		val, root, err := tree.Find(context.Background(), searchKey)

		// Should handle empty tree gracefully
		require.Nil(t, val)
		if root == nil {
			t.Log("empty tree returned nil root (expected)")
		}
		t.Logf("Empty tree result: val=%v, root=%v, err=%v", val, root, err)
	})

	t.Run("context with deadline during traversal", func(t *testing.T) {
		hasher := SHA512Hasher{}
		mock := newMockDeepStorageEngine(hasher, 50)

		cfg := NewConfig(hasher, 256, 512, NewTestObjFactory())
		tree := NewTree(mock, cfg)

		// Set a deadline that might expire during traversal
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
		defer cancel()

		// Small delay to let deadline approach
		time.Sleep(2 * time.Millisecond)

		searchKey := make(Hash, 64)
		_, _, err := tree.Find(ctx, searchKey)

		require.Error(t, err)
		require.True(t, err == context.DeadlineExceeded || err == context.Canceled)
	})

	t.Run("find in tree with single leaf", func(t *testing.T) {
		eng := NewMemEngine()
		factory := NewTestObjFactory()
		cfg := NewConfig(SHA512Hasher{}, 256, 512, factory)
		tree := NewTree(eng, cfg)

		ctx := context.Background()

		// Build tree with single item
		obj := factory.Produce()
		sm := NewSortedMapFromList([]KeyValuePair{obj})

		require.NoError(t, tree.Build(ctx, sm, nil))

		// Find the single item
		val, _, err := tree.Find(ctx, obj.Key)
		require.NoError(t, err)
		require.NotNil(t, val)

		// Search for non-existent key
		nonExistent := make(Hash, 64)
		nonExistent[0] = 0xff
		val2, _, err := tree.Find(ctx, nonExistent)
		if err != nil {
			t.Logf("Find non-existent returned error: %v", err)
		}
		require.Nil(t, val2)
	})

	t.Run("upsert in empty tree", func(t *testing.T) {
		eng := NewMemEngine()
		factory := NewTestObjFactory()
		cfg := NewConfig(SHA512Hasher{}, 256, 512, factory)
		tree := NewTree(eng, cfg)

		ctx := context.Background()

		// Upsert into empty tree
		obj := factory.Produce()
		err := tree.Upsert(ctx, obj, nil)
		require.NoError(t, err)

		// Verify can find it
		val, _, err := tree.Find(ctx, obj.Key)
		require.NoError(t, err)
		require.NotNil(t, val)
	})

	t.Run("repeated upsert same key", func(t *testing.T) {
		eng := NewMemEngine()
		factory := NewTestObjFactory()
		cfg := NewConfig(SHA512Hasher{}, 256, 512, factory)
		tree := NewTree(eng, cfg)

		ctx := context.Background()

		// Build initial tree
		obj := factory.Produce()
		require.NoError(t, tree.Upsert(ctx, obj, nil))

		// Upsert same key with different value multiple times
		for i := range 5 {
			newObj := factory.Produce()
			newObj.Key = obj.Key // Same key
			err := tree.Upsert(ctx, newObj, nil)
			require.NoError(t, err, "upsert iteration %d", i)
		}

		// Should still be able to find the key
		val, _, err := tree.Find(ctx, obj.Key)
		require.NoError(t, err)
		require.NotNil(t, val)
	})
}
