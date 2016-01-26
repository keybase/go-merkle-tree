
package merkleTree

import (
	"testing"
	"math/big"
	"bytes"
	"crypto/rand"
)

func slowBitslice(b []byte, begin int, end int) []byte {
	z := big.NewInt(0).SetBytes(b)
	z = z.Rsh(z, uint(len(b)*8-end))
	modulus := big.NewInt(1)
	modulus = modulus.Lsh(modulus, uint(end - begin))
	z = z.Mod(z, modulus)
	ret := z.Bytes()

	padlen := (end + 7 - begin)/8 - len(ret)
	pad := make([]byte, padlen)
	ret = append(pad, ret...)
	return ret
}

func testBitslice(t *testing.T, size, s1, s2, l1, l2 int) {
	buf := make([]byte, size)
	rand.Read(buf)
	for s := s1; s < s2; s++ {
		for l := l1; l < l2; l++ {
			a1 := slowBitslice(buf, s, s+l)
			a2 := bitslice(buf, ChildIndex(s), ChildIndex(s+l))
			if !bytes.Equal(a1, a2) {
				t.Fatalf("Failed on buf=%v; a1=%v; a2=%v; s=%d; l=%d", buf, a1, a2, s, l)
			}
		}
	}
}

func TestBitslice(t *testing.T) {
	testBitslice(t, 100, 0, 100, 0, 100)
}
