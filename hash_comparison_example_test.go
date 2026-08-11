package merkletree

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestHashComparisonBehavior documents the historical ordering while checking
// that inputs which previously panicked now return safely.
func TestHashComparisonBehavior(t *testing.T) {
	t.Run("fixed-length hashes - behavior unchanged", func(t *testing.T) {
		// Most common case: 64-byte SHA-512 hashes
		h1 := make(Hash, 64)
		h2 := make(Hash, 64)
		h1[0] = 0x01
		h2[0] = 0x02

		// Both old and new: h1 < h2
		require.True(t, h1.Less(h2))
		t.Log("✓ Fixed-length comparison works as before")
	})

	t.Run("leading zeros preserve historical behavior", func(t *testing.T) {
		h1 := Hash{0x00, 0x00, 0x01, 0x02}
		h2 := Hash{0x01, 0x02}

		require.True(t, h1.Less(h2))

		t.Logf("h1.Len() = %d (after stripping zeros)", h1.Len())
		t.Logf("h2.Len() = %d (after stripping zeros)", h2.Len())
	})

	t.Run("all zeros comparison - now safe", func(t *testing.T) {
		h1 := Hash{0x00, 0x00}
		h2 := Hash{0x00}

		// The old implementation panicked in this direction. Its reverse
		// comparison returned equal, so equality is the compatible extension.
		result := h1.cmp(h2)
		require.Zero(t, result)
	})

	t.Run("empty hash comparison", func(t *testing.T) {
		h1 := Hash{}
		h2 := Hash{0x00}

		require.True(t, h1.Eq(h2))

		h3 := Hash{}
		require.True(t, h1.Eq(h3))
	})
}

// ExampleHash_comparison demonstrates the historical comparison behavior.
func ExampleHash_comparison() {
	// Fixed-length hashes (most common case)
	h1 := Hash{0x01, 0x02, 0x03}
	h2 := Hash{0x01, 0x02, 0x04}

	fmt.Printf("h1 < h2: %v\n", h1.Less(h2))
	fmt.Printf("h1 == h2: %v\n", h1.Eq(h2))

	// Variable-length with leading zeros
	h3 := Hash{0x00, 0x01}
	h4 := Hash{0x01}

	fmt.Printf("h3 < h4: %v (leading zeros matter)\n", h3.Less(h4))

	// Output:
	// h1 < h2: true
	// h1 == h2: false
	// h3 < h4: true (leading zeros matter)
}

// TestHashSortingStability verifies that hash sorting is stable and consistent
func TestHashSortingStability(t *testing.T) {
	hashes := []Hash{
		{0x03},
		{0x01},
		{0x00, 0x02},
		{0x02},
		{0x00, 0x01},
	}

	// Create sorted map to see ordering
	kvps := make([]KeyValuePair, len(hashes))
	for i, h := range hashes {
		kvps[i] = KeyValuePair{Key: h, Value: fmt.Sprintf("val-%d", i)}
	}

	sm := NewSortedMapFromList(kvps)

	// Verify ordering is consistent
	for i := ChildIndex(0); i < sm.Len()-1; i++ {
		curr := sm.at(i)
		next := sm.at(i + 1)
		require.True(t, curr.Key.Less(next.Key) || curr.Key.Eq(next.Key), "sorting violation at index %d", i)
	}

	expectedOrder := []Hash{
		{0x00, 0x01},
		{0x00, 0x02},
		{0x01},
		{0x02},
		{0x03},
	}

	for i := ChildIndex(0); i < sm.Len(); i++ {
		require.Equal(t, expectedOrder[i], sm.at(i).Key)
	}
}

// TestVariableLengthKeyFinding tests that variable-length keys can be found
// after being inserted (most important practical test)
func TestVariableLengthKeyFinding(t *testing.T) {
	eng := NewMemEngine()
	factory := NewTestObjFactory()
	cfg := NewConfig(SHA512Hasher{}, 256, 512, factory)
	tree := NewTree(eng, cfg)

	ctx := context.Background()

	// Create test objects with various key lengths and leading zeros
	testKeys := []Hash{
		{0x01},             // 1 byte
		{0x00, 0x01},       // leading zero
		{0x01, 0x02},       // 2 bytes
		{0x00, 0x00, 0x01}, // two leading zeros
	}

	kvps := make([]KeyValuePair, len(testKeys))
	for i, key := range testKeys {
		// Use the test object factory to create proper values
		obj := factory.Produce()
		kvps[i] = KeyValuePair{Key: key, Value: obj.Value}
	}

	// Insert all
	for _, kvp := range kvps {
		require.NoError(t, tree.Upsert(ctx, kvp, nil))
	}

	// Verify we can find all of them
	for _, kvp := range kvps {
		val, _, err := tree.Find(ctx, kvp.Key)
		require.NoError(t, err)
		require.NotNil(t, val, "key %#v", []byte(kvp.Key))
	}
}
