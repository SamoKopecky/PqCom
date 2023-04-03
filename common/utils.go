package common

import (
	"golang.org/x/crypto/sha3"
)

const n = 256

func PMod(input, mod int) int {
	res := input % mod
	if res < 0 {
		return res + mod
	}
	return res
}

func Kdf(input []byte, len int) (output []byte) {
	output = make([]byte, len)
	sha3.ShakeSum256(output, input)
	return
}

func BytesToBits(bytes []byte) (bits []byte) {
	var i, I, j int
	var shiftedByte byte
	bits = make([]byte, len(bytes)*8)
	for i = 0; i < len(bytes); i++ {
		I = i << 3
		shiftedByte = bytes[i]
		for j = 0; j < 8; j++ {
			bits[I+j] = shiftedByte & 0x1
			shiftedByte >>= 1
		}
	}
	return
}

func PolyToBits(poly []int, coefSize int) (bits []byte) {
	var i, I, j, shiftedByte int
	bits = make([]byte, n*coefSize)
	for i = 0; i < n; i++ {
		I = i * coefSize
		shiftedByte = poly[i]
		for j = 0; j < coefSize; j++ {
			bits[I+j] = byte(shiftedByte) & 0x1
			shiftedByte >>= 1
		}
	}
	return
}

func BytesEqual(a, b []byte) (equal bool) {
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
