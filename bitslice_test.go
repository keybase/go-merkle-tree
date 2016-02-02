package merkleTree

import (
	"bytes"
	"testing"
)

func TestBitslice(t *testing.T) {

	data := []byte{
		0x6b, 0xa6, 0x8c, 0xee, 0x3e, 0x51, 0x63, 0x92, 0x94, 0xeb, 0xb4,
		0xb7, 0x51, 0xa9, 0x6f, 0x33, 0x92, 0xc4, 0x94, 0x26, 0x9c, 0x9d,
		0x00, 0xfe, 0x01, 0xa5,
	}

	testVectors := []struct {
		input   []byte
		output  []byte
		numBits int
		level   int
		index   uint32
	}{
		{data, []byte{0x8c, 0xee}, 16, 1, uint32(0x8cee)},
		{data, []byte{0xa6}, 8, 1, uint32(0xa6)},
		{data, []byte{0x0a}, 4, 2, uint32(0x0a)},
		{data, []byte{0x00, 0xe3}, 9, 3, uint32(0xe3)},
	}

	for i, v := range testVectors {
		computed, index := bitslice(v.input, v.numBits, v.level)
		if !bytes.Equal(v.output, computed) {
			t.Fatalf("failure in vector %d: got %v; wanted %v", i, computed, v.output)
		}
		if index != v.index {
			t.Fatalf("wrong index: %d != %d", index, v.index)
		}
	}

}
