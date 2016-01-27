package merkleTree

import (
	"math/big"
)

func byteslice(h []byte, numBytes int, level int) []byte {
	start := numBytes * level
	end := numBytes * (level + 1)
	if start > len(h) {
		start = len(h)
	}
	if end > len(h) {
		end = len(h)
	}
	return h[start:end]
}

func bitslice(h []byte, numBits int, level int) Prefix {
	if numBits%8 == 0 {
		return byteslice(h, numBits/8, level)
	}

	// Set a begin and end bit index
	begin := numBits * level
	end := numBits * (level + 1)

	// No sense in using the large tail of the string, just
	// chop it off here.
	b := h[:((end + 7) >> 3)]

	z := big.NewInt(0).SetBytes(b)
	z = z.Rsh(z, uint(len(b)*8-end))
	modulus := big.NewInt(1)
	modulus = modulus.Lsh(modulus, uint(end-begin))
	z = z.Mod(z, modulus)
	ret := z.Bytes()

	padlen := (end+7-begin)/8 - len(ret)
	pad := make([]byte, padlen)
	ret = append(pad, ret...)
	return ret
}
