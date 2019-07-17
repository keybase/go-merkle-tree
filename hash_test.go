package merkleTree

import "testing"

func TestHashLen(t *testing.T) {
	var tests = []struct {
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
		tt := tt
		t.Run(tt.desc, func(t *testing.T) {
			actual := tt.hash.Len()
			if actual != tt.expected {
				t.Errorf("(%s): expected %d, actual %d", tt.desc, tt.expected, actual)
			}

		})
	}

}
