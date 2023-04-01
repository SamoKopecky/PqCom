package common

import (
	"golang.org/x/crypto/sha3"
)

const n = 256

// func PMod(input, mod int) (output int) {
// 	temp := input % mod
// 	if temp > 0 && temp < mod {
// 		return temp
// 	}
// 	return (temp + mod) % mod
// }

func PMod(input, mod int) (output int) {
	return (input%mod + mod) % mod
}

func Kdf(input []byte, len int) (output []byte) {
	output = make([]byte, len)
	sha3.ShakeSum256(output, input)
	return
}

func BytesToBits(bytes []byte) (bits []byte) {
	var i, I, j int
	bits = make([]byte, len(bytes)*8)
	for i = 0; i < len(bytes); i++ {
		I = i * 8
		for j = 0; j < 8; j++ {
			bits[I+j] = byte((bytes[i] >> j) & 0x1)
		}
	}
	return
}

func PolyToBits(poly []int, coefSize int) (bits []byte) {
	var i, I, j int
	bits = make([]byte, n*coefSize)
	for i = 0; i < n; i++ {
		I = i * coefSize
		for j = 0; j < coefSize; j++ {
			bits[I+j] = byte((poly[i] >> j) & 0x1)
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
