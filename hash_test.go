package merkletree

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHashLen(t *testing.T) {
	tests := []struct {
		desc     string
		hash     Hash
		expected int
	}{
		{"basic", []byte{1, 2, 3}, 3},
		{"leading", []byte{0, 0, 0, 1, 2, 3}, 3},
		{"trailing", []byte{1, 2, 3, 0, 0, 0}, 6},
		{"leading+trailing", []byte{0, 0, 1, 2, 3, 0, 0}, 5},
		{"nil", []byte{}, 0},
		{"middle", []byte{0, 1, 2, 3, 0, 3, 2, 1, 0}, 8},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(_ *testing.T) {
			actual := tt.hash.Len()
			require.Equal(t, tt.expected, actual)
		})
	}
}

// TestHashCmpNoPanic verifies that Hash.cmp does not panic when comparing
// hashes of different lengths, including the case where hashes have the
// same stripped length but different actual lengths.
func TestHashCmpNoPanic(t *testing.T) {
	tests := []struct {
		desc string
		h1   Hash
		h2   Hash
	}{
		{
			desc: "same stripped length, different actual (all zeros)",
			h1:   Hash{0x00, 0x00},
			h2:   Hash{0x00},
		},
		{
			desc: "same stripped length, different actual (prefix match)",
			h1:   Hash{0x01, 0x02, 0x03},
			h2:   Hash{0x01, 0x02},
		},
		{
			desc: "different stripped and actual lengths",
			h1:   Hash{0x00, 0x01},
			h2:   Hash{0x01},
		},
		{
			desc: "one empty",
			h1:   Hash{0x01},
			h2:   Hash{},
		},
		{
			desc: "both empty",
			h1:   Hash{},
			h2:   Hash{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(_ *testing.T) {
			// The test passes if these don't panic
			_ = tt.h1.cmp(tt.h2)
			_ = tt.h2.cmp(tt.h1)
			_ = tt.h1.Eq(tt.h2)
			_ = tt.h1.Less(tt.h2)
		})
	}
}

func TestHashCmpPreservesHistoricalOrdering(t *testing.T) {
	tests := []struct {
		name string
		h1   Hash
		h2   Hash
		want int
	}{
		{
			name: "significant length wins over first byte",
			h1:   Hash{0xff},
			h2:   Hash{0x01, 0x00},
			want: -1,
		},
		{
			name: "reverse significant length comparison",
			h1:   Hash{0x01, 0x00},
			h2:   Hash{0xff},
			want: 1,
		},
		{
			name: "raw bytes break equal significant lengths",
			h1:   Hash{0x00, 0x01},
			h2:   Hash{0x01},
			want: -1,
		},
		{
			name: "former panic extends reverse equality",
			h1:   Hash{0x00, 0x00},
			h2:   Hash{0x00},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.h1.cmp(tt.h2))
		})
	}

	// This is the persisted order produced by the old comparator. Binary
	// search must continue finding both entries after an upgrade.
	sm := NewSortedMapFromSortedList([]KeyValuePair{
		{Key: Hash{0xff}, Value: "first"},
		{Key: Hash{0x01, 0x00}, Value: "second"},
	})
	for _, key := range []Hash{{0xff}, {0x01, 0x00}} {
		require.NotNil(t, sm.find(key))
	}
}

func TestHashCmpMatchesEveryDefinedHistoricalResult(t *testing.T) {
	values := []byte{0x00, 0x01, 0xff}
	hashes := []Hash{nil}
	for _, a := range values {
		hashes = append(hashes, Hash{a})
		for _, b := range values {
			hashes = append(hashes, Hash{a, b})
			for _, c := range values {
				hashes = append(hashes, Hash{a, b, c})
			}
		}
	}

	historicalCmp := func(h1, h2 Hash) (result int, defined bool) {
		defined = true
		defer func() {
			if recover() != nil {
				defined = false
			}
		}()
		if h1.Len() < h2.Len() {
			return -1, true
		}
		if h1.Len() > h2.Len() {
			return 1, true
		}
		for i, b := range h1 {
			if b < h2[i] {
				return -1, true
			}
			if b > h2[i] {
				return 1, true
			}
		}
		return 0, true
	}

	for _, h1 := range hashes {
		for _, h2 := range hashes {
			want, defined := historicalCmp(h1, h2)
			if !defined {
				continue
			}
			require.Equal(t, want, h1.cmp(h2))
		}
	}
}

// TestHashComparisonBoundaries tests hash comparison at extreme values
// and boundary conditions to ensure robustness.
func TestHashComparisonBoundaries(t *testing.T) {
	t.Run("maximum length difference", func(t *testing.T) {
		h1 := make(Hash, 1)
		h1[0] = 0x01
		h2 := make(Hash, 1000)
		h2[0] = 0x01

		// Should handle large length differences gracefully
		result := h1.cmp(h2)
		require.Negative(t, result)
	})

	t.Run("all zeros of different lengths", func(t *testing.T) {
		// Test all combinations of zero-filled hashes up to length 10
		for i := 1; i <= 10; i++ {
			for j := 1; j <= 10; j++ {
				if i == j {
					continue
				}
				h1 := make(Hash, i)
				h2 := make(Hash, j)

				// Should never panic regardless of zero-filled lengths
				_ = h1.cmp(h2)
				_ = h1.Eq(h2)
				_ = h1.Less(h2)

				// Preserve the only result the old comparator returned for
				// these encodings: all-zero hashes compare equal.
				require.True(t, h1.Eq(h2), "expected zero[%d] == zero[%d]", i, j)
			}
		}
	})

	t.Run("single byte hash exhaustive", func(t *testing.T) {
		// Test a representative sample of single-byte comparisons
		testValues := []byte{0x00, 0x01, 0x7f, 0x80, 0xfe, 0xff}
		for _, i := range testValues {
			for _, j := range testValues {
				h1 := Hash{i}
				h2 := Hash{j}

				result := h1.cmp(h2)
				expected := 0
				if i < j {
					expected = -1
				}
				if i > j {
					expected = 1
				}

				require.Equal(t, expected, result, "cmp(%#x, %#x)", i, j)
			}
		}
	})

	t.Run("lexical ordering verification", func(t *testing.T) {
		// Verify lexical byte-order semantics
		tests := []struct {
			h1       Hash
			h2       Hash
			expected int // -1 if h1 < h2, 0 if equal, 1 if h1 > h2
		}{
			{Hash{0x00}, Hash{0x01}, -1},
			{Hash{0x01}, Hash{0x00}, 1},
			{Hash{0x01}, Hash{0x01}, 0},
			{Hash{0x00, 0xff}, Hash{0x01, 0x00}, -1},
			{Hash{0xff}, Hash{0x00, 0x00}, 1},
			{Hash{0x01, 0x02}, Hash{0x01, 0x02}, 0},
		}

		for _, tt := range tests {
			result := tt.h1.cmp(tt.h2)
			require.Equal(t, tt.expected, result, "cmp(%v, %v)", tt.h1, tt.h2)
		}
	})
}
