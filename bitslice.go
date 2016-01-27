
package merkleTree

import (
	"math/big"
)

func byteslice(h Hash, numBytes int, level int) Prefix {
	start := numBytes*level
	end := numBytes*(level+1)
	if end > len(h) {
		end = len(h)
	}
	return Prefix(h[start:end])
}

func genericBitslice(h Hash, numBits ChildIndex, level Level) Prefix {
	return Prefix{}
}

func bitslice(h Hash, numBits ChildIndex, level Level) Prefix {
	if numBits % 8 == 0 {
		return byteslice(h, int(numBits / 8), int(level))
	}

	// Set a begin and end bit index
	begin := int(numBits)*int(level)
	end := int(numBits)*int(level + 1)

	// Only consider the bytes in question; round the left side up and
	// the right side down
	b := h[(begin >> 3):((end + 7) >> 3)]

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

