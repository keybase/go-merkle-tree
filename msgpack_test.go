package merkletree

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecodeFromBytesLargeNodeCompatibility(t *testing.T) {
	// Large nodes were accepted before decoder hardening and must remain
	// readable. In particular, do not add a decode-only size limit which lets a
	// write commit successfully and then makes its root unreadable.
	value := make([]byte, 10*1024*1024+1)
	node := Node{
		Type: NodeTypeLeaf,
		Leafs: []KeyValuePair{
			{Key: Hash{0x01}, Value: value},
		},
	}
	encoded, err := encodeToBytes(node)
	require.NoError(t, err)

	var decoded Node
	require.NoError(t, decodeFromBytes(&decoded, encoded))
	require.Len(t, decoded.Leafs, 1)
	decodedValue, ok := decoded.Leafs[0].Value.([]byte)
	require.True(t, ok, "decoded value has type %T, want []byte", decoded.Leafs[0].Value)
	require.Len(t, decodedValue, len(value))
}

func TestDecodeClaimedHugeCollection(t *testing.T) {
	// array32 with 2^32-1 claimed elements and no contents. MaxInitLen must
	// prevent the header alone from causing an allocation proportional to the
	// claimed length.
	malicious := []byte{0xdd, 0xff, 0xff, 0xff, 0xff}
	var decoded []any
	require.Error(t, decodeFromBytes(&decoded, malicious))
}

func TestDecodeMaxDepth(t *testing.T) {
	// Test that deeply nested structures (beyond go-codec's default depth
	// limit of 100) are rejected.
	var nested any = "leaf"
	for range 101 {
		nested = []any{nested}
	}
	encoded, err := encodeToBytes(nested)
	require.NoError(t, err)

	var decoded any
	err = decodeFromBytes(&decoded, encoded)
	require.Error(t, err)
	require.Contains(t, strings.ToLower(err.Error()), "depth")
}

func TestDecodeLargeLeafCollection(t *testing.T) {
	// Test that collections larger than go-codec's default MaxInitLen (4096)
	// can still be decoded successfully.
	leafs := make([]KeyValuePair, 4097)
	for i := range leafs {
		leafs[i] = KeyValuePair{
			Key:   Hash{byte(i), byte(i >> 8)},
			Value: i,
		}
	}
	node := Node{Type: NodeTypeLeaf, Leafs: leafs}
	encoded, err := encodeToBytes(node)
	require.NoError(t, err)

	var decoded Node
	require.NoError(t, decodeFromBytes(&decoded, encoded))
	require.Len(t, decoded.Leafs, len(leafs))
}

func TestDeepEqualNoPanic(t *testing.T) {
	tests := []struct {
		name string
		i1   any
		i2   any
	}{
		{"both nil", nil, nil},
		{"one nil", nil, "test"},
		{"empty strings", "", ""},
		{"different types", "string", 123},
		{"nested structures", map[string]any{"a": []int{1, 2, 3}}, map[string]any{"a": []int{1, 2, 3}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			_ = deepEqual(tt.i1, tt.i2)
		})
	}
}
