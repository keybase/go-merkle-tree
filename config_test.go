package merkletree

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestNewConfigValidation verifies that NewConfig rejects configurations which
// cannot make progress while retaining historically accepted fanouts.
func TestNewConfigValidation(t *testing.T) {
	hasher := SHA512Hasher{}
	factory := NewTestObjFactory()

	t.Run("valid powers of two", func(_ *testing.T) {
		validM := []ChildIndex{2, 4, 8, 16, 32, 64, 128, 256, 512, 1024}
		for _, m := range validM {
			// Should not panic
			_ = NewConfig(hasher, m, 10, factory)
		}
	})

	t.Run("zero m panics", func(t *testing.T) {
		require.PanicsWithValue(t, "invalid config: m must be at least 2, got 0", func() {
			NewConfig(hasher, 0, 10, factory)
		})
	})

	t.Run("one child panics", func(t *testing.T) {
		require.PanicsWithValue(t, "invalid config: m must be at least 2, got 1", func() {
			NewConfig(hasher, 1, 10, factory)
		})
	})

	t.Run("zero leaf capacity panics", func(t *testing.T) {
		require.PanicsWithValue(t, "invalid config: n must be positive", func() {
			NewConfig(hasher, 256, 0, factory)
		})
	})

	t.Run("historically accepted fanouts remain accepted", func(t *testing.T) {
		historicalM := []ChildIndex{3, 5, 6, 7, 9, 10, 100, 255, 257, 1 << 26, 1 << 31}
		for _, m := range historicalM {
			require.NotPanics(t, func() {
				NewConfig(hasher, m, 10, factory)
			}, "historically accepted m=%d", m)
		}
	})
}

func TestPrefixAtLevelBeyondKey(t *testing.T) {
	cfg := NewConfig(SHA512Hasher{}, 2, 10, NewTestObjFactory())
	key := make(Hash, 64)

	prefix, index := cfg.PrefixAndIndexAtLevel(512, key)
	require.Nil(t, prefix)
	require.Zero(t, index)
	require.Nil(t, cfg.PrefixAtLevel(Level(^uint(0)), key))
}
