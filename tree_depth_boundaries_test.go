package merkletree

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestDepthTraversal is a diagnostic test to understand tree traversal depth
func TestDepthTraversal(t *testing.T) {
	hasher := SHA512Hasher{}

	depths := []Level{10, 63, 64, 65, 100}

	for _, depth := range depths {
		testName := fmt.Sprintf("depth_%d", depth)
		t.Run(testName, func(t *testing.T) {
			mock := newMockDeepStorageEngine(hasher, depth)
			cfg := NewConfig(hasher, 256, 512, NewTestObjFactory())
			tree := NewTree(mock, cfg)

			// Use an all-zero key which should map to index 0 at every level
			searchKey := make(Hash, 64) // SHA-512 is 64 bytes

			_, _, err := tree.Find(context.Background(), searchKey)

			effectiveLimit := cfg.maxKeyLevels(searchKey)
			if depth > effectiveLimit {
				require.Error(t, err)
				require.Contains(t, err.Error(), "max tree depth")
			} else {
				require.False(t, err != nil && strings.Contains(err.Error(), "max tree depth"),
					"depth %d below key depth %d returned %v", depth, effectiveLimit, err)
			}
		})
	}
}
